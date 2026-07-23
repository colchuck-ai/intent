---
title: "Requirement vs task — condition, or work item?"
summary: "A requirement is a condition that must hold; a task is work to do. The tree holds requirements, not a backlog."
---
# Requirement vs task

**The question.** You're adding a requirement. Is it a **condition that must be
true**, or a **task** (a unit of work to schedule)?

**The test.** A requirement is a testable state of the world — you could write
an acceptance check for it. A task is an action with a doer and a "done"; it
belongs in an issue tracker, not the intent tree.

**Good (requirements).**

- "Every declared reference must resolve to a real element."
- "Intent must assign every element a stable logical ID."

**Bad (tasks).**

- "Write the reference resolver." (work to do)
- "Add tests for the validator." (a backlog item)

**Rule of thumb.** Write it as *"\<subject\> must \<condition\>."* If it reads
naturally as *"\<verb\> the \<thing\>"* (an imperative to a developer), it's a
task. A requirement should also `mitigate` a real risk in its outcome — if you
can't name the risk, question whether the requirement is real.

**Failure mode.** A requirements list that's actually a to-do list. It goes
stale the moment the work is done, and it buries the durable conditions among
transient chores.
