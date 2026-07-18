# Phạm vi sản phẩm

> Tách từ `system-design.md` — mục 2.

## 2. Phạm vi sản phẩm

### 2.1 Người dùng chính

| Persona | Nhu cầu |
| --- | --- |
| Khách Hàn/Trung/Anh/Nga | Tìm review đúng ngôn ngữ, so sánh giá, hỏi cộng đồng, nhận cảnh báo và hỗ trợ tại chỗ |
| Người Việt/expat địa phương | Chia sẻ kinh nghiệm, xác minh giá, trả lời câu hỏi và cảnh báo rủi ro |
| Chủ cơ sở kinh doanh | Claim hồ sơ, cập nhật thông tin/giá, phản hồi review và kháng nghị công bằng |
| Nhân viên khách sạn/hướng dẫn viên | Tra playbook, hỗ trợ tạo incident card, xác minh thông tin địa phương |
| Moderator/Data reviewer | Duyệt UGC, bằng chứng giá, scam pattern, hotline, appeal và snapshot |
| Editor/Content steward | Quản trị nguồn báo/video, quyền sử dụng, phiên bản, takedown và audit |

### 2.2 Phân chia theo kênh

| Kênh | Có | Không có |
| --- | --- | --- |
| Responsive web | Full forum/feed/search, place page, post/review, so sánh giá, scam map, nội dung báo/video, tài khoản và monetization | Live Voice/continuous camera |
| Android app | Toàn bộ forum, notification/nearby alert, offline pack, Tourtect Live, Tourtect Lens, price check và SOS | Phiên dịch bên trong cuộc gọi PSTN; bản iOS |
| Zalo Mini App | Feed lite, tìm place, chụp ảnh/price report nhanh, text scam report, incident card và hotline qua openPhone | Full-duplex audio/video, WebRTC, AI nghe cuộc gọi |
| Admin/Moderation web | UGC queue, appeal, source/rights registry, dataset/scam/hotline versioning, ads eligibility, metrics và rollback | Xem raw media khi không có quyền và mục đích hợp lệ |

> Quyết định demo: forum là nền tảng chính; Android là mobile client duy nhất và iOS tạm ngoài phạm vi. Tourtect Live/Lens là mô-đun hỗ trợ; kiến trúc model vẫn adaptive, còn deployment demo chọn profile server-only bằng env.

### 2.3 Trong phạm vi pilot

- Nơi thử nghiệm: Một vài quận/huyện ở Hà Nội
- Forum đa ngôn ngữ, place page, review, price report, scam report, Q&A, feed và search.
- Connector pilot cho nguồn chính thức Việt Nam, RSS/Atom/báo allowlist, YouTube Data API/embed, OSM Vietnam extract/diff và URL do người dùng gửi.
- Query pack Việt/Anh/Hàn/Trung giản thể/Nga để phát hiện nội dung liên quan đến scam, gián đoạn và giá tại Việt Nam; chỉ publish sau rights gate và moderation.
- Năm nhóm giá: taxi/đưa đón, đổi tiền, ăn uống, tour và bán lẻ thiết yếu/đường phố có item chuẩn hóa (ví dụ nước đóng chai, SIM, áo mưa, souvenir phổ biến).
- UI và dịch hai chiều giữa tiếng Việt với:
  - <code>ko-KR</code>
  - <code>zh-Hans</code>
  - <code>en</code>
  - <code>ru-RU</code>
- Scam pattern seed: taxi, đổi tiền, ghost tour, ép thanh toán và tình huống có nguy cơ leo thang.
- Hotline khẩn cấp toàn quốc và hotline du lịch địa phương đã được xác minh.
- Android app native; iOS không build, test hoặc phát hành trong demo.

### 2.4 Ngoài phạm vi

- Cầu nối PSTN/SIP/WebRTC ba bên.
- Nhân viên phiên dịch trực 24/7.
- Full AI offline trên mọi thiết bị.
- Live Voice hoặc Live Camera trên web/Zalo Mini App.
- Tự gọi, tự gửi vị trí, tự gửi incident card hoặc tự liên hệ trusted contact.
- Chấm điểm hoặc công khai danh sách người bán bị tố cáo.
- Dùng dữ liệu hóa đơn thuế hoặc dữ liệu riêng của đối tác khi chưa có quyền hợp pháp.
- Sao chép toàn văn bài báo, tải lại video hoặc vượt paywall/login/robots/access control.
- Browser-scrape GrabFood, ShopeeFood, Shopee, Reddit, TikTok, Facebook, Instagram, Naver Cafe, Telegram hoặc nền tảng tương tự khi chưa có API/partner permission phù hợp.
- iOS app, App Store release, iOS-specific authentication/media/background behavior và test device iPhone/iPad.
- Cho nhà quảng cáo mua vị trí trong organic ranking, mua điểm review hoặc trả tiền để xóa cảnh báo hợp lệ.
- Mở marketplace/booking checkout riêng trong pilot; V1 chỉ dùng referral/affiliate rõ disclosure.

---
