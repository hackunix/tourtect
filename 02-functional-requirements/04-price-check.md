# Price Check

> Tách từ `system-design.md` — mục 4.4.

### 4.4 Price Check

| ID | Yêu cầu |
| --- | --- |
| PC-01 | Nhận ảnh menu, hóa đơn, bảng giá hoặc đồng hồ taxi |
| PC-02 | Trích item, amount, currency, unit và các thuộc tính so sánh |
| PC-03 | Cho người dùng sửa mọi trường OCR có confidence thấp |
| PC-04 | Chọn reference snapshot theo vùng, thời điểm và vertical |
| PC-05 | Trả reference range, alert level, confidence, freshness và giải thích |
| PC-06 | Không gắn nhãn người bán hoặc công khai bằng chứng thô |
| PC-07 | Không dùng submission hiện tại để thay đổi snapshot đang so sánh |
| PC-08 | Cho phép feedback và contribution bằng consent riêng |
| PC-09 | OCR ưu tiên local khi capability/confidence đạt ngưỡng; chỉ offload server sau consent hoặc khi profile cấu hình yêu cầu |
| PC-10 | Reference cohort được chọn theo phiên bản đơn vị hành chính, pricing zone, phân khúc, loại điểm bán và bối cảnh giao dịch |
| PC-11 | Khi cohort quá mỏng, mở rộng theo fallback hierarchy có disclosure; không trộn phân khúc chỉ để đủ mẫu |
| PC-12 | Donation/solicitation không có đơn vị hàng hóa rõ không được tạo kết luận “giá hợp lý”; hành vi gây áp lực chuyển Scam/Safety Engine |
