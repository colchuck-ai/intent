---
title: "ref: requirement"
summary: "Fields and edges of a requirement element."
---
# ref: requirement

Address: `product.jobs.<job>.outcomes.<outcome>.requirements.<req>`.

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `statement` | required | "\<subject\> must \<condition\>." |
| `mitigates` | required, ≥1 | dotted-path list — risks in the same outcome |
| `dependsOn` | optional | dotted-path list — requirement leaves, anywhere |
| `acceptance_criteria` | optional | list of strings |
| `type` | optional | `inline` \| `document` |
| `detail` | optional | freeform markdown |

Add with (edge flags `--mitigates` / `--depends-on` are repeatable):

```
intent add requirement <outcome> <key> --name "…" --statement "…" \
  --mitigates <risk-suffix>
```

Manage edges after creation with `intent link` / `intent unlink`.
