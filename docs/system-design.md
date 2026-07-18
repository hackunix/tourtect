# Thiết kế hệ thống — Trợ lý An toàn Du lịch Mobile Live AI

| Thuộc tính | Giá trị |
| --- | --- |
| Trạng thái | System Design v1 — hồ sơ hackathon |
| Cập nhật | 18/07/2026 |
| Sản phẩm chính | Mobile app iOS/Android |
| Kênh bổ trợ | Desktop/quick web, Zalo Mini App, Admin web |
| Địa bàn pilot | Một vài quận/huyện ở Hà Nội |
| Ngôn ngữ | Tiếng Việt và tiếng Hàn, Trung giản thể, Anh, Nga |
| Mục tiêu tài liệu | Chốt product scope, kiến trúc, dữ liệu, API, an toàn, quyền riêng tư và tiêu chí nghiệm thu |

> Đây là tài liệu thiết kế, không phải cam kết pháp lý hay kết luận rằng một cá nhân hoặc cơ sở kinh doanh đã lừa đảo. Hệ thống đánh giá giao dịch đang được hỏi, luôn hiển thị độ tin cậy và có quyền từ chối kết luận khi dữ liệu chưa đủ.

---

## 1. Tóm tắt điều hành

### 1.1 Khoảng trống cần giải quyết

Khách quốc tế thường ở thế bất lợi thông tin khi gặp:

- Giá taxi không hợp lý, đồng hồ hoặc lộ trình có dấu hiệu bất thường.
- Menu, bảng giá, hóa đơn hoặc phụ phí khó hiểu.
- Đổi tiền với tỷ giá thực nhận thấp hoặc phí không được nói trước.
- Ghost tour, ép thanh toán và các mô-típ lừa đảo phổ biến.
- Sự cố cần trao đổi ngay với người Việt, khách sạn, cơ quan du lịch hoặc lực lượng khẩn cấp.

Các công cụ dịch thông thường không biết mặt bằng giá địa phương. Công cụ kiểm tra giá theo ảnh lại không hiểu diễn biến hội thoại. Hệ thống này kết hợp hai luồng nhưng cố ý tách quyết định:

1. AI hỗ trợ nhìn, nghe, dịch và trích dữ kiện.
2. Price Engine và Safety Engine quyết định dựa trên dữ liệu có phiên bản và luật đã duyệt.

### 1.2 Đề xuất giá trị

Mobile app cung cấp hai trải nghiệm giống một cuộc gọi với AI:

- **Live Voice:** hai bên dùng push-to-talk; ứng dụng dịch hai chiều và âm thầm nhận biết món/dịch vụ, số tiền, tiền tệ hoặc dấu hiệu nguy hiểm.
- **Live Camera:** ứng dụng nhìn menu, món ăn, đồ vật, bảng giá hoặc đồng hồ taxi; kết hợp VLM để hiểu bối cảnh với ảnh tĩnh độ phân giải cao để xác nhận OCR và giá.

Khi nhận thấy một mức giá đáng chú ý, ứng dụng hiển thị **thẻ kín và rung nhẹ** cho khách. Kết quả không tự phát qua loa trước người bán. Khi có nguy hiểm, Safety Engine ưu tiên hướng dẫn rời khỏi tình huống và mở dialer thật; AI không tự gọi hoặc tự gửi vị trí.

### 1.3 Liên hệ tiêu chí chấm

| Tiêu chí | Bằng chứng trong thiết kế |
| --- | --- |
| Độ chính xác cảnh báo giá | Bốn mức kết quả, ngưỡng đỏ nghiêm ngặt, xác nhận OCR, so sánh theo ngành và cho phép “không đủ dữ liệu” |
| Tỷ lệ báo động giả | Red alert precision mục tiêu ≥ 95%, false-positive rate ≤ 2%, calibration theo vùng/ngành |
| Dịch và xử lý khẩn cấp | PTT theo vai, glossary dữ kiện quan trọng, rule engine đứng trước LLM, safety pack offline |
| Cập nhật dữ liệu vùng | Evidence ledger, quarantine, human review, snapshot có phiên bản và chống poisoning |
| Quyền riêng tư | Không cần tài khoản, consent theo mục đích, media mặc định chỉ tồn tại trong bộ nhớ, contribution opt-in |

### 1.4 Nguyên tắc sản phẩm

1. **Safety before intelligence:** xử lý dấu hiệu nguy hiểm trước khi bàn về giá.
2. **Translation must not wait for price reasoning:** luồng dịch và luồng intelligence chạy song song.
3. **No evidence, no accusation:** không gắn nhãn người bán; không cảnh báo đỏ khi chưa xác nhận dữ kiện.
4. **Quiet by default:** price insight riêng tư, không kích thích đối đầu.
5. **Consent by action:** xin camera, micro, vị trí và contribution đúng thời điểm sử dụng.
6. **Graceful abstention:** dữ liệu yếu phải trả “không đủ dữ liệu”.
7. **Emergency remains useful offline:** hotline, phrasebook và incident card không phụ thuộc AI online.

---

## 2. Phạm vi sản phẩm

### 2.1 Người dùng chính

| Persona | Nhu cầu |
| --- | --- |
| Khách Hàn/Trung/Anh/Nga | Dịch tại chỗ, hiểu giá, tránh đối đầu, biết gọi ai |
| Người Việt đồng hành | Giúp khách kiểm tra nhanh qua Zalo, giải thích kết quả và mở hotline |
| Nhân viên khách sạn/hướng dẫn viên | Tra playbook, hỗ trợ tạo incident card, xác minh thông tin địa phương |
| Data reviewer | Duyệt bằng chứng giá, scam pattern, hotline và snapshot |
| Data steward/admin | Quản trị nguồn, quyền truy cập, phiên bản, rollback và audit |

### 2.2 Phân chia theo kênh

| Kênh | Có | Không có |
| --- | --- | --- |
| Mobile app | Live Voice, Live Camera, price check, scam assistant, SOS, dịch tại chỗ, safety pack | Phiên dịch bên trong cuộc gọi PSTN |
| Desktop/quick web | Upload/chụp một ảnh, nhập mô tả, xem khoảng giá, playbook, hotline | Realtime call, continuous camera/audio |
| Zalo Mini App | Chụp/chọn ảnh, Camera Assist theo snapshot, text scam report, incident card, hotline qua openPhone | Full-duplex audio/video, WebRTC, AI nghe cuộc gọi |
| Admin web | Review queue, source registry, dataset/scam/hotline versioning, metrics, rollback | Xem raw media khi không có quyền và mục đích hợp lệ |

> Quyết định pivot đã được khóa: yêu cầu voice/camera realtime ban đầu được chuyển sang mobile app; website chỉ còn vai trò tra nhanh. Đây là thay đổi phạm vi có chủ đích, không phải giới hạn chưa được giải quyết.

### 2.3 Trong phạm vi pilot

- Nơi thử nghiệm: Một vài quận/huyện ở Hà Nội
- Bốn nhóm giá: taxi/đưa đón, đổi tiền, ăn uống, tour.
- UI và dịch hai chiều giữa tiếng Việt với:
  - <code>ko-KR</code>
  - <code>zh-Hans</code>
  - <code>en</code>
  - <code>ru-RU</code>
- Scam pattern seed: taxi, đổi tiền, ghost tour, ép thanh toán và tình huống có nguy cơ leo thang.
- Hotline khẩn cấp toàn quốc và hotline du lịch địa phương đã được xác minh.

### 2.4 Ngoài phạm vi

- Cầu nối PSTN/SIP/WebRTC ba bên.
- Nhân viên phiên dịch trực 24/7.
- Full AI offline trên mọi thiết bị.
- Live Voice hoặc Live Camera trên web/Zalo Mini App.
- Tự gọi, tự gửi vị trí, tự gửi incident card hoặc tự liên hệ trusted contact.
- Chấm điểm hoặc công khai danh sách người bán bị tố cáo.
- Dùng dữ liệu hóa đơn thuế hoặc dữ liệu riêng của đối tác khi chưa có quyền hợp pháp.

---

## 3. Trải nghiệm người dùng

### 3.1 Kiến trúc thông tin mobile

Màn hình Home giữ ba hành động lớn, thao tác được bằng một tay:

1. **Live Voice**
2. **Live Camera**
3. **SOS / Emergency**

Các hành động phụ:

- Kiểm tra một ảnh.
- Mô tả tình huống bằng text.
- Xem hotline và phrasebook offline.
- Xem/xóa lịch sử cục bộ.
- Quản lý ngôn ngữ, quyền và đóng góp dữ liệu.

### 3.2 Wireflow tổng thể

