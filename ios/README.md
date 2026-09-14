# iOS (Phase 2)

Not yet started. Planned: Swift + SwiftUI, mirroring the Android app's MVVM + Clean Architecture layering (see [ADR-0023](../adrs/0023-mobile-architecture.md)) against the same `backend/api/openapi.yaml` contract.

Before implementation begins, add ADRs for anything Phase-1 ADRs didn't already cover platform-neutrally — at minimum: iOS barcode-scanning approach (AVFoundation vs. Vision framework, analogous to [ADR-0025](../adrs/0025-barcode-scanning-library.md) for Android) and Sign in with Apple integration (required by App Store guidelines once this ships, anticipated in [ADR-0014](../adrs/0014-authentication-strategy.md)).
