# Tourtect Live

> Tách từ `system-design.md` — mục 4.5.

### 4.5 Tourtect Live

| ID | Yêu cầu |
| --- | --- |
| LT-01 | Hai hướng dịch được khóa bởi SpeakerRole |
| LT-02 | PCM chỉ được capture trong lúc giữ PTT |
| LT-03 | Translation lane không chờ Price Engine |
| LT-04 | TTS và caption đều có nút phát lại/báo lỗi |
| LT-05 | Bảo toàn số, phủ định, đơn vị, tiền tệ và thực thể khẩn cấp |
| LT-06 | Có fallback caption và phrase audio khi TTS thiếu locale |
| LT-07 | Model router chọn local/server theo execution policy, capability, consent và confidence; mọi lựa chọn được ghi trong ModelTrace |