~~~mermaid
flowchart TD
    HOME["Home mobile"]
    VOICE["Live Voice"]
    CAMERA["Live Camera"]
    SOS["SOS / Emergency"]
    PHOTO["Kiểm tra một ảnh"]
    TEXT["Mô tả bằng text"]

    HOME --> VOICE
    HOME --> CAMERA
    HOME --> SOS
    HOME --> PHOTO
    HOME --> TEXT

    VOICE --> CONSENT_MIC{"Đã đồng ý dùng micro?"}
    CONSENT_MIC -- "Chưa" --> ASK_MIC["Giải thích ngắn và xin quyền"]
    ASK_MIC --> SESSION_V["Phiên PTT hai vai"]
    CONSENT_MIC -- "Rồi" --> SESSION_V
    SESSION_V --> PRIVATE["Thẻ giá kín / rung"]
    SESSION_V --> ESCALATE["Safety escalation"]

    CAMERA --> CONSENT_CAM{"Đã đồng ý dùng camera?"}
    CONSENT_CAM -- "Chưa" --> ASK_CAM["Giải thích frame sampling và xin quyền"]
    ASK_CAM --> SESSION_C["Preview + PTT + vision"]
    CONSENT_CAM -- "Rồi" --> SESSION_C
    SESSION_C --> CAPTURE["Bấm chụp để xác nhận OCR"]
    CAPTURE --> RESULT["Kết quả giá có confidence"]

    SOS --> TRIAGE["Chọn loại sự cố"]
    TRIAGE --> DIAL["Mở dialer tối đa 2 thao tác"]
    TRIAGE --> CARD["Incident card song ngữ"]
    TRIAGE --> PHRASE["Phrasebook / silent mode"]
~~~

### 3.3 Live Voice

#### Giao diện

- Header: ngôn ngữ du khách, trạng thái mạng, thời lượng, Hold và End.
- Khu transcript: câu gốc và bản dịch, tách theo vai.
- Hai nút lớn:
  - “Du khách nói”
  - “Người Việt nói”
- Nút loa/tai nghe, phát lại, báo dịch sai.
- Price insight drawer mặc định thu gọn.
- Nút SOS luôn thấy, không chỉ biểu đạt bằng màu đỏ.

#### Luồng một lượt nói

1. Người dùng giữ nút đúng vai.
2. App phát tín hiệu rung ngắn, bắt đầu lấy PCM.
3. PCM được xử lý VAD/noise reduction và stream dần.
4. Khi nhả nút, app gửi <code>ptt.ended</code>.
5. ASR trả transcript có confidence.
6. Translation lane tạo bản dịch và trả ngay cho client.
7. TTS hệ điều hành phát bản dịch theo audio route hiện tại.
8. Intelligence lane độc lập trích giá, món/dịch vụ, scam signal và critical safety signal.
9. Nếu có price candidate đủ rõ, Price Engine trả insight kín.
10. Nếu có red flag, Safety Engine ưu tiên escalation card nhưng không làm mất transcript/dịch.

#### Dual-lane sequence

~~~mermaid
sequenceDiagram
    actor Tourist as Du khách
    participant App as Mobile App
    participant RT as Realtime Gateway
    participant ASR as ASR Router
    participant MT as Translation Service
    participant Intel as Intelligence Extractor
    participant Price as Price Engine
    participant Safety as Safety Engine

    Tourist->>App: Giữ nút PTT và nói
    App->>RT: ptt.started + PCM chunks
    Tourist->>App: Nhả nút
    App->>RT: ptt.ended
    RT->>ASR: Finalize utterance
    ASR-->>RT: transcript.final + confidence

    par Lane dịch
        RT->>MT: source, target, transcript
        MT-->>RT: translation.ready
        RT-->>App: Text dịch
        App-->>Tourist: TTS hệ điều hành + caption
    and Lane intelligence
        RT->>Intel: Transcript đã chuẩn hóa
        Intel->>Safety: Critical/scam signals
        Safety-->>RT: safe action hoặc no-op
        Intel->>Price: PriceCandidate nếu có
        Price-->>RT: PriceInsight
        RT-->>App: Thẻ kín + haptic
    end
~~~

#### Quy tắc âm thanh

- PTT-only, không mở micro liên tục.
- Mỗi utterance chỉ có một <code>SpeakerRole</code>.
- Không tự phát price insight qua loa.
- Nếu TTS locale không tồn tại:
  1. Hiển thị caption chữ lớn.
  2. Cho phép người dùng đưa màn hình cho bên kia đọc.
  3. Dùng audio phrase đã đóng gói cho câu khẩn cấp đã duyệt.
- Khi dùng tai nghe, người dùng có thể bấm “Đọc riêng” để nghe giải thích giá.

### 3.4 Live Camera

#### Cách hoạt động

- Preview camera chạy cục bộ ở frame rate tự nhiên.
- Client chỉ sample frame tối đa 1 FPS; mặc định 0.5 FPS.
- Context frame có cạnh dài tối đa khoảng 768 px và chỉ phục vụ hiểu bối cảnh.
- Qwen3-VL xác định loại đối tượng: menu, món ăn, đồng hồ taxi, hóa đơn, bảng giá, quầy đổi tiền hoặc không liên quan.
- Khi phát hiện vùng chữ/số ổn định, app đề nghị “Giữ máy yên và chụp để xác nhận”.
- Chỉ ảnh do người dùng bấm chụp mới đi vào pipeline high-resolution.
- PP-OCRv5 chạy trên thiết bị trước; user xác nhận item, giá, tiền tệ và đơn vị.
- Qwen3-VL fallback chỉ hỗ trợ trích xuất/định vị; không được tự quyết cảnh báo.

#### Camera sequence

~~~mermaid
sequenceDiagram
    actor User as Người dùng
    participant App as Mobile App
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

### 3.5 Scam Assistant và safety triage

Thứ tự xử lý:

1. Rule kiểm tra red flag xác định.
2. Trích dữ kiện có schema.
3. Đối chiếu scam pattern đã duyệt.
4. Hỏi tối đa 1–3 câu phân luồng nếu cần.
5. Trả playbook ngắn và tùy chọn escalation.

Mức ưu tiên:

| Mức | Ví dụ | Hành vi |
| --- | --- | --- |
| Critical | Vũ khí, thương tích, bị giữ/nhốt, không thể rời đi | Hiện SOS ngay, hành động an toàn trước |
| Urgent | Ép thanh toán, tài xế không cho xuống, tình hình leo thang | Hướng dẫn rời nơi an toàn, bảo toàn bằng chứng sau |
| Non-emergency | Nghi bị tính cao, ghost tour sau khi đã an toàn | Kiểm tra giá/pattern, hướng dẫn khiếu nại |
| Information | Hỏi cách phòng tránh | Checklist và phrasebook |

Không dùng câu “đây chắc chắn là lừa đảo”. Cách diễn đạt chuẩn: “Tình huống này có một số dấu hiệu giống mẫu rủi ro …”.

### 3.6 Emergency

- Nút SOS không yêu cầu đăng nhập.
- Tối đa hai thao tác từ SOS tới dialer.
- Các lựa chọn ban đầu:
  - Tôi đang bị đe dọa/không thể rời đi.
  - Có người bị thương.
  - Cháy, tai nạn hoặc cần cứu hộ.
  - Tôi an toàn nhưng muốn báo gian lận.
  - Chế độ im lặng.
- V1 dùng:
  - <code>112</code>: tổng đài khẩn cấp quốc gia.
  - <code>113</code>: công an.
  - <code>114</code>: chữa cháy/cứu nạn.
  - <code>115</code>: cấp cứu y tế.
- Không dùng số khẩn cấp cho tranh chấp giá thông thường.
- Hotline địa phương là dữ liệu động, không hard-code ngoài safety pack đã ký.
- Incident card song ngữ tách:
  - Dữ kiện người dùng đã xác nhận.
  - Nhận định do AI đề xuất.
- Không ghi âm cuộc gọi, không nghe audio của cuộc gọi PSTN và không nói rằng app làm được việc này.

### 3.7 Desktop/quick web

Website là công cụ “mì ăn liền”:

- Không yêu cầu cài app hoặc tài khoản.
- Upload/chụp một ảnh.
- Nhập item/giá thủ công khi OCR lỗi.
- Nhập mô tả scam bằng text.
- Xem khoảng giá, confidence, freshness và hotline.
- Hiện QR/đường dẫn chính thức tới mobile app.
- Không xin micro, không chạy Live Voice, không chạy Live Camera.

### 3.8 Zalo Mini App

Đối tượng ưu tiên là người Việt đồng hành, khách sạn, hướng dẫn viên; khách quốc tế đã có Zalo vẫn dùng được.

V1:

- Camera preview và “Bấm để phân tích”; không gửi stream liên tục.
- Chụp/chọn một ảnh menu/hóa đơn/bảng giá.
- Text scam assistant.
- Hotline bằng <code>openPhone</code>.
- Incident card và phrasebook.
- Chọn thành phố thủ công; GPS là tùy chọn đúng ngữ cảnh.
- Dùng Zalo ID giả danh cho session, không xin profile hoặc số điện thoại khi onboarding.

