# Ingestion JSON Contract

This document specifies the provider-independent JSON data contract for the Tourtect ingestion pipeline. It allows Tourtect to safely extract, normalize, and validate travel-critical data from various source documents without coupling the downstream database schemas to any specific LLM provider.

## Contracts Overview

All contracts are defined in JSON Schema Draft 2020-12 and live under `backend/internal/ingestion/contracts/`.

| Contract / Schema | Purpose |
| --- | --- |
| [common.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/common.schema.json) | Shared structures such as `monetary_amount`, `coordinates`, and metadata envelopes (source, quality, provenance, lifecycle). |
| [discovery-candidate.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/discovery-candidate.schema.json) | Represent URLs, feeds, or sitemaps discovered by the crawler prior to any fetching. |
| [source-policy.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/source-policy.schema.json) | Rules governing allowed content types, rate limits, robots directives, and data access modes for a host. |
| [fetch-request.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/fetch-request.schema.json) | Bounded request detailing timeouts, priorities, and policies to the sandbox fetch proxy. |
| [fetch-result.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/fetch-result.schema.json) | Structured results from the fetcher, capturing latency, cache hits, redirect chains, and robots decisions. |
| [raw-source-document.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/raw-source-document.schema.json) | Bounded visible text divided into block fragments for granular claim referencing. |
| [llm-extraction-request.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/llm-extraction-request.schema.json) | Prompt parameters and constraints sent to the external provider wrapper. |
| [llm-extraction-response.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/llm-extraction-response.schema.json) | Discriminated union of raw extracted records matching the permitted content types. |
| [normalized-record.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/normalized-record.schema.json) | Schema holding deterministic unit conversions, timezone normalizations, and entity resolution results. |
| [validation-result.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/validation-result.schema.json) | Results of the pipeline rule checks, containing machine-readable reason codes for rejected or quarantined records. |
| [review-decision.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/review-decision.schema.json) | Audit trail of manual human editor or trusted authority approval / rejection actions. |
| [publishable-record.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/publishable-record.schema.json) | Redacted, clean public-facing payload mapped to standard Tourtect domain entities. |
| [refresh-state.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/refresh-state.schema.json) | Adaptive poll trackers to manage ETag check intervals and change frequencies. |
| [deletion-event.schema.json](file:///home/hanitav/Documents/GitHub/tourtect/backend/internal/ingestion/contracts/deletion-event.schema.json) | Takedown or expiration sync requests triggered by source changes or legal demands. |

---

## Discriminator Design (Polymorphic Content)

To prevent a massive schema with hundreds of nullable fields, the content payloads are modeled as a discriminated union under `content` based on `record_type`.

Supported content types:
- `official_alert`: Weather warnings, safety closures, and regional threats issued by verified authorities.
- `place`: Name, category, address, coordinates, and contact directories.
- `place_status`: Closure state (temporarily closed, restricted, permanently closed).
- `price_observation`: Raw transactions capturing amounts, currencies, units, and contexts.
- `commerce_offer`: platform offers with platform fees, delivery conditions, and vouchers.
- `scam_pattern`: Objective behavioral sequences, preconditions, and recommended tourist safety actions.
- `emergency_directory_entry`: Country and regional hotlines verified by a trusted authority.
- `news_article` / `video`: Travel-critical snippets and metadata (no full re-hosting without license).

---

## Forward-Compatible Migration Guidance

As ingestion needs grow, contracts can be versioned safely:
1. **Backward-Compatible Changes**: Add new fields as optional, and avoid changing existing field types.
2. **Major Upgrades**: If breaking changes are required (e.g. changing currency representations), increment the schema `version` attribute and register a deterministic transformation function inside the pipeline.
