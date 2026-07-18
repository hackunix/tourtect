# Tourtect — Travel Price Transparency & Community Safety Shield

Bộ tài liệu này được cấu trúc để phát triển hệ thống Tourtect: giải pháp quét dị thường giá cả và bảo vệ an toàn cho khách du lịch.

---

## 1. Bản Đồ Tài Liệu Hệ Thống (System Design)

Tài liệu thiết kế gốc được phân rã thành các phần chuyên biệt tại thư mục gốc để giảm context tokens cho tác vụ AI:

| Thư mục | Nội dung |
| --- | --- |
| `00-overview/` | Tóm tắt, mục tiêu và phạm vi sản phẩm |
| `01-product-experience/` | Trải nghiệm người dùng (UX) theo feature và kênh |
| `02-functional-requirements/` | Danh sách Requirement ID dùng nghiệm thu |
| `03-architecture/` | Kiến trúc module, boundary, model routing và runtime |
| `04-data/` | Knowledge graph, dữ liệu giá/scam và thuật toán |
| `05-api/` | Endpoint, public type, WebSocket và response mẫu |
| `06-operations-safety/` | Offline, privacy, security, metric và testing |
| `07-delivery/` | Demo, roadmap, risk, assumption và reference |

*Chi tiết cách nạp tài liệu cho AI, xem tại [AGENTS.md](file:///home/hanitav/Documents/GitHub/tourtect/AGENTS.md) và [CONTEXT_MAP.md](file:///home/hanitav/Documents/GitHub/tourtect/CONTEXT_MAP.md).*

---

## 2. Kiến Trúc Go Backend Monolith (`/backend`)

Tourtect sử dụng kiến trúc Modular Monolith viết bằng Go, tối ưu hiệu năng và độ tin cậy.

### Cấu Trúc Thư Mục
* `cmd/`: Các điểm chạy chính (Binaries)
  * `api/`: REST API Server phục vụ Web/Android (port `8080`).
  * `realtime/`: WebSocket Server xử lý Push-To-Talk streaming (port `8081`).
  * `worker/`: Background worker xử lý các tác vụ outbox bất đồng bộ (port `8082`).
* `internal/`: Logic nghiệp vụ đóng gói theo module miền
  * `places/`: Quản lý danh mục địa điểm du lịch, tích hợp truy vấn PostGIS geofencing.
  * `content/`: Diễn đàn cộng đồng, quy trình nháp & xuất bản bài viết cảnh báo.
  * `pricing/`: **Price Engine** — công cụ phân tích độ lệch giá và xác định alert level (typical/elevated/high_risk) hoàn toàn deterministic dựa trên phân vị cohort (P10, P50, P90).
  * `safety/`: **Safety Engine** — công cụ phân tích sự cố an ninh khẩn cấp dựa trên quy tắc ưu tiên (Rule-First) để trả về hành động an toàn và hotline cứu hộ.
  * `platform/`: Thư viện lõi (Config, Database pool pgx, slog Logger lọc nhạy cảm, Middleware stack).
* `adapters/fptai/`: Adapter tích hợp với FPT AI và fake client phục vụ kiểm thử.
* `db/`: Migrations quản lý bằng `goose` và queries biên dịch bằng `sqlc`.
* `generated/`: Code Go tự động sinh từ OpenAPI schema và SQL queries.

---

## 3. Web Client (`/web`)

Dự án Next.js App Router (TypeScript) đóng vai trò là Client tiêu thụ REST API của Go backend.
* **Quy tắc**: Không truy cập trực tiếp Postgres, không tự xử lý Price/Safety engine logic.
* **Giao diện**: Glassmorphic UI hiện đại, phối màu Neon Obsidian (HSL), tối ưu hóa responsive di động.
* **Tích hợp API**: Tự động tạo UUIDv4 tracing trên header `X-Request-ID` cho mọi request.

---

## 4. Android Client (`/android`)

Dự án Kotlin Native được tổ chức theo kiến trúc đa module hiện đại:
* `:app`: Module khởi chạy và tích hợp điều hướng.
* `:core:network`: Tích hợp OkHttp/Retrofit giao tiếp với Go backend.
* `:core:database`: Bộ nhớ đệm Room DB phục vụ ngoại tuyến (Offline-First).
* `:core:designsystem`: Bộ components giao diện Jetpack Compose dùng chung.
* `:core:model` & `:core:security`: Dữ liệu dùng chung và lưu trữ khóa mật mã (Keystore).
* `:feature-forum` & `:feature-safety`: Các màn hình tính năng giao diện người dùng.

---

## 5. Quy Trình Chạy & Phát Triển Cục Bộ (Local Playbook)

### 1. Chuẩn bị Môi trường
Cần chuẩn bị sẵn Go (1.26+), Node.js, và Podman (hoặc Docker).

Cài đặt các công cụ biên dịch:
```bash
make bootstrap
```

### 2. Khởi chạy Hạ tầng (Databases)
```bash
# Khởi chạy Postgres, Redis, MinIO
make infra-up

# Kiểm tra trạng thái container
make infra-status
```

### 3. Migrations & Seeding Dữ liệu mẫu Hanoi
```bash
# Chạy database migrations
make db-migrate

# Nạp dữ liệu Hanoi giả lập
make db-seed
```

### 4. Biên dịch & Kiểm thử Go Backend
```bash
# Chạy bộ test suite (22+ tests cho Price/Safety engines)
make test

# Biên dịch ra 3 binaries
make api
make realtime
make worker
```

### 5. Khởi chạy Ứng dụng
* **Chạy API Backend**: `./backend/bin/api` (port `8080`)
* **Chạy Web Frontend**: `cd web && npm run dev` (port `3000`)
* **Kiểm tra sức khỏe Backend**:
  ```bash
  make verify-all
  ```
