---
title: "guide: trace context"
summary: "Pick the cheapest command that answers your question: show, affects, or trace."
---
# guide: trace context

**Goal:** understand how an element connects to the rest of the tree. Three
commands form a cost/scope gradient — use the cheapest that answers you.

| Question | Command | Scope |
|---|---|---|
| "What is this?" | `intent show <addr>` | the element + its outgoing edges |
| "What breaks if I change this?" | `intent affects <addr>` | incoming refs (backlinks) |
| "Give me the context." | `intent trace <addr>` | the full neighborhood |

**`trace` knobs** (priced by how much you ask for):

- `--depth N` — how many hops out (default is the immediate neighborhood).
- `--up` — only upstream (things this depends on / mitigates / fulfills).
- `--down` — only downstream (things pointing at this).
- `--edges <kind>` — restrict to one edge kind, e.g. `--edges dependsOn`.

Start with `show`; reach for `affects` before a risky edit; reach for `trace`
only when you need the wider picture.
