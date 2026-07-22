# Prototype: structured source + generated docs

A throwaway spike answering: *can a schema'd data model replace the tier / backlink /
relative-path bookkeeping that makes the Intent skill unwieldy?*

Answer from this spike: **yes, for structure and edges — not for prose or judgment.**

## What's here

```
schema/intent.schema.json   the contract (edges by ID; NO tiers/slugs/paths/maps in it)
model/*.yaml                the source of truth — the SAME content as ../docs, hand-written
intentgen.py                validate (schema + semantic) -> render markdown
out/docs/...                generated output
```

## Run it

```bash
python3 intentgen.py            # validate + generate into out/
python3 intentgen.py --check    # validate only
```

## What the spike proves

1. **Round-trips the real docs.** `out/` is byte-for-byte identical to `../docs/`
   (`diff` clean on all 4 files) — so nothing readable is lost.

2. **Every edge is declared once; the rest is derived.** The `model/` source contains
   zero Risk-Requirement Maps, zero Requirement-Component Maps, zero `See Also` blocks,
   and zero relative `.md` links. All of them are generated from three edge fields:
   `requirement.mitigates`, `component.fulfills`, `record.affects`. Backlink
   bidirectionality is impossible to break because backlinks aren't stored.

3. **The 486-line linter mostly evaporates.** Format/uniqueness/no-reuse → schema
   regex. "Requirement without a risk" → `minItems: 1` on `mitigates` (structurally
   impossible). Empty headings, CR-location, slug/body consistency → the generator
   can't emit them wrong. What's left is a ~40-line semantic check (do ID references
   resolve). Try it: break a `mitigates`/`fulfills` target and run `--check`.

4. **Tier is a render decision, not authored state.** Add a `body:` to a component and
   it auto-promotes to its own file with a generated reference block + correct relative
   link in the parent. No file move, no path edit, no backlink edit.

## What it deliberately does NOT solve

- **Prose stays prose.** Job stories, decision rationale, component behavior are free
  text — kept as plain string fields, not decomposed. (The story regex is a loose
  sanity anchor, not a straitjacket.)
- **Judgment stays human/LLM.** The schema enforces "a requirement names a risk"; it
  cannot tell a real risk from a disguised feature request. The elements.md rubric and
  the draft/judge loop survive unchanged.

## Cost this introduces

Requires the generator + validator to exist and run. The current pure-markdown design
degrades gracefully (hand-write a doc, it still reads); this does not. For a large tree,
keep the model sharded per element (references are IDs, not paths) to preserve selective
context-loading — do NOT collapse to one monolithic YAML (bad merges, bad review, loads
everything).
