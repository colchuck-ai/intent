---
title: "Altitude — is this at the right level of detail?"
summary: "Each element type sits at a fixed altitude; don't smuggle low-level detail into a high-level element or vice versa."
---
# Altitude

**The question.** Is this element pitched at the right level — or is it a
low-level detail wearing a high-level type (or the reverse)?

**The test.** Read the type's job and check the prose matches it:

- **job** — highest: a durable user goal.
- **outcome** — a measurable result under a job.
- **requirement** — a condition that must hold for an outcome.
- **component** — an engineering unit with one responsibility.

Each level explains *why* the level below it exists. If an element answers a
question from a different level, its altitude is off.

**Good.** Outcome: "Minimize the time to trace an element to its reason."
Requirement under it: "Every element must have a stable logical ID." (the
requirement is one rung down, concrete but not implementation)

**Bad.**

- A requirement that says "Use a hash map keyed by address." (implementation —
  too low; that's a component's `behavior` or an ADR)
- A component whose `responsibility` is "Make the product good." (too high — an
  outcome, not a charter)

**Also: does it deserve its own page?** `type: document` is altitude in the
output. Promote only when the element has enough structured content to stand
alone; otherwise keep it `inline`.

**Failure mode.** Detail sinking down or floating up until types stop meaning
anything — a "requirement" you can't test, a "component" you can't build.
