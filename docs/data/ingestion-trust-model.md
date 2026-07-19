# Ingestion Trust Model

To protect Tourtect community safety and ensure high reliability, the ingestion pipeline enforces six clear trust levels. 

## Trust Levels Matrix

```
[Level 0: Discovered] ──> [Level 1: Fetched] ──> [Level 2: AI Extracted] ──> [Level 3: Machine Validated] ──> [Level 4: Verified] ──> [Level 5: Published]
```

| Trust Level | Pipeline Stage | Description |
| --- | --- | --- |
| **Level 0** | Discovered | URL candidate identified by crawler. Not trusted. |
| **Level 1** | Fetched | Raw text/meta extracted and preserved locally in a staging sandbox. |
| **Level 2** | AI Extracted | Factual entities parsed into generic extraction structures by an LLM. |
| **Level 3** | Machine Validated | Deterministic checks pass (impossible coordinates, outlier price filter, robots check). |
| **Level 4** | Verified | Curated and signed off by a human editor or a trusted partner authority. |
| **Level 5** | Published | Made available to the Tourtect mobile client and organic feed query index. |

---

## Architectural Rule: Level 2 Strict Quarantine Gate

No raw model output is ever trusted or exposed directly to tourists. 

A Level 2 AI Extraction is **automatically quarantined** and prevented from reaching Level 5 when it contains:
1. **Accusations**: Claims of fraud, cheating, or theft against specific named businesses or individuals (e.g. "Restaurant X is a scam"). Scam patterns must focus purely on objective behaviors.
2. **Safety-Critical Phone Numbers**: Tourist emergency numbers or police hotlines. Hallucinations could lead to tourists calling incorrect numbers during a crisis.
3. **Legal or Immigration Rules**: Visa requirements or border procedures.
4. **Emergency Instructions**: Evacuation routes or severe weather instructions.
5. **PII**: Phone numbers, emails, addresses of private individuals.

---

## Database Schema Mapping Flow

To preserve provenance and separate staging data from production truth, records flow through distinct database boundaries:

```
                  Ingestion Crawler / Sandbox Fetch
                               ↓
                       [source_documents] (L1)
                               ↓
                      [extraction_runs]
                               ↓
                     [extracted_records] (L2)
                               ↓
                    [normalized_records] (L3)
                               ↓
     ┌─────────────────────────┴─────────────────────────┐
     ▼ (If High-Risk/Needs Review)                       ▼ (If Safe & Low-Risk)
[review_decisions] (L4)                             [machine_validation_pass]
     └─────────────────────────┬─────────────────────────┘
                               ↓
                     [published_records] (L5)
                               ↓
                    Tourtect Core Domain Tables (places, price_observations, safety_directory)
```

1. **Staging Tables**: `source_documents`, `extracted_records`, `normalized_records`, and `validation_results` are used for pipeline tracing. They contain fields for confidence scores, prompt details, and raw JSONB payloads.
2. **Review tables**: `review_decisions` holds editorial audit trails (who, when, what was edited).
3. **Core Domain Tables**: `places`, `price_observations`, and `safety_directory_entries` hold the final validated truth. 
