# ADR-0016: Structured JSON Logging to stdout/stderr

## Status
Accepted

## Context
The backend runs as containerized replicas (see [[0019-compute-and-deployment]]) with no persistent local disk assumed. Logs need to be correlatable across a request's lifecycle (a single barcode scan might touch product lookup, cache, and DB) and machine-parseable for later aggregation.

## Problem
How does the backend emit logs, and where do they go?

## Alternatives
- **Unstructured plain-text logs** — human-readable in a terminal, but expensive to parse/query reliably once aggregated, no consistent field extraction.
- **Direct-to-external-service logging (app pushes logs straight to a SaaS logging API)** — couples the app to a specific vendor's SDK and adds a failure mode (logging call itself can fail/block) inside the request path.
- **Structured (JSON) logs to stdout/stderr, collected by the container runtime and shipped to CloudWatch Logs by the platform, not by the app** — follows the 12-factor app model: the application's only responsibility is to write structured lines to stdout; the platform (ECS, see [[0019-compute-and-deployment]]) handles collection and shipping via the awslogs driver. Every log line includes a `request_id` (generated at the edge, propagated through context) enabling full request tracing across log lines and, later, distributed traces.

## Decision
Use Go's standard library `log/slog` for structured JSON logging to stdout/stderr. Every HTTP request is assigned a `request_id` (from an incoming header if present, generated otherwise) attached to the request's `context.Context` and included in every log line emitted during that request, plus returned to the client in a response header for support/debugging correlation. ECS task definitions ship container stdout/stderr to CloudWatch Logs via the `awslogs` driver; no custom log-shipping code in the application.

## Consequences
- Zero logging infrastructure code in the application beyond structuring log lines — collection is the platform's job.
- CloudWatch Logs Insights (or a future log aggregation tool) can query/filter structured fields directly.
- `request_id` propagation gives request-level traceability without needing full distributed tracing on day one.

## Future Considerations
If/when the system grows beyond the modular monolith (see [[0003-modular-monolith-backend]]) into multiple services, adopt OpenTelemetry tracing with the same `request_id` promoted to a trace ID, and consider shipping logs to a dedicated aggregation platform (e.g., Grafana Loki, Datadog) if CloudWatch Logs Insights becomes limiting.
