---
title: "guide: record a decision or change"
summary: "Decide the record type by domain, then record it with its affects targets."
---
# guide: record a decision or change

**Goal:** capture a decision or change so a future reader knows why.

1. **Is it worth recording?** Only if there was a real choice or a real change
   with lasting consequence — `intent help judgment:worth-recording`.
2. **Pick the type** — `intent help judgment:which-domain-owns-this-decision`:
   - product *what/why* → **PDR** (`affects` ⊆ `product.*`)
   - engineering *how* → **ADR** (`affects` ⊆ `engineering.*`)
   - a modification event, possibly cross-domain → **CR**
3. **Record it:**
   ```
   # a product decision
   intent record pdr <key> --name "…" --summary "…" \
     --affects <product-target> --option "…" --decision "…" --consequence "…"

   # an engineering decision
   intent record adr <key> --name "…" --summary "…" --affects <eng-target>

   # a change (may cross domains)
   intent record cr <key> --name "…" --change "…" --rationale "…" \
     --affects <target> --affects <another-target>
   ```
4. **Domain scope is enforced.** A PDR/ADR whose `affects` leaves its domain is
   `E004` — route a cross-domain link through the trace graph or a CR.
