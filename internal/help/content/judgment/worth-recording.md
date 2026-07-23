---
title: "Worth recording — does this decision/change need a record?"
summary: "Record a decision or change only when it chose among real options with lasting consequences someone will later ask about."
---
# Worth recording

**The question.** You're about to `record` a PDR, ADR, or CR. Is this worth a
permanent record, or is it routine?

**The test.** Record it if a future reader would ask *"why is it this way?"* and
be worse off without the answer. That happens when there was a **real choice
among options** with **consequences that outlive the moment**.

**Good.**

- ADR: "Ship a single Go binary instead of a Python script." (options existed;
  consequences persist)
- CR: "Added the resolvable-references requirement and wired the validator to
  it." (a real change to intent, with a rationale)

**Bad.**

- "Renamed a variable." (no lasting consequence, no option worth weighing)
- "Fixed a typo in a statement." (routine edit — just make it)

**Decision vs change.** A **decision** record (PDR/ADR) captures a choice and
its reasoning. A **change** record (CR) captures a modification event and why it
happened. If nothing was *chosen* and nothing *changed* that others must trace,
don't record anything.

**Failure mode.** A record log full of noise. The signal decisions — the ones
you actually revisit — drown, and people stop reading the log at all.
