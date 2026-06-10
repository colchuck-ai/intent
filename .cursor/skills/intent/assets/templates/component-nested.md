# C<NNN> - <Component Name>

<Responsibility>

<Optional paragraph: scope, boundaries, or context when the one-liner is not enough.>

## Data model

<Owned or exposed entities, schemas, invariants. Link to supplementary files in this folder if the model is large. Omit section if trivial.>

## Interfaces

<APIs, events, protocols, or formats — inputs, outputs, versioning. Link to OpenAPI specs, event catalogs, or similar in this folder if needed. Omit section if trivial.>

## Behavior

<Non-obvious logic: ordering, retries, idempotency, conflict handling, or core algorithms. Omit section if straightforward.>

## Edge cases

- <Condition>: <Expected handling>
- <Condition>: <Expected handling>

## Relationships

- **C<NNN>**: <How this component communicates or depends on the other>
- **C<NNN>**: <How this component communicates or depends on the other>

## Success criteria

- <Engineering-level check — e.g. SLO, invariant, load or correctness assumption>
- <Engineering-level check>

## Notes

<Optional. Constraints, runbooks, or operational detail not covered above.>

## See Also

### Architectural Decision Records

- [<ADR Name>](adrs/<ADR filename>)

### Change Records

- [<CR Name>](crs/<CR filename>)

### Other Materials

- [<Material Name>](<path>)
