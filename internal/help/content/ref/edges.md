---
title: "ref: edges"
summary: "The declared edge kinds, what they connect, and how they're stored."
---
# ref: edges

Edges are the trace graph — **only what you declare here**. Prose links never
create edges. Every edge target is a dotted-path address; a target that doesn't
resolve is `E001`.

| Edge | On | Points to | Storage |
|---|---|---|---|
| `mitigates` | requirement | a risk in the same outcome | flat list |
| `dependsOn` | requirement | requirement leaves, anywhere | flat list |
| `fulfills` | component | a product requirement | flat list |
| `affects` | pdr / adr / cr | what the record touches | flat list |
| `relationships` | component | another element (with a note) | keyed map `addr: note` |

- A whole-domain target is the bare root literal (`product` / `engineering`).
- Flat lists dedup structurally — the same target twice is `E003`.
- `relationships` carries a payload (the note), so it stays a keyed map.
- Record `affects` is domain-scoped: PDR ⊆ `product.*`, ADR ⊆ `engineering.*`
  (CR exempt); a breach is `E004`.

Manage edges with `intent link <kind> <from> <to>` and `intent unlink` — never
hand-build the list.
