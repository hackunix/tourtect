# Tourtect System Design — bản tách module

Bộ tài liệu này được tách từ `system-design.md` để AI coding agent chỉ đọc đúng ngữ cảnh cần thiết, thay vì nạp toàn bộ tài liệu khoảng 165 KB mỗi tác vụ.

## Cách dùng nhanh

1. Đưa cho agent `AGENTS.md` và yêu cầu tuân thủ phạm vi đọc.
2. Chọn tác vụ trong `CONTEXT_MAP.md`.
3. Chỉ nạp các file được chỉ định; thường 2–4 file là đủ.
4. Dùng `MANIFEST.md` khi cần đối chiếu nguồn hoặc kiểm tra phần nào đã được tách.

## Cấu trúc

| Thư mục | Nội dung |
| --- | --- |
| `00-overview/` | Tóm tắt, mục tiêu và phạm vi sản phẩm |
| `01-product-experience/` | UX theo feature và kênh |
| `02-functional-requirements/` | Requirement ID dùng để nghiệm thu |
| `03-architecture/` | Kiến trúc, boundary, model routing và runtime |
| `04-data/` | Knowledge graph, dữ liệu giá/scam và thuật toán |
| `05-api/` | Endpoint, public type, WebSocket và response mẫu |
| `06-operations-safety/` | Offline, privacy, security, metric và testing |
| `07-delivery/` | Demo, roadmap, risk, assumption và reference |

## Nguyên tắc cập nhật

Một thay đổi feature thường phải kiểm tra đồng thời: trải nghiệm người dùng, requirement ID, kiến trúc/data và API contract. `CONTEXT_MAP.md` đã chỉ ra các nhóm file tương ứng.
