# Runbook vận hành Tourtect Frontend + Backend

**Requirement IDs liên quan:** FO-01–FO-08, PC-05–PC-07, SE-02–SE-06  
**Phạm vi:** môi trường local/development cho Next.js Web và Go API. Đây không phải topology production.

## 1. Thành phần và cổng

| Thành phần | Cổng mặc định | Bắt buộc cho Web community |
| --- | ---: | --- |
| Next.js Web | `3000` | Có |
| Go REST API | `8080` | Có |
| PostgreSQL/PostGIS | `5432` | Có |
| Redis | `6379` | Stack local, chưa nằm trên critical path của community feed |
| MinIO API/Console | `9000/9001` | Stack local, cần cho luồng media tương lai |
| Realtime server | `8081` | Không, chỉ Live Voice/Lens |
| Worker | `8082` | Không cho read/write community trực tiếp |
| OpenSearch | `9200` | Không; profile `search` là tùy chọn, search hiện tại dùng PostgreSQL |

## 2. Chuẩn bị lần đầu

### Công cụ

- Linux `x86_64` hoặc `arm64` cho script bootstrap; distro khác vẫn có thể setup thủ công.
- Go `1.26.5+`.
- Node.js/npm tương thích Next.js `16.2.x`.
- Podman và compose provider.
- `curl`; Chromium chỉ cần cho Playwright E2E.

### Bootstrap Linux tự động

```bash
./scripts/setup-linux.sh
```

Script hỗ trợ `apt`, `dnf`, `pacman`, `zypper`; cài base packages qua `sudo`, cài Go/Node thiếu phiên bản vào `~/.local`, chạy `npm ci`, cài `goose/sqlc/oapi-codegen` và chỉ tạo `.env` khi file chưa tồn tại. Download Go/Node được kiểm tra SHA-256 và không dùng `curl | bash`.

Các chế độ hữu ích:

```bash
./scripts/setup-linux.sh --check-only
./scripts/setup-linux.sh --install-browser
./scripts/setup-linux.sh --skip-system-packages
./scripts/setup-linux.sh --start
```

Sau setup, mở shell mới hoặc chạy `source ~/.config/tourtect/env.sh`.

### Setup thủ công

```bash
cp .env.example .env
make bootstrap
cd web && npm ci && cd ..
```

Rà `.env` trước khi chạy stack:

- Quick-start local có thể giữ credential demo nếu các cổng chỉ bind trên máy phát triển; không expose stack đó ra mạng công khai.
- Với môi trường chia sẻ hoặc credential khác, đổi `POSTGRES_PASSWORD`, `REDIS_PASSWORD`, `MINIO_ROOT_PASSWORD` và truyền DSN tương ứng cho migration/test.
- `FPT_AI_API_KEY` chỉ cần cho luồng gọi provider và chỉ được cấp cho backend; community feed, Price Engine deterministic và Safety Engine rule-first không cần key để khởi động.
- Không commit `.env` và không chuyển secret thành `NEXT_PUBLIC_*`.

Lưu ý: các target database trong `Makefile` đang dùng DSN local mặc định `tourtect/change_me_postgres`. Nếu đổi password/port, không dùng nguyên trạng `make db-migrate`; chạy migration bằng DSN tương ứng:

```bash
~/go/bin/goose -dir backend/db/migrations postgres "$DATABASE_URL" up
```

## 3. Khởi động hạ tầng

### Launcher đề xuất

```bash
./scripts/run-local.sh
```

Launcher khởi động dependencies, chờ PostgreSQL healthy, chạy migration/seed idempotent, build API và chạy Next.js dev server. `Ctrl+C` chỉ dừng app; containers và named volumes được giữ lại.

```bash
./scripts/run-local.sh --production     # next build + next start
./scripts/run-local.sh --backend-only   # chỉ API
./scripts/run-local.sh --no-seed
./scripts/run-local.sh --skip-migrate
./scripts/run-local.sh --with-search
```

### Chạy từng bước

```bash
make infra-up
make infra-status
```

