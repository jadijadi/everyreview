# ADR-0001: Record Architecture Decisions

## Status
Accepted

## Context
This project ("EveryReview" — barcode scan → product reviews) is expected to grow from a small MVP (backend + Android) into a multi-client, multi-feature system (iOS, search, AI moderation, recommendations, etc.) over years, touched by multiple contributors and, per project instructions, potentially by other LLM agents generating or modifying architecture over time. Without a durable record of *why* decisions were made, future contributors (human or LLM) will either violate past decisions unknowingly or re-litigate settled questions repeatedly.

## Problem
How do we record architectural decisions so they survive contributor turnover and are discoverable by both humans and LLM agents working on this codebase?

## Alternatives
- **Wiki / external doc tool** — discoverable, but drifts out of sync with code, not versioned with the repo, not visible in PR diffs.
- **Comments in code** — tied to implementation, not to the decision; scattered, hard to find the "why" for cross-cutting concerns.
- **ADRs in-repo (Michael Nygard style), one file per decision** — versioned with the code, reviewable in PRs, greppable, sequential and immutable (superseded rather than edited).

## Decision
Use lightweight ADRs (Nygard style) stored in `adrs/`, one Markdown file per decision, named `NNNN-kebab-case-title.md`, sequential and never renumbered. Use the template in `adrs/template.md`. Every ADR has a `Status`. Every significant architectural decision — language, framework, data store, API style, auth, infra, mobile architecture, third-party library with lock-in risk — must have a corresponding ADR before implementation begins in that area. Implementation must not contradict an `Accepted` ADR; to change direction, write a new ADR that supersedes the old one (mark the old one `Superseded by ADR-XXXX`, leave its content intact for history).

This applies equally to decisions proposed by other LLM agents: an agent-authored change to architecture must land as a new or superseding ADR with `Proposed` status before code implementing it is merged, and a human (or an explicitly authorized review process) must move it to `Accepted`.

## Consequences
- Decision history is versioned, reviewable, and lives next to the code it governs.
- Adds process overhead: every significant decision needs a short document before code.
- Numbering is sequential and global (not per-category), so ADR numbers alone don't indicate topic — use the index in `adrs/README.md`.

## Future Considerations
If the ADR count grows large (100+), consider grouping into subdirectories by domain (`adrs/backend/`, `adrs/mobile/`, `adrs/infra/`) with a generated index — but keep flat numbering across all of them to preserve chronology.