Ràng buộc:

- Camera/micro dừng khi Mini App xuống background.
- <code>openPhone</code> chỉ mở màn hình gọi native.
- Không coi WebSocket/WebRTC trong WebView là hợp đồng nền tảng nếu tài liệu Zalo chưa cam kết.
- Không dùng <code>localStorage</code>, <code>sessionStorage</code> hoặc cookie; token gửi qua Authorization header và cache qua Zalo storage API.

---

## 4. Yêu cầu chức năng

### 4.1 Price Check

| ID | Yêu cầu |
| --- | --- |
| PC-01 | Nhận ảnh menu, hóa đơn, bảng giá hoặc đồng hồ taxi |
| PC-02 | Trích item, amount, currency, unit và các thuộc tính so sánh |
| PC-03 | Cho người dùng sửa mọi trường OCR có confidence thấp |
| PC-04 | Chọn reference snapshot theo vùng, thời điểm và vertical |
| PC-05 | Trả reference range, alert level, confidence, freshness và giải thích |
| PC-06 | Không gắn nhãn người bán hoặc công khai bằng chứng thô |
| PC-07 | Không dùng submission hiện tại để thay đổi snapshot đang so sánh |
| PC-08 | Cho phép feedback và contribution bằng consent riêng |

### 4.2 Live Translation

| ID | Yêu cầu |
| --- | --- |
| LT-01 | Hai hướng dịch được khóa bởi SpeakerRole |
| LT-02 | PCM chỉ được capture trong lúc giữ PTT |
| LT-03 | Translation lane không chờ Price Engine |
| LT-04 | TTS và caption đều có nút phát lại/báo lỗi |
| LT-05 | Bảo toàn số, phủ định, đơn vị, tiền tệ và thực thể khẩn cấp |
| LT-06 | Có fallback caption và phrase audio khi TTS thiếu locale |

### 4.3 Live Vision

| ID | Yêu cầu |
| --- | --- |
| LV-01 | Context frame không vượt 1 FPS |
| LV-02 | Tự giảm FPS khi mạng yếu, pin yếu hoặc thiết bị nóng |
| LV-03 | High-resolution capture luôn cần thao tác người dùng |
| LV-04 | Cảnh báo đỏ cần OCR/user confirmation và Price Engine |
| LV-05 | Dừng capture ngay khi background, revoke permission hoặc end session |

### 4.4 Scam và Emergency

| ID | Yêu cầu |
| --- | --- |
| SE-01 | Rule engine phát hiện critical signal trước LLM |
| SE-02 | Chỉ dùng playbook đã duyệt, có nguồn và ngày review |
| SE-03 | SOS dùng được offline và không cần tài khoản |
| SE-04 | Hotline không được LLM tạo hoặc sửa |
| SE-05 | Không tự gọi/chia sẻ vị trí/incident |
| SE-06 | Chế độ im lặng tắt TTS, rung và flash không cần thiết |

---

## 5. Kiến trúc hệ thống

### 5.1 System context

~~~mermaid
flowchart LR
    Tourist["Du khách"]
    Helper["Người Việt đồng hành"]
    Reviewer["Reviewer / Data steward"]
    Authority["Nguồn chính thức / khảo sát"]

    Mobile["Mobile App\nLive Voice + Live Camera"]
    Web["Desktop / Quick Web"]
    Zalo["Zalo Mini App Lite"]
    Admin["Admin Web"]
    Platform["Tourist Safety Platform"]

    Tourist --> Mobile
    Tourist --> Web
    Helper --> Zalo
    Reviewer --> Admin

    Mobile --> Platform
    Web --> Platform
    Zalo --> Platform
    Admin --> Platform
    Authority --> Platform
~~~

### 5.2 Container architecture

~~~mermaid
flowchart TB
    subgraph Clients["Client applications"]
        Mobile["React Native + Expo Development Build"]
        QuickWeb["Next.js Quick Web"]
        ZMA["Zalo Mini App + ZaUI/ZMP SDK"]
        Admin["Next.js Admin Web"]
    end

    subgraph Edge["Edge/API"]
        Gateway["API Gateway / WAF"]
        Realtime["Realtime Gateway\nWebSocket + session state"]
        Media["Signed Capture API"]
    end

    subgraph Core["FastAPI modular monolith"]
        Session["Session Orchestrator"]
        Router["Adaptive Model Router"]
        Translation["Translation Service"]
        Vision["Vision Service"]
        Price["Price Intelligence Engine"]
        Scam["Scam / Safety Engine"]
        Emergency["Emergency Directory"]
        Consent["Consent / Privacy Service"]
        Contribution["Contribution / Review API"]
    end

    subgraph Inference["Self-hosted inference"]
        ASR["Qwen3-ASR"]
        MT["MADLAD / Qwen3"]
        VLM["Qwen3-VL"]
        Extractor["Qwen3 constrained extraction"]
    end

    subgraph Data["Data platform"]
        PG["PostgreSQL + PostGIS + pgvector"]
        Redis["Redis"]
        Object["MinIO / encrypted object storage"]
        Queue["Worker queue"]
        Snapshot["Versioned datasets + signed safety packs"]
    end

    Mobile --> Gateway
    QuickWeb --> Gateway
    ZMA --> Gateway
    Admin --> Gateway
    Mobile --> Realtime
    Mobile --> Media

    Gateway --> Session
    Gateway --> Price
    Gateway --> Scam
    Gateway --> Emergency
    Gateway --> Consent
    Gateway --> Contribution
    Realtime --> Session
    Session --> Router
    Router --> Translation
    Router --> Vision
    Router --> Price
    Router --> Scam
    Translation --> ASR
    Translation --> MT
    Vision --> VLM
    Scam --> Extractor

    Session --> Redis
    Price --> PG
    Scam --> PG
    Emergency --> PG
    Consent --> PG
    Contribution --> PG
    Media --> Object
    Contribution --> Queue
    Queue --> Snapshot
    Snapshot --> PG
~~~

### 5.3 Lý do chọn modular monolith

Pilot cần phát triển nhanh và giữ transaction/provenance dễ kiểm soát. Core API là modular monolith với boundary rõ:

- Realtime và inference được tách process vì tải GPU/WebSocket khác HTTP CRUD.
- Data worker được tách vì xử lý bất đồng bộ và cần retry/quarantine.
- Các module dùng contract nội bộ, có thể tách service sau khi có số liệu tải thật.
- Không dùng microservice cho từng model trong giai đoạn hackathon.

### 5.4 Trách nhiệm thành phần

| Thành phần | Trách nhiệm | Không được làm |
| --- | --- | --- |
| Realtime Gateway | Session, PCM/frame transport, event ordering, backpressure | Ra quyết định giá |
| Session Orchestrator | State machine, locale, role, mode, resume | Lưu raw media dài hạn |
| Adaptive Model Router | Chọn device/server/model fallback theo capability/confidence | Thay đổi consent |
| Translation Service | ASR, MT, glossary, critical token validation | Tư vấn an toàn tự do |
| Vision Service | Context classification, ROI, OCR fallback | Kết luận gian lận |
| Price Engine | Chuẩn hóa, chọn cohort, anomaly/confidence | Dùng LLM làm phép so sánh cuối |
| Scam/Safety Engine | Triage, pattern retrieval, safe action, escalation | Tạo hotline |
| Emergency Directory | Trả số đã xác minh theo vùng/incident | Tự gọi số |
| Consent Service | Consent ledger, TTL, deletion/export | Gộp consent contribution với xử lý phiên |
| Data Platform | Evidence, review, snapshot, rollback | Cho submission tác động ngay production |

### 5.5 Stack model

| Nhiệm vụ | Model chính | Fallback | Nơi chạy pilot |
| --- | --- | --- | --- |
| Menu/hóa đơn/bảng giá | PP-OCRv5 mobile | Qwen3-VL-4B-Instruct | Device, fallback server |
| Món ăn/đồ vật/bối cảnh | Qwen3-VL-4B-Instruct | Qwen3-VL-8B | Server |
| ASR | Qwen3-ASR-0.6B | Qwen3-ASR-1.7B | Device đủ chuẩn hoặc server |
| Dịch trực tiếp | MADLAD-400-3B | Qwen3-4B-Instruct | Server |
| Trích dữ kiện/playbook | Qwen3-4B-Instruct | Qwen3-8B-Instruct | Server |
| TTS | iOS/Android system TTS | Phrase audio đã duyệt | Device |
| Giá | Rule + robust statistics + LightGBM/CatBoost | Không LLM | Server |
| Emergency | Rule + dữ liệu đã duyệt | Không LLM tự quyết | Device/server |

Tất cả model được bọc bằng interface:

- <code>ASRProvider</code>
- <code>TranslationProvider</code>
- <code>VisionProvider</code>
- <code>ExtractorProvider</code>
- <code>PriceModelProvider</code>

Không gọi trực tiếp model từ controller hoặc UI.

