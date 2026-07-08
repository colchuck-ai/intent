# T11 Verification Report — intent skill review suite

Date: 2026-07-07  
Skill path: `intent/`  
Description length: 872 chars (post-T9)

## 1. validate_skill.py

```bash
python3 ~/.cursor/skills/agent-skill-review/scripts/validate_skill.py intent/
```

**Result: PASS** — `passed: true`, zero errors, zero warnings.

## 2. Cold-maintainer test (fresh subagent, read-only)

Subagent read: `SKILL.md`, `references/elements.md`, `references/structure.md`, `references/workflows.md`, `scripts/lint_intent_tree.py`, `assets/templates/cr-document.md`.

| Question | Result |
|---|---|
| Q1 — purpose & when it fires | Correct |
| Q2a — CR title placeholders | `assets/templates/cr-document.md` H1 | Correct |
| Q2b — jobs have no template | `references/structure.md#naming` | Correct |
| Q2c — `## Dependencies` meaning | `references/elements.md` Requirement section | Correct |
| Q3 — why "Create produces no CR" | Correct rationale (CR = modification log) | Correct |
| Q4 — `lint_intent_tree.py` purpose | Correct (7 checks, JSON stdout, Judge first pass) | Correct |

**Q5 residual (self-discounted — out of assigned read scope):**

- `context-tracing.md` not in minimum read set (referenced but not assigned)
- Inline job form lives in `product.md` (not in minimum read set)
- Tier 2 CR wording in `cr-directory.md` (not in minimum read set)

**Q5 residual (actionable legibility gaps — file as follow-ups):**

1. **Trivial vs material CR threshold** — workflows/elements say trivial edits need no CR but give no rubric.
2. **Judge checklist scatter** — draft/judge/review steps point to rubrics across `elements.md` rather than one enumerated checklist.
3. **Backlink author spec** — linter heuristics documented in code; author-facing format for every link type not fully specified in references.

**Cold-maintainer verdict: PASS** (Q1–Q4 correct; Q5 contains only self-discounted observations plus three minor legibility gaps suitable for follow-up beads).

## 3. Activation eval (T9 description, 20 labeled queries)

Method: labeled queries (~50% should-trigger, near-miss-heavy). Judged whether the current `description` frontmatter would cause correct skill selection.

Precision on near-misses (should **not** trigger): **10/10**  
Recall on should-trigger set: **10/10**

| # | Query | Label | Judgment |
|---|---|---|---|
| 1 | Write a product spec for checkout under docs/product/ | trigger | ✓ anchored to tree |
| 2 | Define the customer's job for O001 in our intent docs | trigger | ✓ |
| 3 | Draft an ADR for our caching decision in docs/engineering/ | trigger | ✓ |
| 4 | Log this change to requirement O001-R002 with a CR | trigger | ✓ |
| 5 | Where does this outcome doc go / what should I name it? | trigger | ✓ |
| 6 | Trace requirement O001-R003 to its component | trigger | ✓ |
| 7 | Check these intent docs for coherence | trigger | ✓ |
| 8 | Run the draft/judge/review workflow on this requirement | trigger | ✓ |
| 9 | What template for a Tier 2 outcome directory? | trigger | ✓ |
| 10 | Review alignment between parent outcome and child requirements | trigger | ✓ |
| 11 | Write a product spec (no docs/tree context) | no-trigger | ✓ T9 anchor requires tree context |
| 12 | Define the customer's job for a sales pitch deck | no-trigger | ✓ not intent-tree work |
| 13 | Update the README install section | no-trigger | ✓ explicit skip |
| 14 | Fix the auth login bug in src/ | no-trigger | ✓ code change skip |
| 15 | Write marketing copy for the landing page | no-trigger | ✓ explicit skip |
| 16 | Jot a personal note that won't link into the intent tree | no-trigger | ✓ ad-hoc note skip |
| 17 | Write a product spec for API endpoint documentation (OpenAPI) | no-trigger | ✓ API docs ≠ intent framework |
| 18 | Record this decision — pick lodash vs ramda in package.json | no-trigger | ✓ code dependency, not ADR/PDR/CR |
| 19 | Capture why we chose X in Slack, not in docs/ | no-trigger | ✓ no intent-tree placement |
| 20 | Draft an ADR for our office lunch policy | no-trigger | ✓ ADR named but not engineering/product intent |

**Activation verdict: PASS** — near-miss precision holds; generic triggers remain anchored to intent-tree context per T9.

## Summary

| Check | Verdict |
|---|---|
| validate_skill.py | PASS |
| Cold-maintainer (Q1–Q4 + Q5) | PASS (3 minor follow-ups) |
| Activation eval (20 queries) | PASS |

Recommended follow-up beads (optional, Minor):

- Define trivial-vs-material CR rubric
- Consolidate or index Judge-loop checklist
- Document author-facing backlink conventions beyond linter heuristics
