#!/usr/bin/env python3
"""Intent generator spike: structured model -> validated -> rendered markdown.

Reads all *.yaml under model/, validates against schema/intent.schema.json plus a
thin semantic layer, then GENERATES the docs/ tree. Everything the current
framework asks an author to hand-maintain is derived here instead:

  - backlinks (See Also / inline record lists)  <- from record.affects
  - Risk-Requirement Map                         <- from requirement.mitigates
  - Requirement-Component Map                    <- from component.fulfills
  - relative link paths                          <- from element location
  - slugs                                        <- from id + name
  - tier (inline vs own file)                    <- from presence of body content

Usage:
  python3 intentgen.py [--check] [MODEL_DIR] [OUT_DIR]

  --check   validate only; do not write files
Exit codes: 0 ok, 1 validation error, 2 usage error.
"""
import json
import os
import re
import sys

try:
    import yaml
except ImportError:
    sys.stderr.write("PyYAML required: pip install pyyaml\n")
    sys.exit(2)

try:
    import jsonschema
except ImportError:
    jsonschema = None

HERE = os.path.dirname(os.path.abspath(__file__))


# --------------------------------------------------------------------------- #
# Load + validate
# --------------------------------------------------------------------------- #
def load_model(model_dir):
    model = {"jobs": [], "components": [], "records": []}
    for fname in sorted(os.listdir(model_dir)):
        if not fname.endswith((".yaml", ".yml")):
            continue
        with open(os.path.join(model_dir, fname)) as fh:
            data = yaml.safe_load(fh) or {}
        for key, val in data.items():
            if isinstance(model.get(key), list) and isinstance(val, list):
                model[key].extend(val)
            else:
                model[key] = val
    return model


def slugify(text):
    s = re.sub(r"[^a-z0-9]+", "-", text.lower()).strip("-")
    return re.sub(r"-{2,}", "-", s)


def semantic_errors(model):
    errs = []
    risk_by_outcome = {}
    req_ids = set()
    req_name = {}

    for job in model.get("jobs", []):
        for o in job.get("outcomes", []):
            oid = o["id"]
            risk_by_outcome[oid] = {r["id"] for r in o.get("risks", [])}
            for req in o.get("requirements", []):
                q = "%s-%s" % (oid, req["id"])
                req_ids.add(q)
                req_name[q] = req["name"]
                for rsk in req.get("mitigates", []):
                    if rsk not in risk_by_outcome[oid]:
                        errs.append("%s mitigates unknown risk %s (not under %s)" % (q, rsk, oid))

    comp_ids = {c["id"] for c in model.get("components", [])}
    for c in model.get("components", []):
        for q in c.get("fulfills", []):
            if q not in req_ids:
                errs.append("%s fulfills unknown requirement %s" % (c["id"], q))
        for rel in c.get("relationships", []):
            if rel["to"] not in comp_ids:
                errs.append("%s relationship to unknown component %s" % (c["id"], rel["to"]))

    known = {"product", "architecture"} | comp_ids | req_ids
    for job in model.get("jobs", []):
        for o in job.get("outcomes", []):
            known.add(o["id"])
    for rec in model.get("records", []):
        for a in rec.get("affects", []):
            if a not in known:
                errs.append("%s affects unresolvable element %s" % (rec["id"], a))
    return errs, req_name


def validate(model):
    errs = []
    if jsonschema is not None:
        with open(os.path.join(HERE, "schema", "intent.schema.json")) as fh:
            schema = json.load(fh)
        v = jsonschema.Draft202012Validator(schema)
        for e in sorted(v.iter_errors(model), key=lambda e: list(e.path)):
            loc = "/".join(str(p) for p in e.path) or "<root>"
            errs.append("schema: %s: %s" % (loc, e.message))
    else:
        errs.append("(jsonschema not installed — schema check skipped)")
    sem, req_name = semantic_errors(model)
    return errs + sem, req_name


# --------------------------------------------------------------------------- #
# Render
# --------------------------------------------------------------------------- #
def relpath(from_file, to_file):
    return os.path.relpath(to_file, os.path.dirname(from_file))


