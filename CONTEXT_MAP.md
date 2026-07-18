# Context Map — chọn file theo tác vụ

| Tác vụ | File nên đọc trước | File bổ sung khi cần |
| --- | --- | --- |
| Hiểu toàn bộ sản phẩm ở mức ngắn | `00-overview/01-executive-summary.md`, `00-overview/02-product-scope.md` | `07-delivery/01-hackathon-demo.md` |
| Xây forum/feed/place page | `01-product-experience/02-forum-place-feed.md`, `02-functional-requirements/01-forum-discovery.md` | `03-architecture/03-modular-monolith-components.md`, `05-api/02-types-core-community.md` |
| Review/reputation/moderation | `02-functional-requirements/02-review-moderation.md`, `01-product-experience/02-forum-place-feed.md` | `04-data/01-community-knowledge-graph.md` |
| Android app foundation | `03-architecture/02-container-android.md`, `02-functional-requirements/08-android-client.md` | `06-operations-safety/01-offline-degraded-mode.md` |
| Live Voice/PTT | `01-product-experience/03-live-voice.md`, `02-functional-requirements/05-live-voice.md` | `03-architecture/04-model-stack-routing.md`, `05-api/03-types-realtime-vision.md`, `05-api/06-websocket-protocol.md` |
| Live Camera/Lens | `01-product-experience/04-live-camera.md`, `02-functional-requirements/06-live-camera.md` | `03-architecture/02-container-android.md`, `05-api/03-types-realtime-vision.md`, `05-api/06-websocket-protocol.md` |
| Price Check UI/API | `02-functional-requirements/04-price-check.md`, `05-api/04-types-price.md`, `05-api/07-price-response-example.md` | `04-data/04-vertical-normalization.md`, `04-data/05-alert-algorithm.md` |
| Price dataset/RAG/reference price | `04-data/02-price-layers-sources.md`, `04-data/03-price-safety-models.md` | `04-data/04-vertical-normalization.md`, `04-data/06-cold-start-poisoning.md` |
| Scam Assistant/SOS | `01-product-experience/05-scam-emergency.md`, `02-functional-requirements/07-scam-emergency.md` | `05-api/05-types-scam.md`, `06-operations-safety/03-threat-model-safety-compliance.md` |
| External crawler/connectors | `01-product-experience/07-external-content.md`, `02-functional-requirements/03-external-content-monetization.md` | `03-architecture/03-modular-monolith-components.md`, `06-operations-safety/03-threat-model-safety-compliance.md` |
| Web/Zalo client | `01-product-experience/06-web-zalo.md` | `02-functional-requirements/01-forum-discovery.md`, `02-functional-requirements/03-external-content-monetization.md` |
| Authentication/Google OAuth | `01-product-experience/09-authentication.md`, `05-api/01-conventions-endpoints.md` | `05-api/02-types-core-community.md`, `06-operations-safety/03-threat-model-safety-compliance.md` |
| Monetization/business tools | `01-product-experience/08-monetization.md`, `02-functional-requirements/03-external-content-monetization.md` | `06-operations-safety/04-metrics-release-gates.md` |
| Podman/local development | `03-architecture/05-local-runtime-podman.md` | `03-architecture/04-model-stack-routing.md` |
| API CRUD chung | `05-api/01-conventions-endpoints.md` | Type file đúng domain trong `05-api/` |
| Realtime protocol | `05-api/06-websocket-protocol.md`, `05-api/03-types-realtime-vision.md` | Live Voice hoặc Lens experience file |
| Privacy/retention | `06-operations-safety/02-consent-retention.md` | `06-operations-safety/03-threat-model-safety-compliance.md` |
| Test/release gate | `06-operations-safety/04-metrics-release-gates.md`, `06-operations-safety/05-test-matrix-golden-data.md` | Requirement file của feature đang test |
| Chuẩn bị demo/pitch | `07-delivery/01-hackathon-demo.md`, `00-overview/01-executive-summary.md` | `07-delivery/03-risks-assumptions.md` |
