# Container Architecture và Android Client

> Tách từ `system-design.md` — mục 5.2.

### 5.2 Container architecture

~~~mermaid
flowchart TB
    subgraph Clients["Client applications"]
        Android["Native Android\nKotlin + Jetpack Compose\nCameraX + AudioRecord/Media3"]
        Web["Next.js Responsive Web / SEO"]
        ZMA["Zalo Mini App + ZaUI/ZMP SDK"]
        Admin["Next.js Admin / Moderation"]
    end

    subgraph Edge["Edge/API"]
        Gateway["API Gateway / WAF"]
        Realtime["Realtime Gateway\nWebSocket + session state"]
        Media["Signed Capture API"]
    end

    subgraph CommunityCore["Community core — modular monolith"]
        Identity["Identity / Profile"]
        Content["Post / Comment / Review"]
        Place["Place / Topic / Price Report"]
        Reputation["Reputation"]
        Moderation["Moderation / Appeal"]
        Feed["Feed / Organic Ranking"]
        Search["Search API"]
        Notify["Notification"]
    end

    subgraph SafetyAI["Safety & AI modules"]
        Session["Session Orchestrator"]
        Router["Adaptive Model Router"]
        Translation["Translation Service"]
        Vision["Vision Service"]
        Price["Price Intelligence Engine"]
        Scam["Scam / Safety Engine"]
        Emergency["Emergency Directory"]
        Consent["Consent / Privacy Service"]
    end

    subgraph ContentRevenue["Content & revenue boundary"]
        Ingest["External Content Ingestion"]
        Rights["Rights / Source Registry"]
        Scheduler["Adaptive Crawl Scheduler\nper-host budget + backoff"]
        Fetcher["Sandbox Fetch Proxy\nconditional GET + SSRF guard"]
        Ads["Ad Eligibility / Decision"]
        Billing["Subscription / Business / Entitlement"]
        Affiliate["Affiliate Disclosure / Events"]
    end

    subgraph DeviceInference["Android on-device — khi policy/capability cho phép"]
        LocalOCR["Signed OCR model pack"]
        LocalASR["Signed ASR model pack"]
        SystemTTS["Android system TTS"]
    end

    subgraph Inference["Server inference providers"]
        FPTAPI["OpenAI-compatible API\nmkp-api.fptcloud.com"]
        ASR["STT model\nchọn trong Marketplace"]
        MT["Qwen3 text / translation"]
        VLM["Qwen VL"]
        TTS["TTS model\nchọn trong Marketplace"]
        Extractor["Qwen3 constrained extraction"]
    end

    subgraph Data["Data platform"]
        PG["PostgreSQL + PostGIS + pgvector"]
        OS["OpenSearch"]
        Redis["Redis"]
        Object["MinIO / encrypted object storage"]
        Queue["Event / worker queue"]
        Snapshot["Versioned datasets + signed safety packs"]
    end

    Android --> Gateway
    Web --> Gateway
    ZMA --> Gateway
    Admin --> Gateway
    Android --> Realtime
    Android --> Media
    Android --> LocalOCR
    Android --> LocalASR
    Android --> SystemTTS

    Gateway --> Identity
    Gateway --> Content
    Gateway --> Place
    Gateway --> Feed
    Gateway --> Search
    Gateway --> Moderation
    Gateway --> Price
    Gateway --> Scam
    Gateway --> Emergency
    Gateway --> Consent
    Gateway --> Billing
    Realtime --> Session
    Session --> Router
    Router -. "execution policy" .-> Android
    Router --> Translation
    Router --> Vision
    Router --> Price
    Router --> Scam
    Translation --> FPTAPI
    Vision --> FPTAPI
    Scam --> FPTAPI
    FPTAPI --> ASR
    FPTAPI --> MT
    FPTAPI --> VLM
    FPTAPI --> TTS
    FPTAPI --> Extractor

    Content --> Moderation
    Content --> Reputation
    Content --> Queue
    Place --> PG
    Content --> PG
    Identity --> PG
    Reputation --> PG
    Moderation --> PG
    Feed --> Redis
    Search --> OS
    Queue --> OS
    Queue --> Notify
    Ingest --> Rights
    Rights --> Scheduler
    Scheduler --> Fetcher
    Fetcher --> Queue
    Ads --> Moderation
    Billing --> PG
    Affiliate --> PG
    Session --> Redis
    Price --> PG
    Scam --> PG
    Emergency --> PG
    Consent --> PG
    Media --> Object
    Queue --> Snapshot
    Snapshot --> PG
~~~

#### 5.2.1 Android client cho demo

Android app là native Kotlin + Jetpack Compose để giảm rủi ro bridge ở camera/audio/realtime trong demo một nền tảng. Không tạo shared UI/business module cho iOS ở giai đoạn này.

| Layer/module | Trách nhiệm |
| --- | --- |
| <code>app</code> | Navigation, dependency wiring, build config, deep link và lifecycle toàn app |
| <code>feature-forum</code> | Feed, search, place, post/review/report và draft |
| <code>feature-live</code> | PTT state machine, AudioRecord, WebSocket event ordering, transcript/audio playback |
| <code>feature-lens</code> | CameraX preview, frame sampler, manual crop, capture/upload và candidate confirmation |
| <code>feature-safety</code> | SOS, incident card, hotline dialer và signed offline safety pack |
| <code>core-network</code> | HTTP/WebSocket, auth refresh, retry policy, certificate/TLS handling và trace ID |
| <code>core-data</code> | Repository, Room cache/draft, DataStore preference và WorkManager sync |
| <code>core-security</code> | Android Keystore-backed token storage, consent state và sensitive screen policy |
| <code>core-ui</code> | Compose design system, accessibility, locale/RTL-safe layout và error/degraded state |

Nguyên tắc Android:

- Camera preview, manual crop, VAD/noise suppression nhẹ và audio playback là xử lý client. Kiến trúc cho phép signed model pack on-device, nhưng flavor demo tắt đường chạy này bằng <code>AI_EXECUTION_MODE=server_only</code>.
- PTT dùng foreground UI session, không giữ micro ở background. CameraX unbind khi app background/permission bị thu hồi/end session.
- OkHttp/WebSocket hoặc transport tương đương gửi PCM/frame theo session token ngắn hạn; backpressure phải bỏ frame cũ trước, không làm nghẽn audio/control event.
- Android system TTS là đường chạy mặc định của thiết kế; router có thể nhận audio TTS server. Offline luôn có phrase audio đã duyệt và đóng gói sẵn.
- Room chỉ giữ post/draft/safety pack cần thiết. Raw PCM/frame/capture nằm trong memory hoặc encrypted cache có TTL; WorkManager dọn rác và đồng bộ deletion receipt.
- FCM chỉ mang opaque notification ID; app fetch nội dung sau khi xác thực. Không đưa incident, transcript hoặc cảnh báo nhạy cảm vào notification payload.
- Build demo có product flavor <code>demo</code>; endpoint và execution policy lấy từ generated BuildConfig, còn provider secret chỉ nằm ở backend. Không đóng gói <code>FPT_AI_API_KEY</code> trong APK ở bất kỳ flavor nào.
- Cleartext <code>http://10.0.2.2</code>/<code>ws://10.0.2.2</code> chỉ được allowlist trong debug network security config để nối Android Emulator tới máy dev; release build bắt buộc HTTPS/WSS và không cho cleartext.
