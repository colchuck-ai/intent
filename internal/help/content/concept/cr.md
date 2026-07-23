---
title: "CR — change record"
summary: "A record of a modification event; the one record type allowed to cross domains."
---
# CR — change record

A **CR** records a change — an atomic modification event to the intent tree,
with what changed and why. It is the one record type that may **cross domains**:
a single change can affect a product requirement and an engineering component at
once.

- **Lives under:** `change_records.<key>` (top-level, cross-cutting).
- **Required edge:** `affects` (≥1) — may **mix** domains and levels freely, and
  may include a bare root literal (`product` / `engineering`).
- **Required fields:** `change` (what changed) and `rationale` (why).
- Carries no `kind`, no `type`, no `summary` — location implies it's a CR.

**Why it's top-level and cross-cutting.** Change events don't decompose by
domain the way decisions do — a change is one indivisible event. Splitting it by
domain would lie about it, so the CR stays whole and lives above both trees.

**Get it right:** `intent help judgment:worth-recording` and
`intent help guide:record-decision`.
