---
title: "guide: edit an element"
summary: "Find it, read it, change it, confirm it — the safe edit loop."
---
# guide: edit an element

**Goal:** change an existing element without breaking anything that points at
it.

1. **Find it.** `intent find <words>` or `intent tree` → the address (a suffix
   is fine).
2. **Read it first.** `intent show <addr>` — its prose and outgoing edges.
   `intent affects <addr>` — who points back at it (what your change might
   ripple to).
3. **Change a field.** `intent set <addr> <field> "<value>"` — e.g.
   `intent set stable-logical-ids statement "…"`. Edges are not fields — use
   `link`/`unlink` (see `guide:add-requirement`).
4. **The write is guarded.** `set` runs `validate` first and re-serializes
   canonically; an edit that would dangle a reference or break containment is
   refused with an `E0NN` code.
5. **Confirm coherence.** Re-read with `intent show <addr>` and check the
   neighborhood still hangs together — `intent help judgment:coherence`.

Renaming the key or moving the element is a different job — see `guide:move`
(`intent mv`), which cascades every referring edge for you.
