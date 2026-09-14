# ADR-0024: Android Tech Stack — Kotlin, Jetpack Compose, Hilt, Coroutines/Flow

## Status
Accepted

## Context
Android is the first mobile client (Phase 1), implementing the MVVM + Clean Architecture pattern from [[0023-mobile-architecture]], and must support the core scan → view → review flow plus an offline-tolerant local cache.

## Problem
What are the concrete Android libraries/tools for language, UI, DI, concurrency, and local storage?

## Alternatives considered per concern
- **Language**: Kotlin (vs. Java) — Kotlin is Google's recommended, modern default for Android: null-safety, coroutines, concise syntax. No real alternative justified for a new project in 2026.
- **UI toolkit**: Jetpack Compose (vs. XML/View system) — Compose is the current standard for new Android UI, declarative (matches the MVVM state-driven model naturally), and actively developed; the View system is legacy-maintenance mode. Chosen.
- **Dependency Injection**: Hilt (vs. Koin, manual DI) — Hilt is Google's recommended DI for Android, compile-time-checked (catches wiring errors at build time, not runtime like Koin's service-locator style), integrates with Jetpack components (ViewModel, WorkManager).
- **Async**: Kotlin Coroutines + Flow (vs. RxJava) — idiomatic Kotlin, first-class language support, `Flow` maps naturally onto reactive UI state in Compose; RxJava is more powerful but unnecessary complexity for this app's needs and not idiomatic Kotlin.
- **Local persistence/offline cache**: Room (vs. raw SQLite, DataStore-only) — Room gives type-safe SQL access over SQLite with coroutine/Flow support, used for the offline cache of recently viewed products/reviews (per the long-term "offline caching" requirement); `DataStore` (not Room) handles simple key-value app preferences (auth tokens reference — see note below, last-used barcode, etc).
- **Networking**: Retrofit + OkHttp, with client models generated from `backend/api/openapi.yaml` via `openapi-generator`'s Kotlin client generator — keeps the Android network layer mechanically in sync with the backend contract (see [[0006-api-style-rest]]) rather than hand-maintained.
- **Camera/barcode**: see [[0025-barcode-scanning-library]] (separate ADR given its significance).
- **Image loading**: Coil — Kotlin-first, Compose-native image loading with built-in caching, simplest fit for Compose over Glide (View-system-oriented) or Picasso (less actively developed).

## Decision
Kotlin + Jetpack Compose + Hilt + Coroutines/Flow + Room + Retrofit/OkHttp (OpenAPI-generated client) + Coil, structured per the MVVM + Clean Architecture layering in [[0023-mobile-architecture]]. Auth tokens (see [[0014-authentication-strategy]]) are stored via `EncryptedSharedPreferences`/Android Keystore-backed storage, never in Room or plain DataStore.

## Consequences
- Fully aligned with current (2026) Google-recommended Android development practices, maximizing long-term library support and hiring pool familiarity.
- Compose's relative newness compared to the View system means some third-party libraries have less mature Compose support — evaluated case by case (barcode scanning UI overlay is the main risk area, addressed in [[0025-barcode-scanning-library]]).
- Room-based offline cache requires an explicit cache-invalidation/sync strategy (documented in `android/README.md` once implemented) to avoid showing stale reviews indefinitely.

## Future Considerations
Revisit Kotlin Multiplatform adoption per [[0023-mobile-architecture]]'s future-considerations once iOS exists; the domain/data layer here is already structured to make that low-risk.
