---
name: Feature request
about: Propose a change to the public behavior or contributor experience
title: "[Feature] "
labels: "enhancement"
assignees: ""
---

## Problem

What user or contributor problem would this solve?

## Proposed behavior

Describe the smallest useful behavior and how it fits the public 0.2.0
standalone skill contract for explicit `$next` use from
`${CODEX_HOME:-$HOME/.codex}/skills/next`.

## Scope check

Explain why the proposal does not require composer insertion, network access,
telemetry, transcript reads, persistence, or automatic continuation.

For Codex 0.149.1, preserve description-driven activation when the user types
`$next`; do not add an automatic footer or metadata that hides the explicit
catalog entry.
