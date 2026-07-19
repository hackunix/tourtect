# Ingestion Schema Assessment

## Documents Selected
1. `AGENTS.md`
2. `CONTEXT_MAP.md`
3. `README.md`
4. `backend/api/openapi.yaml`
5. `compose.yaml`
6. `.env.example`
7. `01-product-experience/07-external-content.md`
8. `02-functional-requirements/03-external-content-monetization.md`
9. `03-architecture/03-modular-monolith-components.md`
10. `04-data/03-price-safety-models.md`
11. `02-functional-requirements/07-scam-emergency.md`
12. `backend/db/migrations/003_places.sql`, `005_pricing.sql`, `006_safety.sql`, `009_ai_intelligence.sql`

## Existing ingestion code
- Found `backend/internal/ingestion/outbox.go`. No significant external ingestion pipeline contracts or code exist yet.

## Existing database tables
- **Places**: `places`, `place_aliases`
- **Pricing**: `price_snapshots`, `price_observations`
- **Safety**: `safety_directory_versions`, `safety_directory_entries`, `consent_records`
- **AI Intelligence**: `assistant_feedback`, `assistant_confirmation_audit`, `assistant_model_traces`

## Existing price and place models
- `places` table requires UUID, name, category, region_id, coordinates (PostGIS Point), etc.
- `price_observations` captures transactions with `vertical`, `region_id`, `amount_minor`, `currency`, `exponent`, `unit`, `service_segment`, `venue_type`, `transaction_context`, and `extraction_confidence`.
- `price_snapshots` aggregates price data with robust intervals (p10, p50, p90) and `effective_from/to`.

## Existing safety models
- `safety_directory_entries` for verified hotlines (`police`, `ambulance`, `fire`, `embassy`, `hotline`, `crisis_center`, `tourist_police`).

## Existing external-source models
- None currently modeled in the database schema. Need separate tables for `source_policies`, `ingestion_jobs`, `source_documents`, `extracted_records`, etc.

## Existing crawl policy rules
- Defined in `01-product-experience/07-external-content.md`. Must respect `robots.txt`, Terms of Service, API policies, and `Retry-After`. Must not circumvent CAPTCHA, logins, or private endpoints. Adaptive polling with Jitter. Token bucket rate limit per host.

## Missing contracts
- All JSON schema contracts requested: `DiscoveryCandidate`, `SourcePolicy`, `FetchRequest`, `FetchResult`, `RawSourceDocument`, `LLMExtractionRequest`, `LLMExtractionResponse`, `NormalizedRecord`, `ValidationResult`, `DeduplicationResult`, `EntityLinkResult`, `ReviewDecision`, `PublishableRecord`, `RefreshState`, `DeletionEvent`, `IngestionJob`, `IngestionTrace`.

## Contract conflicts
- The `price_observations` table stores raw user reports and confirmed items. External source data must be kept in a separate staging area before being merged into `price_observations` to prevent untested AI output from polluting production data.
- The `safety_directory_entries` require rigorous verification; AI extractions cannot be directly inserted without an approval boundary.

## Reusable components
- Common fields (identity, source, provenance, extraction, quality, policy, lifecycle) can be factored into `$defs` in a `common.schema.json` to be referenced by all other schemas.
- `Amount` and `Currency` models can be standardized across pricing records.

## Required migrations
- Need new SQL migrations to store ingestion staging data (e.g., `ingestion_jobs`, `source_policies`, `source_documents`, `extraction_runs`, `extracted_records`, `normalized_records`, `evidence_claims`, `entity_links`, `review_decisions`, `published_records`, `deletion_events`, `ingestion_traces`), mapping the JSON contracts into the DB.
