# Tourtect

Thiết kế và scaffold hạ tầng local cho demo Tourtect.

## Chạy hạ tầng bằng Podman

Yêu cầu: Podman 4.7+ và một Compose provider cho `podman compose` (ví dụ `podman-compose`). Có thể kiểm tra bằng `podman compose version`.

~~~bash
cp .env.example .env
# Điền FPT_AI_API_KEY và đổi các mật khẩu trong .env
podman compose up -d
podman compose ps
~~~

Các container local gồm PostgreSQL/PostGIS, Redis, MinIO và OpenSearch. Model AI không chạy trong container local; backend gọi FPT AI Factory qua API tương thích OpenAI. Các biến `FPT_AI_*` trong `.env` sẽ được truyền vào container backend khi service này được thêm vào manifest.

Riêng deployment demo dùng Android native và đặt `AI_EXECUTION_MODE=server_only`: STT, OCR, VLM, LLM và TTS đều chạy qua backend/FPT AI Factory. Đây chỉ là env override; thiết kế tổng thể vẫn hỗ trợ `adaptive` và `local_only`. Điền model ID STT/TTS mà API key được cấp quyền trong `.env`; không đưa FPT API key vào APK. Với Android Emulator, URL backend local mặc định dùng `10.0.2.2`; iOS tạm ngoài phạm vi.

Kiểm tra nhanh:

~~~bash
podman compose exec postgres pg_isready -U tourtect -d tourtect
curl http://localhost:9200/_cluster/health
~~~

Dừng stack mà vẫn giữ dữ liệu:

~~~bash
podman compose down
~~~

Chỉ dùng `podman compose down -v` khi muốn xóa toàn bộ volume dữ liệu local.

Xem [thiết kế hệ thống](docs/system-design.md) để biết kiến trúc và phạm vi demo.
