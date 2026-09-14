# ADR-0019: AWS ECS Fargate for Backend Compute (Kubernetes Deferred)

## Status
Accepted

## Context
The backend is a single modular-monolith Go service (see [[0003-modular-monolith-backend]]) packaged as a Docker container (see [[0004-backend-language-go]]). It needs to run reliably, scale with load, and support zero-downtime deploys, for a small team without dedicated platform/SRE staff.

## Problem
What runs the backend container(s) in production?

## Alternatives
- **AWS Lambda (serverless)** — scales to zero and requires no server management, attractive for unpredictable early traffic, but cold starts hurt latency-sensitive endpoints (barcode scan is the core interaction and must feel instant), and long-lived DB connection pooling is awkward in a Lambda execution model; better suited to bursty/event-driven workloads (e.g., a future image-processing or OCR pipeline) than the main API.
- **Kubernetes (EKS)** — maximum flexibility and portability, the most "standard" choice for microservices at scale, but a single-service modular monolith doesn't need Kubernetes's orchestration complexity (custom controllers, cluster upgrades, networking/CNI, ingress controllers) yet; running EKS well is itself a part-time job for a small team. Explicitly deferred per ADR-0001's instruction to justify each infra choice, not default to Kubernetes.
- **AWS ECS on EC2** — cheaper at large sustained scale, but requires managing the underlying EC2 fleet (patching, capacity planning) that Fargate abstracts away.
- **AWS ECS on Fargate** — serverless containers: define a task (the Go container), ECS handles placement, scaling (via Application Auto Scaling on CPU/request count), and replacement on failure; no EC2 fleet to manage; integrates natively with an Application Load Balancer, CloudWatch, and Secrets Manager; supports rolling/blue-green deploys out of the box.

## Decision
ECS on Fargate, fronted by an Application Load Balancer (ALB) that terminates TLS and routes to the ECS service's target group. The service runs a minimum of 2 tasks (across 2 Availability Zones) for availability even at low traffic, auto-scaling on CPU utilization and ALB request count. Deploys are rolling (ECS default), with the ALB health check gating traffic shift — a bad deploy is automatically rolled back by failing health checks before it receives production traffic.

## Consequences
- No server/cluster management burden; the team operates at the "container + task definition" level, not the "EC2 fleet" or "Kubernetes cluster" level.
- Less flexible than Kubernetes for complex multi-service orchestration needs — acceptable because the backend is one service today (see [[0003-modular-monolith-backend]]).
- Cost is roughly proportional to running tasks; less optimal than Lambda's scale-to-zero for very bursty, low-traffic-most-of-the-time workloads (a future async job, e.g. image thumbnailing, may be a better fit for Lambda than for the main ECS service).

## Future Considerations
If the modular monolith is later split into multiple independently-scaled services (per [[0003-modular-monolith-backend]]'s extraction path) and the number of services grows past what ECS's simpler orchestration comfortably handles, migrate to EKS. Because everything is already containerized and defined in Terraform, this is an orchestration-layer swap, not an application rewrite. Event-driven/bursty workloads (thumbnail generation, OCR, notification fan-out) should use Lambda or ECS scheduled tasks as they're introduced, independent of this decision for the main API.