`infra-up` chỉ bật PostgreSQL, Redis và MinIO. Khi thực sự cần profile OpenSearch:

```bash
podman compose --profile search up -d opensearch
```

Chờ PostgreSQL báo `healthy`, sau đó:

```bash
make db-migrate
make db-seed
```

Seed chỉ dùng dữ liệu tổng hợp Hanoi. Có thể chạy lại an toàn nhờ các khóa/`ON CONFLICT` đã định nghĩa.

### Chạy toàn bộ bằng container

Không cần Go/Node trên host nếu dùng profile `app`:

```bash
cp .env.example .env
podman compose --profile app up --build -d
podman compose --profile app ps
```

Profile này build backend/Web image; API startup wrapper chạy migration và synthetic seed idempotent trước khi listen, sau đó Web chờ API healthy. Image app chạy bằng user không đặc quyền; tất cả host port chỉ bind `127.0.0.1`. Dùng `podman compose --profile app logs -f api web` để theo dõi.

Khi cần provider key, không truyền `FPT_AI_API_KEY` qua build arg. Tạo Podman secret và dùng override:

```bash
./scripts/podman-secret.sh create /secure/path/fpt-api-key
podman compose -f compose.yaml -f compose.secrets.yaml --profile app up --build -d
```

Chi tiết Quadlet, backup/restore và rotation secret nằm trong [podman-quadlet.md](podman-quadlet.md).

## 4. Chạy Backend

### Cách build binary

```bash
make api
./backend/bin/api
```

Hoặc chạy trực tiếp khi phát triển:

```bash
cd backend
go run ./cmd/api
```

Backend đọc các biến chính:

| Biến | Mặc định | Ý nghĩa |
| --- | --- | --- |
| `PORT` | `8080` | Cổng REST API |
| `DATABASE_URL` | Ghép từ các biến `POSTGRES_*` | PostgreSQL DSN |
| `POSTGRES_HOST` | `localhost` | Host database khi không có DSN |
| `POSTGRES_DB/USER/PASSWORD/PORT` | `tourtect/tourtect/change_me_postgres/5432` | Thành phần DSN local |
| `REDIS_HOST/PORT/PASSWORD` | `localhost/6379/empty` | Redis local |
| `LOG_LEVEL` | `info` | Mức log |
| `FPT_AI_BASE_URL/API_KEY` | FPT endpoint/empty | Provider server-side |
| `FPT_AI_API_KEY_FILE` | empty | File secret mount; biến `FPT_AI_API_KEY` có giá trị sẽ được ưu tiên |

Nếu muốn nạp `.env` vào shell trước khi chạy binary:

```bash
set -a
source .env
set +a
./backend/bin/api
```

Không in `.env`, bearer token, cookie hoặc API key vào log/troubleshooting output.

### Health và smoke test

```bash
curl --fail http://localhost:8080/health/live
curl --fail http://localhost:8080/health/ready
curl --fail 'http://localhost:8080/v1/feed?mode=latest&limit=2'
curl --fail 'http://localhost:8080/v1/feed?mode=nearby&region_id=hanoi-hoan-kiem&limit=2'
curl --fail 'http://localhost:8080/v1/search?q=taxi&tab=safety'
```

`Nearby` yêu cầu `region_id` hoặc input vị trí opt-in hợp lệ; server không tự lấy precise location. Các mutation account-dependent hiện dùng seed identity từ `AuthBoundary`, chỉ dành cho vertical slice.

## 5. Chạy Frontend

```bash
cd web
npm install
npm run dev
```

Localhost mặc định dùng API `http://localhost:8080`. Nếu backend chạy ở URL khác:

```bash
API_URL=http://127.0.0.1:8080 \
NEXT_PUBLIC_API_URL=http://127.0.0.1:8080 \
npm run dev
```

- `API_URL` phục vụ Server Components và không được expose tự động ra browser.
- `NEXT_PUBLIC_API_URL` phục vụ Client Components; chỉ chứa base URL công khai, không chứa secret.
- API phải cho phép origin Web qua CORS. Mỗi request Web tự gắn `X-Request-ID`.

