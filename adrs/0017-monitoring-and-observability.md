# ADR-0017: CloudWatch + OpenTelemetry Metrics, with a Clear Path to a Dedicated Observability Stack

## Status
Accepted

## Context
The system needs visibility into request rates, latencies, error rates, and business metrics (scans/min, reviews/min, cache hit rate) from day one, deployed on AWS (see [[0019-compute-and-deployment]]), with a small team that shouldn't be operating a heavyweight observability stack before there's meaningful traffic.

## Problem
How is the system monitored (metrics, alerting, dashboards)?

## Alternatives
- **No metrics, logs only** — cheapest, but makes it impossible to alert on latency/error-rate regressions proactively; only reactive (post-incident log digging). Rejected.
- **Self-hosted Prometheus + Grafana** — powerful and vendor-neutral, but adds a nontrivial operational burden (running, scaling, and backing up a metrics stack) disproportionate to a small team's initial needs.
- **Third-party SaaS APM (Datadog, New Relic)** — excellent product, but real recurring cost and a new vendor relationship before there's evidence it's needed.
- **AWS-native: application emits metrics via OpenTelemetry SDK → CloudWatch (custom metrics + Container Insights for ECS-level CPU/memory), with alarms on key thresholds (p95 latency, 5xx rate, DB connection saturation)** — no new infrastructure to run, uses the same AWS account/billing as everything else, and OpenTelemetry instrumentation in the app code is vendor-neutral, so switching the *backend* (metrics destination) later doesn't require re-instrumenting the app.

## Decision
Instrument the backend with the OpenTelemetry Go SDK for metrics (request latency histograms per route, error counters, cache hit/miss counters, business counters like reviews-created) and traces (spans per request, per DB query, per external call), exported to AWS CloudWatch/X-Ray via the OpenTelemetry Collector (run as an ECS sidecar). ECS Container Insights covers infra-level CPU/memory/network. CloudWatch Alarms cover: p95/p99 latency thresholds, 5xx error rate, DB CPU/connections, ElastiCache memory. Alarms route to an SNS topic (initially emailing the on-call developer; extensible to PagerDuty/Slack later).

## Consequences
- No new infrastructure to operate beyond what AWS already provides.
- Because instrumentation is via OpenTelemetry (not AWS-SDK-specific calls sprinkled through the code), the export destination can change without re-instrumenting application code.
- CloudWatch's query/dashboard experience is less polished than Grafana or Datadog; accepted trade-off until traffic/team size justifies the switch.

## Future Considerations
Migrate the OpenTelemetry Collector's export target to a dedicated observability platform (self-hosted Grafana stack, or a SaaS APM) once traffic, team size, or incident-response needs outgrow CloudWatch's dashboarding/alerting ergonomics. Because instrumentation is OTel-based, this is a Collector config change, not an application rewrite.
