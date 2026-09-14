# ADR-0025: CameraX + ML Kit Barcode Scanning for Android

## Status
Accepted

## Context
Barcode scanning is the core, first-touch interaction of the entire app (see project overview: "user scans a barcode using the phone camera" is step 2 of the primary flow). It must be fast, reliable across common barcode formats (UPC-A/E, EAN-8/13, and ideally QR/Code128 for future-proofing), and work well within a Compose UI (see [[0024-android-tech-stack]]).

## Problem
What camera and barcode-decoding libraries does the Android app use?

## Alternatives
- **Camera2 API directly** — maximum control, but low-level and verbose; CameraX exists specifically to wrap it with a simpler, lifecycle-aware API without giving up necessary control (frame analysis access for barcode decoding). No reason to drop to raw Camera2.
- **CameraX + ZXing (`zxing-android-embedded`)** — mature, widely used, works offline with no ML model download, but decode speed/accuracy on difficult scans (poor lighting, angled/damaged barcodes — realistic in-store conditions) is generally weaker than ML-based decoders, and the project (a wrapper library) is less actively maintained than Google's own stack.
- **CameraX + ML Kit Barcode Scanning API** — Google's on-device ML Kit barcode scanner integrates directly with CameraX's `ImageAnalysis` use case, runs fully on-device (no network round-trip, no privacy concern from sending camera frames off-device), supports all common 1D/2D formats out of the box, and is actively maintained by Google alongside Compose/CameraX/Android itself — the most cohesive fit with the rest of the [[0024-android-tech-stack]] choices.

## Decision
CameraX for camera lifecycle/preview management, with ML Kit's on-device Barcode Scanning API wired into a CameraX `ImageAnalysis` use case for real-time decode. Scanning restricts to relevant formats (`EAN_13`, `EAN_8`, `UPC_A`, `UPC_E`, plus `CODE_128` and `QR_CODE` for forward compatibility with non-retail-barcode products) rather than scanning all formats, both for decode performance and to reduce false-positive matches. On successful decode, the raw barcode value is sent to `GET /v1/products/{barcode}` (see `backend/api/openapi.yaml`).

## Consequences
- On-device decoding means no camera frames ever leave the device for the scan step itself — good for privacy and works with poor/no connectivity right up until the lookup call.
- ML Kit's barcode model adds to app size versus a pure-code decoder like ZXing, an accepted trade-off for decode accuracy on the app's single most important interaction.
- Tight coupling to Google's ML Kit / CameraX APIs is acceptable since this is Android-specific; iOS will independently use its own native equivalent (AVFoundation's built-in `AVCaptureMetadataOutput` barcode detection or Apple's Vision framework), decided in its own ADR when iOS work begins (Phase 2), consistent with [[0023-mobile-architecture]]'s "shared pattern, not shared code" approach to platform-specific concerns like camera access.

## Future Considerations
When OCR of product labels is added (long-term roadmap item), evaluate whether it reuses the same CameraX pipeline (as a second `ImageAnalysis` use case or mode toggle) or is a separate capture flow — write a new ADR at that time rather than assuming reuse now.
