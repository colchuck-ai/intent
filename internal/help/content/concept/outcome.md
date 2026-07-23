---
title: "Outcome"
summary: "A measurable result a user wants from a job, stated without a solution."
---
# Outcome

An **outcome** is a measurable result the user wants from a job — a direction on
something they experience, with no mention of how you'd achieve it.

- **Prose field:** `statement` — a declarative sentence, usually starting with a
  metric verb ("Minimize the time to ...").
- **Lives under:** `product.jobs.<job>.outcomes.<outcome>`.
- **Contains:** `risks` (how it could fail) and `requirements` (conditions that
  secure it).
- **Promotable:** optional `type: inline | document` (default `inline`).

**Why it exists.** Outcomes separate the result the user cares about from the
design that delivers it, so the design can change without rewriting intent.

**Get it right:** `intent help judgment:outcome-vs-solution` and
`intent help judgment:altitude`.
