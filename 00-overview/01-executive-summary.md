# Tóm tắt điều hành

> Tách từ `system-design.md` — mục 1.

## 1. Tóm tắt điều hành

### 1.1 Khoảng trống cần giải quyết

Khách quốc tế thường ở thế bất lợi thông tin khi gặp:

- Giá taxi không hợp lý, đồng hồ hoặc lộ trình có dấu hiệu bất thường.
- Menu, bảng giá, hóa đơn hoặc phụ phí khó hiểu.
- Đổi tiền với tỷ giá thực nhận thấp hoặc phí không được nói trước.
- Ghost tour, ép thanh toán và các mô-típ lừa đảo phổ biến.
- Sự cố cần trao đổi ngay với người Việt, khách sạn, cơ quan du lịch hoặc lực lượng khẩn cấp.

Thông tin giúp khách tránh rủi ro đang phân tán giữa forum, review, báo chí, video mạng xã hội và kinh nghiệm địa phương. Nội dung thường khác ngôn ngữ, nhanh lỗi thời, thiếu bằng chứng và khó tìm đúng lúc. Tourtect giải quyết khoảng trống này bằng một **community knowledge graph theo địa điểm**; AI hỗ trợ chuẩn hóa thông tin nhưng không thay cộng đồng hoặc cơ chế kiểm chứng.

Hệ thống tách ba lớp trách nhiệm:

1. Cộng đồng tạo post, review, price report, scam report và thảo luận.
2. Nền tảng tổng hợp nội dung công khai đúng quyền sử dụng, liên kết nội dung với địa điểm/chủ đề và xếp hạng theo độ hữu ích, độ mới, mức bằng chứng.
3. AI dịch, OCR, tóm tắt và trích dữ kiện; Price Engine, Safety Engine và moderation policy mới quyết định cảnh báo hoặc hành động nền tảng.

### 1.2 Đề xuất giá trị

**Tourtect là forum an toàn du lịch**, nơi khách có thể tìm một địa điểm và thấy trong cùng một màn hình:

- Review có cấu trúc về minh bạch giá, chất lượng dịch vụ và cảm giác an toàn.
- Khoảng giá tham khảo, lịch sử price report và bằng chứng hóa đơn đã khử thông tin cá nhân.
- Scam report, cảnh báo chính thức, bài báo, video và thảo luận liên quan.
- Bản dịch đa ngôn ngữ, tóm tắt “cần biết gì trước khi đến” và câu hỏi cho cộng đồng.

Android app bổ sung hai công cụ AI tại chỗ:

- **Live Voice:** hai bên dùng push-to-talk; ứng dụng dịch hai chiều và âm thầm nhận biết món/dịch vụ, số tiền, tiền tệ hoặc dấu hiệu nguy hiểm.
- **Live Camera:** ứng dụng nhìn menu, món ăn, đồ vật, bảng giá hoặc đồng hồ taxi; kết hợp VLM để hiểu bối cảnh với ảnh tĩnh độ phân giải cao để xác nhận OCR và giá.

Kết quả từ Live/Lens có thể được người dùng chủ động chuyển thành draft price report hoặc scam report sau khi đã an toàn. Không tự đăng transcript, ảnh hoặc vị trí. Khi có nguy hiểm, Safety Engine ưu tiên hướng dẫn rời khỏi tình huống và mở dialer thật; AI không tự gọi hoặc tự gửi vị trí.

### 1.3 Liên hệ tiêu chí chấm

| Tiêu chí | Bằng chứng trong thiết kế |
| --- | --- |
| Sức mạnh cộng đồng | Place page, post/review có cấu trúc, reputation, moderation, appeal và feed đa ngôn ngữ |
| Độ chính xác cảnh báo giá | Bốn mức kết quả, ngưỡng đỏ nghiêm ngặt, xác nhận OCR, so sánh theo ngành và cho phép “không đủ dữ liệu” |
| Tỷ lệ báo động giả | Red alert precision mục tiêu ≥ 95%, false-positive rate ≤ 2%, calibration theo vùng/ngành |
| Dịch và xử lý khẩn cấp | PTT theo vai, glossary dữ kiện quan trọng, rule engine đứng trước LLM, safety pack offline |
| Cập nhật dữ liệu vùng | Price report cộng đồng, external connector, evidence ledger, quarantine, human review và snapshot có phiên bản |
| Quyền riêng tư | Đọc ẩn danh, consent theo mục đích, media mặc định chỉ tồn tại trong bộ nhớ, đăng bài là hành động xác nhận riêng |
| Kinh doanh bền vững | Quảng cáo, affiliate, Plus, business tools và B2B insights nằm sau commercial trust firewall |

### 1.4 Nguyên tắc sản phẩm

1. **Community first, AI assists:** giá trị cốt lõi đến từ con người và dữ liệu có provenance; AI chỉ hỗ trợ khám phá, dịch và chuẩn hóa.
2. **Safety before engagement:** cảnh báo khẩn cấp quan trọng hơn thời gian xem, lượt tương tác hoặc doanh thu.
3. **No evidence, no accusation:** đánh giá giao dịch/hành vi, không tự gắn nhãn một cá nhân hay doanh nghiệp là “lừa đảo”.
4. **Trust is not for sale:** tiền quảng cáo hoặc gói business không tác động review, moderation, alert threshold hay organic ranking.
5. **Right to reply and appeal:** cơ sở kinh doanh được phản hồi, báo sai và kháng nghị nhưng không được trả tiền để xóa nội dung hợp lệ.
6. **Consent by action:** xin camera, micro, vị trí, đăng công khai và đóng góp dataset ở các bước riêng biệt.
7. **Graceful abstention:** dữ liệu yếu phải trả “không đủ dữ liệu”.
8. **Emergency remains useful offline:** hotline, phrasebook và incident card không phụ thuộc AI online hay gói trả phí.

---
