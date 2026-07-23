---
title: "guide: draft → judge → review"
summary: "Draft an element quickly, judge it against the right test, then review coherence before finalizing."
---
# guide: draft → judge → review

**Goal:** produce elements that are not just valid but *good*. The CLI
guarantees validity; this loop adds judgment.

**1. Draft.** Get it down with `intent add …`. Don't over-polish — you're
capturing intent, not prose. Default `type: inline`; promote later only if it
earns a page (`intent help judgment:altitude`).

**2. Judge.** Apply the test for what you just wrote (the add footer names it):

| You added a… | Test |
|---|---|
| job | `judgment:job-vs-activity` |
| outcome | `judgment:outcome-vs-solution` |
| risk | `judgment:risk-vs-feature` |
| requirement | `judgment:requirement-vs-task` |
| component | `judgment:altitude` |
| any record | `judgment:worth-recording` |

Also sanity-check the prose shape: `judgment:writing-statements`.

**3. Review.** Zoom out and check the change hangs together with its neighbors —
`intent show`, `intent affects`, then `intent help judgment:coherence`. Fix, or
`intent rm` and redraft, until the neighborhood tells one story.

Nothing here is enforced — it's the discipline that keeps a valid tree also
worth reading.