### 5.6 Adaptive hybrid routing

~~~mermaid
flowchart TD
    START["Bắt đầu utterance/capture"]
    CONSENT{"Đã chấp nhận chia sẻ\ndữ liệu thiết yếu cho server?"}
    PRIVATE_GATE{"Local sẵn sàng?\nModel hợp lệ, đáp ứng tác vụ\nvà thiết bị chịu được tải"}
    HYBRID_GATE{"Local sẵn sàng?\nModel hợp lệ, đáp ứng tác vụ\nvà thiết bị chịu được tải"}
    PRIVATE_LOCAL["Chạy on-device\n(local-only)"]
    HYBRID_LOCAL["Chạy on-device\n(hybrid)"]
    PRIVATE_CONF{"Local confidence đạt ngưỡng?"}
    HYBRID_CONF{"Local confidence đạt ngưỡng?"}
    SERVER["Chạy server model chính"]
    SERVER_CONF{"Server confidence đạt ngưỡng?"}
    FALLBACK["Chạy server fallback"]
    DEGRADED["Degraded: caption/manual/offline pack"]
    RESULT["Trả kết quả + ModelTrace"]

    START --> CONSENT
    CONSENT -- "Không" --> PRIVATE_GATE
    CONSENT -- "Có" --> HYBRID_GATE
    PRIVATE_GATE -- "Có" --> PRIVATE_LOCAL
    PRIVATE_GATE -- "Không" --> DEGRADED
    PRIVATE_LOCAL --> PRIVATE_CONF
    PRIVATE_CONF -- "Có" --> RESULT
    PRIVATE_CONF -- "Không" --> DEGRADED
    HYBRID_GATE -- "Có" --> HYBRID_LOCAL
    HYBRID_GATE -- "Không" --> SERVER
    HYBRID_LOCAL --> HYBRID_CONF
    HYBRID_CONF -- "Có" --> RESULT
    HYBRID_CONF -- "Không" --> SERVER
    SERVER --> SERVER_CONF
    SERVER_CONF -- "Có" --> RESULT
    SERVER_CONF -- "Không" --> FALLBACK
    FALLBACK --> RESULT
    SERVER -- "Network/server lỗi" --> DEGRADED
    FALLBACK -- "Fallback lỗi" --> DEGRADED
    DEGRADED --> RESULT
~~~

Router áp dụng nguyên tắc **consent trước khi offload**. Dữ liệu thiết yếu là phần dữ liệu tối thiểu phải gửi tới server để hoàn thành tác vụ (ví dụ: đoạn âm thanh, frame ảnh hoặc văn bản của lượt hiện tại), không bao gồm contribution hay dữ liệu dùng để huấn luyện. Nếu người dùng không chấp nhận điều khoản chia sẻ này, router tuyệt đối không gửi dữ liệu lên server:

- Chỉ chọn local khi model đã tải, chữ ký hợp lệ, đáp ứng được tác vụ và thiết bị vượt qua capacity gate.
- Nếu local không đáp ứng tác vụ, thiết bị không đủ tải hoặc confidence không đạt ngưỡng, chuyển sang chế độ degraded/manual; không tự động đổi sang server.
- Nếu người dùng đã consent, router vẫn ưu tiên local khi đủ điều kiện và chỉ offload sang server khi local không khả dụng hoặc confidence không đạt ngưỡng.
- Trạng thái consent và kết quả từng gate phải được ghi trong <code>ModelTrace</code>, nhưng không ghi raw media.

Gate on-device ASR:

- Model pack checksum/signature hợp lệ.
- Startup benchmark đạt real-time factor mục tiêu.
- Trước mỗi lần chạy, capacity gate kiểm tra tải CPU/GPU/NPU hiện tại và khả năng cấp đủ tài nguyên cho model; không khởi chạy nếu không giữ được resource budget.
- Peak memory không gây memory pressure.
- Thermal state không ở mức serious/critical.
- Pin không dưới ngưỡng an toàn, trừ khi đang sạc.
- User không thu hồi quyền hoặc tắt model download.

PP-OCRv5 mobile vẫn chạy local trên đa số thiết bị. Qwen3-VL và MADLAD server-first trong pilot.

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

---

## 6. Nền tảng dữ liệu giá và scam

### 6.1 Ba lớp dữ liệu

~~~mermaid
flowchart LR
    Sources["Nguồn chính thức\nKhảo sát\nMerchant\nContribution"]
    Intake["Consent + provenance validation\nPII redaction + evidence hash"]
    Ledger["Redacted Evidence Ledger\nMetadata và bằng chứng đã khử định danh"]
    Quarantine["Candidate / Quarantine\nDedup + redaction + moderation"]
    Review["Human Review\n2 reviewer khi rủi ro cao"]
    Model["Robust estimation\nVertical-specific models"]
    Snapshot["Reference Snapshot\nVersioned + signed"]
    Runtime["Price Engine Runtime"]
    Feedback["Appeal / feedback"]

    Sources --> Intake
    Intake --> Ledger
    Ledger --> Quarantine
    Quarantine --> Review
    Review --> Model
    Model --> Snapshot
    Snapshot --> Runtime
    Runtime --> Feedback
    Feedback --> Quarantine
~~~

Submission của phiên hiện tại không bao giờ thay đổi snapshot đang dùng trong cùng phiên.

Raw media không đi vào ledger bất biến. Intake chỉ ghi ledger sau khi đã kiểm tra consent/provenance, khử định danh và tạo hash; raw media tiếp tục tuân theo TTL/xóa. Ledger là append-only ở cấp sự kiện, nhưng payload evidence nằm ngoài ledger: khi có yêu cầu xóa, object và source link bị xóa, ledger chỉ nhận một tombstone cùng non-reversible hash chứng minh thao tác xóa đã xảy ra.

### 6.2 Phân tầng nguồn

| Tầng | Nguồn | Quyền ảnh hưởng |
| --- | --- | --- |
| A | Biểu phí/trần giá/tỷ giá/cảnh báo có thẩm quyền | Luật xác định hoặc nguồn neo |
| B | Khảo sát thực địa, mystery shopper, trung tâm hỗ trợ, đối tác được kiểm toán | Nguồn chính của reference model |
| C | Menu/website/bảng giá/OTA/hãng vận tải công khai | Giá chào bán sau chuẩn hóa |
| D | Contribution có ảnh và consent | Quarantine, cần xác nhận chéo |
| E | Báo chí, forum, review/social | Chỉ phát hiện xu hướng/pattern mới |

Tầng E không trực tiếp thay đổi giá hoặc playbook.

### 6.3 Mô hình dữ liệu

#### PriceObservation

| Trường | Ý nghĩa |
| --- | --- |
| <code>observation_id</code> | UUID |
| <code>canonical_item_id</code> | Item/dịch vụ chuẩn hóa |
| <code>vertical</code> | taxi, exchange, food, tour |
| <code>amount</code>, <code>currency</code>, <code>unit</code> | Giá đã chuẩn hóa; amount lưu bằng database NUMERIC hoặc integer minor units, không dùng float |
| <code>attributes</code> | Thuộc tính so sánh theo vertical |
| <code>region_id</code>, <code>geo_precision</code> | Vùng, không mặc định lưu GPS chính xác |
| <code>observed_at</code> | Thời điểm giao dịch/niêm yết |
| <code>source_tier</code>, <code>source_ref</code> | Provenance |
| <code>evidence_hash</code> | Dedup, không phải URL công khai |
| <code>ocr_confidence</code> | Chất lượng trích xuất |
| <code>moderation_status</code> | pending, accepted, rejected, quarantined |
| <code>consent_grant_id</code> | Bắt buộc với contribution |
| <code>lineage</code> | Các bước transform/model version |

#### ReferencePriceSnapshot

| Trường | Ý nghĩa |
| --- | --- |
| <code>dataset_version</code> | Phiên bản phát hành bất biến |
| <code>canonical_item_id</code> | Nhóm so sánh |
| <code>currency</code>, <code>unit</code> | Đơn vị tiền và đơn vị so sánh |
| <code>cohort_attributes</code> | Thuộc tính bắt buộc tạo nhóm so sánh |
| <code>region_id</code> | Vùng có hiệu lực |
| <code>valid_from</code>, <code>valid_to</code> | Khoảng thời gian |
| <code>p10</code>, <code>p50</code>, <code>p90</code> | Khoảng robust |
| <code>effective_sample_size</code> | Cỡ mẫu sau weighting/dedup |
| <code>independent_source_count</code> | Số nguồn độc lập |
| <code>source_mix</code> | Phân bổ tầng nguồn |
| <code>freshness</code> | Độ mới |
| <code>confidence</code> | Calibrated confidence |
| <code>model_version</code> | Rule/statistical model |
| <code>normalization_version</code> | Phiên bản chuẩn hóa vertical/đơn vị |
| <code>threshold_config_version</code> | Phiên bản materiality và alert gate |
| <code>rule_provenance</code> | Luật/trần giá chính thức nếu có |
| <code>published_at</code>, <code>approved_by</code> | Governance |

