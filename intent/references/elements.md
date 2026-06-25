# Elements

The subsections below define each element in the intent hierarchy — **product**, **engineering**, and **records** — and how they relate to each other.

For naming, paths, and tiers, see [structure.md](structure.md). For how to read, create, update, and delete documents, see [workflows.md](workflows.md).

## Alignment check

Run these yes/no checks against any element. Any "no" is a fail. Read "parent" and "descendants" along the spine: product → job → outcome → risk → requirement → architecture → component.

Vertical alignment:

- Follows from parent: does the element follow from its parent without contradicting it? "No" means it pulls against what its parent claims.
- Honored by descendants: do its descendants honor it without contradiction? "No" means something downstream drifts from or contradicts it.

Sibling scoping (run across the siblings under one parent):

- Minimal overlap: is overlap between siblings minimal? "No" means two siblings cover the same ground.
- Agree on contact: where siblings touch the same situation, do they agree? "No" means siblings conflict on a shared point.
- Cover the parent: do the siblings together cover the parent's scope without gaps? "No" means part of that scope is unserved.

Coherence:

- Consistent links: do laterally linked elements, and any records that name this element, agree with it on the same facts? "No" means a lateral link or record contradicts the element.

## Product

A product is a solution people hire to get a job done.

It is defined not by its features but by:

- The job it serves
- The outcome customers use to measure success
- The risks that threaten those outcomes
- The requirements that constrain how it delivers

The product has no ancestors above it in this framework. Its descendants—jobs, outcomes, risks, and requirements—should read as one story with the product: nothing downstream should contradict what the product claims.

### Job

A job is a goal someone is trying to achieve in a specific situation — the progress they want to make, not the product they use. It is functional, situational, and stable across changes to tools and features.

