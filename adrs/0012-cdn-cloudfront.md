# ADR-0012: CloudFront as CDN for Media and Cacheable API Responses

## Status
Accepted

## Context
Product images and, potentially, cacheable GET API responses (product lookups) need to be served with low latency to a global user base without hitting origin (S3 or the backend) on every request.

## Problem
What CDN sits in front of object storage (and optionally the API) for read traffic?

## Alternatives
- **No CDN, serve directly from S3 / backend** — simplest, but higher latency for geographically distant users and higher origin load/cost as traffic grows.
- **Third-party CDN (Cloudflare, Fastly)** — excellent products, but adds a second cloud vendor relationship and IAM/billing surface when the rest of the infrastructure is AWS-native (see [[0018-infrastructure-as-code]] context); no specific feature here justifies leaving the AWS perimeter.
- **Amazon CloudFront** — integrates natively with S3 (origin access control, no public bucket needed) and with AWS WAF/Shield for basic abuse protection, single vendor/billing relationship alongside the rest of the AWS infrastructure, supports caching both S3 objects and, later, specific cacheable API GET routes behind the same distribution.

## Decision
CloudFront in front of S3 for all media (product/review images), using Origin Access Control so the S3 bucket itself is never public. Cache-Control headers set at upload time based on content type (long TTL + content-addressed/immutable object keys for images, since a new upload gets a new key rather than overwriting). Revisit adding CloudFront in front of specific cacheable API GET endpoints once traffic patterns justify it.

## Consequences
- Low-latency global media delivery, reduced S3 request costs and backend load.
- Cache invalidation is a non-issue for images because object keys are immutable/content-addressed (edits create new objects, not overwrites).
- One more piece of infrastructure to define in Terraform (see [[0018-infrastructure-as-code]]).

## Future Considerations
If API response caching at the edge becomes valuable (e.g., popular product summaries), add a second CloudFront behavior routing specific GET paths to the backend origin with short TTLs, rather than introducing a separate CDN product.
