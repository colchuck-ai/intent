# Develop requirement

**`DEV-REQ`**

You are writing the **full requirement** under `docs/product/outcomes/…/requirements/`—simple file or complex folder with `README.md`—using [`requirement-simple`](../../assets/requirement-simple.md) or [`requirement-complex`](../../assets/requirement-complex.md). Capability, acceptance, and edge cases live here. Work on the product side is chunked **by outcome**, then **by requirement**; one pass of this workflow develops **one** requirement artifact (simple or complex shape as fits). **In scope:** that requirement’s body only. **Out of scope:** sibling requirements, the engineering root index, or component specs (**Enumerate components**, **Develop component**).

**Records:** **PDRs** for **new** product decisions scoped to this requirement (`…/pdrs/` or a parent `pdrs/` if the decision spans more). **CRs** for meaningful edits under this requirement’s `crs/`.

---

### Draft

Bring the outcome README, the one-liner row for this requirement, and any peer requirements you might conflict with. Write so implementation could be tested against acceptance without leaking full design that belongs in engineering. If a new strategic fork appears, capture a **PDR** at the right scope.

### AI-assisted review

Check product language, trace back to the outcome, and consistency with sibling requirements. On **traceability to engineering**, treat **`docs/engineering/README.md`** as the place where requirement ↔ component links live in full; this requirement can name components when it helps the reader, but it should not become a parallel registry of rows. **CRs** when the draft session or this pass moves substantive text.

### Independent review

Human sign-off on whether the requirement is buildable and bounded—and that mapping discipline still holds if the text touched implementation edges. **PDRs** for late strategic recognition; **CRs** for remaining meaningful edits.