Các route vận hành chính:

| Route | Chức năng |
| --- | --- |
| `/` | Feed, composer, Price/Safety context |
| `/search` | Search Places/Posts/Price reports/Safety |
| `/saved` | Bài đã lưu của demo identity |
| `/notifications` | Notification của demo identity |
| `/places/{id}` | Place detail và content liên quan |
| `/profile` | Nhãn trạng thái demo identity |

Web không hiển thị fake posts khi backend lỗi. Error state phải có retry và Request ID nếu API trả về.

## 6. Kiểm thử và code generation

### Backend

```bash
make test
make lint
```

Price/Safety integration tests cần PostgreSQL đang chạy, đã migrate và nhận đúng `DATABASE_URL` nếu credential khác mặc định. Sau khi sửa `backend/api/openapi.yaml`, migration hoặc `backend/db/queries`:

```bash
make generate
git diff --check
```

Commit generated OpenAPI/sqlc files cùng source. Không sửa trực tiếp generated code.

### Frontend

```bash
cd web
npm run lint
npm test
npm run build
```

E2E + accessibility:

```bash
npx playwright install chromium  # một lần trên máy mới
npm run test:e2e
```

Playwright tự chạy Next.js ở `127.0.0.1:3000`. Để kiểm tra data path đầy đủ, giữ API/PostgreSQL hoạt động; khi API tắt, test shell vẫn xác nhận designed error state và không dùng dữ liệu giả.

## 7. Quy trình tắt và phục hồi

Dừng tiến trình FE/BE bằng `Ctrl+C`. Để dừng container nhưng giữ dữ liệu, dùng:

```bash
podman compose --profile app down
```

`make infra-down` và `make podman-down` giữ named volumes. Chỉ `make infra-reset` mới truyền `-v` và xóa dữ liệu local.

Phục hồi môi trường local sạch chỉ khi chủ động chấp nhận mất dữ liệu:

1. Dừng FE/BE.
2. Chạy `make infra-reset`.
3. Chạy lại `make infra-up`, `make db-migrate`, `make db-seed`.

## 8. Troubleshooting

### API readiness trả `503`

- Kiểm tra `make infra-status` và PostgreSQL health.
- Xác nhận DSN/password/port của API trùng với compose.
- Chạy migration còn thiếu rồi thử `/health/ready` lại.

### Feed trả `422`

- `mode` chỉ nhận `following`, `nearby`, `latest`, `trending`, `safety`.
- `nearby` cần vùng do người dùng chọn; dùng `region_id=hanoi-hoan-kiem` cho seed local.

### Web chỉ hiện unavailable state

- Kiểm tra `API_URL` và `NEXT_PUBLIC_API_URL`.
- Gọi trực tiếp `/health/ready` và `/v1/feed`.
- Mở technical details của error state và đối chiếu Request ID với backend log.

### Next.js build không tải được Noto Sans

`next/font` tải font ở build time rồi self-host trong output. Cho phép outbound HTTPS tới Google Fonts trong môi trường build hoặc cung cấp font local được cấp phép; không quay lại `@import` font runtime.

### Go test không kết nối được PostgreSQL

- Khởi động/migrate database trước.
- Nếu sandbox chặn socket local, chạy suite trong môi trường có quyền truy cập `localhost:5432`.
- Không biến lỗi kết nối thành test pass giả; integration suite phải chạy lại trước khi release.

## 9. Safety và dữ liệu vận hành

- Hotline chỉ đến từ safety directory có version; không nhập số tùy ý vào UI hoặc seed ngoài quy trình duyệt.
- Scam report sau xác nhận chuyển `pending`, không được coi là kết luận gian lận.
- Price result `high_risk` không phải bằng chứng pháp lý; `insufficient_data` phải giữ trạng thái trung tính.
- Không đưa advertiser spend, affiliate hoặc business tier vào query/ranking organic.
- Không lưu raw incident/media hoặc precise location ngoài consent/retention policy tương ứng.
