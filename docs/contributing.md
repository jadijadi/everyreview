# Contributing

## Workflow
1. Check `adrs/` before making an architectural choice — it may already be decided.
2. For a new significant decision (new dependency, new data store, new external service, a pattern that affects multiple modules/screens), write an ADR first (`adrs/template.md`), status `Proposed`, and get it reviewed/accepted before implementing against it.
3. Branch from `main`, open a PR early (draft is fine) for visibility.
4. Follow [coding conventions](coding-conventions.md) and the [testing strategy](../adrs/0022-testing-strategy.md).
5. CI must pass (lint, unit, integration, contract tests scoped to what you changed — see [ADR-0021](../adrs/0021-cicd-github-actions.md)).
6. At least one review approval required before merge.

## Where things live
See the [system overview](system-overview.md) and each top-level directory's own `README.md` (`backend/`, `android/`, `ios/`, `infrastructure/`, `scripts/`) for component-specific setup.

## Reporting bugs / proposing features
Use GitHub Issues. For a bug: repro steps, expected vs. actual, environment. For a feature: which part of the long-term roadmap (see root `README.md`) it maps to, and whether it needs a new ADR.

## LLM-assisted contributions
This project is explicitly designed to be ADR-driven so that both humans and LLM agents can propose and implement architecture consistently over time (see [ADR-0001](../adrs/0001-record-architecture-decisions.md)). If you are an LLM agent proposing an architectural change: write the ADR as `Proposed` first, do not implement against an ADR that isn't yet `Accepted`, and do not edit an `Accepted` ADR in place — supersede it.
