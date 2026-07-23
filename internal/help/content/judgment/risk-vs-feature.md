---
title: "Risk vs feature — is this actually a risk?"
summary: "A risk is a way the outcome could fail, not a feature you want, phrased as its absence."
---
# Risk vs feature

**The question.** You're about to add a risk. Is it really a risk, or a feature
you want, dressed up as "the thing that's missing"?

**The test.** A risk names a way the outcome could **fail** — an event or
condition with a plausible cause. If you can only state it as "we don't have
X yet", it's a missing feature, and the real element is a requirement.

**Good (risks).**

- "A reference points at something that does not resolve." (a failure that can
  happen)
- "Frequent edits leave linked elements out of sync." (a degradation with a
  cause)

**Bad (features wearing a risk costume).**

- "We don't have a search command." → that's a requirement, not a risk.
- "The tool should be fast." → an outcome/requirement, not a risk.

**Rule of thumb.** Risk → *requirement that mitigates it*. If you find yourself
writing the fix into the risk, split them: state the failure as the risk, state
the guard as a requirement with `mitigates` pointing at it.

**Failure mode.** A risk list that's really a backlog. It hides the actual
failure modes (so nothing is mitigated deliberately) and inflates scope with
wishes.
