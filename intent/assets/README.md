# Template index (`assets/`)

Markdown shells for **`docs/product/`** and **`docs/engineering/`**. Trees, IDs, and naming: [`../references/guides/structure.md`](../references/guides/structure.md). **Tone:** fill templates with **[plain language](../references/guides/plain-language.md)**. **Workflows** (activities and sessions) live under [`../references/workflows/`](../references/workflows/README.md).

**Product root (`docs/product/README.md`):** choose **simple** or **JTBD**—see [`../SKILL.md`](../SKILL.md) defaults. One style per product tree.

| Template | Typical path |
|----------|----------------|
| [`product-root-simple.md`](product-root-simple.md) | `docs/product/README.md` (who / value / outcomes) |
| [`product-root-jtbd.md`](product-root-jtbd.md) | `docs/product/README.md` (job executor, core job + optional step groups, lifecycle, buyer—omit unused sections) |
| [`outcome.md`](outcome.md) | `docs/product/outcomes/O<NNN>-<slug>/README.md` |
| [`requirement-simple.md`](requirement-simple.md) | `…/requirements/O<NNN>-R<NNN>-<slug>.md` |
| [`requirement-complex.md`](requirement-complex.md) | `…/requirements/O<NNN>-R<NNN>-<slug>/README.md` |
| [`pdr-simple.md`](pdr-simple.md) | `…/pdrs/PDR<NNN>-….md` (simple at any scope) |
| [`pdr-complex.md`](pdr-complex.md) | `…/pdrs/PDR<NNN>-…/README.md` |
| [`cr-simple.md`](cr-simple.md) | Root: `docs/product/crs/PROD-CR<NNN>-<kebab>.md`, `docs/engineering/crs/ENG-CR<NNN>-<kebab>.md`; nested: `…/crs/<scope>-CR<NNN>-<kebab>.md` ([`structure.md`](../references/guides/structure.md)) |
| [`cr-complex.md`](cr-complex.md) | Same scopes as `cr-simple.md`; complex: `…/PROD-CR<NNN>-<kebab>/README.md`, `…/ENG-CR<NNN>-<kebab>/README.md`, or nested `…/<kebab>/README.md` under `crs/` |
| [`adr-simple.md`](adr-simple.md) | `docs/engineering/adrs/ADR<NNN>-….md` |
| [`adr-complex.md`](adr-complex.md) | `…/adrs/ADR<NNN>-…/README.md` |
| [`engineering-root.md`](engineering-root.md) | `docs/engineering/README.md` |
| [`component-simple.md`](component-simple.md) | `docs/engineering/components/C<NNN>-<slug>.md` |
| [`component-complex.md`](component-complex.md) | `docs/engineering/components/C<NNN>-<slug>/README.md` |
