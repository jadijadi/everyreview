# ADR-0005: HTTP Framework — `net/http` + `chi` Router

## Status
Accepted

## Context
Go offers everything from the bare standard library to full-featured frameworks (Gin, Echo, Fiber). The backend is a modular monolith (see [[0003-modular-monolith-backend]]) exposing a REST API (see [[0006-api-style-rest]]).

## Problem
What HTTP routing/framework should the backend use?

## Alternatives
- **Gin / Echo / Fiber (full frameworks)** — batteries-included (binding, validation, middleware ecosystems), but bring their own conventions and abstractions over `net/http` that add a dependency surface and a learning curve, and some (Fiber) aren't even built on `net/http`, limiting compatibility with the standard middleware ecosystem.
- **Bare `net/http`** — zero dependencies, maximal control, but requires hand-rolling route parameters, middleware chaining, and grouping, which is enough boilerplate to slow down a small team.
- **`net/http` + `chi`** — `chi` is a lightweight, idiomatic router built directly on `net/http`'s `Handler` interface (no framework lock-in), adding only routing, URL parameters, middleware chaining, and route grouping (useful for API versioning, see [[0021-api-versioning]]). Everything else (JSON encoding, validation, request binding) stays plain Go, keeping the codebase close to the standard library.

## Decision
Use the standard library `net/http` for the server, with `chi` as the router/middleware layer. Request validation uses explicit Go structs with manual or `go-playground/validator`-based checks; JSON via `encoding/json`.

## Consequences
- Minimal dependency footprint; anything built on `chi` composes with any other `net/http`-compatible middleware.
- Slightly more boilerplate per handler than a batteries-included framework (manual JSON binding/validation), which is an accepted trade-off for long-term maintainability and low framework lock-in.

## Future Considerations
If handler boilerplate becomes a measurable velocity problem, introduce a thin internal helper package for common binding/response patterns rather than switching frameworks.