#### ScamPattern

- <code>pattern_id</code>, vùng/jurisdiction.
- Tín hiệu đa ngôn ngữ.
- Giải thích lành tính có thể xảy ra.
- Câu hỏi phân biệt.
- Mức rủi ro và điều kiện escalation.
- “Làm ngay”, “Không nên làm”, “Giữ bằng chứng”, “Báo ở đâu”.
- Nguồn xác minh, ngày review tiếp theo, phiên bản dịch.

#### EmergencyService

- Vùng và loại sự cố.
- Số ngắn/số quốc tế.
- Giờ hoạt động.
- Ngôn ngữ đã xác minh.
- URL nguồn chính thức.
- <code>verified_at</code>, reviewer, status.
- Kênh fallback và last-known-good version.

### 6.4 Chuẩn hóa theo vertical

| Vertical | Thuộc tính bắt buộc |
| --- | --- |
| Taxi/đưa đón | Khoảng cách, thời gian, loại xe, giờ cao điểm, sân bay, cầu đường, phí chờ |
| Đổi tiền | Cặp tiền, timestamp, số đưa/nhận, phí công khai, effective rate |
| Ăn uống | Món/SKU, khẩu phần, số lượng, thuế, service charge, loại địa điểm |
| Tour | Thời lượng, private/group, ngôn ngữ, lịch trình, vé/xe/bữa ăn, mùa |

### 6.5 Thuật toán cảnh báo

Đầu vào bắt buộc:

- Item/dịch vụ đã chuẩn hóa.
- Giá, tiền tệ và đơn vị đã được xác nhận.
- Vùng và thời điểm.
- Các thuộc tính có ảnh hưởng lớn theo vertical.
- Reference snapshot còn hiệu lực.

Các thành phần confidence:

- Cỡ mẫu hiệu dụng.
- Số nguồn/người bán độc lập.
- Source mix.
- Freshness.
- Chất lượng OCR/user confirmation.
- Product/attribute match.
- Độ chi tiết địa lý.
- Độ ổn định của model.

Logic khởi tạo:

1. Nếu vi phạm một biểu phí/trần giá chính thức đã được xác minh, trả <code>high_risk</code> với rule provenance.
2. Nếu không đủ dữ liệu, không đủ thuộc tính hoặc confidence thấp, trả <code>insufficient_data</code>.
3. Nếu giá nhỏ hơn <code>p10</code>, trả <code>typical</code> kèm lý do “thấp hơn khoảng thường gặp”; hệ thống chỉ phát hiện overcharge, không coi giá thấp là scam.
4. Nếu giá nằm từ <code>p10</code> đến <code>p90</code>, trả <code>typical</code>.
5. Nếu giá vượt <code>p90</code> nhưng chưa thỏa red gate, trả <code>elevated</code>; UI phân biệt “cao nhẹ” với chênh lệch đã vượt materiality threshold.
6. Chỉ trả <code>high_risk</code> khi đồng thời:
   - Giá vượt upper reference bound.
   - Chênh lệch vượt materiality threshold theo ngành.
   - OCR/đơn vị/thuộc tính đã xác nhận.
   - Confidence ≥ ngưỡng được calibration để đạt precision gate.

Ngưỡng không hard-code chung cho bốn vertical. Chúng nằm trong versioned configuration, được chọn trên validation set và phát hành cùng model card.

### 6.6 Cold start không phụ thuộc đối tác

- Seed dữ liệu có provenance từ nguồn công khai chính thức.
- Khảo sát nhỏ phân tầng tại ba cụm pilot.
- Demo data synthetic phải gắn nhãn rõ, không trộn vào production snapshot.
- Taxi và đổi tiền được ưu tiên vì dễ chuẩn hóa hơn.
- Tour có quyền abstain cao vì gói khó so sánh.
- Không giả định có API từ cơ quan du lịch, cơ quan thuế hoặc đối tác.

### 6.7 Chống poisoning

- Rotating pseudonymous device token và rate limit.
- Perceptual hash ảnh, content fingerprint và kiểm tra trùng người bán/thời gian.
- Giới hạn trọng số tối đa của một nguồn/thiết bị.
- Robust estimator, leave-one-source-out và drift detection.
- Cụm submission đột ngột bị quarantine.
- Merchant-declared data có nhãn riêng.
- Audit log bất biến, review RBAC và rollback snapshot.
- Không thưởng contributor dựa trên số “giá cao” tìm được.
- Khiếu nại của merchant/du khách đi qua review, không sửa dữ liệu trực tiếp.

---

## 7. API và realtime protocol

### 7.1 Quy ước

- Prefix: <code>/v1</code>.
- JSON dùng <code>snake_case</code>.
- Thời gian ISO 8601 UTC.
- ID là UUID/opaque ID.
- Lỗi HTTP theo <code>application/problem+json</code>.
- Mọi response AI có <code>trace_id</code>, model/dataset version và confidence nếu áp dụng.
- Tourist session không yêu cầu account; dùng short-lived session token.
- Admin dùng SSO/RBAC và audit log.
- Raw media upload dùng signed URL, content type allowlist và size limit.

### 7.2 Endpoint

| Method | Path | Mục đích |
| --- | --- | --- |
| POST | <code>/v1/realtime/sessions</code> | Tạo Live Voice/Camera session |
| POST | <code>/v1/realtime/sessions/{id}/resume</code> | Cấp token resume ngắn hạn |
| WS | <code>/v1/realtime/sessions/{id}/events</code> | PCM, context frame và realtime event |
| POST | <code>/v1/realtime/sessions/{id}/captures</code> | Tạo capture_id và signed PUT URL |
| PUT | <code>{signed_upload_url}</code> | Upload trực tiếp ảnh đã redaction vào object storage |
| POST | <code>/v1/realtime/sessions/{id}/captures/{capture_id}/finalize</code> | Khóa object, kiểm tra hash/redaction và chuyển capture sang ready |
| DELETE | <code>/v1/realtime/sessions/{id}</code> | End và xóa state/media tạm |
| POST | <code>/v1/price-checks</code> | Tạo price check từ capture_id đã finalize hoặc manual input |
| GET | <code>/v1/price-checks/{id}</code> | Lấy trạng thái/kết quả |
| POST | <code>/v1/scam-assessments</code> | Đánh giá tình huống text/transcript |
| GET | <code>/v1/emergency-services</code> | Hotline theo vùng/incident |
| GET | <code>/v1/safety-packs/{region}</code> | Gói offline có version/signature |
| POST | <code>/v1/contributions</code> | Opt-in contribution |
| DELETE | <code>/v1/privacy/sessions/{id}</code> | Yêu cầu xóa dữ liệu phiên |

### 7.3 Public types

~~~typescript
type Locale = "vi-VN" | "ko-KR" | "zh-Hans" | "en" | "ru-RU";

type Channel = "mobile" | "quick_web" | "zalo_mini_app" | "admin_web";
type LiveMode = "voice" | "camera";
type SpeakerRole = "tourist" | "local";
type ExecutionLocation = "device" | "edge" | "server";
type Vertical = "taxi" | "exchange" | "food" | "tour";

type AlertLevel =
  | "typical"
  | "elevated"
  | "high_risk"
  | "insufficient_data";

type ConsentScope =
  | "process_microphone"
  | "process_camera"
  | "precise_location"
  | "share_incident"
  | "contribute_redacted_data";

interface RealtimeSessionConfig {
  channel: "mobile";
  mode: LiveMode;
  tourist_locale: Locale;
  local_locale: "vi-VN";
  region_id: string;
  execution_policy: "auto";
  enabled_scopes: ConsentScope[];
}

interface ModelTrace {
  provider: string;
  model: string;
  model_version: string;
  execution_location: ExecutionLocation;
  latency_ms: number;
  confidence?: number;
  fallback_reason?: string;
}

interface VisionObservation {
  observation_id: string;
  scene_type:
    | "menu"
    | "food"
    | "taxi_meter"
    | "receipt"
    | "price_board"
    | "exchange_counter"
    | "unknown";
  regions_of_interest: Array<{
    kind: "text" | "price" | "object";
    box: [number, number, number, number];
    confidence: number;
  }>;
  requires_capture: boolean;
  model_trace: ModelTrace;
}

interface Money {
  amount_minor: string;
  currency: string;
  exponent: number;
}

interface PriceCandidateBase {
  canonical_item_id?: string;
  raw_item: string;
  money: Money;
  unit: string;
  region_id: string;
  observed_at: string;
  extraction_confidence: number;
  user_confirmed: boolean;
  comparison_readiness: "ready" | "needs_confirmation" | "insufficient";
  missing_required_fields: string[];
}

interface TaxiAttributes {
  distance_m?: number;
  duration_s?: number;
  vehicle_class?: string;
  airport_trip?: boolean;
  tolls_included?: boolean;
  waiting_fee_included?: boolean;
}

