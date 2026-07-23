---
title: "ref: addressing"
summary: "The dotted-path address, suffix resolution, and how the CLI echoes addresses back."
---
# ref: addressing

**One address encoding — the dotted path — everywhere.** An element's key is its
identity, its path segment, and its slug. Keys are snake_case and only locally
unique; position qualifies them.

```
product.jobs.<job>.outcomes.<outcome>.risks.<risk>
product.jobs.<job>.outcomes.<outcome>.requirements.<req>
engineering.components.<component>
product.decision_records.<key>            # PDR
engineering.decision_records.<key>        # ADR
change_records.<key>                        # CR
```

**Suffix input.** Commands accept the **shortest unambiguous suffix**:
`intent show stable_logical_ids` resolves the full path. If a suffix matches more
than one element, the CLI lists the candidates so you can disambiguate.

**Full addresses out.** The CLI always echoes the full dotted path back, and
anything persisted (prose `{{ }}`, edge lists) uses the full path — suffixes are
input ergonomics only.

Find an address with `intent find <words>` or `intent tree`.
