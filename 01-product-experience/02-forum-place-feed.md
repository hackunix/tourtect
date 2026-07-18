# Forum, Place Page và Feed

> Tách từ `system-design.md` — mục 3.3.

### 3.3 Forum, place page và feed

#### Loại nội dung

| Loại | Trường bắt buộc | Cơ chế tin cậy |
| --- | --- | --- |
| Discussion/Q&A | Tiêu đề, nội dung, topic hoặc place | Reputation, vote hữu ích, moderation |
| Review | Place, thời điểm trải nghiệm, rating có cấu trúc, disclosure quan hệ | Merchant reply, evidence badge, appeal |
| Price report | Item, amount, currency, unit, vùng, thời điểm | Xác nhận OCR/manual, provenance, quarantine trước khi vào Price Engine |
| Scam report | Mô tả hành vi, vùng, thời điểm, mức an toàn hiện tại | Safety triage, khử PII, không tự nêu danh tính cá nhân |
| Tip/itinerary | Khu vực, thời gian hiệu lực, nội dung hướng dẫn | Freshness, saved/share, community correction |
| Official alert | Nguồn có thẩm quyền, ngày hiệu lực, canonical URL | Editor verification và expiry |
| External article/video | Canonical URL, nguồn, tác giả, ngày xuất bản | Rights/policy gate, attribution, dedup, takedown sync |

#### Place page

Mỗi địa điểm/cơ sở/dịch vụ có một trang chuẩn hóa gồm:

- Tên đa ngôn ngữ, loại hình, khu vực, bản đồ với độ chính xác phù hợp và thông tin do chủ cơ sở claim.
- Điểm tổng hợp và phân bố review; các chiều riêng: minh bạch giá, chất lượng, an toàn và mức đáng tiền.
- Khoảng giá hiện tại, lịch sử thay đổi, cỡ mẫu, freshness và nút thêm price report.
- Scam/safety signal theo hành vi đã được kiểm chứng; không hiển thị “scam score” kết luận pháp lý.
- Post, câu hỏi, bài báo và video liên quan; câu trả lời chính thức/merchant reply được gắn nhãn.
- Nút lưu, theo dõi, chia sẻ, báo nội dung, đề xuất chỉnh sửa và kháng nghị.

#### Feed và ranking

Feed gồm Following, Nearby, Latest, Trending và Safety Alerts. Người dùng chọn vùng thủ công trước; vị trí chính xác là opt-in. Ranking hữu cơ chỉ dùng relevance, freshness, evidence, source diversity, community usefulness và safety priority. Các tín hiệu sau bị cấm trong organic ranker: chi tiêu quảng cáo, gói business, affiliate commission và yêu cầu của sales.

Tourtect hiển thị bản dịch theo locale người đọc nhưng luôn cho xem bản gốc. Nội dung AI dịch/tóm tắt có nhãn, nút báo lỗi và không thay thế nguyên văn của tác giả.

#### Reputation và moderation

- Reputation tăng từ contribution được xác minh, sửa lỗi hữu ích, lịch sử ổn định và đánh giá của moderator; không tăng chỉ vì đồng thuận hoặc nhiều upvote.
- Tách reputation theo năng lực: local knowledge, price evidence, translation, safety và merchant representative.
- Phát hiện spam, review exchange, brigading, sockpuppet/Sybil, copy-paste và bất thường theo thời gian/thiết bị/graph.
- Nội dung rủi ro cao đi qua human review; người đăng và đối tượng bị review có quy trình notice, phản hồi và appeal.
- Không có “pay to remove”. Nội dung có thể bị gỡ vì sai, vi phạm, thiếu căn cứ hoặc theo quy trình pháp lý, và quyết định phải có audit trail.
