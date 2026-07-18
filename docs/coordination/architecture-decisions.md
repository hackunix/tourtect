# Architecture Decisions — Locked

These decisions are locked for the first implementation round.
Do not replace without documenting a critical technical blocker and receiving explicit approval.

## AD-001: Go Backend

The primary backend is a Go modular monolith using `net/http`.
Three binaries (`api`, `realtime`, `worker`) share internal domain packages.
No Gin, Fiber, NestJS, Fastify, or Next.js Route Handlers as primary backend.
No separately deployed microservices per domain module.

## AD-002: Kotlin Android

Native Android using Kotlin, Jetpack Compose, Hilt.
No Flutter, React Native, Kotlin Multiplatform, or WebView shell.

## AD-003: Next.js Web Client

Next.js App Router with TypeScript strict mode.
Web is a client of the Go backend — not a shadow backend.
No direct PostgreSQL access from web. No Price/Safety Engine logic in web.

## AD-004: OpenAPI as Contract Source of Truth

`backend/api/openapi.yaml` is the single source of truth.
Generation direction: OpenAPI → Go server types → TypeScript client → Kotlin client.
Clients must not manually redefine shared domain types.

## AD-005: PostgreSQL as System of Record

PostgreSQL 16 with PostGIS is the runtime source of truth.
All runtime endpoints must query PostgreSQL — no mock arrays.

## AD-006: PostgreSQL-First Search

Full-text search, `pg_trgm`, `unaccent`, PostGIS distance search.
OpenSearch remains optional until PostgreSQL-backed search works.

## AD-007: WebSocket for PTT

`coder/websocket` for realtime transport.
JSON text frames for control, binary frames for PCM audio.
No WebRTC in V1.

## AD-008: No Production Python Request Service

Python allowed only under `tools/ml/` for offline tasks.
Runtime AI adapters implemented in Go.

## AD-009: Modular Monolith

Single Go module. Shared internal packages. Three binaries.
No Kafka, RabbitMQ, NATS, or Kubernetes.
PostgreSQL outbox pattern for background jobs.

## AD-010: Backend Gate Before Frontend

`BACKEND_READINESS_GATE=PASS` required before web/Android implementation.
Frontend must never silently fall back to mock data.

## AD-011: Container Runtime

Podman 6.0.1 replaces Docker. `podman compose` used for all infrastructure.

## AD-012: Provider Secrets Server-Only

FPT AI API keys read only by the Go backend.
Never in Android APK, never in web client, never in logs.
