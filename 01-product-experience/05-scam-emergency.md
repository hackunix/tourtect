# Scam Assistant và Emergency

> Tách từ `system-design.md` — mục 3.6–3.7.

### 3.6 Scam Assistant và safety triage

Thứ tự xử lý:

1. Rule kiểm tra red flag xác định.
2. Trích dữ kiện có schema.
3. Đối chiếu scam pattern đã duyệt.
4. Hỏi tối đa 1–3 câu phân luồng nếu cần.
5. Trả playbook ngắn và tùy chọn escalation.

Mức ưu tiên:

| Mức | Ví dụ | Hành vi |
| --- | --- | --- |
| Critical | Vũ khí, thương tích, bị giữ/nhốt, không thể rời đi | Hiện SOS ngay, hành động an toàn trước |
| Urgent | Ép thanh toán, tài xế không cho xuống, tình hình leo thang | Hướng dẫn rời nơi an toàn, bảo toàn bằng chứng sau |
| Non-emergency | Nghi bị tính cao, ghost tour sau khi đã an toàn | Kiểm tra giá/pattern, hướng dẫn khiếu nại |
| Information | Hỏi cách phòng tránh | Checklist và phrasebook |

Không dùng câu “đây chắc chắn là lừa đảo”. Cách diễn đạt chuẩn: “Tình huống này có một số dấu hiệu giống mẫu rủi ro …”.

Phân loại theo **hành vi đã quan sát**, không theo định kiến về người bán. Quán bình dân, hàng rong, người bán lưu động, người xin hỗ trợ hoặc người có hoàn cảnh khó khăn không tự động làm tăng scam score. Các pattern cần theo dõi gồm: đổi giá sau khi đã thỏa thuận; cố ý nhập nhằng đơn vị/tiền tệ; khai phí bắt buộc hoặc danh nghĩa “chính thức” không có căn cứ; giao hàng/dịch vụ không được yêu cầu rồi ép trả tiền; tráo hàng/khẩu phần; giữ tiền thừa/giấy tờ; dàn dựng hư hại, va chạm hoặc cáo buộc để gây áp lực đòi bồi thường. Giá cao đơn thuần chỉ đi qua Price Engine; coercion, đe dọa hoặc không cho rời đi mới nâng Safety urgency.

### 3.7 Emergency

- Nút SOS không yêu cầu đăng nhập.
- Tối đa hai thao tác từ SOS tới dialer.
- Các lựa chọn ban đầu:
  - Tôi đang bị đe dọa/không thể rời đi.
  - Có người bị thương.
  - Cháy, tai nạn hoặc cần cứu hộ.
  - Tôi an toàn nhưng muốn báo gian lận.
  - Chế độ im lặng.
- V1 dùng:
  - <code>112</code>: tổng đài khẩn cấp quốc gia.
  - <code>113</code>: công an.
  - <code>114</code>: chữa cháy/cứu nạn.
  - <code>115</code>: cấp cứu y tế.
- Không dùng số khẩn cấp cho tranh chấp giá thông thường.
- Hotline địa phương là dữ liệu động, không hard-code ngoài safety pack đã ký.
- Incident card song ngữ tách:
  - Dữ kiện người dùng đã xác nhận.
  - Nhận định do AI đề xuất.
- Không ghi âm cuộc gọi, không nghe audio của cuộc gọi PSTN và không nói rằng app làm được việc này.
