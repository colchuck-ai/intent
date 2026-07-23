---
title: "PDR — product decision record"
summary: "A record of a product (what/why) decision; its affects targets stay under product.*"
---
# PDR — product decision record

A **PDR** records a product decision — a choice about *what* to build or *why*
it matters — with the options weighed and the consequences accepted.

- **Lives under:** `product.decision_records.<key>`.
- **Required edge:** `affects` (≥1) — targets **must** resolve under `product.*`
  (domain-pure; a stray non-product target is `E004`).
- **Required field:** `summary` (the decision's core).
- **Optional fields:** `context`, `options` (list), `decision`,
  `consequences` (list).
- Carries no `kind`, no `type`, no `detail` — location implies it's a PDR.

**Why it's domain-pure.** Reasoning decomposes by domain: a product decision is
the *what*, an engineering decision the *how*. Keeping them separate keeps each
record honest and lets the trace graph carry any cross-domain link.

**Get it right:** `intent help judgment:which-domain-owns-this-decision`,
`intent help judgment:worth-recording`, and `intent help guide:record-decision`.
