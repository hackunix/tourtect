# Offline và chế độ suy giảm

> Tách từ `system-design.md` — mục 8.

## 8. Offline và chế độ suy giảm

### 8.1 Safety pack

Mỗi vùng có gói ký số:

- Hotline toàn quốc và địa phương đã xác minh.
- <code>verified_at</code>, giờ hoạt động, nguồn và fallback.
- Scam playbook rút gọn.
- Câu khẩn cấp song ngữ và audio đã duyệt.
- Incident card template.
- Reference price cache rút gọn, kèm freshness/version.
- Public key/version metadata để client xác minh.

Client giữ last-known-good pack. Pack hết hạn vẫn hiển thị hotline toàn quốc nhưng phải gắn nhãn dữ liệu địa phương đã cũ.

### 8.2 Ma trận suy giảm

| Trạng thái | Live Voice | Live Camera | Price | Emergency |
| --- | --- | --- | --- | --- |
| Mạng tốt | Full server PTT/STT/MT/TTS | 0.5–1 FPS + server VLM/OCR | Realtime server | Full directory |
| Mạng yếu | Audio ưu tiên | 0.2 FPS hoặc tắt | Manual/cached | Full offline pack |
| Server vision/OCR lỗi | Không ảnh hưởng dịch | Manual crop/input, không AI | Nếu user nhập đủ dữ kiện | Không ảnh hưởng |
| Server MT lỗi | Caption transcript gốc/phrase | Camera vẫn chạy | Vẫn kiểm tra được input manual | Phrasebook |
| Offline | Phrase/audio pack | Preview local, không AI vision | Cache/manual, không red alert nếu stale | Hotline + incident card |
| Thiết bị nóng/pin yếu | Giảm audio effect, vẫn server ASR | Giảm/tắt frame sampling | Manual | Không ảnh hưởng |

Quy tắc quan trọng: offline/stale data không phát cảnh báo đỏ mới.

---
