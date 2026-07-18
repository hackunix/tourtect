# Kịch bản Demo Hackathon

> Tách từ `system-design.md` — mục 11.

## 11. Demo hackathon

### 11.1 Demo 1 — Forum đa ngôn ngữ theo địa điểm

Khách Hàn tìm một khu phố Hà Nội, đọc place page bằng tiếng Hàn và thấy review, khoảng giá, price report, scam pattern, câu hỏi cộng đồng và nguồn báo/video. Người dùng chuyển về nguyên bản tiếng Việt, xem evidence/freshness và theo dõi địa điểm. Demo chứng minh forum là lõi, không phải màn hình phụ của AI.

### 11.2 Demo 2 — Price report tạo knowledge có kiểm soát

Khách chụp hóa đơn bằng Tourtect Lens, sửa một trường OCR và xác nhận phường/pricing zone, phân khúc bình dân cùng loại điểm bán. Engine chọn đúng cohort, công khai việc fallback địa lý nếu có và trả so sánh riêng tư. Người dùng chủ động tạo draft price report; app khử PII, yêu cầu disclosure/phạm vi công khai, moderation gắn evidence badge và chỉ đưa observation vào quarantine. Place page cập nhật post ngay khi hợp lệ nhưng reference snapshot chỉ đổi sau pipeline review/versioning.

### 11.3 Demo 3 — Tổng hợp báo và video đúng quyền

Editor nhập một RSS item và một YouTube URL. Connector kiểm tra quyền/embed, lấy metadata, canonicalize, dedupe, gắn với place/scam topic và tạo external card có attribution. Khi mô phỏng video bị disable, card ngừng phát và search index được cập nhật.

### 11.4 Demo 4 — Tourtect Live là trợ lý của forum

Khách Trung dùng PTT với người bán Việt Nam. Translation trả trước, intelligence lane tạo price candidate kín. Sau khi kết thúc, app chỉ đề nghị “Lưu riêng” hoặc “Tạo báo cáo”; không tự đăng transcript/media. SOS vẫn mở dialer không cần tài khoản hoặc gói Plus.

### 11.5 Demo 5 — Doanh thu không mua được niềm tin

Feed hiển thị một quảng cáo contextual có nhãn, một affiliate CTA có disclosure và một place đã mua Verified Business. Bật/tắt spend/business tier không làm đổi thứ tự organic, điểm review, price alert hay moderation. Chuyển sang incident critical thì toàn bộ quảng cáo biến mất.

### 11.6 Dữ liệu demo

- Dữ liệu thật có URL/provenance và ngày xác minh.
- Dữ liệu synthetic gắn badge “Demo data”.
- Synthetic data không được đưa vào snapshot pilot/production.
- Không gọi thử hotline khẩn cấp trong demo.

---
