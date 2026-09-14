# Android

Kotlin + Jetpack Compose, MVVM + Clean Architecture. See [ADR-0023](../adrs/0023-mobile-architecture.md), [ADR-0024](../adrs/0024-android-tech-stack.md), [ADR-0025](../adrs/0025-barcode-scanning-library.md).

## Layout (planned)
```
android/
  app/
    src/main/java/app/everyreview/
      presentation/<feature>/     Compose screens + ViewModels
      domain/<feature>/            Use cases, plain Kotlin, no Android imports
      data/<feature>/              Repositories, Retrofit API, Room DAOs
      di/                          Hilt modules
    src/test/                      Unit tests (domain, ViewModels)
    src/androidTest/               Instrumented UI tests
```

## Local development
See [`docs/development-setup.md`](../docs/development-setup.md). Requires Android Studio, JDK 17, an emulator with Google Play services (ML Kit dependency).

## Generating the API client
Network models/interfaces are generated from `backend/api/openapi.yaml` — do not hand-edit generated code. Run `scripts/generate-api-clients.sh` after the spec changes.

## Testing
Unit tests for domain/ViewModel layers on every PR; instrumented Compose UI tests for the golden path (scan → view → submit review) run against an emulator in CI. See [ADR-0022](../adrs/0022-testing-strategy.md).

## MVP scope (current state)

This is a first, intentionally small implementation to test the core scan → view → submit loop end to end, matching `backend/README.md`'s "MVP scope" section. It keeps the architecture/pattern from the ADRs above (MVVM + Clean Architecture, Hilt, Room, Retrofit, CameraX + ML Kit) but trims a few things, flagged here rather than silently diverged, per `docs/coding-conventions.md`:

- **No login/auth screens.** There's no account system on the backend yet (see `backend/README.md`). The nickname used on submitted reviews is just a locally remembered `DataStore` preference (`data/prefs/NicknamePrefs.kt`), not an authenticated identity.
- **No generated API client.** `data/remote/EveryReviewApi.kt` is hand-written Retrofit against the actual (trimmed) backend endpoints, not generated via `scripts/generate-api-clients.sh` from `backend/api/openapi.yaml` — deferred until the API stabilizes past MVP.
- **Room cache is intentionally thin.** Two entities (`ProductEntity`, `ReviewEntity`) cache the last-fetched product/reviews as a network-first, cache-as-fallback read path — no sync/invalidation strategy beyond that (ADR-0024 flags this as a later design item).
- **Moshi codegen (KSP) instead of `openapi-generator`-produced models** — DTOs in `data/remote/dto/` are hand-written and annotated with `@JsonClass(generateAdapter = true)`.
- Only four screens exist: Scan → Product → Write Review / Add Product Details. No instrumented Compose UI tests yet (unit tests for ViewModels/use cases only) — add those once the manual flow is validated on a device/emulator.

### Screens
- **Scan** (`presentation/scan/`) — CameraX preview + ML Kit on-device barcode detection (EAN-13/8, UPC-A/E, Code128, QR per [ADR-0025](../adrs/0025-barcode-scanning-library.md)). On first successful decode, navigates to Product.
- **Product** (`presentation/product/`) — looks up the product by barcode (creating a placeholder server-side if unknown), shows name/brand/rating summary and the review list, refreshing reviews whenever the screen resumes.
- **Add Product Details** (`presentation/addproduct/`) — offered on Product when the backend returned a placeholder (nobody has described this barcode yet). Name is mandatory; brand, manufacturer, price + currency (prefilled from the device locale), description and a photo (camera via `FileProvider`, or the system photo picker) are optional. `data/media/PhotoPreparer.kt` downscales the photo to ≤1600 px JPEG before it goes into the multipart `POST /v1/products/{id}/details`. First writer wins server-side; a `409` is shown as "someone already added details".
- **Write Review** (`presentation/review/`) — nickname (remembered), optional 1–5 star rating, review text; submits and pops back to Product.

### Quick start
1. Start the backend first (see `backend/README.md`'s quick start) — the emulator reaches it at `http://10.0.2.2:8080/v1/`, which is the default `API_BASE_URL` baked into `app/build.gradle.kts` if you don't override it.
2. Copy `local.properties.example` to `local.properties` and set `sdk.dir` for your machine (a `local.properties` is already present in this checkout for the machine it was built on).
3. Open `android/` in Android Studio, or from the command line:
   ```bash
   cd android
   ./gradlew :app:assembleDebug
   ./gradlew :app:testDebugUnitTest
   ```
4. Run on an emulator with Google Play services (required for ML Kit) or a physical device to exercise the actual scan → view → submit flow — this requires a camera and hasn't been exercised headlessly.

### Versioning
- `versionName` comes from `APP_VERSION` in `android/gradle.properties` (semver, bump it by hand for a release).
- `versionCode` is the git commit count (`git rev-list --count HEAD`), so it increases monotonically with every commit.
- Debug builds get a `-dev.<short sha>` suffix, and `BuildConfig.GIT_SHA` is always available. The Scan screen shows the logo plus `v<versionName> (<versionCode>)` in the top-left corner (`presentation/common/AppBranding.kt`).
- The launcher icon is an adaptive icon (`res/mipmap-anydpi-v26/`, foreground/monochrome vectors in `res/drawable/`); `res/drawable/ic_logo.xml` is the same artwork for in-app use.

To point at a different backend (e.g. a physical device on the same Wi-Fi, or a deployed dev environment), pass `-PAPI_BASE_URL=http://<host>:8080/v1/` to Gradle.