Run the [Alignment check](#alignment-check): a job's parent is the product and its descendants are outcomes; sibling jobs together cover the product's scope.

Minimal Pattern:

> When [Situation], I want to [Goal], so I can [Outcome].

- Situation: when, where, or why the job arises — the context or trigger that creates the need.
- Goal: what the person is trying to accomplish — the progress they want to make.
- Outcome: why it matters — the benefit or result they expect from achieving the goal.

Examples:

**Good**

> When I'm comparing options before a big purchase, I want to feel confident I'm not missing a better alternative, so I can buy without second-guessing.

Why it's good: it's functional (comparing options), situational (before a big purchase), and solution-free (no product named). The situation provides the trigger, the goal states the progress, and the outcome explains why it matters.

**Bad**

> I want a product comparison app with filters and ratings.

Why it's bad: it's a solution, not a job. Name is a product. Skips the situation and outcome entirely.

### Outcome

An outcome is a measurable way the customer judges success while (or after) getting a job done. It answers: "What would 'better' look like?" in terms they care about — not in terms of your product. It is stated from the customer's perspective, measurable in principle, and free of solution naming.

Run the [Alignment check](#alignment-check): an outcome's parent is its owning job and its descendants are risks; sibling outcomes together cover the jobs they serve. Apply the sibling-scoping checks to cousin outcomes too, so peers and cousins together stay complete without contradiction.

Minimal Pattern:

> [Verb] [Unit of Measure] [Object]

- Verb: minimize or maximize
- Unit of Measure: time, number, or likelihood — likelihood covers probability and confidence
- Object: the specific thing being measured

Examples:

**Good**

- Minimize the time to identify the most relevant alternatives.
- Minimize the likelihood of overlooking a meaningful difference between options.
- Maximize the confidence that the chosen option meets the most important criteria.

Why they're good: each is customer-centric (stated from their perspective), measurable in principle (time, likelihood, confidence), and solution-free (no feature or technology named). They follow the minimal pattern with a clear verb, unit of measure, and object.

**Bad**

- Show me a side-by-side comparison table
- Give me at least 10 results

Why they're bad: they're feature requests, not measurable success metrics. They name a solution and aren't stated from the customer's perspective of what "better" looks like.

### Risk

A risk is a condition or event that negatively impacts a desired outcome. It is tied to that outcome, measurable by likelihood and impact on the outcome, stated from the customer's perspective, and free of product or feature naming.

Run the [Alignment check](#alignment-check): a risk's parent is its outcome and its descendants are the requirements that mitigate it. For sibling scoping, minimal overlap means no two risks duplicate the same failure mode, and where two risks bear on the same requirement they agree.

Minimal Pattern:

> [Condition/Event] [Negative Impact on Outcome]

- Condition/Event: what might go wrong — a situation, trigger, or circumstance that could occur.
- Negative Impact on Outcome: how it hurts — the specific way it degrades a desired outcome.

Examples:

**Good**

- Conflicting reviews increase the time to identify the most relevant alternatives.
- Time pressure increases the likelihood of overlooking a meaningful difference between options.
- Post-purchase discovery of a better option reduces confidence that the chosen option meets the most important criteria.

Why they're good: each is outcome-linked (directly references a specific desired outcome), measurable by likelihood × impact, customer-centric (stated as conditions the customer faces), and solution-free (no product or feature named). The condition/event is distinct from the negative impact.

**Bad**

- Returning filers take too long to complete their tax return.
- Without a document-import feature, users abandon their return before submitting.

Why they're bad: the first only states a negative result — it restates an outcome (minimize filing time) without naming the condition or event that causes it, so there is nothing whose likelihood × impact you could weigh. The second is solution-aware: it frames the risk as the absence of a named feature (document import), which is really a missing requirement, not a customer-facing condition tied to an outcome.

### Requirement

A requirement is what the product must do or respect to mitigate one or more risks. It is tied to specific risks, states what the product must do (not how it is built), is verifiable, and exists only to mitigate those risks.

Run the [Alignment check](#alignment-check): a requirement's parents are its outcome and the risks it mitigates, and its descendants are architecture and components; sibling requirements together cover the mitigations the outcomes need. Apply the sibling-scoping checks to cousin requirements traced through shared architecture, so their demands stay compatible.

Minimal Pattern:

> [Product/Solution] must [Capability/Constraint]

- Product/Solution: the product or solution being constrained.
- Capability/Constraint: what it must do or respect.

Examples:

**Good**

- The product must surface contradictions across reviews.
- The product must allow saving and resuming a comparison.
- The product must confirm no better-matching options exist at time of decision.

Why they're good: each is risk-linked (directly references one or more specific risks), solution-aware but not implementation-specific (says what the product must do, not how), and verifiable (you can test whether it's met).

**Bad**

- Build a review aggregation microservice.
- Add a save button to the comparison page.

Why they're bad: they're implementation details, not requirements. They specify *how* to build, not *what* the product must do. They aren't tied to any specific risks.

## Engineering

Engineering is how the product is built. It is defined by the architecture that shapes it and the components that compose it.

Run the [Alignment check](#alignment-check): engineering's parents are the product's requirements and outcomes, and its descendants are the architecture and components that realize them.

Per the coherence check, product-facing and engineering-facing descriptions of the same requirement, component, or boundary must agree.

### Architecture

An architecture is the set of decisions that define how components are organized, communicate, and constrain each other. It describes structure and relationships, not individual features in isolation, and is driven by product requirements.

Run the [Alignment check](#alignment-check): the architecture's parents are the product, outcomes, and the requirements it must carry, and its descendants are the components that realize it. Beyond the check, structure, principles, and technology choices within the architecture must not contradict each other.

Examples:

**Good**

> The system separates review analysis from comparison management, communicating through a shared option model. Both are stateless; session persistence is handled by the Comparison Session Store.

Why it's good: structural (describes relationships and boundaries), requirement-driven (the separation traces to distinct product requirements).

**Bad**

> We use microservices and React.

Why it's bad: names technologies, not structure. Doesn't describe how parts relate or why.

### Component

A component is a distinct, implementable part of the system that fulfills one or more requirements. It has a clear scope, maps to those requirements, and is concrete enough to build, test, and deploy.

Run the [Alignment check](#alignment-check): a component's parents are the architecture and the requirements it fulfills; sibling components together cover the architecture's obligations.

Minimal Pattern:

> [Name] — [Responsibility]

- Name: what the component is called.
- Responsibility: what it does.

Examples:

**Good**

- Review Analyzer — surfaces contradictions across reviews.
- Comparison Session Store — persists and restores in-progress comparisons.

Why they're good: bounded (clear scope), requirement-linked (maps to a specific product requirement), and named by what they do.

**Bad**

- Trip Manager — coordinates everything about a ride, from matching and routing to fare calculation and receipts.
- Notification Service — sends notifications to users.

Why they're bad: the first has a plausible name and a stated responsibility, but its scope is unbounded and overlaps several concerns (matching, routing, pricing, receipts) that other components would also own — you could not build, test, or trace it to a single requirement. The second names a real responsibility but states it circularly ("sends notifications") and links to no requirement, so there is no risk it exists to mitigate and no way to tell when it is done.

## Records

Records capture the decisions and changes made to product and engineering elements over time.

Run the [Alignment check](#alignment-check): a record's parents are the elements it concerns, and it must stay coherent with those elements' ancestors and descendants on the trace.

Per the coherence check, records at related scope must not contradict each other on the same fact without an explicit superseding record, and where two records touch the same interface or decision surface they agree.

### Change Record (CR)

A change record documents a modification to a product or engineering element — what changed, why, and what it affects. It is scoped to a specific element, explains what changed and why, and reads as a point-in-time modification. Trivial edits need no record. See [structure.md](structure.md#choosing-a-tier) for when to capture a change inline, as a file, or as a folder.

Run the [Alignment check](#alignment-check): a change record's parent is the element it names, and it must stay coherent with that element's ancestors and descendants on the trace.

Per the coherence check, other change records that bear on the same fact or touch the same interface must not tell conflicting stories, and where they affect shared boundaries they agree on facts.

Examples:

**Good**

> Updated the Review Analyzer component to normalize star ratings across sources. Previously, ratings were passed through raw, causing contradictions to be missed when sources used different scales. This change affects the Review Analyzer component and the "surface contradictions across reviews" requirement.

Why it's good: scoped (names the element changed), traceable (explains what changed and why), and temporal (describes a before/after).

**Bad**

> Updated the scheduling service to retry failed appointment bookings up to three times.

Why it's bad: it reads as a real change and even names a behavior, but it gives no rationale (why retries were added or what problem they solve) and traces to no element — it does not say which component or requirement it modifies, or what else the change affects, so the record cannot be checked for coherence against the trace.

### Product Decision Record (PDR)

A product decision record documents a product decision — the context, the options considered, and the choice made. It states the situation, alternatives, rationale, and consequences (including costs and follow-ups).

Run the [Alignment check](#alignment-check): a product decision record's parents are the product outcomes and requirements it affects, and its consequences must read through to downstream elements.

Per the coherence check, product decision records at overlapping or related scope must not contradict without an explicit superseding record, and where decisions overlap they agree on facts and intent.

Examples:

**Good**

> We decided to limit comparisons to three options at a time. Users comparing more than three reported decision fatigue and longer completion times. We considered unlimited comparisons and a five-option cap. Three balances thoroughness with cognitive load. Consequence: users comparing more than three options must run multiple sessions.

Why it's good: context-driven, option-aware, justified, and states the consequence.

**Bad**

> We're only showing three options.

Why it's bad: no context, no alternatives considered, no justification, no consequences.

### Architectural Decision Record (ADR)

An architectural decision record documents an engineering decision — the context, the options considered, and the choice made. It states the situation, alternatives, rationale, and consequences (including costs and follow-ups).

Run the [Alignment check](#alignment-check): an architectural decision record's parents are the architecture and components it constrains, and its tradeoffs must trace to the requirements and product intent above.

Per the coherence check, architectural decision records at overlapping or related scope must not contradict without an explicit superseding record, and where decisions overlap they agree on facts and intent.

Examples:

**Good**

> We decided to make the Review Analyzer stateless, storing no session data. The Comparison Session Store already handles persistence, and duplicating state would create sync risks. We considered a stateful analyzer with local caching. Consequence: the Review Analyzer must re-fetch review data on every request.

Why it's good: context-driven, option-aware, justified, and states the consequence.

**Bad**

> Review Analyzer is stateless.

Why it's bad: states the decision but not the context, alternatives, justification, or consequences.
