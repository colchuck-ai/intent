<!--
  Section optionality — omit empty sections per references/structure.md ("Never emit an empty heading").

  - Title + Responsibility: required
  - Scope paragraph (below Responsibility): omittable when the one-liner is enough
  - ## Data model: omittable when trivial
  - ## Interfaces: omittable when trivial
  - ## Behavior: omittable when straightforward
  - ## Edge cases: omittable when none
  - ## Relationships: omittable when no relationships to other components
  - ## Success criteria: required
  - ## Notes: omittable
  - ## See Also: omittable when every subsection would be empty (each subsection omittable individually when empty)
-->
# C<NNN> - <Component Name>

<Responsibility>

<!-- omittable: scope, boundaries, or context when the one-liner is not enough -->

## Data model

<!-- omittable when trivial -->

<Owned or exposed entities, schemas, invariants.>

## Interfaces

<!-- omittable when trivial -->

<APIs, events, protocols, or formats — inputs, outputs, versioning.>

## Behavior

<!-- omittable when straightforward -->

<Non-obvious logic: ordering, retries, idempotency, conflict handling, or core algorithms.>

## Edge cases

<!-- omittable when none -->

- <Condition>: <Expected handling>
- <Condition>: <Expected handling>

## Relationships

<!-- omittable when no relationships to other components -->

- **C<NNN> - <Component Name>**: <How this component communicates or depends on the other>
- **C<NNN> - <Component Name>**: <How this component communicates or depends on the other>

## Success criteria

<!-- required -->

- <Engineering-level check — e.g. SLO, invariant, load or correctness assumption>
- <Engineering-level check>

## Notes

<!-- omittable -->

<Constraints, runbooks, or links not covered above.>

## See Also

<!-- omittable when every subsection would be empty; omit each subsection individually when empty. Records live in central directories under docs/engineering/; paths below are relative to this file at docs/engineering/components/C<NNN>-<name>.md. -->

### Architectural Decision Records

- [ADR<NNN> - <ADR Name>](../drs/ADR<NNN>-<name>.md)

### Change Records

- [CR<NNN> - <CR Name>](../crs/CR<NNN>-<name>.md)
