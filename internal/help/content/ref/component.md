---
title: "ref: component"
summary: "Fields and edges of a component element."
---
# ref: component

Address: `engineering.components.<component>`.

| Field | Req? | Shape |
|---|---|---|
| `name` | required | display text |
| `responsibility` | required | one-line charter |
| `fulfills` | required, ≥1 | dotted-path list — product requirements |
| `relationships` | optional | keyed map `<address>: note` |
| `data_model` | optional | structured prose |
| `interfaces` | optional | structured prose |
| `behavior` | optional | structured prose (may interpolate `{{ }}`) |
| `edge_cases` | optional | list of strings |
| `success_criteria` | optional | list of strings |
| `type` | optional | `inline` \| `document` |
| `detail` | optional | freeform markdown |

Add with:

```
intent add component engineering <key> --name "…" --responsibility "…" \
  --fulfills <req-suffix>
```
