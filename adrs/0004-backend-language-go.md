# ADR-0004: Go as the Backend Programming Language

## Status
Accepted

## Context
The backend needs to serve a REST API with moderate concurrency (barcode lookups, review CRUD, image upload handling), be easy to deploy as a single artifact, and be approachable for future contributors and LLM-assisted code generation.

## Problem
Which language should the backend service be written in?

## Alternatives
- **Node.js/TypeScript** — great ecosystem, shares a language with any web admin frontend, but weaker concurrency primitives for CPU-bound work (image processing, future OCR), and runtime type safety depends entirely on discipline.
- **Python** — excellent for AI/ML-adjacent code (future OCR, summarization glue code), but weaker performance/concurrency story for a high-throughput API, and deployment artifacts are heavier (interpreter + deps) than a single binary.
- **Java/Kotlin (Spring Boot)** — mature, strong typing, good concurrency, but heavier runtime footprint, slower cold starts, more ceremony for a small team's MVP; also Kotlin here would blur the line with the Android codebase's language without actually sharing code (JVM ecosystems for server and Android diverge in practice).
- **Go** — compiles to a single static binary (simple deployment, small container images), strong built-in concurrency (goroutines) well suited to many concurrent I/O-bound requests, a small and stable standard library that keeps codegen and maintenance predictable, and a straightforward learning curve for new contributors.

## Decision
Go. The combination of simple deployment (single binary, small Docker images), strong concurrency support, and a stable, boring standard library minimizes operational and maintenance surface area for a small team running a modular monolith (see [[0003-modular-monolith-backend]]).

## Consequences
- Fast container builds and small images, which keeps CI/CD and Kubernetes/Fargate cold-start costs low.
- Smaller ecosystem for certain domains (advanced ML tooling) than Python — any future AI/ML feature that needs a Python-only library (OCR models, embedding generation) is implemented as a separate small service called over the network, not bolted into the Go monolith.
- Team members need Go familiarity; mitigated by Go's small language surface.

## Future Considerations
If an AI/ML-heavy feature (recommendation engine, OCR) needs a Python-only library, extract it as a standalone microservice (per [[0003-modular-monolith-backend]]'s extraction path) rather than introducing a second general-purpose language into the monolith.
