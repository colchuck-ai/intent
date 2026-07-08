#!/usr/bin/env python3
# /// script
# requires-python = ">=3.10"
# dependencies = []
# ///
"""Lint intent-tree documents for mechanical coherence and backlink rules.

Scans markdown under a docs/ root (default: docs/) and reports violations for:
  - logical-ID format, uniqueness, and slug/body consistency
  - backlink bidirectionality between linked intent documents
  - empty headings (structure.md "Never emit an empty heading")
  - requirements with no linked risk (feature-request smell)
  - change records stored on the wrong product/engineering tree side

JSON report on stdout; human-readable diagnostics on stderr.

Usage:
  python3 scripts/lint_intent_tree.py [DOCS_ROOT]
  python3 scripts/lint_intent_tree.py --help

Exit codes:
  0  no errors
  1  one or more lint violations
  2  usage error / docs root missing
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from dataclasses import asdict, dataclass, field
from pathlib import Path
from typing import Iterable

# Match longest/specific IDs first (risk before requirement).
ID_PATTERNS: list[tuple[str, re.Pattern[str]]] = [
    ("risk", re.compile(r"\b(O\d{3,}-RSK\d{3,})\b")),
    ("requirement", re.compile(r"\b(O\d{3,}-R\d{3,})\b")),
    ("outcome", re.compile(r"\b(O\d{3,})\b")),
    ("job", re.compile(r"\b(J\d{3,})\b")),
    ("component", re.compile(r"\b(C\d{3,})\b")),
    ("pdr", re.compile(r"\b(PDR\d{3,})\b")),
    ("adr", re.compile(r"\b(ADR\d{3,})\b")),
    ("cr", re.compile(r"\b(CR\d{3,})\b")),
]

STRICT_ID = {
    "job": re.compile(r"^J\d{3,}$"),
    "outcome": re.compile(r"^O\d{3,}$"),
    "risk": re.compile(r"^O\d{3,}-RSK\d{3,}$"),
    "requirement": re.compile(r"^O\d{3,}-R\d{3,}$"),
    "component": re.compile(r"^C\d{3,}$"),
    "pdr": re.compile(r"^PDR\d{3,}$"),
    "adr": re.compile(r"^ADR\d{3,}$"),
    "cr": re.compile(r"^CR\d{3,}$"),
}

SUSPICIOUS_ID = re.compile(
    r"\b(?:O\d+-R\d+|O\d+-RSK\d+|J\d+|C\d+|PDR\d+|ADR\d+|CR\d+)\b"
)

LINK_RE = re.compile(r"\[([^\]]+)\]\(([^)]+)\)")
HEADING_RE = re.compile(r"^(#{1,6})\s+(.+?)\s*$")
H1_ID_RE = re.compile(
    r"^(?:"
    r"(?P<risk>O\d{3,}-RSK\d{3,})"
    r"|(?P<requirement>O\d{3,}-R\d{3,})"
    r"|(?P<outcome>O\d{3,})"
    r"|(?P<component>C\d{3,})"
    r"|(?P<pdr>PDR\d{3,})"
    r"|(?P<adr>ADR\d{3,})"
    r"|(?P<cr>CR\d{3,})"
    r")\s*-\s*.+$"
)

PRODUCT_ID = re.compile(r"\b(?:J|O\d{3,}|O\d{3,}-RSK\d{3,}|O\d{3,}-R\d{3,}|PDR\d{3,})\b")
ENGINEERING_ID = re.compile(r"\b(?:C\d{3,}|ADR\d{3,})\b")

FILENAME_PATTERNS: list[tuple[str, re.Pattern[str]]] = [
    ("outcome", re.compile(r"^product/outcomes/(O\d{3,})-[^/]+(?:/README\.md|\.md)$")),
    (
        "requirement",
        re.compile(
            r"^product/outcomes/O\d{3,}-[^/]+/requirements/(R\d{3,})-[^/]+(?:/README\.md|\.md)$"
        ),
    ),
    ("component", re.compile(r"^engineering/components/(C\d{3,})-[^/]+(?:/README\.md|\.md)$")),
    ("pdr", re.compile(r"^product/drs/(PDR\d{3,})-[^/]+(?:/README\.md|\.md)$")),
    ("adr", re.compile(r"^engineering/drs/(ADR\d{3,})-[^/]+(?:/README\.md|\.md)$")),
    ("cr", re.compile(r"^product/crs/(CR\d{3,})-[^/]+(?:/README\.md|\.md)$")),
    ("cr", re.compile(r"^engineering/crs/(CR\d{3,})-[^/]+(?:/README\.md|\.md)$")),
]


@dataclass
class Violation:
    check: str
    severity: str  # error | warning
    file: str
    message: str
    detail: dict = field(default_factory=dict)


@dataclass
class Doc:
    path: Path
    rel: str
    text: str
    h1_id: str | None
    h1_kind: str | None
    links: list[tuple[str, Path]]  # (link_text, resolved_path)
    ids: set[str]


def eprint(*args: object) -> None:
    print(*args, file=sys.stderr)


def classify_id(token: str) -> str | None:
    for kind, pattern in ID_PATTERNS:
        if STRICT_ID[kind].fullmatch(token):
            return kind
    return None


def extract_ids(text: str) -> set[str]:
    """Extract logical IDs, preferring longer tokens (risk before outcome)."""
    found: set[str] = set()
    for _, pattern in ID_PATTERNS:
        found.update(pattern.findall(text))
    return found


def canonical_owner_id(doc: Doc) -> str | None:
    """Return the logical ID this file canonically owns, if any."""
    if doc.h1_id:
        return doc.h1_id
    rel = doc.rel.replace("\\", "/")
    for kind, pattern in FILENAME_PATTERNS:
        m = pattern.match(rel)
        if not m:
            continue
        token = m.group(1)
        if kind == "requirement":
            # Requirement files carry R<NNN> in the path; full ID is on the H1.
            continue
        return token
    return None


def filename_number_keys(doc: Doc) -> list[tuple[str, str]]:
    rel = doc.rel.replace("\\", "/")
    keys: list[tuple[str, str]] = []
    for kind, pattern in FILENAME_PATTERNS:
        m = pattern.match(rel)
        if m:
            keys.append((kind, m.group(1)))
    return keys


def parse_h1(text: str) -> tuple[str | None, str | None]:
    for line in text.splitlines():
        if line.startswith("# ") and not line.startswith("##"):
            m = H1_ID_RE.match(line[2:].strip())
            if not m:
                return None, None
            for kind in ("risk", "requirement", "outcome", "component", "pdr", "adr", "cr"):
                if m.group(kind):
                    return m.group(kind), kind
            return None, None
    return None, None


def is_skippable_line(line: str) -> bool:
    s = line.strip()
    return not s or s.startswith("<!--") or s.endswith("-->") or s == "<!--"


def section_content(lines: list[str], start: int, level: int) -> list[str]:
    body: list[str] = []
    for line in lines[start:]:
        hm = HEADING_RE.match(line)
        if hm and len(hm.group(1)) <= level:
            break
        body.append(line)
    return body


def has_meaningful_content(body: Iterable[str]) -> bool:
    for line in body:
        if is_skippable_line(line):
            continue
        if line.strip():
            return True
    return False


def resolve_link(source: Path, target: str, docs_root: Path) -> Path | None:
    if target.startswith(("http://", "https://", "mailto:")):
        return None
    raw = (source.parent / target.split("#", 1)[0]).resolve()
    try:
        rel = raw.relative_to(docs_root.resolve())
    except ValueError:
        return None
    if raw.is_dir():
        raw = raw / "README.md"
    if not raw.exists() and raw.with_suffix(".md").exists():
        raw = raw.with_suffix(".md")
    if raw.exists():
        return docs_root / rel if raw.is_file() else docs_root / rel
    # Normalize to docs-relative file path when possible.
    candidate = docs_root / rel
    if candidate.exists():
        return candidate
    if candidate.with_suffix(".md").exists():
        return candidate.with_suffix(".md")
    return None


def collect_docs(docs_root: Path) -> list[Doc]:
    docs: list[Doc] = []
    for path in sorted(docs_root.rglob("*.md")):
        text = path.read_text(encoding="utf-8")
        rel = str(path.relative_to(docs_root))
        h1_id, h1_kind = parse_h1(text)
        links: list[tuple[str, Path]] = []
        for link_text, target in LINK_RE.findall(text):
            resolved = resolve_link(path, target, docs_root)
            if resolved is not None:
                links.append((link_text, resolved))
        docs.append(
            Doc(
                path=path,
                rel=rel,
                text=text,
                h1_id=h1_id,
                h1_kind=h1_kind,
                links=links,
                ids=extract_ids(text),
            )
        )
    return docs


def check_logical_ids(docs: list[Doc], docs_root: Path) -> list[Violation]:
    violations: list[Violation] = []
    owner_locations: dict[str, list[str]] = {}
    slug_numbers: dict[tuple[str, str], set[str]] = {}

    for doc in docs:
        for token in SUSPICIOUS_ID.findall(doc.text):
            if classify_id(token) is None:
                violations.append(
                    Violation(
                        check="logical_id_format",
                        severity="error",
                        file=doc.rel,
                        message=f"Malformed logical ID '{token}' (expected zero-padded segments per structure.md#naming)",
                    )
                )

        if doc.h1_id and doc.h1_kind and not STRICT_ID[doc.h1_kind].fullmatch(doc.h1_id):
            violations.append(
                Violation(
                    check="logical_id_format",
                    severity="error",
                    file=doc.rel,
                    message=f"H1 logical ID '{doc.h1_id}' does not match expected {doc.h1_kind} format",
                )
            )

        owner = canonical_owner_id(doc)
        if owner:
            owner_locations.setdefault(owner, []).append(doc.rel)

        for key in filename_number_keys(doc):
            slug_numbers.setdefault(key, set()).add(doc.rel)

    for logical_id, locations in sorted(owner_locations.items()):
        unique = sorted(set(locations))
        if len(unique) > 1:
            violations.append(
                Violation(
                    check="logical_id_uniqueness",
                    severity="error",
                    file=unique[0],
                    message=f"Logical ID '{logical_id}' is canonical owner in multiple documents",
                    detail={"locations": unique},
                )
            )

    for key, locations in slug_numbers.items():
        if len(locations) > 1:
            kind, num = key
            violations.append(
                Violation(
                    check="logical_id_no_reuse",
                    severity="error",
                    file=sorted(locations)[0],
                    message=f"Number {num} for {kind} reused across distinct paths (retired numbers must not be reused)",
                    detail={"locations": sorted(locations)},
                )
            )

    return violations


def check_empty_headings(docs: list[Doc]) -> list[Violation]:
    violations: list[Violation] = []
    for doc in docs:
        lines = doc.text.splitlines()
        for idx, line in enumerate(lines):
            hm = HEADING_RE.match(line)
            if not hm:
                continue
            level = len(hm.group(1))
            body = section_content(lines, idx + 1, level)
            if not has_meaningful_content(body):
                violations.append(
                    Violation(
                        check="empty_heading",
                        severity="error",
                        file=doc.rel,
                        message=f"Empty heading '{hm.group(2).strip()}' (omit empty sections per structure.md)",
                        detail={"line": idx + 1},
                    )
                )
    return violations


def check_requirement_risks(docs: list[Doc]) -> list[Violation]:
    violations: list[Violation] = []
    for doc in docs:
        is_requirement = "/requirements/" in doc.rel.replace("\\", "/") or (
            doc.h1_kind == "requirement"
        )
        if not is_requirement:
            continue
        mitigates = re.search(r"^## Mitigates\s*$", doc.text, re.MULTILINE)
        risk_ids = {i for i in doc.ids if classify_id(i) == "risk"}
        if mitigates:
            start = mitigates.end()
            section = doc.text[start : doc.text.find("\n## ", start)]
            risk_ids = {i for i in extract_ids(section) if classify_id(i) == "risk"}
        if not risk_ids:
            violations.append(
                Violation(
                    check="requirement_without_risk",
                    severity="error",
                    file=doc.rel,
                    message="Requirement has no linked risk in ## Mitigates (feature-request smell)",
                )
            )
    return violations


def check_cr_tree_side(docs: list[Doc]) -> list[Violation]:
    violations: list[Violation] = []
    for doc in docs:
        rel = doc.rel.replace("\\", "/")
        if "/crs/" not in rel:
            continue
        body = doc.text
        has_product = bool(PRODUCT_ID.search(body))
        has_engineering = bool(ENGINEERING_ID.search(body))
        in_product = rel.startswith("product/")
        in_engineering = rel.startswith("engineering/")
        if in_product and has_engineering and not has_product:
            violations.append(
                Violation(
                    check="cr_tree_side",
                    severity="error",
                    file=doc.rel,
                    message="CR under docs/product/crs/ references engineering elements only; move to docs/engineering/crs/",
                )
            )
        if in_engineering and has_product and not has_engineering:
            violations.append(
                Violation(
                    check="cr_tree_side",
                    severity="error",
                    file=doc.rel,
                    message="CR under docs/engineering/crs/ references product elements only; move to docs/product/crs/",
                )
            )
    return violations


def check_backlinks(docs: list[Doc], docs_root: Path) -> list[Violation]:
    violations: list[Violation] = []
    by_path = {d.path.resolve(): d for d in docs}

    def doc_ids(doc: Doc) -> set[str]:
        ids = set(doc.ids)
        if doc.h1_id:
            ids.add(doc.h1_id)
        return ids

    for src in docs:
        for link_text, target in src.links:
            tgt = by_path.get(target.resolve())
            if tgt is None:
                continue
            if tgt.path.resolve() == src.path.resolve():
                continue
            tgt_ids = doc_ids(tgt)
            src_ids = doc_ids(src)
            forward = extract_ids(link_text) or (tgt_ids & src_ids)
            if not forward and tgt.h1_id:
                forward = {tgt.h1_id}
            if not forward:
                continue
            # Reverse: target should mention source logical ID or link back.
            reverse_ids = tgt_ids
            reverse_links = {p.resolve() for _, p in tgt.links}
            if src.path.resolve() in reverse_links:
                continue
            if src.h1_id and src.h1_id in reverse_ids:
                continue
            if src_ids & reverse_ids:
                continue
            violations.append(
                Violation(
                    check="backlink_bidirectionality",
                    severity="error",
                    file=src.rel,
                    message=(
                        f"Link to '{tgt.rel}' missing reverse backlink "
                        f"(expected reference to '{src.h1_id or src.rel}' in target See Also/Affects)"
                    ),
                    detail={"target": tgt.rel, "link_text": link_text},
                )
            )
    return violations


def run_lint(docs_root: Path) -> dict:
    docs = collect_docs(docs_root)
    violations: list[Violation] = []
    violations.extend(check_logical_ids(docs, docs_root))
    violations.extend(check_empty_headings(docs))
    violations.extend(check_requirement_risks(docs))
    violations.extend(check_cr_tree_side(docs))
    violations.extend(check_backlinks(docs, docs_root))

    errors = [v for v in violations if v.severity == "error"]
    warnings = [v for v in violations if v.severity == "warning"]

    for v in violations:
        eprint(f"{v.severity.upper()} [{v.check}] {v.file}: {v.message}")

    report = {
        "docs_root": str(docs_root),
        "files_scanned": len(docs),
        "summary": {"errors": len(errors), "warnings": len(warnings)},
        "violations": [asdict(v) for v in violations],
        "checks": {
            "logical_id_format": True,
            "logical_id_uniqueness": True,
            "logical_id_no_reuse": True,
            "backlink_bidirectionality": True,
            "empty_heading": True,
            "requirement_without_risk": True,
            "cr_tree_side": True,
        },
    }
    for v in violations:
        report["checks"][v.check] = False
    return report


def main(argv: list[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        description="Lint intent-tree documents for mechanical coherence and backlink rules."
    )
    parser.add_argument(
        "docs_root",
        nargs="?",
        default="docs",
        help="Root directory containing product/ and engineering/ trees (default: docs)",
    )
    args = parser.parse_args(argv)

    docs_root = Path(args.docs_root)
    if not docs_root.is_dir():
        eprint(f"error: docs root not found: {docs_root}")
        return 2

    report = run_lint(docs_root.resolve())
    print(json.dumps(report, indent=2))
    return 1 if report["summary"]["errors"] else 0


if __name__ == "__main__":
    sys.exit(main())
