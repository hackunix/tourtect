# Rủi ro, phương án và giả định đã khóa

> Tách từ `system-design.md` — mục 13–14.

## 13. Rủi ro và phương án

| Rủi ro | Ảnh hưởng | Phương án |
| --- | --- | --- |
| Forum ít nội dung/cold start | Cao | Seed place/topic có kiểm chứng, local contributor program, Q&A theo vùng và external card đúng quyền |
| Review giả/brigading/Sybil | Cao | Reputation đa chiều, graph/rate anomaly, evidence badge, quarantine, human review và rollback |
| Cáo buộc gây hại/merchant retaliation | Cao | Review hành vi/giao dịch, PII redaction, right-to-reply, notice/appeal, legal escalation và audit |
| Fetch nội dung vi phạm quyền hoặc bị xóa | Cao | Official API/RSS/partner, rights registry, metadata/embed tối thiểu, refresh/takedown sync |
| Platform khóa bot/API hoặc đổi điều khoản | Cao | Deny-by-default source policy, policy expiry, kill switch, partner route và không phụ thuộc một platform |
| Crawler gây tải/bị block | Cao | Stable identity/IP, per-host budget, conditional GET, jitter/backoff, circuit breaker và abuse contact |
| Giá commerce bị hiểu là giá tại chỗ | Cao | CommerceOfferObservation riêng; tách promo/fee/location và cần evidence độc lập trước Price Engine |
| Địa giới/địa chỉ cũ sau sắp xếp | Cao | Official versioned snapshot, temporal region graph, legacy alias/redirect và spatial QA |
| OSM POI sai/cũ hoặc license violation | Trung bình | Provenance/version, merchant/official verification, diff QA, attribution và ODbL review |
| Social signal thiên lệch theo ngôn ngữ/platform | Cao | Query pack bản ngữ, source diversity cap, dedupe cross-post, calibration và human review |
| Quảng cáo làm mất niềm tin | Cao | Commercial trust firewall, disclosure, ranking feature allowlist, no-ads safety surface và trust-health gate |
| Một nguồn doanh thu chi phối | Trung bình | Theo dõi revenue concentration, đa dạng Plus/business/B2B/grant và governance độc lập |
| Model FPT AI Factory dịch sai ngữ cảnh khẩn cấp | Cao | Glossary, critical validator, phrase pack, benchmark theo model ID/version và rule đứng trước LLM |
| FPT AI Factory lỗi/rate-limit/latency cao | Cao | Timeout, retry có backoff, circuit breaker, budget/rate cap và degraded/manual mode |
| Profile demo server-only làm Live/Lens phụ thuộc mạng | Cao | PTT/frame nhỏ, backpressure, timeout/circuit breaker, manual input, caption/phrase và offline safety pack |
| STT không đạt realtime hoặc không hỗ trợ đủ locale | Trung bình | Benchmark model Marketplace đã chọn, giới hạn utterance, caption/manual và phrase pack |
| VLM đọc sai số | Cao | Không dùng VLM làm source of truth; high-res OCR + user confirmation |
| Dữ liệu giá mỏng | Cao | Abstain, mở rộng cohort có disclosure, khảo sát tập trung |
| Crowdsourcing bị đầu độc | Cao | Quarantine, source cap, dedup, review, rollback |
| Live làm nóng máy/tốn data | Trung bình | PTT, adaptive FPS, burst video, thermal/network policy |
| Android fragmentation/lifecycle làm rò camera/micro | Cao | CameraX lifecycle binding, foreground-only capture, permission revoke test và device matrix |
| TTS không có đủ locale | Trung bình | Caption và audio phrase đã đóng gói |
| Người dùng hiểu “Live” là gọi người thật | Trung bình | Tên “Live Voice/Live Camera với AI”, onboarding rõ |
| Zalo Mini App runtime khác Android native | Trung bình | Lite scope, test Android thật, không cam kết realtime |
| Hotline thay đổi | Cao | Registry có nguồn/ngày xác minh, signed pack, expiry |
| Cảnh báo làm xung đột leo thang | Cao | Private card, haptic, safe wording, no public TTS |

---

## 14. Giả định đã khóa

- Tên sản phẩm là **Tourtect**.
- Forum/community knowledge graph là nền tảng; AI phiên dịch, Lens, price/scam intelligence là mô-đun hỗ trợ.
- Responsive web và Android app đều là sản phẩm đầy đủ cho forum; Android có thêm Live/Lens/offline/SOS. iOS tạm ngoài phạm vi demo.
- Zalo Mini App là kênh lite cho discovery, quick report, snapshot và hotline.
- Đọc public content không cần tài khoản; đăng/tương tác cần account. SOS và AI session riêng tư không bắt buộc account.
- V1 hỗ trợ account email/password và Google Sign-In; Google chỉ dùng authentication scope <code>openid email profile</code>, không mặc định xin quyền Google API khác.
- External content chỉ qua official API, RSS/Atom, partner feed hoặc shared URL trong phạm vi được phép; không sao chép toàn văn/re-host video mặc định.
- Việt Nam là geo scope và nguồn xác minh chính; Trung/Hàn/Nga trước mắt là market/locale discovery về trải nghiệm tại Việt Nam, chưa phải mở rộng POI toàn cầu.
- Không crawl GrabFood/ShopeeFood/Shopee hoặc public social bằng browser automation nếu chưa có chấp thuận; commerce/social ưu tiên partner API, creator opt-in và user-submitted URL.
- Google Custom Search JSON API không phải dependency; search API nếu bổ sung chỉ làm discovery và không thay rights gate.
- Địa giới dùng danh mục/mã chính thức Việt Nam có version; OSM là base map/POI/cross-check có attribution, không là nguồn pháp lý.
- Organic ranking, review, moderation, Price Engine và Safety Engine độc lập hoàn toàn với spend/commission/business tier.
- Không paywall SOS, hotline, báo cáo an toàn và cảnh báo thiết yếu; không có pay-to-remove.
- “Call” nghĩa là phiên nói với AI trong app, không phải cuộc gọi PSTN.
- Live Voice chỉ dùng push-to-talk theo vai.
- Live Camera là AI camera assist, không kết nối video tới người khác.
- Model server-side của profile demo gọi FPT AI Factory bằng API tương thích OpenAI; Gemini Live chỉ là mẫu tham khảo.
- Kiến trúc model là adaptive device/server. Riêng env demo đặt <code>AI_EXECUTION_MODE=server_only</code>; khi API không khả dụng Android dùng degraded/manual/offline safety pack.
- Dependency local chạy bằng Podman Compose; không yêu cầu Docker daemon hoặc GPU local.
- Android system TTS là mặc định của thiết kế; demo có thể ép server TTS theo execution policy và offline luôn có phrase audio đã duyệt.
- Price insight luôn riêng tư bằng card/rung.
- Pilot có thể chạy bằng dữ liệu seed/cộng đồng mà không phụ thuộc đối tác thương mại; connector ngoài phải có quyền truy cập hợp lệ.
- Contribution chỉ có hiệu lực khi user opt-in riêng.
- Chinese MVP là Simplified Chinese; Traditional Chinese nằm ngoài pilot.
- Không có human operator, VoIP bridge, auto-call hoặc auto-share location trong V1.

---
