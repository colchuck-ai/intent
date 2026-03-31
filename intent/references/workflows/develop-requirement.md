# Develop requirement

**`DEV-REQ`**

You are writing the **full requirement** under `docs/product/outcomes/…/requirements/`—simple file or complex folder with `README.md`—using [`requirement-simple`](../../assets/requirement-simple.md) or [`requirement-complex`](../../assets/requirement-complex.md). Capability, acceptance, and edge cases live here. Work on the product side is chunked **by outcome**, then **by requirement**; one pass of this workflow develops **one** requirement artifact (simple or complex shape as fits). **In scope:** that requirement’s body only. **Out of scope:** sibling requirements, the engineering root index, or component specs (**Enumerate components**, **Develop component**).

**Records:** **PDRs** for **new** product decisions scoped to this requirement (`…/pdrs/` or a parent `pdrs/` if the decision spans more). **CRs** for meaningful edits under this requirement’s `crs/`.

---

### Draft

Bring the outcome README, the one-liner row for this requirement, and any peer requirements you might conflict with. Write so implementation could be tested against acceptance without leaking full design that belongs in engineering. Keep the wording **[plain language](../guides/plain-language.md)**—short sentences, common words, and no clever or academic gloss—so a junior engineer can follow the requirement without re-reading. If a new strategic fork appears, capture a **PDR** at the right scope.

### AI-assisted review

Check product language: it should stay plain and traceable ([`plain-language.md`](../guides/plain-language.md))—wired back to the outcome, consistent with sibling requirements. On **traceability to engineering**, treat **`docs/engineering/README.md`** as the place where requirement ↔ component links live in full; this requirement can name components when it helps the reader, but it should not become a parallel registry of rows. **CRs** when the draft session or this pass moves substantive text.

### Independent review

Human sign-off on whether the requirement is buildable and bounded—and that mapping discipline still holds if the text touched implementation edges. The human reader should also confirm the text reads clearly for a junior engineer ([`plain-language.md`](../guides/plain-language.md)): if anything still sounds performative or needlessly dense, fix it in this pass. **PDRs** for late strategic recognition; **CRs** for remaining meaningful edits.
