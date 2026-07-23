---
title: "Requirement"
summary: "A condition that must hold to secure an outcome; mitigates at least one risk."
---
# Requirement

A **requirement** is a condition that must be true for its outcome to hold. It
is testable — you could write an acceptance check for it — and it exists to
mitigate a risk.

- **Prose field:** `statement` — *"\<subject\> must \<condition\>."*
- **Lives under:** `product.jobs.<job>.outcomes.<outcome>.requirements.<req>`.
- **Required edge:** `mitigates` (≥1) — a risk in the same outcome.
- **Optional edges:** `dependsOn` — other requirement leaves, anywhere.
- **Optional field:** `acceptance_criteria` (list).
- **Promotable:** optional `type: inline | document`.
- **Fulfilled by** engineering components (via their `fulfills` edge).

**Why it exists.** Requirements are the durable conditions that connect a
user's outcome to the engineering that satisfies it — the hinge of the trace
graph.

**Get it right:** `intent help judgment:requirement-vs-task` (a requirement is a
condition, not a to-do) and `intent help guide:add-requirement`.
