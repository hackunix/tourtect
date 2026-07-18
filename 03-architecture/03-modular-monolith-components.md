# Modular Monolith và trách nhiệm thành phần

> Tách từ `system-design.md` — mục 5.3–5.4.

### 5.3 Lý do chọn modular monolith

Pilot cần phát triển nhanh và giữ transaction/provenance dễ kiểm soát. Core API là modular monolith với boundary rõ:

- Realtime được tách process vì tải WebSocket khác HTTP CRUD. Profile demo gọi FPT AI Factory và không vận hành GPU container local; kiến trúc vẫn cho phép provider khác hoặc on-device sau demo.
- Search indexing, external ingestion, notification và data worker chạy bất đồng bộ vì cần retry/quarantine.
- Các module dùng contract nội bộ, có thể tách service sau khi có số liệu tải thật.
- Không dùng microservice hoặc self-host container cho từng model trong giai đoạn hackathon; đây là giới hạn deployment demo, không phải khóa thiết kế dài hạn.

### 5.4 Trách nhiệm thành phần

| Thành phần | Trách nhiệm | Không được làm |
| --- | --- | --- |
| Identity/Profile | Registration, email verification, Google OIDC, provider linking, session rotation/revocation và public profile | Tin client token chưa verify, dùng email làm Google identity key hoặc tự cấp verified role |
| Content/Review Service | Post, comment, review, draft, merchant reply và version history | Tự quyết sanction hoặc điểm uy tín |
| Place/Topic Service | Canonical place, alias đa ngôn ngữ, geo, topic và price report | Cho merchant ghi đè community fact |
| Feed/Organic Ranking | Candidate, relevance, freshness, diversity và safety boost | Nhận spend/commission/business tier làm feature |
| Search | Full-text, geo, filter, autocomplete và index version | Index private incident hoặc PII bị redaction |
| Reputation | Điểm theo lĩnh vực, decay, evidence và abuse signal | Dùng upvote đơn thuần làm trust score |
| Moderation/Appeal | Policy engine, queue, sanction, notice, appeal và transparency log | Cho sales/advertiser sửa quyết định |
| External Ingestion | Connector, dedupe, entity linking, refresh và takedown | Bypass access control hoặc re-host trái quyền |
| Rights/Source Registry | Access mode, robots/terms/license, allowed field, retention, owner và kill switch | Mặc định cho phép nguồn chưa review |
| Adaptive Crawl Scheduler | Per-host budget, priority/freshness, jitter, conditional refresh, backoff và circuit breaker | Né block bằng proxy/IP rotation hoặc vượt quota |
| Sandbox Fetch Proxy | Egress allowlist, DNS/IP/MIME/size guard, stable bot identity và HTTP cache | Login, chạy script trình duyệt hoặc truy cập private endpoint |
| Monetization | Ad eligibility, subscription, business entitlement và disclosure | Ghi vào rank/review/alert/moderation |
| Realtime Gateway | Session, PCM/frame transport, event ordering, backpressure | Ra quyết định giá |
| Session Orchestrator | State machine, locale, role, mode, resume | Lưu raw media dài hạn |
| Adaptive Model Router | Chọn device/server/provider fallback theo execution policy, capability, consent và confidence | Thay đổi consent hoặc làm lộ provider credential |
| Translation Service | ASR, MT, glossary, critical token validation | Tư vấn an toàn tự do |
| Vision Service | Context classification, ROI, OCR fallback | Kết luận gian lận |
| Price Engine | Chuẩn hóa, chọn cohort, anomaly/confidence | Dùng LLM làm phép so sánh cuối |
| Scam/Safety Engine | Triage, pattern retrieval, safe action, escalation | Tạo hotline |
| Emergency Directory | Trả số đã xác minh theo vùng/incident | Tự gọi số |
| Consent Service | Consent ledger, TTL, deletion/export | Gộp consent contribution với xử lý phiên |
| Data Platform | Evidence, review, snapshot, rollback | Cho submission tác động ngay production |
