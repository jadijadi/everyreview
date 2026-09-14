# ADR-0021: GitHub Actions for CI/CD

## Status
Accepted

## Context
The monorepo (see [[0002-repository-structure]]) contains backend (Go), Android (Kotlin/Gradle), later iOS (Swift), and infrastructure (Terraform), each needing its own build/test/deploy pipeline, ideally without needing a separate CI product.

## Problem
What runs CI (build/test/lint) and CD (deploy) for this project?

## Alternatives
- **Jenkins (self-hosted)** — maximum flexibility, but is itself infrastructure to run, patch, and secure; disproportionate operational cost for this team size.
- **CircleCI / GitLab CI** — solid products, but the repo is hosted on GitHub (implied by the `.github/` directory in the requested structure), so GitHub Actions avoids a second vendor relationship and integrates natively with PR checks, required-status-checks branch protection, and GitHub's OIDC-to-AWS federation (no long-lived AWS keys stored in CI).
- **GitHub Actions** — native to the hosting platform, path-filtered workflows (only run Android CI when `android/**` changes, only run backend CI when `backend/**` changes — important for monorepo build-time hygiene), first-class macOS runners for future iOS builds, and OIDC federation to AWS IAM roles for deploy jobs (no static AWS credentials in CI secrets).

## Decision
GitHub Actions, with workflows scoped by path filters per top-level directory:
- `backend-ci.yml` — lint (`golangci-lint`), unit tests, `sqlc`/migration validation, build image, on every PR touching `backend/**`.
- `android-ci.yml` — lint, unit tests, instrumented tests (emulator), assemble APK, on every PR touching `android/**`.
- `infrastructure-ci.yml` — `terraform fmt -check`, `terraform validate`, `terraform plan` (posted as a PR comment), on every PR touching `infrastructure/**`.
- `backend-deploy.yml` — on merge to `main` (path-filtered to `backend/**` or `infrastructure/**`): run DB migrations, `terraform apply` per environment, build+push image to ECR, update ECS service. Deploys to `staging` automatically; `prod` requires a manual approval gate (GitHub Environments protection rule).
- `android-release.yml` — triggered manually or on release tag: build signed release APK/AAB, upload to Play Console (internal/beta tracks initially).

AWS access from CI uses GitHub's OIDC provider federated to a scoped IAM role — no long-lived AWS access keys stored as GitHub secrets.

## Consequences
- Fast CI feedback because monorepo path filtering avoids running irrelevant pipelines (a backend-only PR doesn't trigger a 20-minute Android build).
- No CI infrastructure to operate/patch.
- Deploy credentials never leave AWS/GitHub's federated trust relationship, eliminating a class of leaked-static-credential risk.
- Coupling to GitHub Actions specifically; acceptable since the repo is already GitHub-hosted per the requested project structure.

## Future Considerations
If build minutes/cost become significant at scale, consider self-hosted runners for the heaviest jobs (Android instrumented tests, Docker builds) while keeping orchestration in GitHub Actions.
