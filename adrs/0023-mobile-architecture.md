# ADR-0023: Shared Mobile Architecture — MVVM + Clean Architecture, Platform-Native UI

## Status
Accepted

## Context
Android ships first (Phase 1), iOS follows (Phase 2), and both must evolve independently without breaking each other while sharing as much *architecture* (not necessarily code, since Kotlin and Swift don't share a runtime without extra tooling) as possible, per the project's platform requirements.

## Problem
What architectural pattern do the mobile clients follow, and how much actual code is shared between Android and iOS?

## Alternatives
- **Cross-platform framework (Flutter, React Native, KMP with shared UI)** — maximizes code sharing including UI, but the spec explicitly calls for native stacks (Kotlin/Compose for Android, Swift/SwiftUI for iOS) presumably for native camera/barcode-scanning performance and platform-idiomatic UX, which is the harder-to-reverse decision here (see [[0025-barcode-scanning-library]]); ruled out by the given platform requirements rather than re-litigated in this ADR.
- **Kotlin Multiplatform (KMP) for shared business logic** — genuine code sharing (networking, models, use cases) between Android and iOS while keeping native UI on each platform. Attractive, but adds build-tooling complexity (KMP + iOS interop, Xcode integration) before there's a second platform to share code *with* — Phase 1 is Android-only.
- **No code sharing, parallel architecture only** — each platform independently implements the same *pattern* (MVVM + Clean Architecture layering) against the same OpenAPI contract, so the two codebases are structurally analogous and a developer familiar with one can navigate the other, but no runtime code is shared.

## Decision
Each platform uses **MVVM + Clean Architecture**, layered identically in concept:
- **Presentation** (Compose/SwiftUI + ViewModel) — UI state and user intent handling, no business logic.
- **Domain** (Use Cases / Interactors, plain Kotlin/Swift, no framework dependencies) — business rules (e.g., "submitting a review requires non-empty text"), independent of Android/iOS SDKs, fully unit-testable.
- **Data** (Repositories + data sources: remote API client generated from the shared OpenAPI spec, plus local cache/offline store) — implements domain-layer repository interfaces.

Both platforms consume the **same OpenAPI spec** (`backend/api/openapi.yaml`, see [[0006-api-style-rest]]) via generated client code (Retrofit/Moshi models on Android via `openapi-generator`, equivalent Swift codegen on iOS later), which is the actual mechanism of "shared architecture" between platforms — a consistent contract and layering pattern, not shared compiled code. Kotlin Multiplatform is explicitly deferred, not adopted, for Phase 1.

## Consequences
- Domain and data layers are fully unit-testable without Android/iOS framework dependencies.
- A developer moving from the Android codebase to the iOS codebase (once it exists) finds the same layering and naming conventions, reducing ramp-up time.
- No runtime code sharing means some logic (e.g., client-side validation rules) is genuinely duplicated across platforms; each duplication point should reference the domain rule's ADR or the OpenAPI spec as the source of truth to keep both implementations honest.

## Future Considerations
Once iOS exists and both platforms have stabilized, revisit Kotlin Multiplatform for the domain/data layers specifically (business rules, API client, models) if duplicated logic proves to be a recurring source of drift/bugs between platforms. This ADR's layering already puts that logic in framework-independent Kotlin, which is a prerequisite for a low-risk KMP adoption later.
