# Modify requirement

**`MOD-REQ`**

You are **changing** an existing requirement. **Ripple** can reach the **parent outcome README** (index rows, risks), **`docs/engineering/README.md`** (table rows, links), and **component** docs where behavior or scope shifts.

**Context span:** the requirement artifact, its **outcome**, the **engineering root** table, and any **components** linked to that requirement. If product language and engineering reality diverge, fix both in the same narrative.

**Records:** **CRs** at requirement scope (and outcome or `ENG-CR…` when parents move). **PDRs** if the change introduces a **new** strategic decision (rarer than CRs, but real when scope or promise shifts).

---

### Draft

Edit with traceability: update indices and tables when the requirement’s meaning moves. Avoid leaving stale links from engineering back to the wrong text. Keep the requirement text in **[plain language](../guides/plain-language.md)**—plain words, short sentences—so the change is easy to follow.

### AI-assisted review

Check outcome alignment, engineering **table** consistency, and component specs. If the requirement text still points at implementation, keep the **README table** as the canonical mapping list—do not grow a second registry inside the requirement. Tighten wording for plain readability ([`plain-language.md`](../guides/plain-language.md)). **CRs** for substantive edits; **PDRs** for new product decisions.

### Independent review

Human judgment on whether the change is coherent end-to-end and reads clearly for a junior engineer without re-reading ([`plain-language.md`](../guides/plain-language.md)). **CRs** and **PDRs** as needed.
