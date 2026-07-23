---
title: "ref: risk"
summary: "Fields and shape of a risk element."
---
# ref: risk

Address: `product.jobs.<job>.outcomes.<outcome>.risks.<risk>`. Always inline (no
`type`), no outgoing edges.

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `statement` | required | a failure that could happen |
| `detail` | optional | freeform markdown |

Requirements in the same outcome point back at a risk via `mitigates`. Add with:

```
intent add risk <outcome> <key> --name "…" --statement "…"
```
