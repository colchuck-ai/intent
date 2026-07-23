---
title: "ref: cr"
summary: "Fields of a change record."
---
# ref: cr

Address: `change_records.<key>` (top-level). No `kind`, `type`, or `summary`.

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `affects` | required, ≥1 | dotted-path list — may **mix** domains; root literals allowed |
| `change` | required | what changed |
| `rationale` | required | why it changed |

A CR is the only record whose `affects` may cross domains (no `E004`). Create
with:

```
intent record cr <key> --name "…" --change "…" --rationale "…" \
  --affects <target> --affects <another-target>
```
