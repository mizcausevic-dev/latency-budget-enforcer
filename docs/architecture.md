# Latency Budget Enforcer Architecture

## Service Overview

Latency Budget Enforcer models a backend enforcement layer for teams that need to turn service-path latency budgets into ranked operational action.

It represents the sort of service platform and performance teams use to surface:

- sustained latency breach
- p95 drift
- p99 breach
- dependency drag
- error-pressure coupling

## Processing Flow

1. A budget and observation set are loaded into a typed request.
2. Enforcement checks compare expected p95/p99 limits against the measured path.
3. Severity scores are assigned to each breach family.
4. A consolidated report is emitted with evidence and next actions.

## Current Output Modes

- JSON API response
- terminal summary

## Enforcement Families

### Sustained Breach

- repeated path latency above target
- budget overage strong enough to warrant coordinated response

### P95 Drift

- typical-path latency inflation
- user experience erosion under ordinary load

### P99 Breach

- tail spikes beyond safe service behavior
- incident-shaping extreme-path latency

### Dependency Drag

- downstream or infrastructure layers stretching the path
- path contribution that points to the true bottleneck

### Error-Pressure Coupling

- latency and elevated error rate rising together
- user-visible failure pressure instead of pure slowness
