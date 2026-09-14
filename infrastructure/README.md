# Infrastructure

Terraform IaC on AWS. See [ADR-0018](../adrs/0018-infrastructure-as-code-terraform.md), [ADR-0019](../adrs/0019-compute-and-deployment.md), [ADR-0020](../adrs/0020-secrets-management.md).

## Layout
```
infrastructure/
  modules/
    network/          VPC, subnets, security groups
    database/          RDS PostgreSQL
    cache/             ElastiCache Redis
    compute/            ECS cluster, Fargate service, ALB
    storage/            S3 buckets (media)
    cdn/                CloudFront distribution
    secrets/            Secrets Manager resources
  environments/
    dev/                composes modules for dev
    staging/            composes modules for staging
    prod/               composes modules for prod
  bootstrap/            one-time setup: Terraform state S3 bucket + DynamoDB lock table
```

## Usage
```bash
cd infrastructure/environments/dev
terraform init
terraform plan
```
`terraform apply` for `staging`/`prod` runs only via CI ([ADR-0021](../adrs/0021-cicd-github-actions.md)) with a manual approval gate for `prod`. Never apply against shared environments from a local machine.

See [`docs/deployment.md`](../docs/deployment.md) for the full deployment diagram and release process.