def records_for(model, element_token, kinds):
    out = []
    for rec in model.get("records", []):
        if element_token in rec.get("affects", []) and rec["kind"] in kinds:
            out.append(rec)
    return out


def rec_path(rec):
    slug = rec.get("slug") or slugify(rec["name"])
    fname = "%s-%s.md" % (rec["id"], slug)
    if rec["kind"] == "pdr":
        return "docs/product/drs/" + fname
    if rec["kind"] == "adr":
        return "docs/engineering/drs/" + fname
    return "docs/crs/" + fname


REC_HEADING = {"pdr": "Product Decision Records", "adr": "Architectural Decision Records", "cr": "Change Records"}


def render_product(model, files):
    p = model["product"]
    L = ["# %s" % p["name"], "", p["summary"], "", "## Jobs"]
    for job in model["jobs"]:
        L += ["", "### %s" % job["name"], "", "> %s" % job["story"]]
        for o in job.get("outcomes", []):
            oid = o["id"]
            L += ["", "#### %s - %s" % (oid, o["name"]), "", o["statement"]]
            risks = o.get("risks", [])
            reqs = o.get("requirements", [])
            if risks:
                L += ["", "**Risks**", ""]
                for r in risks:
                    L.append("- **%s-%s** - %s: %s" % (oid, r["id"], r["name"], r["statement"]))
            if reqs:
                L += ["", "**Requirements**", ""]
                for req in reqs:
                    L.append("- **%s-%s** - %s: %s" % (oid, req["id"], req["name"], req["statement"]))
            # Risk-Requirement Map — DERIVED from each requirement's `mitigates`.
            map_lines = []
            for r in risks:
                hits = ["%s-%s - %s" % (oid, req["id"], req["name"])
                        for req in reqs if r["id"] in req.get("mitigates", [])]
                if hits:
                    map_lines.append("- **%s-%s - %s**: %s" % (oid, r["id"], r["name"], ", ".join(hits)))
            if map_lines:
                L += ["", "**Risk-Requirement Map**", ""] + map_lines
    # See Also — DERIVED from records that affect the product.
    see = []
    for kind in ("pdr", "cr"):
        recs = records_for(model, "product", {kind})
        if recs:
            see += ["", "### %s" % REC_HEADING[kind], ""]
            for rec in recs:
                see.append("- [%s - %s](%s)" % (rec["id"], rec["name"],
                            relpath("docs/product/README.md", rec_path(rec))))
    if see:
        L += ["", "## See Also"] + see
    files["docs/product/README.md"] = "\n".join(L) + "\n"


def render_engineering(model, files):
    a = model["architecture"]
    L = ["# %s" % a["title"], "", a["summary"]]
    for label, key in (("Principles", "principles"), ("Constraints", "constraints"),
                       ("Technology Choices", "technologyChoices")):
        items = a.get(key, [])
        if items:
            L += ["", "## %s" % label, ""] + ["- %s" % i for i in items]

    comps = model.get("components", [])
    comp_name = {c["id"]: c["name"] for c in comps}
    L += ["", "## Components"]
    for c in comps:
        promoted = bool(c.get("body"))
        L += ["", "### %s - %s" % (c["id"], c["name"]), "", c["responsibility"]]
        if promoted:
            # Tier is a RENDER decision: body present -> own file + reference block.
            slug = "%s-%s" % (c["id"], slugify(c["name"]))
            path = "docs/engineering/components/%s.md" % slug
            L.append("")
            L.append("See [%s - %s](%s)." % (c["id"], c["name"],
                     relpath("docs/engineering/README.md", path)))
            render_component_file(c, comp_name, model, files, path)
            continue
        rels = c.get("relationships", [])
        if rels:
            L += ["", "**Relationships**", ""]
            for rel in rels:
                L.append("- **%s - %s**: %s" % (rel["to"], comp_name[rel["to"]], rel["note"]))
        for kind in ("adr", "cr"):
            recs = records_for(model, c["id"], {kind})
            if recs:
                L += ["", "**%s**" % REC_HEADING[kind], ""]
                for rec in recs:
                    L.append("- [%s - %s](%s)" % (rec["id"], rec["name"],
                             relpath("docs/engineering/README.md", rec_path(rec))))

    # Requirement-Component Map — DERIVED from each component's `fulfills`.
    _, req_name = semantic_errors(model)
    ordered_reqs = []
    for job in model["jobs"]:
        for o in job.get("outcomes", []):
            for req in o.get("requirements", []):
                ordered_reqs.append("%s-%s" % (o["id"], req["id"]))
    map_lines = []
    for q in ordered_reqs:
        hits = ["%s - %s" % (c["id"], c["name"]) for c in comps if q in c.get("fulfills", [])]
        if hits:
            map_lines.append("- **%s - %s**: %s" % (q, req_name[q], ", ".join(hits)))
    if map_lines:
        L += ["", "## Requirement-Component Map", ""] + map_lines
    files["docs/engineering/README.md"] = "\n".join(L) + "\n"


