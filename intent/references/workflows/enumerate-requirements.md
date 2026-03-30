# Enumerate requirements

**`ENUM-REQ`**

**In scope for this workflow:** edits to the outcome **`README.md`** (risks table and requirement one-liner tables). **Out of scope:** creating or editing files under **`requirements/`**—that is **Develop requirement** (`DEV-REQ`).

You are deriving **requirement one-liners** (and usually **risk rows**) on the **outcome README**. In practice that means asking how the outcome is achieved, asking **“how else?”** so you do not lock onto one design, and asking **why** each candidate requirement exists so it traces to the outcome. For risks: what could prevent the outcome from being true, naming each risk on the README and tying it to a **mitigating requirement** (existing or new). The goal is an **index** you trust before you open full requirement files.

**Records:** **PDRs** when framing requirements implies a **new** strategic choice. **CRs** when the **outcome README** changes in a meaningful way (including table work).

### Automation / agents

By default, propose one-liners and risk rows in review or in the outcome README only; do not create `requirements/` files unless the user explicitly asks to run **Develop requirement** or to materialize files.

---

### Draft

Stress-test the outcome: ask what must be true for success, what could go wrong, and whether each line traces back to the outcome. Lock **one-liners** (and risk links) on the outcome README before you invest in long requirement bodies. If the team makes a new strategic call while doing this, record a **PDR** at outcome or product scope as fits.

### AI-assisted review

Look for missing risks, orphan requirements, and requirements that sound like implementation smuggled into product language. Table edits that shift meaning earn **CRs** at outcome scope.

### Independent review

Confirm the index is enough to prioritize and build from. **PDRs** only for decisions that surface at the gate; **CRs** if the reviewer’s edits are still substantial.
