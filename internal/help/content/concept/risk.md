---
title: "Risk"
summary: "A way an outcome could fail — an event or condition, not a missing feature."
---
# Risk

A **risk** names a way its outcome could fail: an event or condition with a
plausible cause. It is what requirements exist to guard against.

- **Prose field:** `statement` — the failure, stated as something that can
  happen ("A reference points at something that does not resolve.").
- **Lives under:** `product.jobs.<job>.outcomes.<outcome>.risks.<risk>`.
- **Always inline** — risks carry no `type` and declare no outgoing edges.
- **Referenced by** the requirements that `mitigate` them (same outcome).

**Why it exists.** Naming failure modes explicitly means each requirement can
point at the risk it addresses, so you can see what is and isn't guarded.

**Get it right:** `intent help judgment:risk-vs-feature` (a risk is a failure,
not a wish for a feature).
