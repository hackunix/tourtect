# Tourtect Lens — Live Camera

> Tách từ `system-design.md` — mục 3.5.

### 3.5 Tourtect Lens (Live Camera)

#### Cách hoạt động

- Preview camera chạy cục bộ ở frame rate tự nhiên.
- Client chỉ sample frame tối đa 1 FPS; mặc định 0.5 FPS.
- Context frame có cạnh dài tối đa khoảng 768 px và chỉ phục vụ hiểu bối cảnh.
- Qwen3-VL xác định loại đối tượng: menu, món ăn, đồng hồ taxi, hóa đơn, bảng giá, quầy đổi tiền hoặc không liên quan.
- Khi phát hiện vùng chữ/số ổn định, app đề nghị “Giữ máy yên và chụp để xác nhận”.
- Chỉ ảnh do người dùng bấm chụp mới đi vào pipeline high-resolution.
- OCR mobile chạy trước khi thiết bị đủ capability; user xác nhận item, giá, tiền tệ và đơn vị.
- VLM/OCR server là fallback sau consent khi local không khả dụng hoặc confidence thấp; không được tự quyết cảnh báo.

#### Camera sequence

~~~mermaid
sequenceDiagram
    actor User as Người dùng
    participant App as Android App
    participant Vision as Vision Service
    participant OCR as On-device OCR
    participant Capture as Capture API
    participant Storage as Object Storage
    participant Price as Price Engine

    User->>App: Bật Live Camera
    loop Tối đa 0.2–1 FPS
        App->>Vision: Frame thấp + session context
        Vision-->>App: vision.observation
    end
    App-->>User: Gợi ý vùng cần giữ ổn định
    User->>App: Bấm chụp để xác nhận
    App->>OCR: Ảnh độ phân giải cao
    OCR-->>App: Text, bounding boxes, confidence
    App-->>User: Xác nhận item/giá/đơn vị
    User->>App: Xác nhận hoặc sửa
    App->>Capture: Tạo capture
    Capture-->>App: capture_id + signed PUT URL
    App->>Storage: PUT ảnh đã redaction
    App->>Capture: Finalize capture + hash/redaction metadata
    Capture-->>App: Capture finalized
    App->>Price: POST price check bằng capture_id + candidate đã xác nhận
    Price-->>App: PriceInsight + range + freshness
    App-->>User: Thẻ kết quả riêng tư
~~~

#### Guardrail thị giác

- Không suy ra “giá hợp lý” chỉ từ vẻ ngoài món ăn hoặc địa điểm.
- Không cảnh báo đỏ từ frame Live chưa xác nhận.
- Không nhận diện danh tính khuôn mặt.
- Khuôn mặt, hộ chiếu, thẻ thanh toán và số điện thoại được làm mờ trước khi contribution.
- Biển số có thể được giữ cục bộ trong incident card; chỉ upload khi người dùng chủ động chia sẻ bằng chứng.
