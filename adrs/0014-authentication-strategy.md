# ADR-0014: Anonymous-First Authentication with JWT, Extensible to OAuth

## Status
Accepted

## Context
The core flow (scan → view reviews) should work with minimal friction, ideally before any signup. Writing a review requires some form of identity to prevent abuse and to let users manage their own reviews later. The spec requires: anonymous users initially, email login initially, OAuth providers later, all without a future rewrite.

## Problem
How do we authenticate/identify users across anonymous browsing, email login, and future OAuth providers, in a way that's consistent across Android and iOS?

## Alternatives
- **Require login for everything** — simplest backend logic, but adds friction to the core "scan and see reviews" value proposition before a user has any reason to trust the app; rejected.
- **Fully anonymous, device-ID-only, forever** — zero friction, but no path to protecting a user's review history across devices/reinstalls, and no foundation for reputation (long-term feature) or moderation accountability.
- **Anonymous device identity by default, upgradeable to a registered identity, backend-issued JWTs, provider-agnostic identity linking table** — every client (even before any login) gets a backend-issued anonymous account (a real `user` row with `auth_provider = 'anonymous'`, tied to a device via a securely stored refresh token) so it can read and even submit reviews immediately. A user can later "upgrade" that same account by linking an email/password or OAuth identity, preserving their review history. Auth providers (email, Google, Apple, etc.) are stored in a separate `identity` table (`user_id`, `provider`, `provider_subject`) rather than columns on `user`, so adding a new OAuth provider is a data change, not a schema change.

## Decision
- Every client obtains an anonymous account on first launch (`POST /v1/auth/anonymous` → backend creates a `user` row + issues tokens). No email/password required to browse or submit reviews initially, keeping the product's friction low while still tying every review to an accountable, rate-limitable identity (mitigating anonymous abuse, see [[0015-rate-limiting]]).
- Auth uses short-lived JWT access tokens (15 min) + long-lived opaque refresh tokens (stored hashed server-side, rotated on use), issued by the backend's own `auth` module — not a third-party auth service — since requirements are simple enough that a managed IdP (Cognito, Auth0) would add cost/lock-in without a corresponding win at this stage.
- Email login adds a password credential (Argon2id-hashed) linked to the same account via the `identity` table; "upgrading" from anonymous to email-authenticated is a linking operation on the existing `user_id`, not a new account.
- OAuth providers (Google, Apple — required for iOS App Store "Sign in with Apple" rules once iOS ships) are added later purely as new rows in the `identity` table plus a provider-specific token-verification adapter behind a common `IdentityProvider` interface; no schema change and no change to how already-issued JWTs work.

## Consequences
- Low-friction onboarding: the core scan-and-read flow needs no signup screen.
- Reviews are always attributable to a stable `user_id` even before "real" login, enabling moderation/rate-limiting from day one.
- The backend owns credential storage and token issuance, meaning it must implement password hashing, token rotation, and OAuth token verification correctly and securely — more responsibility than delegating to a managed IdP, accepted for cost/control reasons at this stage.
- Anonymous accounts that are never upgraded accumulate; a retention/cleanup policy (e.g., purge anonymous accounts with no activity after N months) will be needed.

## Future Considerations
If auth requirements grow substantially (SSO for an eventual admin dashboard, enterprise requirements, more OAuth providers than justify custom adapters), migrate to a managed IdP (AWS Cognito or Auth0) fronting the same `identity`-table model — the `IdentityProvider` interface boundary is designed to make that swap contained.
