# Test Matrix và Golden Data

> Tách từ `system-design.md` — mục 10.3–10.4.

### 10.3 Test matrix

#### Forum, review và moderation

- Đăng ký, email verification hết hạn/dùng lại, login đúng/sai và password reset không làm lộ account.
- Google callback sai <code>state</code>/<code>nonce</code>/<code>aud</code>/<code>iss</code>, token hết hạn, PKCE sai, code replay và redirect URI ngoài allowlist.
- Google email trùng account password, link/unlink sau re-authentication, Google subject đổi email và account còn/không còn login method.
- Refresh-token rotation/reuse, revoke một thiết bị, logout all, session hết hạn và Android Keystore/cookie bị xóa.
- Anonymous merge preview/confirm/rollback; saved item/draft được merge nhưng incident/raw media không bị chuyển.
- Account mới spam hàng loạt, copy post, review exchange và vote ring.
- Brigading theo nhóm/ngôn ngữ, Sybil nhiều account và merchant tự review.
- Review tiêu cực hợp lệ bị report hàng loạt; merchant phản hồi/appeal nhưng không được ép xóa.
- Nội dung có PII, biển số, khuôn mặt, cáo buộc chưa có bằng chứng, harassment và doxxing.
- Sửa/xóa post, merge place, block user, notification race và search index eventual consistency.
- Dịch/tóm tắt sai số tiền, phủ định hoặc biến ý kiến thành sự thật.

#### External content

- Canonical URL có tracking, redirect loop, duplicate/cross-post và cùng sự kiện khác nguồn.
- RSS/API timeout, quota/rate limit, source sửa/xóa, embed bị disable và takedown.
- ETag/Last-Modified đúng/sai/thiếu; 304 không chạy extraction; jitter tránh burst; 429 có/không có Retry-After và 403/451 đóng circuit.
- Robots/ToS thay đổi giữa hai lần chạy, kill switch, abuse complaint và source owner opt-out.
- Commerce offer có voucher, membership, delivery/location fee và giá cá nhân hóa không bị nhập thành giá tại chỗ.
- Query tiếng Việt/Anh/Hàn/Trung/Nga có alias địa danh, phủ định, từ đồng âm và spam SEO; đo precision theo query pack.
- URL trỏ private IP, file quá lớn, MIME giả, prompt injection và HTML/script độc hại.
- Article/video sponsored nhưng metadata thiếu disclosure; classifier không được tự khẳng định tuyệt đối.
- Entity linking nhầm địa điểm trùng tên hoặc gán sai thời gian/khu vực.
- OSM diff mất sequence/replay/trùng; POI merge sai chi nhánh; official boundary đổi nhưng legacy address vẫn resolve.

#### Monetization và trust firewall

- Cùng một query/feed seed phải giữ thứ tự organic khi thay advertiser spend, commission hoặc business tier.
- Ads không xuất hiện ở SOS, critical incident, silent mode hoặc nội dung chưa đủ ad eligibility.
- Affiliate/sponsored card luôn có disclosure trước CTA và screen reader đọc được.
- Merchant hết hạn gói vẫn giữ public reply/history nhưng mất feature trả phí, không làm đổi rating.
- Sales/ad token không gọi được moderation, ranking feature, PriceSnapshot hoặc ScamPattern write API.
- Trust-health gate tắt được một ad surface mà forum, SOS và alert vẫn hoạt động.

#### Live Voice

- Khách Hàn nói số tiền, người Việt nói phụ phí hợp lệ.
- Khách Trung code-switch tên món tiếng Việt.
- Khách Nga nói trong tiếng ồn đường phố.
- Khách Anh hỏi lại và sửa transcript.
- Hai người bấm nhầm vai.
- TTS locale thiếu hoặc audio route đổi sang Bluetooth.
- STT local/server confidence thấp, adaptive fallback, timeout và rate-limit.
- System/server TTS lỗi hoặc thiếu locale; app chuyển caption + phrase audio đã duyệt.
- Mất mạng đúng lúc nhả PTT.

#### Live Camera

- Menu nhiều cột, giá thiếu ký hiệu tiền tệ.
- Hóa đơn mờ, phản sáng, nghiêng.
- Đồng hồ taxi có nhiều số.
- Món ăn nhìn giống nhau nhưng khẩu phần khác.
- Bảng đổi tiền mua/bán hai cột.
- Camera thấy thẻ ngân hàng/khuôn mặt.
- Network giảm từ 4G xuống 2G.
- App background, khóa màn hình hoặc thiết bị quá nhiệt.

