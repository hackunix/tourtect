# AGENTS.md — quy tắc đọc context cho coding agent

## Mục tiêu

Giảm token và tránh agent bị nhiễu bởi phần hệ thống không liên quan. Không đọc toàn bộ cây tài liệu theo mặc định.

## Quy trình bắt buộc

1. Đọc `README.md` và `CONTEXT_MAP.md` trước.
2. Xác định feature/bounded context đang sửa.
3. Chỉ đọc nhóm file được map cho tác vụ đó, ưu tiên tối đa 2–4 file.
4. Chỉ mở thêm file khi gặp dependency cụ thể chưa được mô tả.
5. Dùng tìm kiếm theo Requirement ID, endpoint, interface hoặc tên module trước khi mở cả file.

## Không làm

- Không glob và đọc tất cả file Markdown.
- Không nạp toàn bộ thư mục chỉ để tìm một endpoint hoặc một requirement.
- Không suy luận rằng quảng cáo/business tier được phép tác động organic ranking, review, moderation, PriceSnapshot hoặc Safety Engine.
- Không để LLM tự sinh hotline, kết luận pháp lý hoặc cáo buộc một cá nhân/doanh nghiệp lừa đảo.
- Không đóng gói provider API key trong Android app.

## Khi sửa tài liệu hoặc code

- Thay đổi UX: kiểm tra file experience + requirement tương ứng.
- Thay đổi domain/data: kiểm tra data model + public type + endpoint.
- Thay đổi realtime: kiểm tra Live Voice/Lens + architecture + WebSocket protocol.
- Thay đổi safety/privacy: kiểm tra requirement + operations/safety + API/data liên quan.
- Ghi rõ Requirement ID bị ảnh hưởng trong PR/commit.

## Quy tắc ưu tiên

Safety và privacy > tính năng tiện lợi > engagement > monetization. Khi dữ liệu không đủ, hệ thống phải abstain hoặc trả `insufficient_data`.
