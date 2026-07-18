# Responsive Web và Zalo Mini App

> Tách từ `system-design.md` — mục 3.8–3.9.

### 3.8 Responsive web

Website là bề mặt chính cho discovery, SEO, thảo luận và doanh thu:

- Đọc feed, place page, review, bảng giá, alert và nội dung ngoài mà không cần tài khoản.
- Tìm kiếm theo thành phố, địa điểm, topic, ngôn ngữ, thời gian, loại bằng chứng và mức giá.
- Tài khoản bắt buộc khi đăng, bình luận, vote, theo dõi, báo nội dung hoặc lưu đồng bộ.
- Cho phép upload một ảnh để price check và chuyển kết quả thành draft post sau xác nhận riêng.
- Có trang public indexable cho place/topic; trang cá nhân, incident và dữ liệu nhạy cảm không index.
- Hiện QR/deep link sang Android cho Tourtect Live/Lens; web không chạy continuous camera/audio trong V1.

### 3.9 Zalo Mini App

Đối tượng ưu tiên là người Việt đồng hành, khách sạn, hướng dẫn viên; khách quốc tế đã có Zalo vẫn dùng được.

V1:

- Feed lite theo thành phố, search place và xem review/price/scam card.
- Tạo price report hoặc scam report ngắn; bước xuất bản phải xác nhận nội dung và phạm vi công khai.
- Camera preview và “Bấm để phân tích”; không gửi stream liên tục.
- Chụp/chọn một ảnh menu/hóa đơn/bảng giá.
- Text scam assistant.
- Hotline bằng <code>openPhone</code>.
- Incident card và phrasebook.
- Chọn thành phố thủ công; GPS là tùy chọn đúng ngữ cảnh.
- Dùng Zalo ID giả danh cho session, không xin profile hoặc số điện thoại khi onboarding.

Ràng buộc:

- Camera/micro dừng khi Mini App xuống background.
- <code>openPhone</code> chỉ mở màn hình gọi native.
- Không coi WebSocket/WebRTC trong WebView là hợp đồng nền tảng nếu tài liệu Zalo chưa cam kết.
- Không dùng <code>localStorage</code>, <code>sessionStorage</code> hoặc cookie; token gửi qua Authorization header và cache qua Zalo storage API.
