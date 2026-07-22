<!--
  Tier 2 template. Sibling: requirement-document.md (Tier 1).

  Delta from sibling:
  - `## See Also` comment: paths relative to .../requirements/R<NNN>-<name>/README.md (not .md).
  - PDR/CR backlinks use one extra `../` (../../../../ vs ../../../) for the directory depth.
  - Adds `### Other Materials` under `## See Also`.

  Sync: Title through `## Dependencies` must match requirement-document.md; update both siblings in the same commit.
-->
# O<NNN>-R<NNN> - <Requirement Name>

[Product/Solution] must [Capability/Constraint]

## Mitigates

<!-- List every risk this requirement mitigates; the owning outcome follows from the risk's ID. Omit this section only if there are no linked risks (in which case the entry is a feature request, not a requirement). -->

- **O<NNN>-RSK<NNN>** - <Risk Name>  (serves **O<NNN>** - <Outcome Name>)
- **O<NNN>-RSK<NNN>** - <Risk Name>  (serves **O<NNN>** - <Outcome Name>)

## Detail

<Expanded description of the requirement.>

## Edge Cases

- <Condition>: <Expected Behavior>
- <Condition>: <Expected Behavior>

## Examples

### <Scenario Name>

- Input: <Input>
- Expected Output: <Expected Output>
- Verification: <Verification Step(s)>

## Acceptance Criteria

- [ ] <Concrete, testable condition>
- [ ] <Concrete, testable condition>

## Dependencies

<!-- Prerequisite requirements only: other requirement elements this one cannot be satisfied without (blocking product coupling — not implementation order, sprint sequencing, or shared-risk links). When this requirement assumes a component contract defined elsewhere, list the requirement that owns that contract. Omit this section when there are no prerequisites. -->

- **O<NNN>-R<NNN>** - <Requirement Name>
- **O<NNN>-R<NNN>** - <Requirement Name>

## See Also

<!-- Omit any subsection with no entries; omit this entire `## See Also` if every subsection would be empty. Records about this requirement live in the central directories under docs/product/; paths below are relative to this file at docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/README.md. -->

### Product Decision Records

- [PDR<NNN> - <PDR Name>](../../../../drs/PDR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](../../../../crs/CR<NNN>-<name>.md)

### Other Materials

- [<Material Name>](<path>)
