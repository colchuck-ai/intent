---
title: "ref: adr"
summary: "Fields of an architecture decision record."
---
# ref: adr

Address: `engineering.decision_records.<key>`. No `kind`, `type`, or `detail`.
A technology choice is an ADR (there is no `technology_choices` field).

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `affects` | required, ≥1 | dotted-path list — **must** be under `engineering.*` |
| `summary` | required | the decision's core statement |
| `context` | optional | prose |
| `options` | optional | list of strings |
| `decision` | optional | prose |
| `consequences` | optional | list of strings |

Out-of-domain `affects` is `E004`. Create with:

```
intent record adr <key> --name "…" --summary "…" --affects <engineering-target>
```
