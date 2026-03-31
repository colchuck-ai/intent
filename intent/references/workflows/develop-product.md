# Develop product

**`DEV-PROD`**

You are establishing **`docs/product/README.md`** as a **product statement**: who the product serves, what value means, and the boundaries of the problem space. **In scope:** that product root README only. **Out of scope:** outcome one-liners (**Enumerate outcomes**, `ENUM-OUT`) and outcome folders under `docs/product/outcomes/` (**Develop outcome**, `DEV-OUT`).

Pick **one** root template—[`product-root-simple`](../../assets/product-root-simple.md) or [`product-root-jtbd`](../../assets/product-root-jtbd.md)—and stay with it for the life of that tree. If you use JTBD, optional background lives in [`jtbd-primer`](../guides/jtbd-primer.md); the template itself is allowed to omit sections that do not help.

**Records:** this is where **new product or strategic decisions** surface as **PDRs** when there is a real decision—not for ordinary drafting. There is usually **no CR** until the file has been meaningfully edited in a later pass; the first time you might still emit a **CR** if you are replacing placeholder content with substantive text in one jump.

---

### Draft

Settle enough surrounding context (charter, repo, constraints) that the README is not floating in a vacuum. Write the product root so a stranger understands **intent** without reading engineering yet. If a genuine strategic fork appears, capture it as a **PDR** at product scope under `docs/product/pdrs/`.

### AI-assisted review

Read the draft as someone who was not in the room. Check that language stays **product-level** (avoid sneaking implementation in as requirements). If review surfaces **new** strategic decisions, add **PDRs**; routine tightening rarely needs one.

### Independent review

A human who did not author the draft decides whether the story holds. Same bar for **PDRs**: new decisions, not polish. If this pass changes meaning in a material way, a **CR** at product root scope (`docs/product/crs/`, `PROD-CR…` naming) records what moved.
