---
title: "ref: job"
summary: "Fields and shape of a job element."
---
# ref: job

Address: `product.jobs.<job>`. Always its own document (no `type`).

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `story` | required | JTBD narrative: "When … I want … so I can …" |
| `detail` | optional | freeform markdown |
| `outcomes` | map | child outcomes |

No outgoing edges. Add with:

```
intent add job product <key> --name "…" --story "…"
```
