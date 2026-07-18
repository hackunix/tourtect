# External Content và Monetization

> Tách từ `system-design.md` — mục 4.3.

### 4.3 External content và monetization

| ID | Yêu cầu |
| --- | --- |
| CM-01 | Chỉ fetch qua official API, RSS/Atom, partner feed hoặc URL chia sẻ trong phạm vi được phép |
| CM-02 | Lưu canonical URL, attribution, rights status, retrieved/checked time và deletion state |
| CM-03 | Dedupe, entity-link, fact/opinion/sponsored labeling, freshness và takedown sync |
| CM-04 | Không vượt paywall/login/access control hoặc re-host toàn bài/video khi chưa có license |
| CM-05 | Ad/sponsored/affiliate phải có disclosure; SOS và critical incident không có quảng cáo |
| CM-06 | Business trả phí không được sửa rank/review/moderation/alert; mọi quyền được enforce bằng RBAC và audit |
| CM-07 | B2B chỉ xuất aggregate đủ ngưỡng và không bán raw content, media hoặc định danh người dùng |
| CM-08 | Trust-health gate có thể ngừng monetization surface độc lập với product availability |
| CM-09 | Mỗi connector có access mode, policy/robots version, owner, crawl budget, retention, deletion SLA và kill switch |
| CM-10 | Scheduler dùng per-host token bucket, conditional GET, adaptive backoff và không vượt login/paywall/CAPTCHA/private endpoint |
| CM-11 | Commerce offer tách khỏi giá tại chỗ; promotion/delivery/personalization không được nhập thẳng Price Engine |
| CM-12 | Query pack đa ngôn ngữ có version, editor bản ngữ, place aliases, negative keyword và precision review |
| CM-13 | Search API chỉ discovery URL; mọi kết quả vẫn qua rights/policy gate và không dùng để lách platform restriction |
