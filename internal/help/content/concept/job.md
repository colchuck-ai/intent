---
title: "Job"
summary: "A durable, solution-independent goal a user is trying to accomplish."
---
# Job

A **job** is a goal a user is trying to get done — stated the jobs-to-be-done
way, independent of any particular solution. It is the top of the product tree:
outcomes, risks, and requirements all hang beneath it.

- **Prose field:** `story` — *"When \<situation\>, I want \<motivation\> so I can
  \<expected outcome\>."*
- **Lives under:** `product.jobs.<job>`.
- **Contains:** `outcomes` (each with its own `risks` and `requirements`).
- **Always its own document** — jobs carry no `type`.

**Why it exists.** Anchoring intent to a stable user goal means the tree
survives redesigns: the *how* can change completely while the job stays put.

**Get it right:** `intent help judgment:job-vs-activity` (a job is a goal, not a
step) and `intent help judgment:writing-statements` (the `story` shape).
