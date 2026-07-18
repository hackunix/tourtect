# Tourtect — Travel Price Transparency & Community Safety Shield

Tourtect là nền tảng cộng đồng giúp du khách chia sẻ kinh nghiệm địa phương, kiểm tra khoảng giá tham chiếu và nhận hướng dẫn an toàn đã được backend phê duyệt. Repository gồm Go modular monolith, Next.js Web, Android client và bộ tài liệu thiết kế theo Requirement ID.

## Trạng thái hiện tại

- Web: community feed responsive, search, saved, notifications, place detail, composer draft/confirm, Price Check và Safety Assessment.
- Backend: REST API `:8080`, realtime server `:8081`, worker `:8082`, PostgreSQL/PostGIS, Redis và MinIO.
- Community API: feed Following/Nearby/Latest/Trending/Safety, comments, useful, saved, follows, reports, blocks và notifications.
- Authentication hiện dùng seed identity cho vertical slice; không phải cơ chế đăng nhập production.
- Android không dùng chung mã UI với Web và được vận hành theo module riêng.

## Chạy nhanh FE + BE

Trên Linux `x86_64` hoặc `arm64`, dùng bootstrap tự động:

```bash
./scripts/setup-linux.sh
./scripts/run-local.sh
```

`setup-linux.sh` nhận diện `apt`, `dnf`, `pacman` hoặc `zypper`, cài dependency hệ thống, Go/Node phù hợp và tool codegen. `run-local.sh` bật hạ tầng, migration, seed, API và Web; `Ctrl+C` dừng app nhưng giữ volumes.

Kiểm tra trước mà không cài gì:

```bash
./scripts/setup-linux.sh --check-only
```

Thiết lập thủ công và các tùy chọn production/backend-only nằm trong runbook.

Muốn chạy toàn bộ FE + BE bằng rootless Podman, không cần cài Go/Node trên host:

```bash
cp .env.example .env
make podman-up
```

Profile `app` build image non-root, chạy migration + synthetic seed idempotent rồi bật API/Web. Mọi cổng local bind vào `127.0.0.1`; secret AI có override riêng, không được bake vào image. Với máy demo chạy lâu dài, dùng Quadlet để systemd user quản lý restart/boot.

Mở:

- Web: <http://localhost:3000>
- API liveness: <http://localhost:8080/health/live>
- API readiness: <http://localhost:8080/health/ready>
- MinIO console: <http://localhost:9001>

Hướng dẫn đầy đủ, biến môi trường, kiểm thử và troubleshooting: [docs/operations/frontend-backend-runbook.md](docs/operations/frontend-backend-runbook.md).

## Kiểm thử và build

```bash
# Backend, cần PostgreSQL đã migrate
make test
make lint

# Web
cd web
npm run lint
npm test
npm run build
npx playwright install chromium   # chạy một lần trên máy mới
npm run test:e2e
```

Khi sửa OpenAPI, SQL query hoặc migration:

```bash
make generate
```

Commit cả file generated tương ứng và ghi Requirement ID bị ảnh hưởng.

## Cấu trúc repository

| Đường dẫn | Nội dung |
| --- | --- |
| `backend/` | Go API, realtime, worker, domain modules, migrations và generated code |
| `web/` | Next.js App Router, community UI và browser/unit tests |
| `android/` | Kotlin/Jetpack Compose đa module |
| `00-overview/`–`07-delivery/` | Tài liệu sản phẩm, requirements, architecture, API, safety và delivery |
| `docs/design/` | Design-system semantics dùng chung giữa client |
| `docs/operations/` | Runbook vận hành local và kiểm tra FE/BE |
| `deploy/podman/quadlet/` | Rootless systemd/Quadlet units cho máy demo |
| `scripts/` | Bootstrap Linux, launcher, secret, Quadlet và backup/restore |

Agent/cộng tác viên phải đọc [AGENTS.md](AGENTS.md), [CONTEXT_MAP.md](CONTEXT_MAP.md) và chỉ nạp bounded context liên quan.

## Nguyên tắc bắt buộc

- Safety/privacy ưu tiên hơn tiện lợi, engagement và monetization.
- Web/Android không truy cập trực tiếp Postgres và không tự tính Price/Safety result.
- Không để advertiser spend, affiliate commission hoặc business tier tác động organic ranking.
- Không tự sinh hotline, kết luận pháp lý hoặc cáo buộc cá nhân/doanh nghiệp lừa đảo.
- Không đóng gói provider API key trong Web public bundle hoặc Android app.
- Khi dữ liệu không đủ, trả `insufficient_data` hoặc abstain thay vì dựng dữ liệu giả.

## Tài liệu chính

- [Context map](CONTEXT_MAP.md)
- [API catalog](05-api/01-conventions-endpoints.md)
- [Local Podman runtime](03-architecture/05-local-runtime-podman.md)
- [Frontend design system](docs/design/frontend-design-system.md)
- [FE + BE operations runbook](docs/operations/frontend-backend-runbook.md)
- [Rootless Podman + Quadlet operations](docs/operations/podman-quadlet.md)
