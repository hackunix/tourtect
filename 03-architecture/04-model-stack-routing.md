# Model Stack, Adaptive Routing và Gemini Live

> Tách từ `system-design.md` — mục 5.5–5.7.

### 5.5 Stack model và profile demo

| Nhiệm vụ | Thiết kế adaptive | Fallback | Override trong demo <code>server_only</code> |
| --- | --- | --- | --- |
| Menu/hóa đơn/bảng giá | Signed OCR model pack trên Android | Server VLM/OCR + user sửa manual | <code>FPT_AI_OCR_MODEL</code> |
| Món ăn/đồ vật/bối cảnh | Server VLM | Model VLM khác được cấp quyền | <code>FPT_AI_VISION_MODEL</code> |
| ASR | Signed ASR model pack khi device đủ capability | Server ASR | <code>FPT_AI_STT_MODEL</code> |
| Dịch trực tiếp | Server translation model | Phrase pack đã duyệt | <code>FPT_AI_TEXT_MODEL</code> |
| Trích dữ kiện/playbook | Server structured extraction | Rule-only/degraded | <code>FPT_AI_TEXT_MODEL</code> |
| TTS | Android system TTS | Server TTS/caption/phrase audio | <code>FPT_AI_TTS_MODEL</code> |
| Giá | Rule + robust statistics + LightGBM/CatBoost trên server | Không LLM | Không đổi |
| Emergency | Rule server + dữ liệu đã duyệt | Signed safety pack trên Android | Không đổi |

Profile demo đặt <code>AI_EXECUTION_MODE=server_only</code> và dùng FPT AI Inference qua API tương thích OpenAI tại <code>https://mkp-api.fptcloud.com</code>, xác thực bằng <code>Authorization: Bearer &lt;FPT_AI_API_KEY&gt;</code>. Đây là deployment override, không thay đổi kiến trúc adaptive. Model ID thực tế phụ thuộc region, gói dịch vụ và quyền gắn với API key, vì vậy backend đọc model ID từ env thay vì hard-code.

API key chỉ tồn tại ở backend/server-side secret. Web, Android và Zalo Mini App gọi Tourtect API, tuyệt đối không gọi FPT AI Factory trực tiếp. Adapter phải đặt timeout, retry có backoff cho lỗi tạm thời, rate limit, giới hạn token/ảnh và log <code>trace_id</code> nhưng không log API key hoặc raw media.

Tất cả model được bọc bằng interface:

- <code>ASRProvider</code>
- <code>TranslationProvider</code>
- <code>VisionProvider</code>
- <code>OCRProvider</code>
- <code>TTSProvider</code>
- <code>ExtractorProvider</code>
- <code>PriceModelProvider</code>

Không gọi trực tiếp model từ controller hoặc UI.

### 5.6 Adaptive model routing

~~~mermaid
flowchart TD
    START["Bắt đầu utterance/capture"]
    POLICY{"AI_EXECUTION_MODE?"}
    LOCAL_GATE{"Local model hợp lệ\nvà device đủ capability?"}
    LOCAL["Chạy on-device"]
    LOCAL_CONF{"Local confidence đạt ngưỡng?"}
    CONSENT{"Đã consent gửi dữ liệu\ntối thiểu lên server?"}
    SERVER["Chạy server provider\nchính theo tác vụ"]
    SERVER_CONF{"Server confidence đạt ngưỡng?"}
    FALLBACK["Chạy server fallback"]
    DEGRADED["Degraded: manual/caption\nphrase + safety pack offline"]
    RESULT["Trả kết quả + ModelTrace"]

    START --> POLICY
    POLICY -- "adaptive/local_only" --> LOCAL_GATE
    POLICY -- "server_only" --> CONSENT
    LOCAL_GATE -- "Có" --> LOCAL
    LOCAL_GATE -- "Không + local_only" --> DEGRADED
    LOCAL_GATE -- "Không + adaptive" --> CONSENT
    LOCAL --> LOCAL_CONF
    LOCAL_CONF -- "Có" --> RESULT
    LOCAL_CONF -- "Không + local_only" --> DEGRADED
    LOCAL_CONF -- "Không + adaptive" --> CONSENT
    CONSENT -- "Không" --> DEGRADED
    CONSENT -- "Có" --> SERVER
    SERVER --> SERVER_CONF
    SERVER_CONF -- "Có" --> RESULT
    SERVER_CONF -- "Không" --> FALLBACK
    FALLBACK --> RESULT
    SERVER -- "Network/server lỗi" --> DEGRADED
    FALLBACK -- "Fallback lỗi" --> DEGRADED
    DEGRADED --> RESULT
~~~

Router áp dụng nguyên tắc **consent trước khi offload**. Dữ liệu thiết yếu là phần tối thiểu để hoàn thành lượt hiện tại: PCM của utterance, frame/capture hoặc text; không bao gồm contribution hay quyền dùng để huấn luyện. Nếu không consent, router chỉ được chạy local hoặc chuyển degraded.

- <code>AI_EXECUTION_MODE=adaptive</code> ưu tiên local khi model pack hợp lệ, device đủ capability và confidence đạt ngưỡng; nếu không thì mới xin consent và offload.
- <code>AI_EXECUTION_MODE=local_only</code> không gửi media/text tới inference server; local không đạt thì degraded.
- <code>AI_EXECUTION_MODE=server_only</code> bỏ qua local gate và dùng server sau consent. Chỉ profile demo hiện tại đặt giá trị này trong <code>.env</code>; không coi đây là ràng buộc kiến trúc.
- Model pack on-device phải signed/versioned, qua startup benchmark, memory/thermal/battery gate và có kill switch. Credential của server provider không bao giờ nằm trong APK.
- Mỗi kết quả ghi provider, model/version, latency, confidence, fallback reason và <code>execution_location=device|server</code>; không ghi raw media/prompt nhạy cảm.
- Price comparison và emergency decision vẫn là service/rule độc lập, không giao quyền quyết định cuối cho LLM/VLM.

### 5.7 Gemini Live là tài liệu tham khảo, không phải runtime

Các mẫu học:

- Full-screen call-like UI.
- Trạng thái camera/micro/Hold/End rõ ràng.
- Caption và transcript sau lượt nói.
- Local camera preview, chỉ gửi frame lấy mẫu.
- Phiên có state và hỗ trợ reconnect.

Không sao chép:

- Không gọi trải nghiệm là cuộc gọi người-người.
- Không mặc định micro luôn mở.
- Không để một multimodal LLM vừa dịch vừa tự quyết giá/khẩn cấp.
- Không gửi video 30 FPS tới AI.
