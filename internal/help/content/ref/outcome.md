---
title: "ref: outcome"
summary: "Fields and shape of an outcome element."
---
# ref: outcome

Address: `product.jobs.<job>.outcomes.<outcome>`.

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `statement` | required | declarative sentence (a measurable result) |
| `type` | optional | `inline` \| `document` (default `inline`) |
| `detail` | optional | freeform markdown |
| `risks` | map | child risks |
| `requirements` | map | child requirements |

No outgoing edges. Add with:

```
intent add outcome <job> <key> --name "…" --statement "…"
```
