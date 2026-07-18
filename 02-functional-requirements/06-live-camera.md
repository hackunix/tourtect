# Tourtect Lens

> Tách từ `system-design.md` — mục 4.6.

### 4.6 Tourtect Lens

| ID | Yêu cầu |
| --- | --- |
| LV-01 | Context frame không vượt 1 FPS |
| LV-02 | Tự giảm FPS khi mạng yếu, pin yếu hoặc thiết bị nóng |
| LV-03 | High-resolution capture luôn cần thao tác người dùng |
| LV-04 | Cảnh báo đỏ cần OCR/user confirmation và Price Engine |
| LV-05 | Dừng capture ngay khi background, revoke permission hoặc end session |
| LV-06 | Local inference không upload frame; mọi frame/capture offload chỉ gửi sau consent và raw media server tuân TTL |
