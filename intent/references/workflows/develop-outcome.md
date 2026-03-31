# Develop outcome

**`DEV-OUT`**

You are creating **`docs/product/outcomes/O<NNN>-<slug>/`** with a real **outcome document** (typically `README.md` in a complex folder), using the [`outcome`](../../assets/outcome.md) shape. The product tree is chunked **by outcome**; one pass of this workflow develops **one** outcome. At this stage you may **omit** the risk and requirement tables—or leave them empty—if you have not run **Enumerate requirements** yet. The outcome can be **partial**: problem, success, boundaries, and story first; the index of requirements can arrive in a later workflow. **In scope:** this outcome’s folder and README. **Out of scope:** populating risk and requirement **tables** on that README (**Enumerate requirements**, `ENUM-REQ`) and any files under **`requirements/`** (**Develop requirement**, `DEV-REQ`).

**Records:** **PDRs** for **new** product decisions at **outcome** scope (`outcomes/…/pdrs/`), or at product root if the decision is bigger than one outcome. **CRs** for meaningful edits to **this outcome’s** narrative under that outcome’s `crs/`.

---

### Draft

Load the product root and the one-liner for this outcome so the detail **pulls upward** to the same language. Build the outcome README so someone could explain **why this outcome exists** without reading requirements yet. If you already know major risks, you can sketch them; it is fine to defer full risk→requirement mapping to **Enumerate requirements**.

### AI-assisted review

Check alignment with the product root and sibling outcomes (peek at other outcome READMEs if boundaries might blur). Strategic surprises become **PDRs**; substantive edits to this README become **CRs** at outcome scope.

### Independent review

Human judgment on whether the outcome is coherent and scoped. **PDRs** for late strategic clarity; **CRs** if this review still changes the document in a traceable way.
