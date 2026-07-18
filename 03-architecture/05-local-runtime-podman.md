# Runtime local bằng Podman

> Tách từ `system-design.md` — mục 5.8.

### 5.8 Runtime container local bằng Podman

Demo dùng Podman và <code>podman compose</code>; không yêu cầu Docker daemon. File <code>compose.yaml</code> ở repository root khởi động các dependency stateful:

- PostgreSQL/PostGIS cho source of truth và geo.
- Redis cho cache/session/realtime state.
- MinIO làm object storage tương thích S3 trong local.
- OpenSearch cho full-text/geo/feed retrieval.

Trong profile demo, model không chạy trong stack local: backend gọi FPT AI Factory qua HTTPS. Khi source code backend/web được thêm, các service app sẽ được bổ sung vào cùng manifest hoặc file override, nhận cấu hình từ <code>.env</code> nhưng chỉ backend được nhận <code>FPT_AI_API_KEY</code>.

~~~bash
cp .env.example .env
# Điền API key và thay toàn bộ mật khẩu mẫu trong .env
podman compose up -d
podman compose ps
~~~

<code>.env.example</code> chỉ chứa tên biến và giá trị demo, được commit. <code>.env</code> chứa secret thật, bị loại khỏi Git. Cấu hình hiện tại là single-node dành cho máy phát triển/hackathon, không phải topology production.

---
