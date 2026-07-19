# Tourtect AI-Native Assessment

Assessment date: 2026-07-19  
Baseline commit inspected: `473dcb7` plus the working tree present at assessment time  
Affected requirements: PC-01–PC-12, SE-01–SE-08, new AIN-01–AIN-18 contract set below

## Executive finding

Before this transformation, Tourtect was a verified community, Price, and Safety platform with several AI/realtime adapter-shaped components. It was not AI-native. The default Web experience was the community feed; Android was a disconnected scaffold; each Price or Safety request was stateless and structured; realtime used a separate in-memory session; Lens had no implementation; and missing AI credentials produced plausible fake realtime answers.

The reusable foundation is strong: deterministic Price and Safety engines, versioned price snapshots and safety directory data, places and aliases, public community posts, PostgreSQL, Redis infrastructure, the outbox, Go provider interfaces, and a realtime process already exist. The safe transformation path is an application layer inside the Go modular monolith, not replacement of those systems.

## 1. What was already AI-native?

No complete product flow met the nine-part AI-native definition.

Parts that were architecturally reusable were:

- Backend-only ASR, Translation, Vision, and Extraction provider interfaces.
- A realtime PTT state machine that accepted audio and produced transcript/translation event shapes.
- A deterministic Price Engine and rule-first Safety Engine that can remain authoritative tools.
- Place aliases, public posts, verified price snapshots, and versioned safety-directory records that can ground responses.
- Explicit design rules for consent-before-offload, provider traces, offline safety, and no LLM ownership of price or emergency decisions.

These were enabling components, not a single AI-native assistance experience.

## 2. What was an isolated AI feature?

- Realtime ASR and translation lived in the WebSocket process and did not feed the REST Price/Safety engines, place retrieval, community retrieval, or a durable assistance session.
- Provider adapter methods existed separately from the community/Price/Safety application flow.
- Web exposed deterministic Price and Safety demo actions in `ContextPanel`, using fixed taxi and dispute examples rather than natural input.
- Live and Lens Android modules existed only as empty Gradle modules.
- The so-called real ASR, Vision, and Extraction methods returned hard-coded successful Vietnamese examples; realtime selected fake successful providers when the secret was absent.

## 3. Where was context lost?

- Every REST Price Check and Safety Assessment was independent.
- Realtime session state was held in a process-local map and deleted on disconnect.
- Realtime accepted a caller-provided session ID without binding it to the authenticated user or a REST assistance session.
- There was no shared locale, target locale, current place, recent price candidate, user-confirmed facts, consent state, recent tool result, or current safety state.
- Web held only page/component state. Android had no repository or ViewModel session state.
- No final session summary, quarantined correction, redacted AI trace, or confirmation audit existed.

## 4. Which flows required manual tool choice?

- Web navigation made users start in Community, Search, Price Check, or Safety-like demo controls.
- Price Check required vertical, amount, currency, unit, region, segment, venue, context, and time.
- Safety Assessment required users or demo code to choose structured facts and indicators.
- Realtime required entry into a separate PTT connection and returned only transcript/translation.
- Android exposed only a static welcome screen plus disconnected Forum and Safety composables.
- There was no place where speech, text, a confirmed image candidate, or a structured fallback entered the same orchestration flow.

## 5. Which APIs were missing?

The baseline OpenAPI lacked:

- Assistance-session create/resume/delete.
- Typed multimodal messages.
- Structured assistant responses, evidence, and tool results.
- Server-issued consequential confirmations.
- Feedback/correction quarantine.
- Safe capture consent/finalization.
- A contract linking realtime utterances to owned assistance sessions.

Capture endpoints must remain absent until MinIO storage, redaction, explicit media consent, failure retention, ownership, and cleanup are implemented. Lens must not upload against an invented contract.

## 6. Which modules can be reused?

- `internal/pricing.Engine.Evaluate` behind an allowlisted `evaluate_price` adapter.
- `internal/safety.Engine.Assess` behind an allowlisted `evaluate_safety` adapter.
- `internal/places.Service` for canonical place and alias resolution.
- `internal/content.Service` for bounded public community evidence linked to a place.
- `price_snapshots`, `price_observations`, places, post-place links, public posts, safety-directory versions, approved entries, audit events, and outbox events.
- Redis infrastructure and configuration for active session TTL state.
- Backend-only FPT/OpenAI-compatible translation adapter after consent and configuration checks.
- Existing Web Price/Safety/error components as typed assistant cards.
- Existing realtime state machine after protocol, identity, idempotency, and production-mock hardening.

## 7. What was mock-only or disconnected?

