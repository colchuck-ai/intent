---
title: "Materiality — is this worth capturing at all?"
summary: "Capture an element only if leaving it out would let a real mistake through."
---
# Materiality

**The question.** Should this be its own element (principle, constraint, risk,
requirement), or is it noise that adds surface without adding signal?

**The test.** Imagine it's absent. Would someone plausibly make a wrong
decision *because* it wasn't written down? If yes, it's material. If it's true
but nobody would act differently knowing it, leave it out.

**Good.**

- A constraint: "The CLI ships as one Go binary with no runtime prerequisite."
  (rules out whole design options)
- A principle: "Derive only the mechanical duals of what is declared."
  (settles recurring arguments)

**Bad.**

- "The code should be readable." (nobody argues the opposite; guides nothing)
- A requirement restating a general best practice with no bearing on *this*
  outcome.

**Rule of thumb.** Principles and constraints are axioms that *close off
options*. If it doesn't close anything, it isn't one.

**Failure mode.** A tree padded with truisms. Readers stop trusting that every
element earns its place, and the material ones get lost in the filler.
