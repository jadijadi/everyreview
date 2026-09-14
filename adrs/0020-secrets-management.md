# ADR-0020: AWS Secrets Manager for Application Secrets

## Status
Accepted

## Context
The backend needs database credentials, Redis auth, JWT signing keys, OAuth client secrets (later), and third-party API keys. These must never be committed to the repository and must differ per environment (dev/staging/prod).

## Problem
Where do application secrets live, and how does the running backend obtain them?

## Alternatives
- **`.env` files / environment variables set manually** — fine for local dev, but manually managing production secrets invites human error (accidental commits, stale copies, no audit trail, no rotation).
- **Encrypted values checked into the repo (e.g., SOPS, git-crypt)** — versioned and reviewable, but adds key-management overhead of its own and still requires a separate secret store for the encryption keys.
- **HashiCorp Vault (self-hosted)** — powerful, but another stateful service to operate, disproportionate for this project's needs.
- **AWS Secrets Manager** — native IAM-based access control (ECS task role grants read access to only the specific secrets a task needs), automatic rotation support for RDS credentials, audit trail via CloudTrail, referenced directly from ECS task definitions (injected as environment variables at container start, never baked into images or Terraform state in plaintext).

## Decision
AWS Secrets Manager for all production/staging secrets. Terraform provisions the secret *resources* (names/structure) but never their values inline in `.tf` files; values are set out-of-band (via `aws secretsmanager put-secret-value` in a bootstrap step, or manually for initial setup) so secret values never appear in Terraform state history in plaintext beyond what's unavoidable. ECS task definitions reference secrets by ARN; Fargate injects them as environment variables at container start. Local development uses a `.env.local` file (gitignored, documented in `docs/development-setup.md`) with dummy/dev-only values — never a copy of real secrets.

## Consequences
- No secret values in git history or CI logs.
- IAM scoping means a compromised task can only read the secrets it's explicitly granted, not the entire secret store.
- RDS credential rotation can be automated via Secrets Manager's native rotation Lambda.
- Adds a small amount of Terraform/bootstrap complexity versus plain environment variables.

## Future Considerations
If secret sprawl grows (many services, many environments) enough that Secrets Manager's per-secret cost or organization becomes unwieldy, evaluate a self-hosted Vault — unlikely to be justified before the system has multiple independently-deployed services.
