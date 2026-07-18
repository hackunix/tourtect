# Android Client

> Tách từ `system-design.md` — mục 4.8.

### 4.8 Android client

| ID | Yêu cầu |
| --- | --- |
| AN-01 | App native Kotlin/Jetpack Compose; CameraX cho Lens và AudioRecord cho PTT |
| AN-02 | Không bao giờ đóng gói <code>FPT_AI_API_KEY</code>; flavor demo server-only không đóng gói model, còn model pack tương lai phải signed, versioned và tải theo capability |
| AN-03 | Camera/micro dừng khi background, permission revoke, session end hoặc process bị reclaim |
| AN-04 | Token lưu bằng Android Keystore-backed storage; notification không chứa transcript/incident/raw safety data |
| AN-05 | Room/DataStore chỉ cache dữ liệu cần thiết; raw media cache mã hóa có TTL và cleanup receipt |
| AN-06 | iOS build, shared UI layer và iOS-specific behavior nằm ngoài demo |

---
