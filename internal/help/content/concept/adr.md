---
title: "ADR — architecture decision record"
summary: "A record of an engineering (how) decision; its affects targets stay under engineering.*"
---
# ADR — architecture decision record

An **ADR** records an engineering decision — a choice about *how* to build
something, including technology choices — with the options weighed and the
consequences accepted. (A technology choice *is* an ADR; there is no separate
`technology_choices` field.)

- **Lives under:** `engineering.decision_records.<key>`.
- **Required edge:** `affects` (≥1) — targets **must** resolve under
  `engineering.*` (domain-pure; a stray non-engineering target is `E004`).
- **Required field:** `summary`.
- **Optional fields:** `context`, `options` (list), `decision`,
  `consequences` (list).
- Carries no `kind`, no `type`, no `detail` — location implies it's an ADR.

**Why it's domain-pure.** The *how* belongs to engineering. A decision that
seems to span domains is really a product decision plus an engineering one,
linked through the elements they touch.

**Get it right:** `intent help judgment:which-domain-owns-this-decision` and
`intent help guide:record-decision`.
