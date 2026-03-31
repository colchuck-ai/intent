# Enumerate requirements

**`ENUM-REQ`**

You are deriving **requirement one-liners** (and usually **risk rows**) on the **outcome README**. **In scope:** edits to that README’s tables—risks and requirement one-liners—only. **Out of scope:** creating or editing files under **`requirements/`**—that is **Develop requirement** (`DEV-REQ`). In practice, that means asking how the outcome is achieved, **“how else?”** so you do not lock onto one design, and **why** each candidate requirement exists so it traces to the outcome. For risks: what could prevent the outcome from being true, naming each risk on the README and tying it to a **mitigating requirement** (existing or new). The goal is an **index** you trust before you open full requirement files.

**Records:** **PDRs** when framing requirements implies a **new** strategic choice. **CRs** when the **outcome README** changes in a meaningful way (including table work).

---

### Draft

Stress-test the outcome: ask what must be true for success, what could go wrong, and whether each line traces back to the outcome. Lock **one-liners** (and risk links) on the outcome README before you invest in long requirement bodies. If the team makes a new strategic call while doing this, record a **PDR** at outcome or product scope as fits.

### AI-assisted review

Look for missing risks, orphan requirements, and requirements that sound like implementation smuggled into product language. Table edits that shift meaning earn **CRs** at outcome scope.

### Independent review

Confirm the index is enough to prioritize and build from. **PDRs** only for decisions that surface at the gate; **CRs** if the reviewer’s edits are still substantial.