- Production realtime silently selected `FakeASR` and `FakeTranslation` without a provider secret.
- `RealASR`, `RealVision`, and `RealExtraction` returned hard-coded successful content without calling a provider.
- Realtime binary input was treated as raw PCM, not the documented protobuf `MediaChunk`; it had no media ACK, resume token, durable session, or reliable server event sequence.
- Price and Safety realtime event constants existed but no corresponding lanes ran.
- Web used a handwritten API client and had no Assistant endpoints or session ownership.
- Android DTOs did not reliably match the snake_case/nested OpenAPI data, and no Retrofit client construction connected them.
- Android Live, Lens, Room, security, and design-system modules were mostly empty scaffolds.
- `make android-build` and `make android-test` returned success while skipping because `android/gradlew` was missing.

## 8. What must change without breaking the verified backend?

1. Add a Tourtect Intelligence Layer inside the Go monolith.
2. Keep Price and Safety decisions immutable and deterministic.
3. Store bounded, redacted, user-owned active context in Redis with TTL and an explicit version.
4. Run critical safety phrase rules before normal routing.
5. Restrict tool execution to registered schemas; never expose SQL, arbitrary HTTP, storage, or shell execution.
6. Retrieve bounded public/verified evidence before composing grounded answers.
7. Return typed response objects with evidence IDs, tool-result IDs, dataset versions, factors, missing information, and fallback state.
8. Require a server-issued, expiring confirmation ID for consequential actions and audit the decision.
9. Quarantine corrections and contributions; do not update trusted datasets automatically.
10. Remove successful production mocks and make provider degradation visible.
11. Make `/assistant` the Web default and preserve Community as the evidence/feedback network.
12. Make the native Assistant the Android default while keeping PTT/Lens unavailable until their backend contracts and project build foundation are real.

## Frozen architecture decisions

- The Intelligence Layer is part of the Go modular monolith.
- Redis holds active assistance sessions; PostgreSQL holds feedback, confirmation audit, and redacted model traces.
- Raw audio and images never enter assistance-session JSON.
- Session defaults: 30-minute sliding TTL, schema version 1, 256 KiB maximum serialized size, 20 recent responses, 64 processed message IDs, 32 confirmed facts, and 8 active capture references.
- The initial router is deterministic/rule-first and schema-constrained. A model-assisted candidate router may be added only behind the same schema and safety pre-check.
- Retrieval is limited to public or verified allowed sources, six items, approximately 6,000 characters, locale preference, freshness labels, duplicate removal, and basic PII redaction.
- The response composer uses deterministic templates as the guaranteed fallback. It cannot change engine-owned fields.
- Capture APIs are deliberately deferred; no continuous video, facial identification, or unconsented media upload is permitted.

## New contract requirements

| ID | Requirement |
| --- | --- |
| AIN-01 | Assistant is the primary entry surface on supported clients. |
| AIN-02 | A user may begin with text, voice transcript, confirmed image capture reference, or structured fallback facts. |
| AIN-03 | Assistance sessions are user-owned, versioned, bounded, redacted, and TTL-controlled. |
| AIN-04 | Critical safety rules run before or alongside intent routing. |
| AIN-05 | Router output is schema-constrained and selects only allowlisted tools. |
| AIN-06 | Every tool invocation records a trace ID, duration, status, and bounded output. |
| AIN-07 | Price alert levels come only from the Price Engine. |
| AIN-08 | Safety urgency and emergency numbers come only from the Safety Engine and versioned directory. |
| AIN-09 | Assistant responses preserve evidence and provenance. |
| AIN-10 | The composer cannot mutate deterministic facts and has a non-model fallback. |
| AIN-11 | Consequential actions require an unexpired server-issued confirmation bound to session, user, and action. |
| AIN-12 | Feedback and corrections enter quarantine and never update trusted data directly. |
| AIN-13 | AI traces exclude raw media, credentials, full sensitive transcripts, and chain-of-thought. |
| AIN-14 | Provider failure produces an explicit degraded state, never a plausible fake answer. |
| AIN-15 | PTT final transcripts feed the same owned assistance session and preserve utterance IDs. |
| AIN-16 | Lens upload remains unavailable until consent, capture, finalize, retention, and cleanup contracts exist. |
| AIN-17 | Community remains accessible and supplies supporting public knowledge without controlling deterministic decisions. |
| AIN-18 | Offline/manual Price, Safety, directory, phrasebook, and private-draft paths remain available. |

## Initial verification baseline

- Backend: tests passed with a writable Go build cache and available local dependencies; package tests that require PostgreSQL cannot run in a network-restricted sandbox.
- Web: lint and 3 initial unit tests passed. The initial sandbox build failed only while fetching Google Fonts; a network-enabled verification passed.
- Android: no trustworthy build/test result existed because the Gradle wrapper and essential application files were absent; Make targets skipped and returned false success.

This assessment is intentionally the pre-change baseline. Completion claims belong in the final transformation report and must distinguish implemented vertical slices from contract-gated PTT/Lens limitations.
