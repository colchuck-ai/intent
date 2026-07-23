---
title: "ref: pdr"
summary: "Fields of a product decision record."
---
# ref: pdr

Address: `product.decision_records.<key>`. No `kind`, `type`, or `detail`.

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `affects` | required, ≥1 | dotted-path list — **must** be under `product.*` |
| `summary` | required | the decision's core statement |
| `context` | optional | prose |
| `options` | optional | list of strings |
| `decision` | optional | prose |
| `consequences` | optional | list of strings |

Out-of-domain `affects` is `E004`. Create with:

```
intent record pdr <key> --name "…" --summary "…" --affects <product-target>
```
