<!--
  Section optionality — omit empty sections per references/structure.md ("Never emit an empty heading").

  - Title + intro paragraph: required
  - ## Principles: omittable when none
  - ## Constraints: omittable when none
  - ## Technology Choices: omittable when none
  - ## Components: required (at least one component)
  - ### C<NNN> component blocks: required per component (inline or reference form per child tier — see references/structure.md#child-rendering)
  - **Relationships** under inline components: omittable when none
  - **Architectural Decision Records** / **Change Records** under inline components: omittable when none
  - ## Requirement-Component Map: omittable when no mappings
  - ## See Also: omittable when every subsection would be empty (each subsection omittable individually when empty)
-->
# <Architecture Name>

<How components are organized, communicate, and constrain each other.>

## Principles

<!-- omittable when none -->

- <Principle>
- <Principle>

## Constraints

<!-- omittable when none -->

- <Constraint>
- <Constraint>

## Technology Choices

<!-- omittable when none -->

- <Technology>: <Rationale>
- <Technology>: <Rationale>

## Components

<!-- required (at least one component); render each by tier (see references/structure.md#child-rendering) -->

<!-- Render each component by its tier (see references/structure.md#child-rendering): -->
<!--   - Tier 0 component — inline block: H3 heading + responsibility line + bold-labeled flat list. -->
<!--   - Tier 1+ component — reference block: H3 heading + responsibility one-liner + "See [link]" line. -->
<!-- Pick exactly one form per component — never emit both blocks below for the same child. -->

<!-- Inline form (Tier 0 component) — use this block OR the reference block, not both -->
### C<NNN> - <Component Name>

<!-- required per component -->

<Responsibility>

**Relationships**

<!-- omittable when none -->

- **C<NNN> - <Component Name>**: <How they communicate or depend>
- **C<NNN> - <Component Name>**: <How they communicate or depend>

**Architectural Decision Records**

<!-- omittable when none -->

- [ADR<NNN> - <ADR Name>](drs/ADR<NNN>-<name>.md)

**Change Records**

<!-- omittable when none -->

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)

<!-- Reference form (Tier 1+ component) — use this block OR the inline block above, not both -->
### C<NNN> - <Component Name>

<!-- required per component -->

<Responsibility>

See [C<NNN> - <Component Name>](components/C<NNN>-<name>.md).

## Requirement-Component Map

<!-- omittable when no mappings -->

- **O<NNN>-R<NNN> - <Requirement Name>**: C<NNN> - <Component Name>, C<NNN> - <Component Name>
- **O<NNN>-R<NNN> - <Requirement Name>**: C<NNN> - <Component Name>

## See Also

<!-- omittable when every subsection would be empty; omit each subsection individually when empty. Records about the architecture as a whole live in the central directories below. -->

### Architectural Decision Records

- [ADR<NNN> - <ADR Name>](drs/ADR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](crs/CR<NNN>-<name>.md)
