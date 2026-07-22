<!--
  Tier 2 template. Sibling: outcome-document.md (Tier 1).

  Delta from sibling:
  - `## Requirements`: tier rendering comments (Tier 0 inline vs Tier 1+ reference link).
  - `## See Also` comment: paths relative to .../outcomes/O<NNN>-<name>/README.md (not .md).
  - PDR/CR backlinks use one extra `../` (../../ vs ../) for the directory depth.
  - Adds `### Other Materials` under `## See Also`.

  Sync: Title, [Direction] line, `## Risks`, and `## Risk-Requirement Map` must match outcome-document.md; update both siblings in the same commit.
-->
# O<NNN> - <Outcome Name>

[Direction] [Metric] [Object] [Context]

## Risks

- **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]
- **O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]

## Requirements

<!-- Render each requirement by its tier (see references/structure.md#child-rendering): -->
<!--   - Tier 0 requirement — inline form: bullet with the minimal pattern. -->
<!--   - Tier 1+ requirement — reference form: same bullet shape, but the name links to the child doc. -->
<!-- Pick exactly one form per requirement — never emit both bullets below for the same child. -->

<!-- Inline form (Tier 0 requirement) — use this bullet OR the reference bullet, not both -->
- **O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]

<!-- Reference form (Tier 1+ requirement) — use this bullet OR the inline bullet above, not both -->
- **O<NNN>-R<NNN>** - [<Requirement Name>](requirements/R<NNN>-<name>.md): [Product/Solution] must [Capability/Constraint]

## Risk-Requirement Map

- **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>, O<NNN>-R<NNN> - <Requirement Name>
- **O<NNN>-RSK<NNN> - <Risk Name>**: O<NNN>-R<NNN> - <Requirement Name>

## See Also

<!-- Omit any subsection with no entries; omit this entire `## See Also` if every subsection would be empty. Records about this outcome live in the central directories under docs/product/; paths below are relative to this file at docs/product/outcomes/O<NNN>-<name>/README.md. -->

### Product Decision Records

- [PDR<NNN> - <PDR Name>](../../drs/PDR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](../../crs/CR<NNN>-<name>.md)

### Other Materials

- [<Material Name>](<path>)
