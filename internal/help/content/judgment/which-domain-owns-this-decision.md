---
title: "Which domain owns this decision?"
summary: "Product decisions (the what) are PDRs; engineering decisions (the how) are ADRs; a decision is never both."
---
# Which domain owns this decision?

**The question.** You're recording a decision. Is it a **PDR** (product), an
**ADR** (engineering), or — the trap — "both"?

**The test.** Ask what the decision is *about*:

- About **what** to build or why it matters to a user → **PDR**, under
  `product.decision_records`, `affects` targets under `product.*`.
- About **how** to build it (technology, structure, mechanism) → **ADR**, under
  `engineering.decision_records`, `affects` targets under `engineering.*`.

A decision is never both. What looks cross-domain is almost always a product
decision *plus* an engineering decision, joined by a `fulfills` edge between the
elements they touch.

**Good.**

- PDR: "Prefer explicit declaration over inference." (a product stance)
- ADR: "Ship a single Go binary instead of a Python script." (a build choice)

**Bad.**

- One record that decides both the user-facing behavior and the tech stack —
  split it into a PDR and an ADR.
- A PDR whose `affects` names an `engineering.components.*` target → that's
  `E004`.

**When it's really a change, not a decision.** A modification that touches both
domains at once is a **change record** (CR) — the one record type allowed to
cross domains.

**Failure mode.** Cross-domain records that blur what/how, so neither side owns
the reasoning and the trace graph can't carry the real dependency. See
`intent help E004`.
