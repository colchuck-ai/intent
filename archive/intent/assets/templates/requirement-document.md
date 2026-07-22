<!--
  Section optionality — omit empty sections per references/structure.md ("Never emit an empty heading").

  - Title + minimal pattern: required
  - ## Mitigates: required when risks are linked; omit when none (feature request, not a requirement)
  - ## Detail: omittable when the minimal pattern is sufficient
  - ## Edge Cases: omittable when none
  - ## Examples: omittable when none
  - ## Acceptance Criteria: required
  - ## Dependencies: omittable when no prerequisite requirements
  - ## See Also: omittable when every subsection would be empty (each subsection omittable individually when empty)
-->
# O<NNN>-R<NNN> - <Requirement Name>

[Product/Solution] must [Capability/Constraint]

## Mitigates

<!-- required when risks are linked; omit when none (feature request, not a requirement) -->

- **O<NNN>-RSK<NNN>** - <Risk Name>  (serves **O<NNN>** - <Outcome Name>)
- **O<NNN>-RSK<NNN>** - <Risk Name>  (serves **O<NNN>** - <Outcome Name>)

## Detail

<!-- omittable when the minimal pattern is sufficient -->

<Expanded description of the requirement.>

## Edge Cases

<!-- omittable when none -->

- <Condition>: <Expected Behavior>
- <Condition>: <Expected Behavior>

## Examples

<!-- omittable when none -->

### <Scenario Name>

- Input: <Input>
- Expected Output: <Expected Output>
- Verification: <Verification Step(s)>

## Acceptance Criteria

<!-- required -->

- [ ] <Concrete, testable condition>
- [ ] <Concrete, testable condition>

## Dependencies

<!-- omittable when no prerequisite requirements — blocking product coupling only, not implementation order or shared-risk links -->

- **O<NNN>-R<NNN>** - <Requirement Name>
- **O<NNN>-R<NNN>** - <Requirement Name>

## See Also

<!-- omittable when every subsection would be empty; omit each subsection individually when empty. Records live in central directories under docs/product/; paths below are relative to this file at docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>.md. -->

### Product Decision Records

- [PDR<NNN> - <PDR Name>](../../../drs/PDR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](../../../crs/CR<NNN>-<name>.md)