interface ExchangeAttributes {
  source_currency: string;
  target_currency: string;
  quoted_rate?: string;
  fee?: Money;
  rate_timestamp: string;
}

interface FoodAttributes {
  portion_size?: string;
  quantity: number;
  venue_segment?: string;
  tax_included?: boolean;
  service_charge_included?: boolean;
}

interface TourAttributes {
  duration_minutes: number;
  group_type: "private" | "group";
  guide_language?: Locale;
  inclusions: string[];
  season?: string;
}

type PriceCandidate =
  | (PriceCandidateBase & {
      vertical: "taxi";
      attributes: TaxiAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "exchange";
      attributes: ExchangeAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "food";
      attributes: FoodAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "tour";
      attributes: TourAttributes;
    });

interface CaptureCreateResponse {
  capture_id: string;
  upload_url: string;
  required_headers: Record<string, string>;
  expires_at: string;
}

interface CaptureFinalizeRequest {
  object_etag: string;
  sha256: string;
  media_type: "image/jpeg" | "image/png";
  redaction_applied: boolean;
  client_ocr_model?: string;
}

type PriceCheckRequest =
  | {
      source: "capture";
      capture_id: string;
      candidate: PriceCandidate;
    }
  | {
      source: "manual";
      candidate: PriceCandidate;
    };

interface PriceInsight {
  alert_level: AlertLevel;
  observed: Money;
  reference?: {
    p10_minor: string;
    p50_minor: string;
    p90_minor: string;
    currency: string;
    exponent: number;
    unit: string;
  };
  confidence: number;
  comparison_scope: string;
  freshness: string;
  reasons: string[];
  possible_benign_explanations: string[];
  dataset_version: string;
  trace_id: string;
}

interface ScamAssessment {
  urgency: "critical" | "urgent" | "non_emergency" | "information";
  matched_pattern_ids: string[];
  confirmed_facts: string[];
  ai_inferences: string[];
  safe_actions: string[];
  do_not: string[];
  follow_up_questions: string[];
  escalation?: {
    incident_type: string;
    emergency_service_ids: string[];
  };
  confidence: number;
  playbook_version: string;
  trace_id: string;
}
~~~

Chỉ candidate có <code>comparison_readiness = "ready"</code> mới được Price Engine đánh giá. Candidate thiếu thuộc tính vẫn được nhận để UI hỏi lại, nhưng kết quả phải là <code>insufficient_data</code> cùng <code>missing_required_fields</code>; server luôn validate schema theo vertical, không tin type phía client.

### 7.4 WebSocket messages

Client control:

- <code>ptt.started</code>: role, sequence, audio format.
- Binary <code>MediaChunk</code> loại AUDIO_PCM trong thời gian giữ nút.
- <code>ptt.ended</code>: utterance_id, final media sequence, client timestamp.
- Binary <code>MediaChunk</code> loại CAMERA_JPEG cho context frame, không phải capture bằng chứng.
- <code>session.hold</code>, <code>session.resume</code>, <code>session.end</code>.
- <code>feedback.translation</code> và <code>feedback.price_insight</code>.

WebSocket dùng hai kiểu frame:

1. Text frame JSON cho control và server event.
2. Binary frame Protobuf cho media; không gửi PCM/JPEG trần và không base64 media trong JSON.

~~~protobuf
enum MediaKind {
  MEDIA_KIND_UNSPECIFIED = 0;
  AUDIO_PCM = 1;
  CAMERA_JPEG = 2;
}

message MediaChunk {
  uint32 protocol_version = 1;
  MediaKind kind = 2;
  string media_id = 3;
  optional string utterance_id = 4;
  uint64 sequence = 5;
  int64 captured_at_ms = 6;
  bool is_last = 7;
  bytes payload = 8;
}
~~~

Server trả <code>media.ack</code> chứa <code>media_id</code>, highest contiguous <code>sequence</code> và trạng thái accepted/duplicate/rejected. WebSocket giữ ordering trong một kết nối; sau reconnect, client chỉ gửi lại chunk chưa ACK trong buffer bộ nhớ. <code>kind</code> là discriminator bắt buộc để route PCM và JPEG.

Server event:

- <code>media.ack</code>
- <code>transcript.partial</code>
- <code>transcript.final</code>
- <code>translation.ready</code>
- <code>vision.observation</code>
- <code>price.candidate</code>
- <code>price.insight</code>
- <code>safety.escalation</code>
- <code>network.degraded</code>
- <code>session.expired</code>

Envelope:

~~~typescript
interface RealtimeEvent<T> {
  event_id: string;
  event_type: string;
  session_id: string;
  utterance_id?: string;
  sequence: number;
  occurred_at: string;
  payload: T;
  trace_id: string;
}
~~~

Yêu cầu protocol:

- Event idempotent theo <code>event_id</code>.
- Sequence của event và sequence của media là hai namespace riêng; media sequence tăng đơn điệu trong từng <code>media_id</code>.
- Client bỏ qua event có sequence cũ nhưng vẫn ACK.
- Backpressure làm giảm camera frame trước, không drop <code>ptt.ended</code>.
- Audio format được negotiate khi tạo session; pilot chuẩn hóa PCM mono 16-bit, 16 kHz.
- Reconnect dùng resume token; không resume camera/micro nếu app đang background.

### 7.5 Price check response mẫu

~~~json
{
  "check_id": "pc_01...",
  "status": "completed",
  "candidate": {
    "vertical": "food",
    "raw_item": "pho bo",
    "money": {
      "amount_minor": "180000",
      "currency": "VND",
      "exponent": 0
    },
    "unit": "bowl",
    "region_id": "hcm_d1",
    "observed_at": "2026-07-18T03:15:00Z",
    "attributes": {
      "quantity": 1,
      "portion_size": "regular",
      "venue_segment": "restaurant",
      "tax_included": true,
      "service_charge_included": false
    },
    "extraction_confidence": 0.94,
    "user_confirmed": true,
    "comparison_readiness": "ready",
    "missing_required_fields": []
  },
  "insight": {
    "alert_level": "elevated",
    "observed": {
      "amount_minor": "180000",
      "currency": "VND",
      "exponent": 0
    },
    "reference": {
      "p10_minor": "55000",
      "p50_minor": "75000",
      "p90_minor": "120000",
      "currency": "VND",
      "exponent": 0,
      "unit": "bowl"
    },
    "confidence": 0.78,
    "comparison_scope": "District 1, restaurant segment",
    "freshness": "2026-07-01",
    "reasons": [
      "Giá cao hơn khoảng thường gặp của nhóm so sánh"
    ],
    "possible_benign_explanations": [
      "Khẩu phần hoặc loại thịt có thể khác",
      "Có thể đã gồm thuế hoặc phí phục vụ"
    ],
    "dataset_version": "price-v2026.07.1",
    "trace_id": "tr_01..."
  }
}
~~~

---

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
| Mạng tốt | Full PTT | 0.5–1 FPS + capture | Realtime | Full directory |
| Mạng yếu | Audio ưu tiên | 0.2 FPS hoặc tắt | Manual/cached | Full offline pack |
| Server vision lỗi | Không ảnh hưởng dịch | OCR local + manual | Nếu đủ dữ kiện | Không ảnh hưởng |
| Server MT lỗi | Caption transcript gốc/phrase | Camera vẫn chạy | Vẫn kiểm tra được input manual | Phrasebook |
| Offline | Phrase/audio pack | Preview local, không AI vision | Cache/manual, không red alert nếu stale | Hotline + incident card |
| Thiết bị nóng/pin yếu | ASR server hoặc caption | Tắt frame sampling | Manual | Không ảnh hưởng |

Quy tắc quan trọng: offline/stale data không phát cảnh báo đỏ mới.

---

## 9. Quyền riêng tư, bảo mật và an toàn

### 9.1 Consent model

Consent không gộp:

| Scope | Thời điểm xin | Mặc định |
| --- | --- | --- |
| Microphone processing | Lần đầu mở Live Voice hoặc PTT trong Live Camera | Tắt |
| Camera processing | Lần đầu mở Live Camera | Tắt |
| Precise location | Khi cần chọn hotline/cohort chính xác hơn | Tắt; chọn vùng thủ công trước |
| Share incident | Khi người dùng bấm chia sẻ/xuất | Tắt |
| Contribute redacted data | Sau khi đã nhận kết quả, ở màn hình riêng | Tắt |

Thu hồi consent phải dừng capture ngay và không làm mất hotline/phrasebook.

### 9.2 Retention

