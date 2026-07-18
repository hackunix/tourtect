# Rootless Podman, Quadlet, secret và backup

**Requirement IDs liên quan:** FO-01–FO-08, PC-05–PC-07, SE-02–SE-06  
**Phạm vi:** local full-container và máy demo Linux dùng systemd user. Production cần bổ sung ingress TLS, identity production, observability và quy trình release riêng.

## 1. Vì sao dùng Podman ở đây

- Rootless containers giới hạn blast radius của app local/demo.
- `compose.yaml` phục vụ vòng lặp local và CI tương thích Compose.
- Quadlet biến container/network/volume thành systemd user units có dependency, restart và boot management mà không cần daemon riêng.
- Podman secret mount key lúc runtime; key không nằm trong Containerfile, image layer hay public Web bundle.
- Named volume cố định giúp backup/restore có target rõ ràng giữa Compose và Quadlet.

Không chạy Compose và Quadlet đồng thời: hai workflow dùng chung `tourtect-*-data` volumes và cùng host ports.

## 2. Full stack bằng Compose

```bash
cp .env.example .env
make podman-up
podman compose --profile app ps
podman compose --profile app logs -f api web
```

Profile `app` chờ PostgreSQL healthy, API startup wrapper chạy migrations + synthetic demo seed idempotent, rồi Web chờ API healthy. Redis/MinIO cũng có healthcheck. Cách này tương thích cả `podman compose` và provider `podman-compose`. Web ở `127.0.0.1:3000`, API ở `127.0.0.1:8080`; database/object store không mở ra LAN.

Realtime/worker là tùy chọn ngoài Web scope:

```bash
podman compose --profile app --profile realtime --profile worker up --build -d
```

Dừng nhưng giữ dữ liệu:

```bash
make podman-down
```

`make infra-reset` xóa volumes và chỉ dùng khi chủ động reset dữ liệu local.

## 3. Provider key bằng Podman secret

Tạo từ file quyền hạn chế hoặc stdin:

```bash
./scripts/podman-secret.sh create /secure/path/fpt-api-key
./scripts/podman-secret.sh status
make podman-up-secret
```

`compose.secrets.yaml` mount secret thành `/run/secrets/fpt_ai_api_key`; backend đọc qua `FPT_AI_API_KEY_FILE`. Không cần key để community/Price deterministic/Safety rule-first khởi động.

Rotation có chủ đích:

```bash
./scripts/podman-secret.sh remove
./scripts/podman-secret.sh create /secure/path/new-fpt-api-key
podman compose -f compose.yaml -f compose.secrets.yaml --profile app up -d --force-recreate api
```

Không đưa key vào `.env` trên máy chia sẻ. Nếu cả env và file cùng tồn tại, giá trị env được ưu tiên để giữ tương thích local.

## 4. Máy demo dùng Quadlet

Yêu cầu Podman có Quadlet và systemd user session. Installer build hai image local, tạo file environment mode `0600`, cài units vào `~/.config/containers/systemd` và enable Web cùng toàn bộ dependency:

```bash
./scripts/install-quadlet.sh
```

Với secret đã tạo:

```bash
./scripts/install-quadlet.sh --with-fpt-secret
```

Các lệnh vận hành:

```bash
systemctl --user status tourtect-web.service
journalctl --user -u tourtect-api.service -f
systemctl --user restart tourtect-api.service
systemctl --user stop tourtect-web.service
```

Muốn stack tiếp tục sau logout, quản trị viên máy có thể bật linger cho đúng user bằng `loginctl enable-linger USER`. Đây là thay đổi host-level nên installer không tự thực hiện.

Quadlet demo tự load synthetic seed idempotent. Trước staging có dữ liệu thật, bỏ dependency `tourtect-seed.service` khỏi API unit và áp dụng quy trình migration/seed được phê duyệt riêng.

## 5. Backup và restore

Backup online PostgreSQL dạng custom dump:

```bash
./scripts/podman-backup.sh
```

Cold backup cả ba named volumes yêu cầu dừng Compose lẫn Quadlet trước:

```bash
podman compose --profile app down
systemctl --user stop tourtect-web.service tourtect-api.service tourtect-minio.service tourtect-redis.service tourtect-postgres.service
./scripts/podman-backup.sh --cold-volumes
```

Mỗi backup có `SHA256SUMS`. Restore luôn cần cờ xác nhận rõ ràng:

```bash
# PostgreSQL logical restore; PostgreSQL phải đang chạy
./scripts/podman-restore.sh --from backups/tourtect-YYYYMMDDTHHMMSSZ --confirm-overwrite

# Cold volume restore; tất cả container dùng volume phải dừng
./scripts/podman-restore.sh --from backups/tourtect-YYYYMMDDTHHMMSSZ --volumes --confirm-overwrite
```

Restore dùng `--clean` cho database hoặc thay đúng ba volumes `tourtect-postgres-data`, `tourtect-redis-data`, `tourtect-minio-data`. Luôn smoke-test `/health/ready`, feed, Price `insufficient_data` và Safety abstain trước khi mở traffic.

## 6. Kiểm tra và xử lý lỗi

```bash
curl --fail http://127.0.0.1:8080/health/ready
curl --fail http://127.0.0.1:3000/
podman healthcheck run tourtect-api
```

- Migration/seed one-shot lỗi: xem log `migrate`/`seed` trong Compose hoặc `journalctl --user -u tourtect-migrate.service`.
- Quadlet unit không xuất hiện: chạy `systemctl --user daemon-reload`, rồi xem generator error bằng `systemctl --user status tourtect-web.service`.
- Port đã dùng: đảm bảo workflow local `run-local.sh`, Compose và Quadlet không chạy cùng lúc.
- Không restore khi checksum sai hoặc khi target volume còn được container sử dụng.
