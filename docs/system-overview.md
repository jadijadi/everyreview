# System Overview

EveryReview lets a user scan a product's barcode and immediately see reviews other users have written about it. See the [ADRs](../adrs/README.md) for the reasoning behind every decision referenced here.

## Primary flow

```mermaid
sequenceDiagram
    participant U as User
    participant A as Android App
    participant API as Backend API
    participant DB as PostgreSQL
    participant Cache as Redis

    U->>A: Scan barcode (CameraX + ML Kit)
    A->>API: GET /v1/products/{barcode}
    API->>Cache: lookup cached product+summary
    alt cache hit
        Cache-->>API: product + rating summary
    else cache miss
        API->>DB: SELECT product WHERE barcode = ?
        alt product exists
            DB-->>API: product row
        else not found
            API->>DB: INSERT placeholder product
            DB-->>API: new product row
        end
        API->>Cache: write-through cache
    end
    API-->>A: product + reviews + rating summary
    A-->>U: Show product screen

    U->>A: Submit review
    A->>API: POST /v1/products/{id}/reviews
    API->>DB: INSERT review (txn)
    API->>DB: update rating_summary
    API->>Cache: invalidate product cache entry
    API-->>A: 201 Created
    Note over A,API: Other users' next GET reflects the new review immediately (cache invalidated on write).
```

## High-level architecture

```mermaid
flowchart LR
    subgraph Clients
        AND[Android App]
        IOS[iOS App - Phase 2]
    end

    subgraph AWS
        ALB[Application Load Balancer]
        subgraph ECS[ECS Fargate Service]
            API[Go Backend - Modular Monolith]
        end
        RDS[(PostgreSQL - RDS)]
        REDIS[(Redis - ElastiCache)]
        S3[(S3 - Media)]
        CF[CloudFront CDN]
        SM[Secrets Manager]
    end

    AND -->|REST /v1, JWT| ALB
    IOS -->|REST /v1, JWT| ALB
    ALB --> API
    API --> RDS
    API --> REDIS
    API -->|presigned URLs| S3
    AND -->|image downloads| CF
    CF --> S3
    API -.reads secrets.-> SM
```

## Backend module boundaries (modular monolith — see [ADR-0003](../adrs/0003-modular-monolith-backend.md))

```mermaid
flowchart TD
    subgraph Backend[Go Backend Process]
        Auth[auth module]
        Product[product module]
        Review[review module]
        Media[media module]
        Moderation[moderation module - future]
    end
    Review -->|interface call, in-process| Product
    Review -->|interface call, in-process| Auth
    Media -->|interface call, in-process| Auth
    Moderation -.future.-> Review
```

Each module owns its own tables and exposes only Go interfaces to other modules — no cross-module SQL. This is what makes extracting a module into a standalone service later a contained refactor rather than a rewrite.

## Key documents
- [Database schema](database-schema.md)
- [API spec](../backend/api/openapi.yaml)
- [Deployment](deployment.md)
- [Development setup](development-setup.md)
- [Coding conventions](coding-conventions.md)
- [Contributing](contributing.md)
- [All ADRs](../adrs/README.md)
