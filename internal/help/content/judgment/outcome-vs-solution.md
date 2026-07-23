---
title: "Outcome vs solution — result, or how you'd get it?"
summary: "An outcome is a measurable result the user wants; a solution is a way to achieve it."
---
# Outcome vs solution

**The question.** You're adding an outcome. Does it name a **result** the user
wants, or a **solution** for producing that result?

**The test.** An outcome is measurable and solution-free: it names a direction
(minimize / increase / ensure) on something the user feels, without saying how.
If it names a feature, mechanism, or UI, it's a solution.

**Good (outcomes).**

- "Minimize the time to trace an element back to its reason."
- "Minimize the likelihood that linked elements fall out of sync."

**Bad (solutions).**

- "Add a backlinks panel." (a feature — the *how*)
- "Cache the trace graph in memory." (a mechanism)

**Rule of thumb.** Outcomes start with a metric verb — *minimize, maximize,
reduce, ensure* — applied to something the user experiences. The solution then
lives below it as requirements and components that fulfill those requirements.

**Failure mode.** Outcomes that pin intent to one design. When the design
changes the outcome looks wrong, even though the user's actual goal never moved.
