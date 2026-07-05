# 0001: Deferred Semantic Merge Conflict Resolution

**Status:** Accepted for Phase 1. Revisit in a later phase.

## Problem

Merge conflict resolution exists to combine changes made independently on
different branches into a single consistent result. When two branches modify
the same piece of data differently, the system must decide which changes to
keep, whether they can be combined safely, or whether the conflict requires
more than a mechanical decision. The goal is to preserve as much useful
information as possible while still producing a valid, unambiguous final
state.

## Decision: last-write-wins for v1

Phase 1 resolves every conflict with **last-write-wins (LWW)**: the commit
with the more recent timestamp is kept, and the other commit's payload is
discarded (though preserved in `MergeResult.Discarded`, never silently lost
from the system entirely). This guarantees every merge produces a
deterministic result without requiring manual review or semantic analysis.

## What this gets wrong

LWW treats every conflict as if the newer change is automatically correct,
which is often false. It can discard valuable information from the older
branch even when both changes were independently true and could have been
combined. For example:

> **Branch A:** "The price increased due to inflation."
> **Branch B:** "The price increased due to supply shortages."

LWW keeps whichever commit was created later and completely loses the other
explanation — even though both statements may be correct and, ideally,
should coexist. LWW optimizes for simplicity and determinism, not for
preserving intent or meaning.

## The eventual alternative: semantic / LLM-assisted merge

A semantic or LLM-assisted merge would analyze the *meaning* of the
conflicting changes rather than only their timestamps. Instead of choosing
one version outright, it could determine whether the two changes are
actually compatible, synthesize them into a coherent combined result, or
flag genuinely irreconcilable conflicts for review rather than guessing.

## Why this is deferred, not just unknown

This is deferred because it introduces significantly more complexity,
computational cost, and uncertainty than a deterministic merge strategy. An
LLM-assisted merge requires semantic understanding, prompt design, model
integration, evaluation of correctness, and mechanisms to handle ambiguous
or incorrect outputs — concerns that are outside the scope of Phase 1, whose
goal is a reliable, deterministic version-control core.

The Phase 1 merge algorithm is intentionally simple not because a better
approach is unknown, but because the project's current objective is
correctness, predictability, and a solid foundation. Establishing a simple,
predictable merge strategy first provides the stable ground on which more
advanced semantic conflict resolution can be implemented and evaluated
later — with a known-correct baseline to compare against.

## Revisit when

- Phase 2/3 integration surfaces real conflict patterns from actual agent
  workloads, not synthetic tests.
- There's a concrete evaluation method for judging whether a semantic merge
  is actually *more correct* than LWW, not just more sophisticated.