#### Android app

- Cold start, process recreation, rotation/configuration change và navigation restore không khởi động lại capture ngoài ý muốn.
- Camera/micro permission grant/deny/revoke, app background/foreground, screen lock và session end dừng resource đúng lifecycle.
- Mạng đổi Wi-Fi/4G/2G, WebSocket reconnect/resume, event duplicate/out-of-order và backpressure ưu tiên audio/control hơn frame.
- Android Keystore entry mất/invalidated, refresh-token reuse, Room migration, cache TTL cleanup và WorkManager retry.
- Bluetooth/wired/speaker audio route, audio focus interruption và system/server TTS lỗi.
- Static scan APK/AAB không thấy <code>FPT_AI_API_KEY</code>; flavor demo không có model file/runtime, adaptive flavor chỉ chấp nhận signed model pack/runtime allowlist.

#### Price

- Giá hợp lệ cao vì sân bay/phí cầu đường/service charge.
- Cùng một món nhưng khác phường/pricing zone, budget/premium và street stall/attraction concession; engine không trộn cohort sai phân khúc.
- Pricing zone thiếu mẫu phải fallback sang xã/phường rồi tỉnh với disclosure; không đủ independent source thì abstain.
- Observation dùng địa chỉ quận/huyện cũ phải map qua đúng admin snapshot tại thời điểm giao dịch.
- Bác bỏ/loại trừ hoàn toàn các observation có thời gian trước 01/07/2025 (trước thời điểm sáp nhập hành chính).
- Nước đóng chai/SIM/áo mưa có SKU chuẩn hóa được so sánh; đồ thủ công hoặc “đồ linh tinh” không đủ thuộc tính phải trả <code>insufficient_data</code>.
- Donation/solicitation không tạo fair-price verdict; giao hàng không được yêu cầu rồi đòi tiền chuyển sang Scam/Safety.
- OCR nhầm <code>15</code> thành <code>75</code>.
- VND và USD bị nhầm.
- Dữ liệu vùng thiếu, phải mở rộng cohort.
- Snapshot stale.
- Một merchant/source gửi nhiều submission.
- Tour private bị so nhầm tour group.

#### Scam/Emergency

- Khách bị giữ trong taxi.
- Tranh chấp giá nhưng khách đã an toàn.
- Quán bình dân/hàng rong báo và giữ đúng giá: không tăng scam score chỉ vì venue type.
- Người bán đổi giá sau thỏa thuận, nhập nhằng “mỗi phần/mỗi người/mỗi 100 g” hoặc tự nhận một khoản phí là bắt buộc/chính thức.
- Đưa vòng tay/chụp ảnh/dịch vụ không được yêu cầu rồi ép trả tiền; hệ thống hỏi xác nhận trước khi match pattern.
- Người xin hỗ trợ không gây áp lực: chỉ hiển thị guidance nếu user hỏi, không gắn nhãn scam. Nếu chặn đường, đe dọa hoặc dàn dựng va chạm để đòi tiền thì route theo behavior tương ứng.
- Dàn dựng hư hại/va chạm hoặc cáo buộc để đòi bồi thường; ưu tiên rời nơi an toàn và không hướng dẫn đối đầu/quay đuổi.
- Ghost tour không có đe dọa.
- Người bị thương, không có vũ khí.
- Chế độ im lặng.
- GPS/micro bị từ chối.
- Roaming không gọi được số ngắn.
- Hotline địa phương hết giờ.

#### Zalo Mini App

- Camera permission được duyệt/từ chối trên Zalo Android thật.
- Mini App xuống background khi mở dialer.
- Storage đầy hoặc file tạm bị xóa.
- Không có cookie/localStorage.
- API CORS và Authorization header.
- User không cung cấp profile/số điện thoại.

### 10.4 Golden data

Ground truth không lấy trực tiếp từ nút thích/không thích của cộng đồng.

- Giá: khảo sát mù, mystery shopper, nguồn có thẩm quyền và chuyên gia phân xử.
- Scam: pattern do chuyên gia địa phương và ít nhất hai reviewer duyệt.
- Translation: bilingual reviewer; tập trung dữ kiện khẩn cấp và tiền.
- OCR: ảnh thực tế đã được phép sử dụng, đa ánh sáng/thiết bị.
- Mọi golden set có data card, consent/provenance và version.

---
