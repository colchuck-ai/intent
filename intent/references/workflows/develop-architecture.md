# Develop architecture

**`DEV-ARCH`**

You are writing **`docs/engineering/README.md`** as a **high-level architectural story**: how the system hangs together, main flows, boundaries, and principles. **In scope:** that engineering root README and its narrative. **Out of scope:** component one-liners and the **components table**—those belong in **Enumerate components** (`ENUM-COMP`). It is fine for the README to gesture at major pieces in prose before they have stable `C###` rows.

**Strong recommendation:** have **outcomes and requirement indices** in good shape first so architecture answers real intent rather than inventing it.

**Records:** **ADRs** for **new** architectural decisions at engineering-root scope (`docs/engineering/adrs/`). **CRs** are less common on a first pass unless you are replacing large chunks of placeholder text in one go—use **`ENG-CR…`** at engineering root when the README change is meaningful and needs traceability.

---

### Draft

Read enough product context that the engineering story does not contradict intent. Draft narrative and, if useful, a light diagram path. When a **new** technical decision appears (stack, boundary, pattern), prefer a short **ADR** rather than burying the decision only in prose.

### AI-assisted review

Check for clarity, missing boundaries, and decisions that should be **ADRs**. Meaningful README edits get **CRs** if not already recorded.

### Independent review

Human read for coherence and fit with product. **ADRs** for decisions surfaced late; **CRs** if the review still shifts the engineering root materially.
