# ADR-0011: S3 for Object Storage and Image Uploads

## Status
Accepted

## Context
The initial product model already has an optional product image; reviews will later support photos and video. Images/media need durable storage, must be servable at scale to a global mobile user base, and uploads happen from untrusted mobile clients.

## Problem
Where and how are uploaded images/media stored, and how do untrusted clients upload them safely?

## Alternatives
- **Store binary blobs in PostgreSQL (bytea)** — simplest conceptually, but bloats the database, hurts backup/restore times, and databases are not built to serve high-volume binary reads; ruled out immediately at any real scale.
- **Backend proxies all uploads/downloads through itself to storage** — centralizes validation, but doubles bandwidth costs and load on the backend for every image (upload AND every subsequent view), and doesn't scale independently of the API.
- **Amazon S3 with presigned URLs for direct client upload, served via CloudFront (see [[0012-cdn-cloudfront]])** — clients (mobile apps) upload directly to S3 using a short-lived presigned URL issued by the backend (after the backend validates the request — auth, content-type, size limit), and read traffic is served via CDN, not through the backend at all. The backend never touches raw image bytes on the hot path.

## Decision
S3 for all object storage (product images, review photos, and future video). Upload flow: client requests a presigned upload URL from the backend (backend validates auth, enforces content-type/size constraints, generates a unique object key), client uploads directly to S3, client then confirms completion to the backend which records the object key against the product/review. Downloads are served via CloudFront in front of S3, never directly from S3 and never proxied through the backend (see [[0012-cdn-cloudfront]]). A background job (or synchronous Lambda-on-upload trigger) performs virus/content scanning and generates resized thumbnail variants.

## Consequences
- Upload/download bandwidth never transits the backend, keeping it lightweight and cheap to scale.
- Requires careful presigned-URL scoping (short expiry, constrained content-type/size, unique per-request key) to prevent abuse.
- Orphaned uploads (client gets a presigned URL but never confirms) need a periodic cleanup job (lifecycle rule on a temp prefix, or a scheduled reconciliation job).

## Future Considerations
When AI moderation (per the long-term feature list) is introduced, hook it into the same upload-confirmation flow as an async moderation step before an image is considered "published" and served publicly.