| Dữ liệu | Mặc định | Khi opt-in contribution |
| --- | --- | --- |
| PCM/raw audio | Bộ nhớ tạm, hủy sau ASR | Không giữ raw audio |
| Low-resolution context frame | Bộ nhớ tạm, không lưu | Không dùng làm contribution |
| High-resolution capture | Xóa sau xử lý; lỗi tối đa 24 giờ | Redaction rồi giữ bằng chứng tối đa 30 ngày để audit |
| Transcript | Session memory; lịch sử local nếu user bật | Chỉ lưu text đã redaction cần thiết |
| Incident card | Local | Upload/share chỉ sau xác nhận |
| Normalized observation | Không tạo nếu không opt-in | Giữ tối đa 24 tháng để xây reference; sau đó xóa hoặc chỉ giữ aggregate không còn source/session link |
| Operational telemetry | 30 ngày, không raw content | Không thay đổi |
| Security audit | 90 ngày, metadata tối thiểu | Không raw content |

Contribution được liên kết bằng deletion token ngẫu nhiên, không bằng danh tính. Hash của token được lưu tách khỏi dataset để người dùng có thể rút contribution. Khi rút consent, evidence/source link và candidate bị xóa khỏi pipeline trong tối đa 30 ngày, snapshot tương lai được tái tính; snapshot đã công bố chỉ chứa aggregate không định danh và được supersede bằng phiên bản mới khi cần.

### 9.3 Threat model

| Rủi ro | Biện pháp |
| --- | --- |
| Prompt injection từ menu/transcript | Coi OCR/transcript là dữ liệu, không phải instruction; tool allowlist và schema validation |
| LLM tạo hotline | Hotline chỉ từ Emergency Directory đã ký |
| Data poisoning/Sybil | Quarantine, device token, dedup, source caps, robust estimator |
| URL/file độc hại | Signed upload, MIME sniffing, size limit, malware scan, image re-encode |
| Session hijack | Short-lived token, bind device/session, TLS, rate limit |
| Rò raw media trong log | Structured logging allowlist; cấm payload content |
| Admin lạm dụng | RBAC, least privilege, audit, dual approval cho snapshot/playbook |
| Model hallucination | Rule engine, constrained schema, confidence, provenance, abstention |
| Người bán thấy cảnh báo | Private card, haptic, headphone-only explain |
| Camera thu người ngoài | Explicit indicator, foreground-only, no continuous recording, consent reminder |

### 9.4 Safety wording

- “Trong khoảng thường gặp.”
- “Cao hơn mặt bằng của nhóm so sánh; hãy kiểm tra phụ phí.”
- “Rất bất thường so với dữ liệu hiện có; đây không phải kết luận pháp lý.”
- “Chưa đủ dữ liệu để đánh giá đáng tin.”

Không dùng:

- “Người bán này là lừa đảo.”
- “Tuyệt đối không trả tiền” khi người dùng đang bị giữ/đe dọa.
- “Hãy quay phim/đuổi theo/giằng lại tiền.”

### 9.5 Compliance baseline

- Thiết kế mapping consent, purpose limitation, deletion, export và incident handling theo Luật Bảo vệ dữ liệu cá nhân Việt Nam hiện hành.
- Cần legal review trước pilot công khai, đặc biệt với dữ liệu trẻ vị thành niên, sinh trắc học, location và chia sẻ cho cơ quan chức năng.
- Không tuyên bố “zero retention” nếu hạ tầng/model server chưa chứng minh được.

---

## 10. Quan sát, đánh giá và kiểm thử

### 10.1 Metrics

| Nhóm | Metrics |
| --- | --- |
| OCR/Vision | CER, field extraction F1, price/currency exact match, scene classification F1 |
| ASR | WER theo locale/accent/noise, critical entity exact match |
| Translation | Human adequacy, critical token preservation, latency, semantic parity |
| Price | Red precision, false-positive rate, recall severe overcharge, calibration, abstention |
| Data | Coverage, freshness, source diversity, effective sample size, drift |
| Safety | Critical escalation recall, unsafe-advice rate, playbook parity |
| Realtime | P50/P95 translation latency, disconnect, resume success, dropped frame/audio |
| Device | Battery drain, peak memory, thermal throttling, crash-free sessions |
| Privacy | TTL deletion success, consent errors, raw-content log violations |
| Emergency | Dialer success, hotline lookup availability, offline pack validity |

Metrics được slice theo:

- Ngôn ngữ.
- Thành phố/cụm.
- Vertical.
- Thiết bị/tier.
- On-device/server.
- Network condition.
- Dữ liệu mới/cũ.

### 10.2 Release gates

| Gate | Mục tiêu |
| --- | --- |
| PTT → audio dịch đầu tiên | P95 ≤ 2 giây trên 4G ổn định |
| Price insight từ lời nói | P95 ≤ 3 giây |
| Live Camera observation | P95 ≤ 3 giây |
| OCR + confirmed price result | P95 ≤ 5 giây |
| Camera sampling | Không vượt 1 FPS |
| Red alert precision | ≥ 95% |
| False-positive trên giá hợp lệ | ≤ 2% |
| Critical safety escalation | 100% trên golden safety set |
| Unsafe confrontation advice | 0 trường hợp trên safety set |
| Critical translation fields | Không sai phủ định, thương tích, số tiền, tiền tệ, vị trí, biển số, số người |
| Background capture | Dừng 100% khi background/end/revoke |
| Offline SOS | Mở hotline/incident card tối đa 2 thao tác |
| Poisoning simulation | Attack budget 1–500; p50/p90 dịch ≤ 5% và không lật quyết định high-risk |

Gate giá được đánh giá theo từng vertical × vùng, không chỉ toàn cục. Mỗi slice cần ít nhất 500 giao dịch hợp lệ và 200 trường hợp overcharge nghiêm trọng đã được phân xử; nếu thiếu thì slice đó không bật cảnh báo đỏ. Point estimate và one-sided Wilson 95% confidence bound đều phải đạt gate: lower bound của precision ≥ 95% và upper bound của false-positive rate ≤ 2%. Nếu chưa đạt, chỉ trả <code>elevated</code> hoặc <code>insufficient_data</code>.

Poisoning test chạy với ngân sách 1, 10, 50, 100 và 500 submission phối hợp; đo dịch chuyển <code>p50</code>, <code>p90</code>, confidence và tỷ lệ lật quyết định alert. Gate yêu cầu cả <code>p50</code> và <code>p90</code> không dịch quá 5%, đồng thời không có price case chuẩn nào bị đổi từ <code>typical</code> sang <code>high_risk</code> hoặc ngược lại.

### 10.3 Test matrix

#### Live Voice

- Khách Hàn nói số tiền, người Việt nói phụ phí hợp lệ.
- Khách Trung code-switch tên món tiếng Việt.
- Khách Nga nói trong tiếng ồn đường phố.
- Khách Anh hỏi lại và sửa transcript.
- Hai người bấm nhầm vai.
- TTS locale thiếu hoặc audio route đổi sang Bluetooth.
- ASR local confidence thấp và failover server.
- Mất mạng đúng lúc nhả PTT.

#### Live Camera

- Menu nhiều cột, giá thiếu ký hiệu tiền tệ.
- Hóa đơn mờ, phản sáng, nghiêng.
- Đồng hồ taxi có nhiều số.
- Món ăn nhìn giống nhau nhưng khẩu phần khác.
- Bảng đổi tiền mua/bán hai cột.
- Camera thấy thẻ ngân hàng/khuôn mặt.
- Network giảm từ 4G xuống 2G.
- App background, khóa màn hình hoặc thiết bị quá nhiệt.

#### Price

- Giá hợp lệ cao vì sân bay/phí cầu đường/service charge.
- OCR nhầm <code>15</code> thành <code>75</code>.
- VND và USD bị nhầm.
- Dữ liệu vùng thiếu, phải mở rộng cohort.
- Snapshot stale.
- Một merchant/source gửi nhiều submission.
- Tour private bị so nhầm tour group.

#### Scam/Emergency

- Khách bị giữ trong taxi.
- Tranh chấp giá nhưng khách đã an toàn.
- Ghost tour không có đe dọa.
- Người bị thương, không có vũ khí.
- Chế độ im lặng.
- GPS/micro bị từ chối.
- Roaming không gọi được số ngắn.
- Hotline địa phương hết giờ.

#### Zalo Mini App

- Camera permission được duyệt/từ chối trên Zalo Android và iOS thật.
- Mini App xuống background khi mở dialer.
- Storage đầy hoặc file tạm bị xóa.
- Không có cookie/localStorage.
- API CORS và Authorization header.
- User không cung cấp profile/số điện thoại.

### 10.4 Golden data

Ground truth không lấy trực tiếp từ nút thích/không thích của cộng đồng.

- Giá: khảo sát mù, mystery shopper, nguồn có thẩm quyền và chuyên gia phân xử.
- Scam: pattern do chuyên gia địa phương và ít nhất hai reviewer duyệt.
- Translation: bilingual reviewer; tập trung dữ kiện khẩn cấp và tiền.
- OCR: ảnh thực tế đã được phép sử dụng, đa ánh sáng/thiết bị.
- Mọi golden set có data card, consent/provenance và version.

