---
title: "Coherence — does this edit still hang together?"
summary: "After any change, check that the element and everything it links still tell one consistent story."
---
# Coherence

**The question.** After you add or change an element, does the local
neighborhood still make sense together — parent to child, edge to target,
prose to declared edges?

**The test.** Read the element and each thing it links to, in one sitting. If a
new reader would ask "wait, how do these relate?", it isn't coherent yet.

**Good.** A requirement `mitigates` a risk that is genuinely about the same
outcome; its `statement` describes the same concern the risk raises. A
component `fulfills` a requirement whose wording it plausibly satisfies.

**Bad.** A requirement mitigates a risk from an unrelated outcome. A
component's `responsibility` never touches the requirement it claims to fulfill.
Prose mentions another element (`{{ ... }}`) as a real dependency, but no edge
declares it.

**How to check it fast.**

- `intent show <addr>` — the element and its outgoing edges.
- `intent affects <addr>` — who points back at it; do those still make sense?
- `intent trace <addr>` — the wider neighborhood when you're unsure.

**Failure mode.** Edits that are each locally valid but drift apart over time,
so the tree validates yet no longer describes one coherent intent. Validity is
the CLI's job; coherence is yours.
