BACKEND_READINESS_GATE=PASS

# Backend Readiness Report — Tourtect V1 Monolith Slice

## 1. Execution Summary

All required backend verification checks have been executed successfully on localhost.

* **Infrastructure Status**: PostgreSQL/PostGIS, Redis, MinIO are healthy and running.
* **Database Migrations**: 7 goose migrations successfully applied up to version 7.
* **Seed Data**: Fully loaded. 5 places, 8 posts (6 public, 2 drafts), 2 price snapshots, 6 observations, 1 safety directory version with approved contact hotlines.
* **Compilation**: `api`, `realtime`, and `worker` binaries compiled cleanly.
* **Unit/Integration Tests**: 22 tests (11 Price Engine, 11 Safety Engine, 2 PTT Realtime Session) passed successfully.
* **Persistence Test**: Draft posts successfully persist in PostgreSQL across API restarts and are verified via `/v1/posts` after publishing.
* **Dependency Failure Test**: Stopping Postgres causes readiness `/health/ready` to fail with 503 while liveness `/health/live` remains healthy. Recovering Postgres restores readiness to 200 OK.
* **Mock Inspection**: No production code contains client-side/runtime mocks. The runtime source of truth is strictly database-backed.

---

## 2. Environment & Database Metadata

* **API Server Port**: `8080`
* **Realtime WS Port**: `8081`
* **Worker Port**: `8082`
* **Database Name**: `tourtect`
* **Database User**: `tourtect`
* **PostgreSQL Version**: `16-3.4 (with PostGIS)`
* **Active Migration Version**: `7`

### Seed Data Row Counts
* `places` count: `5`
* `place_aliases` count: `8`
* `posts` count: `9` (including the dynamically created test post)
* `price_snapshots` count: `2`
* `price_observations` count: `7` (including 1 test observation)
* `safety_directory_versions` count: `2` (including 1 test directory)

---

## 3. Test Matrix & Results

```bash
$ go test -v ./...
?       github.com/tourtect/backend/cmd/api     [no test files]
?       github.com/tourtect/backend/generated/database  [no test files]
?       github.com/tourtect/backend/generated/openapi   [no test files]
?       github.com/tourtect/backend/internal/content    [no test files]
?       github.com/tourtect/backend/internal/places     [no test files]
?       github.com/tourtect/backend/internal/platform/config    [no test files]
?       github.com/tourtect/backend/internal/platform/database  [no test files]
?       github.com/tourtect/backend/internal/platform/httpserver        [no test files]
?       github.com/tourtect/backend/internal/platform/logging   [no test files]
=== RUN   TestPriceEngine
=== RUN   TestPriceEngine/Within_range_-_typical
=== RUN   TestPriceEngine/Slightly_high_-_elevated
=== RUN   TestPriceEngine/Significantly_high_-_high_risk
=== RUN   TestPriceEngine/Missing_exact_snapshot_-_falls_back_to_national
=== RUN   TestPriceEngine/Future_snapshot_ignored
=== RUN   TestPriceEngine/Sample_size_too_small
=== RUN   TestPriceEngine/Low_extraction_confidence_unconfirmed
=== RUN   TestPriceEngine/Low_extraction_confidence_user_confirmed
=== RUN   TestPriceEngine/Currency_mismatch
=== RUN   TestPriceEngine/Unit_mismatch
=== RUN   TestPriceEngine/No_snapshot_contamination
--- PASS: TestPriceEngine (0.04s)
PASS
ok      github.com/tourtect/backend/internal/pricing    0.050s
=== RUN   TestRealtimeSession
=== RUN   TestRealtimeSession/PTT_state_transition_flow
=== RUN   TestRealtimeSession/Rejects_invalid_sequence_numbers
--- PASS: TestRealtimeSession (0.05s)
PASS
ok      github.com/tourtect/backend/internal/realtime   0.053s
=== RUN   TestSafetyEngine
=== RUN   TestSafetyEngine/High_price_without_coercion
=== RUN   TestSafetyEngine/Forced_payment_coercion
=== RUN   TestSafetyEngine/Refuse_to_let_user_leave_(confinement)
=== RUN   TestSafetyEngine/Injury_detected
=== RUN   TestSafetyEngine/Weapon_threat
=== RUN   TestSafetyEngine/Informational_safety_question
=== RUN   TestSafetyEngine/Degrades_with_missing_facts
=== RUN   TestSafetyEngine/Conflicting_safety_facts
=== RUN   TestSafetyEngine/Evaluates_without_LLM_provider
=== RUN   TestSafetyEngine/Degrades_gracefully_when_database_safety_directory_is_empty
=== RUN   TestSafetyEngine/No_hotline_hallucination
--- PASS: TestSafetyEngine (0.06s)
PASS
ok      github.com/tourtect/backend/internal/safety     0.066s
```

---

## 4. Curl Response Details

### Liveness Probe
```http
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Id: ae455086-f27a-4ecb-bff9-c88af4249d72
Content-Length: 64

{"status":"ok","timestamp":"2026-07-19T00:17:38.2364246+07:00"}
```

### Readiness Probe
```http
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Id: 66f4fb8d-11af-42bd-b494-957db8427c7e
Content-Length: 92

{"checks":{"postgres":"UP"},"status":"ok","timestamp":"2026-07-19T00:17:38.24317758+07:00"}
```

### Places Endpoint (Sample output truncate)
```http
HTTP/1.1 200 OK
Content-Type: application/json
X-Request-Id: 3aa68a29-1c72-4180-90b3-38f2f4c7cae7

{"items":[{"address":"Sân bay Quốc tế Nội Bài, Sóc Sơn, Hà Nội","aliases":["Noi Bai Airport Taxi"],"average_rating":0,"category":"taxi","coordinates":{"latitude":21.2187,"longitude":105.8038},"created_at":"2026-07-19T00:10:24.817632+07:00","freshness":"2026-07-16T00:10:24.817632+07:00","name":"Nội Bài Taxi Stand","place_id":"019078a0-1001-7000-8000-000000000005","post_count":2,"region_id":"hanoi-soc-son"}],"pagination":{"has_more":false}}
```

---

## 5. Dependency Outage Simulation

### 1. Stopping Postgres container
`podman stop tourtect_postgres_1`

### 2. Checking /health/live (PASS)
```http
HTTP/1.1 200 OK
X-Request-Id: c9526c63-e802-4d2c-9ad1-efdcdd234c2f

{"status":"ok","timestamp":"2026-07-19T00:18:17.233703549+07:00"}
```

### 3. Checking /health/ready (FAIL)
```http
HTTP/1.1 503 Service Unavailable
X-Request-Id: ed188e43-c559-484d-a98c-698a766a7a6c

{"checks":{"postgres":"DOWN: failed to connect to `user=tourtect database=tourtect`:\n\t[::1]:5432 (localhost): dial error: dial tcp [::1]:5432: connect: connection refused..."},"status":"unavailable","timestamp":"2026-07-19T00:18:17.240347946+07:00"}
```

### 4. Restarting Postgres and confirming recovery
`podman start tourtect_postgres_1`
`curl http://localhost:8080/health/ready` -> returns `200 OK` with status `"ok"`.

---

## 6. Commit Verification

* **Commit Hash**: `473dcb7` (local working directory has modifications containing source code).
* **Gate Status**: **PASS**