---

## 11. Demo hackathon

### 11.1 Demo 1 — Live Voice

**Bối cảnh:** khách Hàn nói chuyện với người bán Việt Nam tại TP.HCM.

1. Chọn <code>ko-KR ↔ vi-VN</code>.
2. Khách giữ “Du khách nói”; app dịch và phát tiếng Việt.
3. Người bán giữ “Người Việt nói”, nêu giá.
4. App dịch ngay, đồng thời trích số tiền.
5. Price Engine trả <code>elevated</code>.
6. Khách nhận rung và thẻ kín; người bán không nghe cảnh báo.
7. App gợi ý câu hỏi lịch sự về phí đã bao gồm.

Điểm chứng minh: translation lane độc lập, private insight, không đối đầu.

### 11.2 Demo 2 — Live Camera

**Bối cảnh:** khách Trung hoặc Nga nhìn menu/đồng hồ taxi tại Đà Nẵng–Hội An.

1. Camera nhận biết scene.
2. App đánh dấu vùng chữ/số và yêu cầu bấm chụp.
3. PP-OCRv5 đọc ảnh high-resolution.
4. Người dùng xác nhận giá/đơn vị.
5. Price Engine so snapshot theo khu vực.
6. Kết quả hiện range, freshness, confidence và giải thích.

Điểm chứng minh: VLM hiểu bối cảnh nhưng không tự quyết; OCR/user confirmation giảm false alarm.

### 11.3 Demo 3 — Kênh lite

- Quick web upload cùng ảnh và nhận kết quả không cần tài khoản.
- Zalo Mini App chụp một ảnh, xem incident card và dùng <code>openPhone</code>.
- UI ghi rõ Live chỉ có trong mobile app.

### 11.4 Dữ liệu demo

- Dữ liệu thật có URL/provenance và ngày xác minh.
- Dữ liệu synthetic gắn badge “Demo data”.
- Synthetic data không được đưa vào snapshot pilot/production.
- Không gọi thử hotline khẩn cấp trong demo.

---

## 12. Lộ trình sau tài liệu

### Giai đoạn 0 — System Design

- Hoàn thiện tài liệu này.
- Review safety, privacy, data governance và model licenses.
- Chốt product name/design language sau, không ảnh hưởng API.

### Giai đoạn 1 — Technical spikes

1. Expo Development Build: PCM capture/playback, audio focus, system TTS.
2. Qwen3-ASR-0.6B quantized trên thiết bị đại diện.
3. WebSocket PTT và latency budget end-to-end.
4. Camera sampling + on-device PP-OCRv5.
5. Qwen3-VL context frame + high-resolution capture.
6. MADLAD benchmark cho bốn cặp ngôn ngữ.
7. Zalo CameraContext/openPhone trên thiết bị thật.

### Giai đoạn 2 — Hackathon prototype

- Một vertical slice Live Voice.
- Một vertical slice Live Camera.
- Seed reference snapshot bốn vertical, ba cụm.
- Quick web và Zalo click-through/quick check.
- Admin review queue tối thiểu.

### Giai đoạn 3 — Pilot 90 ngày

- Khảo sát thực địa, golden data và calibration.
- Reviewer workflow, signed safety pack và hotline verification.
- Closed beta với khách sạn/hướng dẫn viên.
- Privacy/security review và incident drill mô phỏng.

### Giai đoạn 4 — Mở rộng có điều kiện

- Partner data/API.
- Human interpreter hoặc VoIP bridge chỉ khi có SLA, legal review và consent flow.
- Thêm <code>zh-Hant</code>.
- Mở rộng địa lý sau khi từng vertical đạt release gate.

---

## 13. Rủi ro và phương án

| Rủi ro | Ảnh hưởng | Phương án |
| --- | --- | --- |
| MADLAD dịch sai ngữ cảnh khẩn cấp | Cao | Glossary, critical validator, phrase pack, benchmark và fallback Qwen3 |
| ASR 0.6B không realtime trên thiết bị phổ thông | Trung bình | Startup benchmark và server fallback |
| VLM đọc sai số | Cao | Không dùng VLM làm source of truth; high-res OCR + user confirmation |
| Dữ liệu giá mỏng | Cao | Abstain, mở rộng cohort có disclosure, khảo sát tập trung |
| Crowdsourcing bị đầu độc | Cao | Quarantine, source cap, dedup, review, rollback |
| Live làm nóng máy/tốn data | Trung bình | PTT, adaptive FPS, burst video, thermal/network policy |
| TTS không có đủ locale | Trung bình | Caption và audio phrase đã đóng gói |
| Người dùng hiểu “Live” là gọi người thật | Trung bình | Tên “Live Voice/Live Camera với AI”, onboarding rõ |
| Zalo runtime khác Android/iOS | Trung bình | Lite scope, device spike, không cam kết realtime |
| Hotline thay đổi | Cao | Registry có nguồn/ngày xác minh, signed pack, expiry |
| Cảnh báo làm xung đột leo thang | Cao | Private card, haptic, safe wording, no public TTS |

---

## 14. Giả định đã khóa

- Mobile app là sản phẩm đầy đủ; web và Zalo là kênh lite.
- Web dùng được trên desktop và mobile browser nhưng không có Live.
- “Call” nghĩa là phiên nói với AI trong app, không phải cuộc gọi PSTN.
- Live Voice chỉ dùng push-to-talk theo vai.
- Live Camera là AI camera assist, không kết nối video tới người khác.
- Stack runtime là mã nguồn mở; Gemini Live chỉ là mẫu tham khảo.
- Adaptive hybrid chỉ ưu tiên ASR/OCR on-device trong pilot; model 3B–8B server-first.
- TTS hệ điều hành là mặc định.
- Price insight luôn riêng tư bằng card/rung.
- Pilot không phụ thuộc đối tác.
- Contribution chỉ có hiệu lực khi user opt-in riêng.
- Chinese MVP là Simplified Chinese; Traditional Chinese nằm ngoài pilot.
- Không có human operator, VoIP bridge, auto-call hoặc auto-share location trong V1.

---

## 15. Nguồn tham khảo chính

### Realtime multimodal

- [Gemini Apps Help — Talk naturally with Gemini Live](https://support.google.com/gemini/answer/15274899?hl=en-GB)
- [Gemini Live API — Capabilities](https://ai.google.dev/gemini-api/docs/live-api/capabilities)
- [Gemini Live API — WebSocket guide](https://ai.google.dev/gemini-api/docs/live-api/get-started-websocket)

### Model mã nguồn mở

- [PaddleOCR — PP-OCRv5 multilingual recognition](https://www.paddleocr.ai/latest/en/version3.x/algorithm/PP-OCRv5/PP-OCRv5_multi_languages.html)
- [Qwen3-ASR official repository](https://github.com/QwenLM/Qwen3-ASR)
- [Qwen3-VL official repository](https://github.com/QwenLM/Qwen3-VL)
- [Qwen3 official repository](https://github.com/QwenLM/Qwen3)
- [Google MADLAD-400-3B-MT model card](https://huggingface.co/google/madlad400-3b-mt)

### Zalo Mini App

- [Zalo Mini App CameraContext](https://docs.zaloplatforms.com/docs/MA/api/media/camera/createCameraContext)
- [Zalo Mini App API reference](https://docs.zaloplatforms.com/docs/MA/api/intro)
- [Zalo Mini App permission guide](https://docs.zaloplatforms.com/docs/MA/intro/request-permission)
- [Zalo Mini App openPhone](https://docs.zaloplatforms.com/docs/MA/api/device/contact/openPhone)
- [Zalo Mini App call API guidance](https://docs.zaloplatforms.com/docs/MA/intro/best-practices/call-restful-api)

### Du lịch, khẩn cấp và quyền riêng tư

- [Tổng đài khẩn cấp quốc gia 112](https://en.mae.gov.vn/hotline-to-receive-disaster-and-emergency-reports-nationwide-9018.htm)
- [Miễn phí các cuộc gọi 113, 114, 115](https://mst.gov.vn/mien-phi-cac-cuoc-goi-den-so-113-114-115-197140697.htm)
- [Trung tâm Hỗ trợ Du khách Đà Nẵng](https://bana.danang.gov.vn/vi/web/dng/w/trung-t%C3%A2m-h%E1%BB%97-tr%E1%BB%A3-du-kh%C3%A1ch-%C4%90%C3%A0-n%E1%BA%B5ng)
- [Cục Du lịch Quốc gia — tiếp nhận phản ánh du khách](https://vietnamtourism.gov.vn/post/34576)
- [Cục Du lịch Quốc gia — cảnh báo scam du lịch](https://vietnamtourism.gov.vn/post/51087)
- [Luật Bảo vệ dữ liệu cá nhân 91/2025/QH15](https://vbpl.moj.gov.vn/bocongan/Pages/vbpq-thuoctinh.aspx?ItemID=179252&Keyword=)