def render_component_file(c, comp_name, model, files, path):
    L = ["# %s - %s" % (c["id"], c["name"]), "", c["responsibility"]]
    b = c.get("body", {})
    if b.get("dataModel"):
        L += ["", "## Data model", "", b["dataModel"]]
    if b.get("interfaces"):
        L += ["", "## Interfaces", "", b["interfaces"]]
    if b.get("behavior"):
        L += ["", "## Behavior", "", b["behavior"]]
    if b.get("edgeCases"):
        L += ["", "## Edge cases", ""] + ["- %s" % e for e in b["edgeCases"]]
    rels = c.get("relationships", [])
    if rels:
        L += ["", "## Relationships", ""]
        for rel in rels:
            L.append("- **%s - %s**: %s" % (rel["to"], comp_name[rel["to"]], rel["note"]))
    if b.get("successCriteria"):
        L += ["", "## Success criteria", ""] + ["- %s" % s for s in b["successCriteria"]]
    files[path] = "\n".join(L) + "\n"


def render_records(model, files):
    for rec in model.get("records", []):
        L = ["# %s - %s" % (rec["id"], rec["name"]), "", rec["summary"]]
        if rec["kind"] in ("pdr", "adr"):
            if rec.get("context"):
                L += ["", "## Context", "", rec["context"]]
            if rec.get("options"):
                L += ["", "## Options", ""] + ["- %s" % o for o in rec["options"]]
            if rec.get("decision"):
                L += ["", "## Decision", "", rec["decision"]]
            if rec.get("consequences"):
                L += ["", "## Consequences", ""] + ["- %s" % c for c in rec["consequences"]]
        else:  # cr
            if rec.get("change"):
                L += ["", "## Change", "", rec["change"]]
            if rec.get("rationale"):
                L += ["", "## Rationale", "", rec["rationale"]]
            if rec.get("affectsNarrative"):
                L += ["", "## Affects", "", rec["affectsNarrative"]]
        files[rec_path(rec)] = "\n".join(L) + "\n"


def main(argv):
    args = [a for a in argv if not a.startswith("--")]
    check_only = "--check" in argv
    model_dir = args[0] if len(args) > 0 else os.path.join(HERE, "model")
    out_dir = args[1] if len(args) > 1 else os.path.join(HERE, "out")

    model = load_model(model_dir)
    errs, _ = validate(model)
    if errs:
        sys.stderr.write("VALIDATION FAILED:\n")
        for e in errs:
            sys.stderr.write("  - %s\n" % e)
        return 1
    sys.stderr.write("validation: OK\n")
    if check_only:
        return 0

    files = {}
    render_product(model, files)
    render_engineering(model, files)
    render_records(model, files)

    for rel, content in sorted(files.items()):
        dest = os.path.join(out_dir, rel)
        os.makedirs(os.path.dirname(dest), exist_ok=True)
        with open(dest, "w") as fh:
            fh.write(content)
        sys.stderr.write("wrote %s\n" % rel)
    sys.stderr.write("generated %d files into %s\n" % (len(files), out_dir))
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
