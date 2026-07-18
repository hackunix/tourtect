# API conventions và endpoint catalog

> Tách từ `system-design.md` — mục 7.1–7.2.

### 7.1 Quy ước

- Prefix: <code>/v1</code>.
- JSON dùng <code>snake_case</code>.
- Thời gian ISO 8601 UTC.
- ID là UUID/opaque ID.
- Lỗi HTTP theo <code>application/problem+json</code>.
- Mọi response AI có <code>trace_id</code>, model/dataset version và confidence nếu áp dụng.
- Đọc public content và dùng SOS/AI session riêng tư không yêu cầu account; dùng anonymous/short-lived session token.
- Đăng bài, bình luận, vote, follow, report, subscription và merchant action yêu cầu account; endpoint có idempotency key và rate limit phù hợp.
- Admin dùng SSO/RBAC và audit log.
- Raw media upload dùng signed URL, content type allowlist và size limit.

### 7.2 Endpoint

| Method | Path | Mục đích |
| --- | --- | --- |
| POST | <code>/v1/auth/registrations</code> | Tạo account email/password và gửi verification |
| POST | <code>/v1/auth/email-verifications</code> | Xác minh email bằng token một lần |
| POST | <code>/v1/auth/sessions</code> | Đăng nhập email/password và tạo Tourtect session |
| POST | <code>/v1/auth/sessions/refresh</code> | Rotate refresh token, phát hiện reuse |
| DELETE | <code>/v1/auth/sessions/{session_id}</code> | Revoke một phiên/thiết bị |
| DELETE | <code>/v1/auth/sessions</code> | Đăng xuất tất cả phiên |
| POST | <code>/v1/auth/password-resets</code> | Yêu cầu email reset, response không làm lộ account tồn tại |
| POST | <code>/v1/auth/password-resets/confirm</code> | Đổi mật khẩu bằng token một lần và revoke session khác |
| POST | <code>/v1/auth/oauth/google/attempts</code> | Tạo state/nonce/PKCE auth attempt và authorization URL |
| GET/POST | <code>/v1/auth/oauth/google/callback</code> | Verify callback/code và tạo/link Tourtect session |
| POST | <code>/v1/auth/identities/google/link</code> | Link Google với account đang đăng nhập sau re-authentication |
| DELETE | <code>/v1/auth/identities/{identity_id}</code> | Unlink provider nếu account còn phương thức đăng nhập khác |
| POST | <code>/v1/auth/anonymous-merge</code> | Preview/confirm merge saved item, draft và preference |
| GET | <code>/v1/account/sessions</code> | Liệt kê phiên gần đây và metadata thiết bị tối thiểu |
| GET | <code>/v1/feed</code> | Feed Following/Nearby/Latest/Trending/Safety với reason code |
| GET | <code>/v1/search</code> | Full-text/geo search place, post, giá, scam và external content |
| POST/GET | <code>/v1/posts</code> | Tạo draft/xuất bản và đọc danh sách post |
| GET/PATCH/DELETE | <code>/v1/posts/{id}</code> | Đọc, sửa có version hoặc yêu cầu xóa post |
| POST/GET | <code>/v1/posts/{id}/comments</code> | Bình luận và thread |
| POST | <code>/v1/posts/{id}/votes</code> | Đánh dấu hữu ích; không đồng nghĩa evidence |
| POST | <code>/v1/posts/{id}/reports</code> | Báo vi phạm/an toàn/PII |
| POST | <code>/v1/follows</code> | Theo dõi place/topic/user |
| GET | <code>/v1/places/{id}</code> | Place page aggregate |
| POST/GET | <code>/v1/places/{id}/reviews</code> | Review có cấu trúc và danh sách review |
| POST | <code>/v1/places/{id}/claim</code> | Merchant claim; không đổi review/ranking |
| POST/GET | <code>/v1/price-reports</code> | Price report công khai và evidence workflow |
| POST/GET | <code>/v1/scam-reports</code> | Scam report, safety triage và moderation |
| GET | <code>/v1/external-content</code> | External card đã qua rights/policy gate |
| POST | <code>/v1/external-content/submissions</code> | Gửi URL để connector kiểm tra; không fetch tùy ý từ client |
| POST/GET | <code>/v1/moderation/appeals</code> | Tạo và theo dõi appeal |
| GET/PATCH | <code>/v1/notifications</code> | Danh sách và trạng thái notification |
| POST/GET | <code>/v1/subscriptions</code> | Tourtect Plus và entitlement |
| GET/PATCH | <code>/v1/business-profiles/{place_id}</code> | Business tools sau claim/verification |
| POST | <code>/v1/affiliate-events</code> | Ghi sự kiện disclosure/click tối thiểu, chống giả mạo |
| POST | <code>/v1/realtime/sessions</code> | Tạo Live Voice/Camera session |
| POST | <code>/v1/realtime/sessions/{id}/resume</code> | Cấp token resume ngắn hạn |
| WS | <code>/v1/realtime/sessions/{id}/events</code> | PCM, context frame và realtime event |
| POST | <code>/v1/realtime/sessions/{id}/captures</code> | Tạo capture_id và signed PUT URL |
| PUT | <code>{signed_upload_url}</code> | Upload capture đã redaction từ local path hoặc capture tạm mã hóa khi dùng server path |
| POST | <code>/v1/realtime/sessions/{id}/captures/{capture_id}/finalize</code> | Kiểm tra hash/MIME/redaction metadata; server path tiếp tục OCR/redaction trước khi capture thành ready |
| DELETE | <code>/v1/realtime/sessions/{id}</code> | End và xóa state/media tạm |
| POST | <code>/v1/price-checks</code> | Tạo price check từ capture_id đã finalize hoặc manual input |
| GET | <code>/v1/price-checks/{id}</code> | Lấy trạng thái/kết quả |
| POST | <code>/v1/scam-assessments</code> | Đánh giá tình huống text/transcript |
| GET | <code>/v1/emergency-services</code> | Hotline theo vùng/incident |
| GET | <code>/v1/safety-packs/{region}</code> | Gói offline có version/signature |
| POST | <code>/v1/contributions</code> | Opt-in contribution |
| DELETE | <code>/v1/privacy/sessions/{id}</code> | Yêu cầu xóa dữ liệu phiên |
