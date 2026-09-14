# ADR-0018: Terraform for Infrastructure as Code, AWS as Cloud Provider

## Status
Accepted

## Context
Infrastructure (compute, database, object storage, CDN, load balancer, networking, secrets, monitoring) must be reproducible across environments (dev, staging, prod) and auditable/reviewable like application code. AWS is specified as the default cloud unless another is clearly justified — nothing in this project's requirements (barcode scanning, review CRUD, image storage) needs a capability unique to another cloud, so AWS is used.

## Problem
How is cloud infrastructure defined, provisioned, and kept reproducible; which IaC tool?

## Alternatives
- **Manual console configuration ("ClickOps")** — fast to start, but not reproducible, not reviewable, and drifts silently; rejected outright for anything beyond a personal experiment.
- **AWS CDK (TypeScript/Python) or CloudFormation** — first-party AWS tooling, but ties IaC to AWS-specific abstractions/templating, and CDK adds a compile step and another language surface (TypeScript/Python) alongside Go/Kotlin already in use.
- **Pulumi** — general-purpose language IaC (could use Go, matching the backend), but a smaller ecosystem/community than Terraform for AWS-specific modules and examples.
- **Terraform (HCL)** — the de facto industry standard for AWS IaC, huge module ecosystem, cloud-agnostic syntax (reduces lock-in risk if a future component needs a different provider), mature state-management and plan/apply review workflow that fits naturally into a PR-based CI/CD process.

## Decision
Terraform, with state stored remotely in an S3 backend with DynamoDB state locking (both provisioned via a small bootstrap Terraform config applied manually once per AWS account). Infrastructure is organized as reusable modules (`infrastructure/modules/{network,database,compute,storage,cdn,secrets}`) composed per environment (`infrastructure/environments/{dev,staging,prod}`), so environments share module code and differ only in variable values. `terraform plan` runs on every PR touching `infrastructure/`; `terraform apply` runs from CI on merge to main, per environment, requiring manual approval for `prod` (see [[0021-cicd]]).

## Consequences
- All infrastructure changes are code-reviewed, versioned, and reproducible from scratch in a new AWS account if needed (disaster recovery, see `docs/deployment.md`).
- Requires discipline to never make manual console changes that cause state drift.
- Terraform state (in S3+DynamoDB) becomes a critical asset requiring its own backup/access-control care.

## Future Considerations
If infrastructure grows complex enough to need per-team or per-service ownership boundaries, split the Terraform root modules per domain (already modularized, so this is a matter of composition, not a rewrite).
