---
title: "Job vs activity — is this a job or a step?"
summary: "A job is a stable goal a user is trying to accomplish; an activity is a step or tactic toward it."
---
# Job vs activity

**The question.** You're adding a job. Is it a real jobs-to-be-done goal, or
just one activity/step on the way to a goal?

**The test.** A job is **stable and solution-independent** — it stays true even
if the tool changes completely. An activity is a step you'd perform; it
disappears when the solution changes.

**Good (jobs).**

- "Understand the rationale behind an element."
- "Change documentation without drift."

**Bad (activities dressed as jobs).**

- "Open the element view." (a step)
- "Click validate." (a UI action)
- "Run the linter nightly." (a tactic)

**Rule of thumb.** Phrase it as *"When \<situation\>, I want \<motivation\> so I
can \<outcome\>."* If the motivation is really a mechanism ("so I can click the
button"), you have an activity. Activities, if they matter, live as requirements
or components under the outcome — not as jobs.

**Failure mode.** A job list that reads like a click-path. It ties intent to
today's solution, so the tree churns every time the implementation moves.
