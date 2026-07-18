# Tourtect — Thiết kế nền tảng cộng đồng bảo vệ khách du lịch

| Thuộc tính | Giá trị |
| --- | --- |
| Trạng thái | System Design v2 — forum-first |
| Cập nhật | 18/07/2026 |
| Sản phẩm chính | Forum du lịch đa ngôn ngữ trên responsive web và Android app |
| Kênh bổ trợ | Zalo Mini App, Admin/Moderation web, Tourtect Live/Lens trên Android |
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

Thông tin giúp khách tránh rủi ro đang phân tán giữa forum, review, báo chí, video mạng xã hội và kinh nghiệm địa phương. Nội dung thường khác ngôn ngữ, nhanh lỗi thời, thiếu bằng chứng và khó tìm đúng lúc. Tourtect giải quyết khoảng trống này bằng một **community knowledge graph theo địa điểm**; AI hỗ trợ chuẩn hóa thông tin nhưng không thay cộng đồng hoặc cơ chế kiểm chứng.

Hệ thống tách ba lớp trách nhiệm:

1. Cộng đồng tạo post, review, price report, scam report và thảo luận.
2. Nền tảng tổng hợp nội dung công khai đúng quyền sử dụng, liên kết nội dung với địa điểm/chủ đề và xếp hạng theo độ hữu ích, độ mới, mức bằng chứng.
3. AI dịch, OCR, tóm tắt và trích dữ kiện; Price Engine, Safety Engine và moderation policy mới quyết định cảnh báo hoặc hành động nền tảng.

### 1.2 Đề xuất giá trị

**Tourtect là forum an toàn du lịch**, nơi khách có thể tìm một địa điểm và thấy trong cùng một màn hình:

- Review có cấu trúc về minh bạch giá, chất lượng dịch vụ và cảm giác an toàn.
- Khoảng giá tham khảo, lịch sử price report và bằng chứng hóa đơn đã khử thông tin cá nhân.
- Scam report, cảnh báo chính thức, bài báo, video và thảo luận liên quan.
- Bản dịch đa ngôn ngữ, tóm tắt “cần biết gì trước khi đến” và câu hỏi cho cộng đồng.

Android app bổ sung hai công cụ AI tại chỗ:

- **Live Voice:** hai bên dùng push-to-talk; ứng dụng dịch hai chiều và âm thầm nhận biết món/dịch vụ, số tiền, tiền tệ hoặc dấu hiệu nguy hiểm.
- **Live Camera:** ứng dụng nhìn menu, món ăn, đồ vật, bảng giá hoặc đồng hồ taxi; kết hợp VLM để hiểu bối cảnh với ảnh tĩnh độ phân giải cao để xác nhận OCR và giá.

Kết quả từ Live/Lens có thể được người dùng chủ động chuyển thành draft price report hoặc scam report sau khi đã an toàn. Không tự đăng transcript, ảnh hoặc vị trí. Khi có nguy hiểm, Safety Engine ưu tiên hướng dẫn rời khỏi tình huống và mở dialer thật; AI không tự gọi hoặc tự gửi vị trí.

### 1.3 Liên hệ tiêu chí chấm

| Tiêu chí | Bằng chứng trong thiết kế |
| --- | --- |
| Sức mạnh cộng đồng | Place page, post/review có cấu trúc, reputation, moderation, appeal và feed đa ngôn ngữ |
| Độ chính xác cảnh báo giá | Bốn mức kết quả, ngưỡng đỏ nghiêm ngặt, xác nhận OCR, so sánh theo ngành và cho phép “không đủ dữ liệu” |
| Tỷ lệ báo động giả | Red alert precision mục tiêu ≥ 95%, false-positive rate ≤ 2%, calibration theo vùng/ngành |
| Dịch và xử lý khẩn cấp | PTT theo vai, glossary dữ kiện quan trọng, rule engine đứng trước LLM, safety pack offline |
| Cập nhật dữ liệu vùng | Price report cộng đồng, external connector, evidence ledger, quarantine, human review và snapshot có phiên bản |
| Quyền riêng tư | Đọc ẩn danh, consent theo mục đích, media mặc định chỉ tồn tại trong bộ nhớ, đăng bài là hành động xác nhận riêng |
| Kinh doanh bền vững | Quảng cáo, affiliate, Plus, business tools và B2B insights nằm sau commercial trust firewall |

### 1.4 Nguyên tắc sản phẩm

1. **Community first, AI assists:** giá trị cốt lõi đến từ con người và dữ liệu có provenance; AI chỉ hỗ trợ khám phá, dịch và chuẩn hóa.
2. **Safety before engagement:** cảnh báo khẩn cấp quan trọng hơn thời gian xem, lượt tương tác hoặc doanh thu.
3. **No evidence, no accusation:** đánh giá giao dịch/hành vi, không tự gắn nhãn một cá nhân hay doanh nghiệp là “lừa đảo”.
4. **Trust is not for sale:** tiền quảng cáo hoặc gói business không tác động review, moderation, alert threshold hay organic ranking.
5. **Right to reply and appeal:** cơ sở kinh doanh được phản hồi, báo sai và kháng nghị nhưng không được trả tiền để xóa nội dung hợp lệ.
6. **Consent by action:** xin camera, micro, vị trí, đăng công khai và đóng góp dataset ở các bước riêng biệt.
7. **Graceful abstention:** dữ liệu yếu phải trả “không đủ dữ liệu”.
8. **Emergency remains useful offline:** hotline, phrasebook và incident card không phụ thuộc AI online hay gói trả phí.

---

## 2. Phạm vi sản phẩm

### 2.1 Người dùng chính

| Persona | Nhu cầu |
| --- | --- |
| Khách Hàn/Trung/Anh/Nga | Tìm review đúng ngôn ngữ, so sánh giá, hỏi cộng đồng, nhận cảnh báo và hỗ trợ tại chỗ |
| Người Việt/expat địa phương | Chia sẻ kinh nghiệm, xác minh giá, trả lời câu hỏi và cảnh báo rủi ro |
| Chủ cơ sở kinh doanh | Claim hồ sơ, cập nhật thông tin/giá, phản hồi review và kháng nghị công bằng |
| Nhân viên khách sạn/hướng dẫn viên | Tra playbook, hỗ trợ tạo incident card, xác minh thông tin địa phương |
| Moderator/Data reviewer | Duyệt UGC, bằng chứng giá, scam pattern, hotline, appeal và snapshot |
| Editor/Content steward | Quản trị nguồn báo/video, quyền sử dụng, phiên bản, takedown và audit |

### 2.2 Phân chia theo kênh

| Kênh | Có | Không có |
| --- | --- | --- |
| Responsive web | Full forum/feed/search, place page, post/review, so sánh giá, scam map, nội dung báo/video, tài khoản và monetization | Live Voice/continuous camera |
| Android app | Toàn bộ forum, notification/nearby alert, offline pack, Tourtect Live, Tourtect Lens, price check và SOS | Phiên dịch bên trong cuộc gọi PSTN; bản iOS |
| Zalo Mini App | Feed lite, tìm place, chụp ảnh/price report nhanh, text scam report, incident card và hotline qua openPhone | Full-duplex audio/video, WebRTC, AI nghe cuộc gọi |
| Admin/Moderation web | UGC queue, appeal, source/rights registry, dataset/scam/hotline versioning, ads eligibility, metrics và rollback | Xem raw media khi không có quyền và mục đích hợp lệ |

> Quyết định demo: forum là nền tảng chính; Android là mobile client duy nhất và iOS tạm ngoài phạm vi. Tourtect Live/Lens là mô-đun hỗ trợ; kiến trúc model vẫn adaptive, còn deployment demo chọn profile server-only bằng env.

### 2.3 Trong phạm vi pilot

- Nơi thử nghiệm: Một vài quận/huyện ở Hà Nội
- Forum đa ngôn ngữ, place page, review, price report, scam report, Q&A, feed và search.
- Connector pilot cho nguồn chính thức Việt Nam, RSS/Atom/báo allowlist, YouTube Data API/embed, OSM Vietnam extract/diff và URL do người dùng gửi.
- Query pack Việt/Anh/Hàn/Trung giản thể/Nga để phát hiện nội dung liên quan đến scam, gián đoạn và giá tại Việt Nam; chỉ publish sau rights gate và moderation.
- Bốn nhóm giá: taxi/đưa đón, đổi tiền, ăn uống, tour.
- UI và dịch hai chiều giữa tiếng Việt với:
  - <code>ko-KR</code>
  - <code>zh-Hans</code>
  - <code>en</code>
  - <code>ru-RU</code>
- Scam pattern seed: taxi, đổi tiền, ghost tour, ép thanh toán và tình huống có nguy cơ leo thang.
- Hotline khẩn cấp toàn quốc và hotline du lịch địa phương đã được xác minh.
- Android app native; iOS không build, test hoặc phát hành trong demo.

### 2.4 Ngoài phạm vi

- Cầu nối PSTN/SIP/WebRTC ba bên.
- Nhân viên phiên dịch trực 24/7.
- Full AI offline trên mọi thiết bị.
- Live Voice hoặc Live Camera trên web/Zalo Mini App.
- Tự gọi, tự gửi vị trí, tự gửi incident card hoặc tự liên hệ trusted contact.
- Chấm điểm hoặc công khai danh sách người bán bị tố cáo.
- Dùng dữ liệu hóa đơn thuế hoặc dữ liệu riêng của đối tác khi chưa có quyền hợp pháp.
- Sao chép toàn văn bài báo, tải lại video hoặc vượt paywall/login/robots/access control.
- Browser-scrape GrabFood, ShopeeFood, Shopee, Reddit, TikTok, Facebook, Instagram, Naver Cafe, Telegram hoặc nền tảng tương tự khi chưa có API/partner permission phù hợp.
- iOS app, App Store release, iOS-specific authentication/media/background behavior và test device iPhone/iPad.
- Cho nhà quảng cáo mua vị trí trong organic ranking, mua điểm review hoặc trả tiền để xóa cảnh báo hợp lệ.
- Mở marketplace/booking checkout riêng trong pilot; V1 chỉ dùng referral/affiliate rõ disclosure.

---

## 3. Trải nghiệm người dùng

### 3.1 Kiến trúc thông tin đa kênh

Navigation chính trên web và Android:

1. **Khám phá:** feed theo địa điểm/chủ đề, nearby và trending.
2. **Tìm kiếm:** place, món/dịch vụ, khoảng giá, scam pattern và nội dung.
3. **Đăng:** thảo luận, review, price report, scam report, câu hỏi hoặc tip.
4. **Đã lưu:** post, place, danh sách và safety pack.
5. **Hồ sơ:** reputation, contribution, notification và cài đặt ngôn ngữ.

Android có action dock luôn truy cập được:

- **Tourtect Live:** gọi AI bằng PTT để dịch hai chiều.
- **Tourtect Lens:** camera assist cho menu, đồ vật, món ăn và đồng hồ taxi.
- **SOS:** hotline, incident card và phrasebook offline.

### 3.2 Wireflow tổng thể

~~~mermaid
flowchart TD
    LANDING["Tourtect Home"] --> FEED["Feed theo vùng/ngôn ngữ"]
    LANDING --> SEARCH["Search place / giá / scam"]
    LANDING --> CREATE["Tạo post"]
    FEED --> PLACE["Place page"]
    SEARCH --> PLACE
    PLACE --> REVIEWS["Review + merchant reply"]
    PLACE --> PRICES["Khoảng giá + price reports"]
    PLACE --> ALERTS["Scam reports + official alerts"]
    PLACE --> SOURCES["Báo / video / post liên quan"]
    CREATE --> TYPE{"Loại nội dung"}
    TYPE --> REVIEW["Review có cấu trúc"]
    TYPE --> PRICE["Price report + bằng chứng"]
    TYPE --> SCAM["Scam report + safety check"]
    TYPE --> DISCUSS["Thảo luận / câu hỏi / tip"]
    PRICE --> MOD["Moderation + evidence level"]
    SCAM --> MOD
    MOD --> PUBLISH["Xuất bản / giới hạn / kháng nghị"]
    PLACE --> MOBILE_AI["Mở Tourtect Live / Lens trên Android"]
    MOBILE_AI --> DRAFT["Kết quả riêng tư"]
    DRAFT --> OPTIN{"Người dùng chủ động đăng?"}
    OPTIN -- "Có" --> CREATE
    OPTIN -- "Không" --> END["Kết thúc và xóa media theo TTL"]
~~~

### 3.3 Forum, place page và feed

#### Loại nội dung

| Loại | Trường bắt buộc | Cơ chế tin cậy |
| --- | --- | --- |
| Discussion/Q&A | Tiêu đề, nội dung, topic hoặc place | Reputation, vote hữu ích, moderation |
| Review | Place, thời điểm trải nghiệm, rating có cấu trúc, disclosure quan hệ | Merchant reply, evidence badge, appeal |
| Price report | Item, amount, currency, unit, vùng, thời điểm | Xác nhận OCR/manual, provenance, quarantine trước khi vào Price Engine |
| Scam report | Mô tả hành vi, vùng, thời điểm, mức an toàn hiện tại | Safety triage, khử PII, không tự nêu danh tính cá nhân |
| Tip/itinerary | Khu vực, thời gian hiệu lực, nội dung hướng dẫn | Freshness, saved/share, community correction |
| Official alert | Nguồn có thẩm quyền, ngày hiệu lực, canonical URL | Editor verification và expiry |
| External article/video | Canonical URL, nguồn, tác giả, ngày xuất bản | Rights/policy gate, attribution, dedup, takedown sync |

#### Place page

Mỗi địa điểm/cơ sở/dịch vụ có một trang chuẩn hóa gồm:

- Tên đa ngôn ngữ, loại hình, khu vực, bản đồ với độ chính xác phù hợp và thông tin do chủ cơ sở claim.
- Điểm tổng hợp và phân bố review; các chiều riêng: minh bạch giá, chất lượng, an toàn và mức đáng tiền.
- Khoảng giá hiện tại, lịch sử thay đổi, cỡ mẫu, freshness và nút thêm price report.
- Scam/safety signal theo hành vi đã được kiểm chứng; không hiển thị “scam score” kết luận pháp lý.
- Post, câu hỏi, bài báo và video liên quan; câu trả lời chính thức/merchant reply được gắn nhãn.
- Nút lưu, theo dõi, chia sẻ, báo nội dung, đề xuất chỉnh sửa và kháng nghị.

#### Feed và ranking

Feed gồm Following, Nearby, Latest, Trending và Safety Alerts. Người dùng chọn vùng thủ công trước; vị trí chính xác là opt-in. Ranking hữu cơ chỉ dùng relevance, freshness, evidence, source diversity, community usefulness và safety priority. Các tín hiệu sau bị cấm trong organic ranker: chi tiêu quảng cáo, gói business, affiliate commission và yêu cầu của sales.

Tourtect hiển thị bản dịch theo locale người đọc nhưng luôn cho xem bản gốc. Nội dung AI dịch/tóm tắt có nhãn, nút báo lỗi và không thay thế nguyên văn của tác giả.

#### Reputation và moderation

- Reputation tăng từ contribution được xác minh, sửa lỗi hữu ích, lịch sử ổn định và đánh giá của moderator; không tăng chỉ vì đồng thuận hoặc nhiều upvote.
- Tách reputation theo năng lực: local knowledge, price evidence, translation, safety và merchant representative.
- Phát hiện spam, review exchange, brigading, sockpuppet/Sybil, copy-paste và bất thường theo thời gian/thiết bị/graph.
- Nội dung rủi ro cao đi qua human review; người đăng và đối tượng bị review có quy trình notice, phản hồi và appeal.
- Không có “pay to remove”. Nội dung có thể bị gỡ vì sai, vi phạm, thiếu căn cứ hoặc theo quy trình pháp lý, và quyết định phải có audit trail.

### 3.4 Tourtect Live (Live Voice)

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
7. TTS hệ điều hành phát bản dịch theo audio route hiện tại; router có thể dùng server TTS khi cần.
8. Intelligence lane độc lập trích giá, món/dịch vụ, scam signal và critical safety signal.
9. Nếu có price candidate đủ rõ, Price Engine trả insight kín.
10. Nếu có red flag, Safety Engine ưu tiên escalation card nhưng không làm mất transcript/dịch.

#### Dual-lane sequence

~~~mermaid
sequenceDiagram
    actor Tourist as Du khách
    participant App as Android App
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
        App-->>Tourist: System TTS + caption
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
- Nếu TTS locale không tồn tại hoặc engine lỗi:
  1. Hiển thị caption chữ lớn.
  2. Cho phép người dùng đưa màn hình cho bên kia đọc.
  3. Dùng audio phrase đã đóng gói cho câu khẩn cấp đã duyệt.
- Khi dùng tai nghe, người dùng có thể bấm “Đọc riêng” để nghe giải thích giá.

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

### 3.6 Scam Assistant và safety triage

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

### 3.7 Emergency

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

### 3.8 Responsive web

Website là bề mặt chính cho discovery, SEO, thảo luận và doanh thu:

- Đọc feed, place page, review, bảng giá, alert và nội dung ngoài mà không cần tài khoản.
- Tìm kiếm theo thành phố, địa điểm, topic, ngôn ngữ, thời gian, loại bằng chứng và mức giá.
- Tài khoản bắt buộc khi đăng, bình luận, vote, theo dõi, báo nội dung hoặc lưu đồng bộ.
- Cho phép upload một ảnh để price check và chuyển kết quả thành draft post sau xác nhận riêng.
- Có trang public indexable cho place/topic; trang cá nhân, incident và dữ liệu nhạy cảm không index.
- Hiện QR/deep link sang Android cho Tourtect Live/Lens; web không chạy continuous camera/audio trong V1.

### 3.9 Zalo Mini App

Đối tượng ưu tiên là người Việt đồng hành, khách sạn, hướng dẫn viên; khách quốc tế đã có Zalo vẫn dùng được.

V1:

- Feed lite theo thành phố, search place và xem review/price/scam card.
- Tạo price report hoặc scam report ngắn; bước xuất bản phải xác nhận nội dung và phạm vi công khai.
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

### 3.10 Thu thập bài báo, video mạng xã hội và post ngoài Tourtect

Tourtect không xây crawler “lấy mọi thứ”. Mỗi connector phải có phương thức truy cập và quyền hiển thị rõ ràng:

1. Ưu tiên API chính thức, RSS/Atom, partner feed hoặc URL do người dùng chia sẻ.
2. Kiểm tra allowlist nguồn, robots/access control, điều khoản API, khả năng embed và loại giấy phép trước khi fetch.
3. Chỉ lưu metadata cần thiết: canonical URL, source, author/channel, thời điểm, thumbnail được phép, đoạn mô tả ngắn được phép và content hash. Không lưu toàn văn báo hoặc re-host video nếu chưa có license.
4. Chuẩn hóa, phát hiện trùng/cùng sự kiện, entity-link tới place/topic/scam pattern và phân loại fact/opinion/sponsored.
5. Chạy moderation; claim rủi ro cao phải chờ editor hoặc hiển thị dưới dạng “nguồn đang đưa tin”, không được biến thành kết luận của Tourtect.
6. Dùng embed/player chính thức khi được phép; nếu không thì chỉ hiển thị link card có attribution.
7. Đồng bộ sửa/xóa/disable từ nguồn, hỗ trợ takedown và đặt TTL tái kiểm tra.

Với YouTube, connector dùng YouTube Data API, lọc video cho phép embed và phát bằng player chính thức. TikTok, Facebook/Instagram hoặc nền tảng khác chỉ tích hợp qua API/embed chính thức đã được phê duyệt; chưa có quyền thì chỉ nhận link do người dùng chủ động chia sẻ. “Public trên trình duyệt” không đồng nghĩa với quyền thu thập hàng loạt hoặc quyền tái sử dụng.

~~~mermaid
flowchart LR
    SRC["Official API / RSS / partner feed / shared URL"] --> RIGHTS{"Rights & policy gate"}
    RIGHTS -- "Không đạt" --> LINK["Link card hoặc từ chối"]
    RIGHTS -- "Đạt" --> FETCH["Fetch metadata / permitted snippet"]
    FETCH --> NORM["Normalize + canonicalize + dedupe"]
    NORM --> ENTITY["Place/topic/entity linking"]
    ENTITY --> SAFETY["Classifier + moderation"]
    SAFETY --> INDEX["Search index + feed candidate"]
    INDEX --> REFRESH["Refresh / deletion / takedown sync"]
~~~

#### 3.10.1 Crawl những gì trước: travel-critical corpus cho Việt Nam

Mục tiêu không phải số URL lớn mà là trả lời được câu hỏi của khách **trước, trong và ngay sau chuyến đi**. Backlog dữ liệu được ưu tiên như sau:

| Mức | Nhóm dữ liệu | Trường tối thiểu | Nguồn ưu tiên | Freshness mục tiêu |
| --- | --- | --- | --- | --- |
| P0 | Khẩn cấp, thiên tai, thời tiết cực đoan, ngập/cháy, dịch bệnh, biểu tình hoặc đóng cửa đột xuất | Vùng ảnh hưởng, hiệu lực, hướng dẫn, nguồn, hotline đã xác minh | Cơ quan trung ương/địa phương, khí tượng, phòng chống thiên tai, y tế, công an, đại sứ quán | 5–15 phút khi có sự kiện; ≤ 60 phút bình thường |
| P0 | Visa/xuất nhập cảnh, cửa khẩu, sân bay, chuyến bay/tàu/xe bị gián đoạn | Điều kiện áp dụng, ngày hiệu lực, tuyến, trạng thái, canonical URL | Cổng xuất nhập cảnh, cảng hàng không, hãng vận chuyển/nhà ga qua API/feed/partner | 5–30 phút cho trạng thái; hằng ngày cho quy định |
| P0 | Scam pattern đang nổi | Hành vi, địa bàn, thời gian, bằng chứng, ngôn ngữ, mức xác minh | Cơ quan chức năng, báo chí allowlist, Tourtect report, social URL đã qua rights gate | 15–60 phút cho discovery; human review trước cảnh báo rộng |
| P1 | Địa điểm thiết yếu | Tên/alias, tọa độ, giờ mở, contact, accessibility, category, trạng thái đóng/mở | OSM, cổng dữ liệu/cơ quan du lịch, merchant claim, khảo sát | 6–24 giờ cho trạng thái; 1–7 ngày cho POI thường |
| P1 | Giá có thể gây tranh chấp | Item/unit, giá niêm yết và giá thực trả, thời điểm, địa bàn, phí/phụ thu, bằng chứng | Bảng giá chính thức, merchant, hóa đơn người dùng, khảo sát; partner commerce nếu có hợp đồng | 6–24 giờ với nguồn động; snapshot theo tuần nếu ổn định |
| P1 | Giao thông và di chuyển | Điểm đón, loại phương tiện, khoảng giá, tuyến, giờ, phụ phí, cách mua vé | Đơn vị vận tải/địa phương, merchant, OSM, khảo sát | 15 phút–24 giờ tùy dữ liệu realtime/tĩnh |
| P1 | Điểm đến, sự kiện, vé, quy tắc ứng xử | Lịch, giá vé, dress code, closure, accessibility, official URL | Cơ quan du lịch/đơn vị vận hành/merchant | 6 giờ khi gần sự kiện; hằng ngày/tuần nếu tĩnh |
| P2 | Review, tip, video và thảo luận | Canonical URL, locale, author/channel, thời gian, permitted snippet/embed, claim/evidence | API chính thức, RSS, URL người dùng gửi | 15–60 phút cho nguồn đã theo dõi; long-tail theo ngày |
| P3 | Nội dung cảm hứng/evergreen | Chủ đề, địa điểm, ngôn ngữ, thời gian hữu ích | Creator/press/partner có quyền | Hằng tuần |

Trong pilot Hà Nội, seed theo hành trình: sân bay Nội Bài → taxi/xe công nghệ/xe buýt → khách sạn → đổi tiền/SIM/eSIM → ăn uống/menu → điểm tham quan/vé → nightlife/mua sắm → tour/day trip → y tế/công an/đại sứ quán. Mỗi hành trình phải có giá/unit, câu hỏi nên hỏi, scam pattern, phương án an toàn và nguồn cập nhật; không seed nội dung chỉ để làm feed trông đông.

#### 3.10.2 Ma trận quyết định theo nguồn

| Nguồn | Giá trị cần lấy | Cách truy cập được chấp nhận | Không làm | Quyết định pilot |
| --- | --- | --- | --- | --- |
| Cổng Chính phủ, Bộ/ngành, UBND, cơ quan du lịch, khí tượng, sân bay/nhà ga | Quy định, cảnh báo, hotline, closure, lịch/tariff chính thức | API/open data/RSS/sitemap; HTML allowlist nếu robots và điều khoản cho phép | Suy diễn từ bản tin cũ; coi social repost là văn bản gốc | **Ưu tiên cao nhất** |
| Báo chí Việt Nam và báo đa ngôn ngữ | Scam case, disruption, thay đổi chính sách, điều tra giá | RSS/Atom/licensed feed; fetch bài mới/đổi với conditional GET; lưu metadata/snippet được phép | Lấy toàn văn, ảnh hoặc paywall; lấy từ aggregator khi đã có canonical publisher | **Allowlist theo tòa soạn** |
| Báo Mới/aggregator tương tự | Discovery URL và cluster sự kiện | Chỉ qua feed/API/partner được phép; resolve về nguồn xuất bản gốc | Dùng aggregator làm nguồn độc lập thứ hai hoặc copy snippet/ảnh lặp | **Discovery-only** |
| Grab/GrabFood | Tên merchant, menu/giá/phí tại thời điểm, ETA | Partner/merchant export hoặc API/hợp đồng bằng văn bản; receipt do user opt-in | Reverse-engineer app/private endpoint, crawl catalog đại trà, dùng giá khuyến mại/cá nhân hóa làm giá phố | **Không crawl; xin partnership** |
| ShopeeFood | Merchant/menu/giá/khuyến mại | Partner feed, merchant export hoặc user receipt có consent | Robot/spider/scrape khi chưa có chấp thuận bằng văn bản | **Blocked-by-default** |
| Shopee marketplace | Giá SIM, travel accessory, ticket/tour listing có liên quan | Shopee Open Platform/affiliate/partner nếu use case và quyền cho phép | Crawl search/product/review; dùng seller listing làm fact an toàn | **Chỉ partner/API** |
| YouTube/Shorts | Video cảnh báo, hướng dẫn, local news, creator experience | YouTube Data API; query theo locale/region/date; <code>videos.list</code> để refresh; official embed/player | Download/re-host video, scrape trang watch/comment | **Connector pilot** |
| TikTok | Video URL do user/creator đưa vào, nội dung của creator đã OAuth | Display API cho account đã cấp quyền, official embed/link card | Dùng Display API để giả lập global search; browser scraping; dùng Research API thương mại | **Submission/creator opt-in** |
| Facebook/Instagram/Reels | Post/reel từ cơ quan, báo, business/creator và URL do user gửi | Graph/Instagram API đúng permission, oEmbed/embed khi được duyệt; partner; Meta Content Library chỉ nếu tổ chức/use case đủ điều kiện nghiên cứu | Quét profile/group/comment public bằng browser; truy cập private/closed group; lưu media | **Official pages + submissions** |
| Reddit | Thread scam/travel theo subreddit và query | Data API sau khi có phê duyệt; Reddit embed/link; deletion sync | Scrape HTML/API không duyệt; dùng thương mại hoặc huấn luyện AI khi chưa có chấp thuận | **Chờ approval; link submission trước** |
| OSM | Road/POI/category/address geometry và alias | Regional PBF + replication diff; API/provider phù hợp; attribution ODbL | Bulk qua public Nominatim/Overpass/tile server; coi OSM là nguồn pháp lý cho địa giới | **Base map/place seed** |
| Google Maps/Places | Cross-check place ID, giờ mở, business status nếu dùng | Places API theo license và storage/display rules | Scrape Maps/Search result hoặc trộn/cache trường bị hạn chế vô thời hạn | **Không cần cho pilot; đánh giá sau** |

Giá trên GrabFood/ShopeeFood/Shopee thường chứa markup, voucher, phí nền tảng, phí giao hàng, membership, vị trí và thời điểm; vì vậy chỉ tạo <code>CommerceOfferObservation</code>, không nhập thẳng vào <code>PriceObservation</code>. Chỉ sau khi tách được item base price, phí, discount, delivery và điều kiện áp dụng, nguồn này mới được dùng làm một tín hiệu tham khảo có nhãn “giá online”, không làm reference truth cho giá tại chỗ.

#### 3.10.3 Trung Quốc, Hàn Quốc và Nga: lấy tín hiệu về Việt Nam

Việt Nam vẫn là entity/geo scope chính; mở rộng thị trường nghĩa là bổ sung **nguồn và query theo ngôn ngữ của khách đang nói về Việt Nam**, không nhập toàn bộ dữ liệu du lịch tại ba quốc gia.

| Thị trường | Query pack ban đầu | Nguồn khả thi | Nguồn bị giới hạn/cách xử lý |
| --- | --- | --- | --- |
| Việt/Anh | <code>lừa đảo du lịch</code>, <code>chặt chém</code>, <code>taxi sân bay</code>, <code>đổi tiền</code>, <code>Vietnam scam</code>, cùng alias địa danh | Cơ quan/báo RSS, YouTube, Tourtect community, OSM, merchant/partner | Facebook/TikTok/commerce theo ma trận quyền |
| Trung giản thể | <code>越南 诈骗</code>, <code>越南 出租车 宰客</code>, <code>河内 换汇</code>, tên Việt/Hán/pinyin của place | Nguồn báo/du lịch có RSS, Baidu Search/Maps API nếu ký use case phù hợp, Bilibili/Douyin creator opt-in, shared URL | Douyin API chỉ lấy video account đã cấp OAuth; Xiaohongshu/WeChat/Weibo không giả định có public search API thương mại—dùng partner, creator opt-in hoặc link submission |
| Hàn | <code>베트남 사기</code>, <code>하노이 택시 바가지</code>, <code>환전</code>, alias Hangul/Latin | Naver/Daum Search API nếu điều khoản cho phép index metadata, YouTube, public news/blog, creator opt-in; VisitKorea/TourAPI là mẫu taxonomy đa ngôn ngữ | Naver Cafe/private post và Kakao group không crawl; chỉ public result/link, owner/partner feed hoặc user submission |
| Nga | <code>Вьетнам мошенничество</code>, <code>такси Ханой обман</code>, <code>обмен валюты</code>, alias Cyrillic/Latin | Yandex Search API với region/language/date, YouTube, VK API nếu được cấp quyền, báo/RSS tiếng Nga | Telegram không dùng để aggregate/AI pipeline nếu chưa có cơ sở quyền rõ; nhận link do user gửi hoặc hợp tác với channel owner |

Mỗi query có <code>query_pack_id</code>, locale, market, intent, place aliases, negative keywords, version và hiệu suất. Editor bản ngữ duyệt query/translation và false-positive. Kết quả từ một cộng đồng/ngôn ngữ không tự động trở thành “scam pattern”: cần entity-link đúng nơi/thời điểm, dedupe cross-post, ít nhất một nguồn độc lập hoặc evidence cộng đồng và human review cho cảnh báo high-risk.

Khi Tourtect thật sự mở rộng địa lý ra ngoài Việt Nam, dùng connector bản địa: Korea Tourism Organization TourAPI/Kakao Local cho Hàn Quốc, Baidu Maps cho Trung Quốc và Yandex Maps/Search cho Nga, sau review license/data residency riêng. Không dùng các API này để thay thế dữ liệu chính thức Việt Nam.

#### 3.10.4 Luôn cập nhật mà không gây tải hoặc bị coi là bot xấu

Tourtect không né bot detection. Crawler dùng User-Agent ổn định dạng <code>TourtectBot/1.0 (+https://tourtect.example/crawler; data@tourtect.example)</code>, IP egress ổn định, trang mô tả mục đích/opt-out và email xử lý abuse. Không xoay proxy/IP, giả browser fingerprint, dùng tài khoản người dùng, vượt CAPTCHA, login, paywall hoặc endpoint private.

Lịch mặc định là **adaptive**, có jitter ±20% để tránh burst đồng hồ:

| Connector/item state | Poll mặc định | Sau khi không đổi | Khi active/breaking | Ghi chú |
| --- | --- | --- | --- | --- |
| Official alert/transport API hoặc feed | 5–15 phút | 30–60 phút | 2–5 phút nếu API/quota cho phép | Ưu tiên webhook/WebSub/stream nếu nguồn cung cấp |
| News RSS/Atom allowlist | 10–15 phút | Giữ 15–30 phút | 5 phút cho feed breaking-news đã thỏa thuận | Chỉ fetch article khi GUID/URL/hash mới hoặc đổi |
| Official HTML không có feed | 30–120 phút | 6 giờ → 24 giờ | 15–30 phút có thời hạn do editor bật | Chỉ allowlist; không deep-crawl toàn domain |
| Social official API query | 15–60 phút | 2–6 giờ | 10–15 phút cho incident query | Tuân quota/rate header; refresh ID đã biết rẻ hơn search lại |
| Commerce/merchant/partner | 6–24 giờ | 1–7 ngày | 1–6 giờ trong campaign/sự kiện | Không poll nếu partner có webhook/export delta |
| OSM | Hourly replication cho pilot; minutely chỉ khi có nhu cầu | Nightly derived-place rebuild | Hourly | Full regional PBF định kỳ; lưu sequence để resume |
| URL active đã index | 2–6 giờ trong 48 giờ đầu | 1 ngày → 7 ngày → 30 ngày | 15–60 phút nếu source báo thay đổi | Xóa/disable được ưu tiên hơn enrichment |
| <code>robots.txt</code>/terms/rights policy | Cache tối đa 24 giờ cho robots; review policy hằng tuần | — | Ngay khi nhận 403/451/complaint | Policy đổi có thể đóng connector tự động |

Per-host scheduler dùng token bucket riêng. Với HTML không có chỉ dẫn khác, bắt đầu ở concurrency 1 và tối đa 1 request/5 giây/host; giảm tốc khi latency hoặc 5xx tăng. Với API, tuân rate-limit/quota chính thức và không tự đặt tốc độ cao hơn. Luôn:

- Dùng <code>If-None-Match</code>/<code>ETag</code>, <code>If-Modified-Since</code>/<code>Last-Modified</code>, gzip và field projection; <code>304</code> không chạy lại extraction.
- Canonicalize trước fetch sâu; dedupe URL/content/perceptual hash và event cluster để không tải cùng nội dung từ tracking URL/cross-post.
- Tôn trọng <code>robots.txt</code> theo RFC 9309, ToS, API policy và <code>Retry-After</code>. <code>401/403/451</code> đóng circuit và đưa vào policy review; <code>429</code> backoff theo header, nếu thiếu thì exponential backoff có jitter.
- Đặt crawl budget theo host/ngày, byte/day và error budget. Dừng tự động khi 429 > 1%, 5xx > 5%, latency P95 gấp 3 baseline, robots đổi sang disallow hoặc có abuse complaint.
- Lưu <code>fetched_at</code>, <code>source_updated_at</code>, <code>etag</code>, <code>last_modified</code>, status, bytes, latency, policy version và next-check; không lưu cookie/token/raw HTML lâu hơn nhu cầu xử lý.
- Có kill switch toàn cục/per-source, dashboard quota, dead-letter queue, replay idempotent và lịch takedown/deletion sync nhanh hơn lịch enrichment.

#### 3.10.5 Có cần Google Search API không?

**Không dùng Google Search API làm dependency của pilot.** Tại thời điểm cập nhật tài liệu, Custom Search JSON API đã đóng với khách hàng mới và khách hàng hiện hữu phải chuyển trước ngày 01/01/2027. Vertex AI Search phù hợp search trên tập domain/data được quản lý; Grounding with Google Search tạo grounded answer trong hệ sinh thái Gemini, không thay thế ingestion API có quyền lưu, refresh và deletion sync, đồng thời không phù hợp quyết định dùng FPT AI Factory cho runtime demo.

Discovery nên đi theo thứ tự: official API/RSS/WebSub → sitemap/allowlisted HTML → platform API (YouTube, search API bản địa) → partner feed → URL do user gửi. Nếu sau pilot vẫn thiếu long-tail discovery, mua một web-search API có điều khoản thương mại rõ, domain/date/language filter và quyền dùng metadata. Search result chỉ tạo <code>DiscoveryCandidate</code>; connector vẫn phải qua rights/robots/policy gate trước khi fetch hoặc publish. Không dùng search API để lách hạn chế của Grab/Shopee/Meta/TikTok/Reddit.

### 3.11 Mô hình doanh thu và commercial trust firewall

| Nguồn doanh thu | Giá trị cung cấp | Guardrail bắt buộc |
| --- | --- | --- |
| Quảng cáo contextual/native | Banner/card theo ngữ cảnh du lịch | Nhãn “Quảng cáo”, ad-eligible UGC, không đặt trên SOS/incident critical, không nhắm mục tiêu theo dữ liệu nhạy cảm |
| Affiliate/referral | Hotel, tour, bảo hiểm, eSIM, vận chuyển hoặc vé từ đối tác | Disclosure cạnh CTA, organic ranking độc lập commission, log click/conversion tối thiểu dữ liệu |
| Tourtect Plus | Không quảng cáo, list/offline pack nâng cao, alert tùy chỉnh, quota AI cao hơn | Không paywall SOS, hotline, báo cáo an toàn hoặc cảnh báo thiết yếu |
| Verified Business | Claim profile, cập nhật menu/giá, trả lời review, analytics và inbox | Không mua điểm, không ẩn review, không giảm alert; badge chỉ xác minh đại diện |
| Sponsored content | Destination guide hoặc topic do đối tác tài trợ | Nhãn tài trợ nổi bật, slot riêng, không trộn với official alert |
| B2B insights/API | Xu hướng aggregate cho cơ quan du lịch, khách sạn, hãng bảo hiểm | Chỉ aggregate đủ ngưỡng, privacy review, không bán raw post/media hoặc hồ sơ cá nhân |
| Grant/public-interest partnership | Tài trợ safety pack, dữ liệu địa phương và nghiên cứu | Công khai nhà tài trợ; không được can thiệp moderation |

**Commercial trust firewall** là ranh giới kỹ thuật và tổ chức:

- Ad/Billing service không có quyền ghi review score, moderation status, PriceSnapshot, ScamPattern hoặc organic ranking features.
- Sponsored inventory có namespace và analytics riêng; không truyền advertiser spend vào ranker hữu cơ.
- Sales không có quyền moderator. Mọi gỡ/sửa nội dung cần policy reason, actor và audit log; high-risk appeal cần dual review.
- Merchant có thể mua công cụ quản lý nhưng không mua “uy tín”. Verified badge chỉ chứng minh quyền đại diện.
- Trust-health gate có quyền tắt một bề mặt kiếm tiền nếu fake-review rate, complaint rate, policy violation hoặc ad-induced safety risk vượt ngưỡng.

Thứ tự triển khai doanh thu: contextual ads/affiliate sau khi moderation đủ chuẩn → Plus → business tools → B2B aggregate. Không tối ưu doanh thu trước khi có content eligibility, consent analytics và cơ chế appeal.

Các mức dưới đây chỉ là **giả thuyết để A/B test**, không phải giá đã chốt:

| Gói | Giả thuyết giá pilot | Cách kiểm chứng |
| --- | --- | --- |
| Tourtect Plus | 59.000–99.000 VND/tháng tại Việt Nam hoặc mức tương đương theo thị trường | Conversion trial→paid, retention 30/90 ngày, mức dùng offline/alert/AI quota |
| Verified Business | 299.000–999.000 VND/tháng theo số place và analytics | Tỷ lệ claim, phản hồi review, menu/giá được cập nhật và churn; tuyệt đối không đo bằng điểm review tăng |
| Contextual ads | CPM/CPC theo mạng quảng cáo hoặc direct campaign | Revenue per eligible session sau khi trừ policy loss và trust complaint |
| Affiliate | CPA/revenue share theo booking hợp lệ | Incremental conversion, cancellation/refund và disclosure comprehension |
| B2B insights | Báo giá theo coverage/SLA, không theo số hồ sơ cá nhân | Renewal, số vùng/vertical đủ aggregate threshold và privacy audit pass |

North-star kinh doanh không phải ad impressions mà là **trusted trip utility**: số người tìm được thông tin hữu ích hoặc tránh một quyết định rủi ro mà không làm tăng false accusation. Mọi dashboard doanh thu phải đặt cạnh retention, search success, fake-review rate, appeal overturn rate và trust complaint rate.

### 3.12 Đăng ký, đăng nhập và Google OAuth

Tourtect hỗ trợ ba trạng thái danh tính:

1. **Anonymous:** đọc nội dung, tìm kiếm, dùng SOS và tạo dữ liệu nháp cục bộ.
2. **Registered:** đăng bài, bình luận, vote, follow, đồng bộ saved list và nhận notification.
3. **Verified role:** moderator/editor hoặc đại diện doanh nghiệp đã qua quy trình xác minh riêng; đăng nhập Google không tự cấp vai trò này.

#### Đăng ký bằng email

- Người dùng nhập email, mật khẩu, display name, locale và chấp nhận Terms/Privacy/Community Guidelines thành các consent record có phiên bản.
- Gửi email verification bằng token dùng một lần, hash ở server và hết hạn sau 15–30 phút.
- Chỉ cho xuất bản nội dung sau khi xác minh email; vẫn cho chuẩn bị draft trong thời gian chờ.
- Mật khẩu tối thiểu 12 ký tự, cho phép password manager/paste, kiểm tra mật khẩu phổ biến/rò rỉ và hash bằng Argon2id với tham số có version.
- Quên mật khẩu dùng token một lần, vô hiệu hóa sau sử dụng; reset mật khẩu thu hồi các refresh session khác và gửi security notification.

#### Đăng nhập với Google

- Web dùng **Google Identity Services (GIS)**; Android dùng Credential Manager hoặc system browser theo OpenID Connect Authorization Code Flow với PKCE.
- Chỉ xin scope <code>openid email profile</code>. Đăng nhập không xin quyền Drive, Contacts, YouTube hoặc quyền đăng thay người dùng.
- Client tạo <code>state</code>, <code>nonce</code> và PKCE <code>code_verifier/code_challenge</code>; backend thực hiện code exchange và xác minh ID token.
- Backend kiểm tra chữ ký/JWKS, <code>iss</code>, <code>aud</code>, <code>exp</code>, <code>iat</code>, <code>nonce</code> và <code>email_verified</code>. Khóa liên kết bằng cặp <code>(issuer, sub)</code>, không dùng email làm định danh Google bất biến.
- Google access/refresh token không được lưu nếu chỉ dùng để đăng nhập. Authentication và authorization tới Google API là hai consent flow riêng.
- Redirect URI là allowlist chính xác theo môi trường; không nhận redirect do client tùy ý truyền.

#### Liên kết tài khoản và phiên

- Anonymous session được merge vào account sau đăng nhập theo preview rõ: saved item, draft và preference nào sẽ được chuyển; incident/private media không tự merge.
- Nếu email Google trùng account email/password đã xác minh, Tourtect yêu cầu người dùng chứng minh quyền kiểm soát account hiện tại trước khi link; không tự gộp chỉ dựa trên email.
- Một account có thể liên kết nhiều identity provider; phải giữ ít nhất một phương thức đăng nhập trước khi unlink.
- Web dùng session cookie <code>HttpOnly</code>, <code>Secure</code>, <code>SameSite=Lax/Strict</code> phù hợp; Android lưu refresh token bằng Android Keystore-backed encrypted storage. Access token sống ngắn, refresh token rotation và reuse detection.
- Trang Security hiển thị thiết bị/phiên gần đây, cho revoke từng phiên hoặc “đăng xuất tất cả”. Logout cục bộ phải xóa session; logout Google không đồng nghĩa xóa account Tourtect.
- Không dùng Google profile photo làm public avatar mặc định nếu người dùng chưa xác nhận phạm vi công khai.

~~~mermaid
sequenceDiagram
    actor U as Người dùng
    participant C as Web/Android
    participant I as Tourtect Identity
    participant G as Google Identity
    participant D as Account DB

    U->>C: Chọn Đăng nhập với Google
    C->>I: Tạo auth attempt + state/nonce/PKCE
    I-->>C: authorization URL đã allowlist
    C->>G: Authorization request
    G-->>C: code + state
    C->>I: callback(code, state, code_verifier)
    I->>I: Kiểm tra state và exchange code
    I->>G: Token/JWKS verification data
    G-->>I: ID token / keys
    I->>I: Verify iss/aud/exp/nonce/email_verified
    I->>D: Find/link bằng issuer + sub
    D-->>I: account + roles
    I-->>C: Tourtect session + rotated refresh credential
    C-->>U: Đăng nhập; hỏi trước khi merge anonymous data
~~~

---

## 4. Yêu cầu chức năng

### 4.1 Forum và discovery

| ID | Yêu cầu |
| --- | --- |
| FO-01 | Đọc public content không cần đăng nhập; đăng/tương tác cần account hoặc session đã xác minh |
| FO-02 | Hỗ trợ discussion, Q&A, review, price report, scam report, tip, official alert và external content card |
| FO-03 | Mọi post liên kết được với place, region, topic, locale và evidence level |
| FO-04 | Feed có Following, Nearby, Latest, Trending và Safety; mỗi item giải thích được lý do xuất hiện |
| FO-05 | Search full-text và filter theo place, vùng, vertical, ngôn ngữ, freshness, giá và evidence |
| FO-06 | Bản dịch AI luôn có nhãn, xem được nguyên bản và có feedback |
| FO-07 | Saved list, follow, comment, mention, notification, share và block/report |
| FO-08 | Organic ranking không nhận advertiser spend, affiliate commission hoặc business tier |
| FO-09 | Đăng ký email cần verification; đăng nhập hỗ trợ email/password và Google OIDC |
| FO-10 | Anonymous data chỉ merge vào account sau preview/confirmation; không merge incident/media mặc định |
| FO-11 | Account có quản lý phiên, revoke từng thiết bị, logout all, password reset và provider linking an toàn |
| FO-12 | OAuth dùng state/nonce/PKCE, redirect allowlist và server-side ID-token validation; không dùng email thay Google subject ID |

### 4.2 Review, reputation và moderation

| ID | Yêu cầu |
| --- | --- |
| RM-01 | Review có rating tổng thể và các chiều price transparency, service, safety, value |
| RM-02 | Disclosure bắt buộc nếu reviewer được mời, nhận ưu đãi, là nhân viên/đối tác hoặc có xung đột lợi ích |
| RM-03 | Evidence badge phân biệt không bằng chứng, metadata, hóa đơn đã xác minh và nguồn chính thức |
| RM-04 | Merchant reply, report và appeal có SLA; không có pay-to-remove |
| RM-05 | Phát hiện spam, duplicate, harassment, PII, fake review, brigading và Sybil trước/sau xuất bản |
| RM-06 | Reputation theo lĩnh vực, có decay và audit; không chỉ dựa vào upvote |
| RM-07 | Nội dung/sanction rủi ro cao cần human review; người dùng được biết lý do và quyền kháng nghị |
| RM-08 | Safety report ưu tiên bảo vệ nạn nhân nhưng tránh công khai danh tính, biển số hoặc cáo buộc chưa xác minh |

### 4.3 External content và monetization

| ID | Yêu cầu |
| --- | --- |
| CM-01 | Chỉ fetch qua official API, RSS/Atom, partner feed hoặc URL chia sẻ trong phạm vi được phép |
| CM-02 | Lưu canonical URL, attribution, rights status, retrieved/checked time và deletion state |
| CM-03 | Dedupe, entity-link, fact/opinion/sponsored labeling, freshness và takedown sync |
| CM-04 | Không vượt paywall/login/access control hoặc re-host toàn bài/video khi chưa có license |
| CM-05 | Ad/sponsored/affiliate phải có disclosure; SOS và critical incident không có quảng cáo |
| CM-06 | Business trả phí không được sửa rank/review/moderation/alert; mọi quyền được enforce bằng RBAC và audit |
| CM-07 | B2B chỉ xuất aggregate đủ ngưỡng và không bán raw content, media hoặc định danh người dùng |
| CM-08 | Trust-health gate có thể ngừng monetization surface độc lập với product availability |
| CM-09 | Mỗi connector có access mode, policy/robots version, owner, crawl budget, retention, deletion SLA và kill switch |
| CM-10 | Scheduler dùng per-host token bucket, conditional GET, adaptive backoff và không vượt login/paywall/CAPTCHA/private endpoint |
| CM-11 | Commerce offer tách khỏi giá tại chỗ; promotion/delivery/personalization không được nhập thẳng Price Engine |
| CM-12 | Query pack đa ngôn ngữ có version, editor bản ngữ, place aliases, negative keyword và precision review |
| CM-13 | Search API chỉ discovery URL; mọi kết quả vẫn qua rights/policy gate và không dùng để lách platform restriction |

### 4.4 Price Check

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
| PC-09 | OCR ưu tiên local khi capability/confidence đạt ngưỡng; chỉ offload server sau consent hoặc khi profile cấu hình yêu cầu |

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

### 4.6 Tourtect Lens

| ID | Yêu cầu |
| --- | --- |
| LV-01 | Context frame không vượt 1 FPS |
| LV-02 | Tự giảm FPS khi mạng yếu, pin yếu hoặc thiết bị nóng |
| LV-03 | High-resolution capture luôn cần thao tác người dùng |
| LV-04 | Cảnh báo đỏ cần OCR/user confirmation và Price Engine |
| LV-05 | Dừng capture ngay khi background, revoke permission hoặc end session |
| LV-06 | Local inference không upload frame; mọi frame/capture offload chỉ gửi sau consent và raw media server tuân TTL |

### 4.7 Scam và Emergency

| ID | Yêu cầu |
| --- | --- |
| SE-01 | Rule engine phát hiện critical signal trước LLM |
| SE-02 | Chỉ dùng playbook đã duyệt, có nguồn và ngày review |
| SE-03 | SOS dùng được offline và không cần tài khoản |
| SE-04 | Hotline không được LLM tạo hoặc sửa |
| SE-05 | Không tự gọi/chia sẻ vị trí/incident |
| SE-06 | Chế độ im lặng tắt TTS, rung và flash không cần thiết |

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

## 5. Kiến trúc hệ thống

### 5.1 System context

~~~mermaid
flowchart LR
    Tourist["Du khách"]
    Community["Cộng đồng địa phương / expat"]
    Merchant["Chủ cơ sở"]
    Reviewer["Moderator / Editor / Data steward"]
    Sources["Official API / RSS / partner / khảo sát"]

    Android["Tourtect Android\nForum + Live + Lens"]
    Web["Tourtect Responsive Web\nForum + Search + Places"]
    Zalo["Zalo Mini App Lite"]
    Admin["Admin / Moderation Web"]
    Platform["Tourtect Platform"]

    Tourist --> Android
    Tourist --> Web
    Community --> Android
    Community --> Web
    Merchant --> Web
    Reviewer --> Admin

    Android --> Platform
    Web --> Platform
    Zalo --> Platform
    Admin --> Platform
    Sources --> Platform
~~~

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

### 5.3 Lý do chọn modular monolith

Pilot cần phát triển nhanh và giữ transaction/provenance dễ kiểm soát. Core API là modular monolith với boundary rõ:

- Realtime được tách process vì tải WebSocket khác HTTP CRUD. Profile demo gọi FPT AI Factory và không vận hành GPU container local; kiến trúc vẫn cho phép provider khác hoặc on-device sau demo.
- Search indexing, external ingestion, notification và data worker chạy bất đồng bộ vì cần retry/quarantine.
- Các module dùng contract nội bộ, có thể tách service sau khi có số liệu tải thật.
- Không dùng microservice hoặc self-host container cho từng model trong giai đoạn hackathon; đây là giới hạn deployment demo, không phải khóa thiết kế dài hạn.

### 5.4 Trách nhiệm thành phần

| Thành phần | Trách nhiệm | Không được làm |
| --- | --- | --- |
| Identity/Profile | Registration, email verification, Google OIDC, provider linking, session rotation/revocation và public profile | Tin client token chưa verify, dùng email làm Google identity key hoặc tự cấp verified role |
| Content/Review Service | Post, comment, review, draft, merchant reply và version history | Tự quyết sanction hoặc điểm uy tín |
| Place/Topic Service | Canonical place, alias đa ngôn ngữ, geo, topic và price report | Cho merchant ghi đè community fact |
| Feed/Organic Ranking | Candidate, relevance, freshness, diversity và safety boost | Nhận spend/commission/business tier làm feature |
| Search | Full-text, geo, filter, autocomplete và index version | Index private incident hoặc PII bị redaction |
| Reputation | Điểm theo lĩnh vực, decay, evidence và abuse signal | Dùng upvote đơn thuần làm trust score |
| Moderation/Appeal | Policy engine, queue, sanction, notice, appeal và transparency log | Cho sales/advertiser sửa quyết định |
| External Ingestion | Connector, dedupe, entity linking, refresh và takedown | Bypass access control hoặc re-host trái quyền |
| Rights/Source Registry | Access mode, robots/terms/license, allowed field, retention, owner và kill switch | Mặc định cho phép nguồn chưa review |
| Adaptive Crawl Scheduler | Per-host budget, priority/freshness, jitter, conditional refresh, backoff và circuit breaker | Né block bằng proxy/IP rotation hoặc vượt quota |
| Sandbox Fetch Proxy | Egress allowlist, DNS/IP/MIME/size guard, stable bot identity và HTTP cache | Login, chạy script trình duyệt hoặc truy cập private endpoint |
| Monetization | Ad eligibility, subscription, business entitlement và disclosure | Ghi vào rank/review/alert/moderation |
| Realtime Gateway | Session, PCM/frame transport, event ordering, backpressure | Ra quyết định giá |
| Session Orchestrator | State machine, locale, role, mode, resume | Lưu raw media dài hạn |
| Adaptive Model Router | Chọn device/server/provider fallback theo execution policy, capability, consent và confidence | Thay đổi consent hoặc làm lộ provider credential |
| Translation Service | ASR, MT, glossary, critical token validation | Tư vấn an toàn tự do |
| Vision Service | Context classification, ROI, OCR fallback | Kết luận gian lận |
| Price Engine | Chuẩn hóa, chọn cohort, anomaly/confidence | Dùng LLM làm phép so sánh cuối |
| Scam/Safety Engine | Triage, pattern retrieval, safe action, escalation | Tạo hotline |
| Emergency Directory | Trả số đã xác minh theo vùng/incident | Tự gọi số |
| Consent Service | Consent ledger, TTL, deletion/export | Gộp consent contribution với xử lý phiên |
| Data Platform | Evidence, review, snapshot, rollback | Cho submission tác động ngay production |

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

### 5.8 Runtime container local bằng Podman

Demo dùng Podman và <code>podman compose</code>; không yêu cầu Docker daemon. File <code>compose.yaml</code> ở repository root khởi động các dependency stateful:

- PostgreSQL/PostGIS cho source of truth và geo.
- Redis cho cache/session/realtime state.
- MinIO làm object storage tương thích S3 trong local.
- OpenSearch cho full-text/geo/feed retrieval.

Trong profile demo, model không chạy trong stack local: backend gọi FPT AI Factory qua HTTPS. Khi source code backend/web được thêm, các service app sẽ được bổ sung vào cùng manifest hoặc file override, nhận cấu hình từ <code>.env</code> nhưng chỉ backend được nhận <code>FPT_AI_API_KEY</code>.

~~~bash
cp .env.example .env
# Điền API key và thay toàn bộ mật khẩu mẫu trong .env
podman compose up -d
podman compose ps
~~~

<code>.env.example</code> chỉ chứa tên biến và giá trị demo, được commit. <code>.env</code> chứa secret thật, bị loại khỏi Git. Cấu hình hiện tại là single-node dành cho máy phát triển/hackathon, không phải topology production.

---

## 6. Nền tảng dữ liệu forum, giá và scam

### 6.1 Community knowledge graph

Các thực thể lõi tạo thành graph <code>User → Post/Review/Report → Place/Topic → PriceObservation/ScamPattern → ExternalSource</code>. PostgreSQL/PostGIS là source of truth; OpenSearch là read model cho full-text/geo/feed retrieval và có thể rebuild từ event log.

#### Post

| Trường | Ý nghĩa |
| --- | --- |
| <code>post_id</code>, <code>author_id</code> | Opaque ID; public profile tách khỏi identity nhạy cảm |
| <code>post_type</code> | discussion, question, review, price_report, scam_report, tip, official_alert, external_link |
| <code>original_locale</code> | Ngôn ngữ gốc; translation là derivative riêng |
| <code>title</code>, <code>body</code> | Nội dung có version history và redaction state |
| <code>place_ids</code>, <code>topic_ids</code>, <code>region_id</code> | Liên kết entity; geo chính xác chỉ khi phù hợp |
| <code>evidence_level</code> | none, metadata, verified_receipt, verified_source |
| <code>commercial_disclosure</code> | none, invited, gifted, affiliate, employee, sponsored |
| <code>moderation_status</code> | draft, pending, published, limited, removed, appealed |
| <code>created_at</code>, <code>updated_at</code>, <code>event_time</code> | Phân biệt thời gian đăng và thời gian sự việc |

#### Review và Place

- <code>Review</code>: place, visit time, overall rating, price transparency, service, safety, value, evidence, disclosure và merchant reply.
- <code>Place</code>: canonical name, alias đa ngôn ngữ, category, region/geo, claim status, contact/public metadata và merge history.
- Điểm place là aggregate có Bayesian shrinkage/minimum count, freshness và distribution; không chỉ hiển thị trung bình đơn giản.
- Xóa review không làm mất audit event nhưng public payload/PII phải được tombstone theo policy.

#### Region, địa giới Việt Nam và OpenStreetMap

Từ 01/07/2025, mô hình chính quyền địa phương Việt Nam vận hành hai cấp: cấp tỉnh và cấp xã; cấp huyện kết thúc hoạt động. Dataset chuẩn phải theo Quyết định <code>19/2025/QĐ-TTg</code> về danh mục/mã đơn vị hành chính, gồm 34 đơn vị cấp tỉnh và 3.321 xã/phường/đặc khu. Tuy vậy, địa chỉ, bài báo, hóa đơn và thói quen người dùng vẫn chứa quận/huyện/tỉnh cũ, nên Tourtect không xóa lịch sử mà quản lý temporal region graph.

<code>Region</code> gồm:

- <code>region_id</code> nội bộ bất biến; <code>official_code</code> theo phiên bản danh mục; <code>admin_level</code> là province/commune/special_zone hoặc legacy_district.
- Tên chính thức, short name và alias tiếng Việt/Anh/Hàn/Trung/Nga; alias không dấu, Hán tự, Hangul, Cyrillic và tên trước sáp nhập.
- <code>valid_from</code>, <code>valid_to</code>, <code>predecessor_ids</code>, <code>successor_ids</code> và loại thay đổi merge/split/rename/boundary_adjustment.
- Geometry/version/source, centroid/representative point, parent hiện hành, legacy parent và mapping confidence.
- <code>osm_relation_id</code>/<code>osm_version</code> chỉ là cross-reference; không dùng OSM ID làm primary key nghiệp vụ.

Luồng cập nhật địa giới:

1. Theo dõi văn bản Chính phủ/Quốc hội/Bộ Nội vụ; ingest danh mục/mã chính thức thành candidate có checksum và effective date.
2. Data steward đối chiếu số lượng, mã, tên, predecessor/successor và ban hành <code>AdminBoundarySnapshot</code> có version.
3. Geometry lấy từ nguồn nền địa lý/địa giới chính thức khi có quyền. OSM dùng để bổ sung/cross-check và phục vụ map UX, không tự ghi đè official code/boundary.
4. Chạy spatial QA: geometry validity, gap/overlap, centroid, containment của place, coverage và diff diện tích. Thay đổi lớn cần dual review.
5. Re-link place/address bất đồng bộ; giữ URL redirect và alias cũ. Search “Hoàn Kiếm”, “quận Hoàn Kiếm” vẫn trả place đúng dù cấp huyện là legacy context.

Chiến lược OSM cho pilot:

| Nhu cầu | Cách làm | Không dùng public service theo cách nào |
| --- | --- | --- |
| Base POI/road/address | Bootstrap từ Vietnam regional <code>.osm.pbf</code>; lọc category du lịch bằng osm2pgsql/ imposm; giữ <code>osm_type/id/version/timestamp/tags</code> | Không enumerate toàn Việt Nam qua Overpass/Nominatim |
| Cập nhật | Apply replication diff theo sequence; pilot hourly, nightly rebuild derived place/search và weekly full consistency check | Không poll từng OSM object |
| Geocoding/search | Self-host Nominatim/Photon hoặc mua provider có SLA và quyền cache | Public Nominatim tối đa 1 req/s, không autocomplete/bulk |
| Map tiles | Provider OSM-derived có SLA hoặc self-host tile stack/CDN | Không dùng <code>tile.openstreetmap.org</code> làm tile backend production/offline prefetch |
| Routing | Self-host OSRM/Valhalla/GraphHopper hoặc provider | Không giả định OSM public server cung cấp routing SLA |
| On-map attribution | Hiển thị <code>© OpenStreetMap contributors</code> và link copyright theo guideline | Không bỏ attribution; review ODbL nếu phát hành derived database |

OSM POI là candidate, không phải verified business record. Place merge dùng name/alias, category, khoảng cách, phone/website và geometry với threshold theo category; trường hợp chuỗi cửa hàng, cổng khác nhau của một điểm du lịch hoặc place đã chuyển địa chỉ phải vào review. Merchant claim/official source có thể sửa business facts nhưng không làm mất provenance hoặc đóng góp OSM.

#### ExternalContent

- <code>external_content_id</code>, source platform, canonical URL và source content ID.
- Author/channel, published time, retrieved/last-checked time và locale.
- Permitted title/snippet/thumbnail/embed metadata, license/rights status và policy version.
- Entity links, duplicate cluster, fact/opinion/sponsored label và moderation status.
- Source state: active, changed, deleted, embed_disabled, takedown hoặc expired.
- Access mode: official_api, rss_atom, webhook, partner_feed, allowlisted_html, user_submitted_url, embed_only hoặc blocked.
- HTTP/source state: ETag, Last-Modified, content/perceptual hash, source updated time, next check, failure/backoff và deletion check time.

<code>SourcePolicy</code> tách khỏi content record: host/platform, owner, terms/robots URL và hash, allowed path/content/fields, display/retention rule, auth/quota, per-host budget, policy reviewer/expiry, complaint contact và kill-switch state. Connector chỉ chạy khi policy còn hiệu lực; “chưa review” tương đương deny.

<code>DiscoveryCandidate</code> lưu URL/query pack/source/first-seen/rank tối thiểu trong thời gian ngắn. Candidate không xuất hiện trong search/feed cho đến khi rights gate, canonicalization, entity-link và moderation hoàn tất.

<code>CommerceOfferObservation</code> lưu platform/merchant/item, displayed base price, promotion, delivery/service fee, membership/location condition, observed time và provenance. Nó không được Price Engine coi là giá tại cửa hàng nếu chưa có mapping/unit và bằng chứng độc lập.

### 6.2 Ba lớp dữ liệu giá

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

### 6.3 Phân tầng nguồn

| Tầng | Nguồn | Quyền ảnh hưởng |
| --- | --- | --- |
| A | Biểu phí/trần giá/tỷ giá/cảnh báo có thẩm quyền | Luật xác định hoặc nguồn neo |
| B | Khảo sát thực địa, mystery shopper, trung tâm hỗ trợ, đối tác được kiểm toán | Nguồn chính của reference model |
| C | Menu/website/bảng giá/OTA/hãng vận tải công khai | Giá chào bán sau chuẩn hóa |
| D | Contribution có ảnh và consent | Quarantine, cần xác nhận chéo |
| E | Báo chí, forum, review/social | Chỉ phát hiện xu hướng/pattern mới |

Tầng E không trực tiếp thay đổi giá hoặc playbook.

### 6.4 Mô hình dữ liệu giá và an toàn

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

### 6.5 Chuẩn hóa theo vertical

| Vertical | Thuộc tính bắt buộc |
| --- | --- |
| Taxi/đưa đón | Khoảng cách, thời gian, loại xe, giờ cao điểm, sân bay, cầu đường, phí chờ |
| Đổi tiền | Cặp tiền, timestamp, số đưa/nhận, phí công khai, effective rate |
| Ăn uống | Món/SKU, khẩu phần, số lượng, thuế, service charge, loại địa điểm |
| Tour | Thời lượng, private/group, ngôn ngữ, lịch trình, vé/xe/bữa ăn, mùa |

### 6.6 Thuật toán cảnh báo

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

### 6.7 Cold start không phụ thuộc đối tác

- Seed dữ liệu có provenance từ nguồn công khai chính thức.
- Khảo sát nhỏ phân tầng tại ba cụm pilot.
- Demo data synthetic phải gắn nhãn rõ, không trộn vào production snapshot.
- Taxi và đổi tiền được ưu tiên vì dễ chuẩn hóa hơn.
- Tour có quyền abstain cao vì gói khó so sánh.
- Không giả định có API từ cơ quan du lịch, cơ quan thuế hoặc đối tác.

### 6.8 Chống poisoning

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
- Đọc public content và dùng SOS/AI session riêng tư không yêu cầu account; dùng anonymous/short-lived session token.
- Đăng bài, bình luận, vote, follow, report, subscription và merchant action yêu cầu account; endpoint có idempotency key và rate limit phù hợp.
- Admin dùng SSO/RBAC và audit log.
- Raw media upload dùng signed URL, content type allowlist và size limit.

### 7.2 Endpoint

| Method | Path | Mục đích |
| --- | --- | --- |
| POST | <code>/v1/auth/registrations</code> | Tạo account email/password và gửi verification |
| POST | <code>/v1/auth/email-verifications</code> | Xác minh email bằng token một lần |
| POST | <code>/v1/auth/sessions</code> | Đăng nhập email/password và tạo Tourtect session |
| POST | <code>/v1/auth/sessions/refresh</code> | Rotate refresh token, phát hiện reuse |
| DELETE | <code>/v1/auth/sessions/{session_id}</code> | Revoke một phiên/thiết bị |
| DELETE | <code>/v1/auth/sessions</code> | Đăng xuất tất cả phiên |
| POST | <code>/v1/auth/password-resets</code> | Yêu cầu email reset, response không làm lộ account tồn tại |
| POST | <code>/v1/auth/password-resets/confirm</code> | Đổi mật khẩu bằng token một lần và revoke session khác |
| POST | <code>/v1/auth/oauth/google/attempts</code> | Tạo state/nonce/PKCE auth attempt và authorization URL |
| GET/POST | <code>/v1/auth/oauth/google/callback</code> | Verify callback/code và tạo/link Tourtect session |
| POST | <code>/v1/auth/identities/google/link</code> | Link Google với account đang đăng nhập sau re-authentication |
| DELETE | <code>/v1/auth/identities/{identity_id}</code> | Unlink provider nếu account còn phương thức đăng nhập khác |
| POST | <code>/v1/auth/anonymous-merge</code> | Preview/confirm merge saved item, draft và preference |
| GET | <code>/v1/account/sessions</code> | Liệt kê phiên gần đây và metadata thiết bị tối thiểu |
| GET | <code>/v1/feed</code> | Feed Following/Nearby/Latest/Trending/Safety với reason code |
| GET | <code>/v1/search</code> | Full-text/geo search place, post, giá, scam và external content |
| POST/GET | <code>/v1/posts</code> | Tạo draft/xuất bản và đọc danh sách post |
| GET/PATCH/DELETE | <code>/v1/posts/{id}</code> | Đọc, sửa có version hoặc yêu cầu xóa post |
| POST/GET | <code>/v1/posts/{id}/comments</code> | Bình luận và thread |
| POST | <code>/v1/posts/{id}/votes</code> | Đánh dấu hữu ích; không đồng nghĩa evidence |
| POST | <code>/v1/posts/{id}/reports</code> | Báo vi phạm/an toàn/PII |
| POST | <code>/v1/follows</code> | Theo dõi place/topic/user |
| GET | <code>/v1/places/{id}</code> | Place page aggregate |
| POST/GET | <code>/v1/places/{id}/reviews</code> | Review có cấu trúc và danh sách review |
| POST | <code>/v1/places/{id}/claim</code> | Merchant claim; không đổi review/ranking |
| POST/GET | <code>/v1/price-reports</code> | Price report công khai và evidence workflow |
| POST/GET | <code>/v1/scam-reports</code> | Scam report, safety triage và moderation |
| GET | <code>/v1/external-content</code> | External card đã qua rights/policy gate |
| POST | <code>/v1/external-content/submissions</code> | Gửi URL để connector kiểm tra; không fetch tùy ý từ client |
| POST/GET | <code>/v1/moderation/appeals</code> | Tạo và theo dõi appeal |
| GET/PATCH | <code>/v1/notifications</code> | Danh sách và trạng thái notification |
| POST/GET | <code>/v1/subscriptions</code> | Tourtect Plus và entitlement |
| GET/PATCH | <code>/v1/business-profiles/{place_id}</code> | Business tools sau claim/verification |
| POST | <code>/v1/affiliate-events</code> | Ghi sự kiện disclosure/click tối thiểu, chống giả mạo |
| POST | <code>/v1/realtime/sessions</code> | Tạo Live Voice/Camera session |
| POST | <code>/v1/realtime/sessions/{id}/resume</code> | Cấp token resume ngắn hạn |
| WS | <code>/v1/realtime/sessions/{id}/events</code> | PCM, context frame và realtime event |
| POST | <code>/v1/realtime/sessions/{id}/captures</code> | Tạo capture_id và signed PUT URL |
| PUT | <code>{signed_upload_url}</code> | Upload capture đã redaction từ local path hoặc capture tạm mã hóa khi dùng server path |
| POST | <code>/v1/realtime/sessions/{id}/captures/{capture_id}/finalize</code> | Kiểm tra hash/MIME/redaction metadata; server path tiếp tục OCR/redaction trước khi capture thành ready |
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

type Channel = "android_app" | "responsive_web" | "zalo_mini_app" | "admin_web";
type PostType =
  | "discussion"
  | "question"
  | "review"
  | "price_report"
  | "scam_report"
  | "tip"
  | "official_alert"
  | "external_link";
type EvidenceLevel = "none" | "metadata" | "verified_receipt" | "verified_source";
type ModerationStatus = "draft" | "pending" | "published" | "limited" | "removed" | "appealed";
type CommercialDisclosure = "none" | "invited" | "gifted" | "affiliate" | "employee" | "sponsored";
type RevenueSurface = "contextual_ad" | "affiliate" | "plus" | "business" | "sponsored" | "b2b";
type LiveMode = "voice" | "camera";
type SpeakerRole = "tourist" | "local";
type ExecutionLocation = "device" | "server";
type ExecutionPolicy = "adaptive" | "server_only" | "local_only";
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

type AccountStatus = "pending_email_verification" | "active" | "suspended" | "scheduled_for_deletion";
type IdentityProvider = "password" | "google";

interface Account {
  account_id: string;
  display_name: string;
  primary_email_masked: string;
  email_verified: boolean;
  status: AccountStatus;
  locale: Locale;
  created_at: string;
}

interface FederatedIdentity {
  identity_id: string;
  account_id: string;
  provider: "google";
  issuer: string;
  subject: string;
  email_at_link_time_masked?: string;
  linked_at: string;
}

interface AccountSession {
  session_id: string;
  account_id: string;
  device_label?: string;
  created_at: string;
  last_seen_at: string;
  expires_at: string;
  revoked_at?: string;
}

interface OAuthAttempt {
  attempt_id: string;
  provider: "google";
  state_hash: string;
  nonce_hash: string;
  pkce_challenge: string;
  redirect_uri_id: string;
  expires_at: string;
  consumed_at?: string;
}

interface Post {
  post_id: string;
  author_id: string;
  post_type: PostType;
  original_locale: Locale;
  title: string;
  body: string;
  place_ids: string[];
  topic_ids: string[];
  region_id?: string;
  evidence_level: EvidenceLevel;
  commercial_disclosure: CommercialDisclosure;
  moderation_status: ModerationStatus;
  created_at: string;
  updated_at: string;
}

interface Review {
  review_id: string;
  post_id: string;
  place_id: string;
  visited_at?: string;
  overall_rating: number;
  price_transparency_rating?: number;
  service_rating?: number;
  safety_rating?: number;
  value_rating?: number;
  evidence_level: EvidenceLevel;
  commercial_disclosure: CommercialDisclosure;
}

interface ExternalContent {
  external_content_id: string;
  platform: string;
  canonical_url: string;
  source_content_id?: string;
  original_locale?: string;
  published_at?: string;
  last_checked_at: string;
  rights_status: "embed_allowed" | "metadata_only" | "partner_licensed" | "blocked";
  source_state: "active" | "changed" | "deleted" | "embed_disabled" | "takedown" | "expired";
  place_ids: string[];
  topic_ids: string[];
}

interface ReputationProfile {
  user_id: string;
  local_knowledge: number;
  price_evidence: number;
  translation: number;
  safety: number;
  last_recalculated_at: string;
}

interface RealtimeSessionConfig {
  channel: "android_app";
  mode: LiveMode;
  tourist_locale: Locale;
  local_locale: "vi-VN";
  region_id: string;
  execution_policy: ExecutionPolicy;
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
| Mạng tốt | Full server PTT/STT/MT/TTS | 0.5–1 FPS + server VLM/OCR | Realtime server | Full directory |
| Mạng yếu | Audio ưu tiên | 0.2 FPS hoặc tắt | Manual/cached | Full offline pack |
| Server vision/OCR lỗi | Không ảnh hưởng dịch | Manual crop/input, không AI | Nếu user nhập đủ dữ kiện | Không ảnh hưởng |
| Server MT lỗi | Caption transcript gốc/phrase | Camera vẫn chạy | Vẫn kiểm tra được input manual | Phrasebook |
| Offline | Phrase/audio pack | Preview local, không AI vision | Cache/manual, không red alert nếu stale | Hotline + incident card |
| Thiết bị nóng/pin yếu | Giảm audio effect, vẫn server ASR | Giảm/tắt frame sampling | Manual | Không ảnh hưởng |

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
| Publish public post/review/report | Khi bấm xuất bản, kèm preview phạm vi công khai | Tắt; không suy ra từ consent xử lý AI |
| Personalized feed/ads | Khi bật cá nhân hóa; contextual feed/ads vẫn dùng được nếu tắt | Tắt |
| Marketing/affiliate analytics | Trước khi dùng dữ liệu ngoài đo lường thiết yếu | Tắt |

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
| Password hash | Trong vòng đời account; xóa khi account bị xóa | Không áp dụng cho account chỉ có Google identity |
| Google identity | <code>issuer + subject</code>, email masked và thời điểm link | Không lưu Google access/refresh token khi chỉ dùng authentication |
| Auth session | Refresh credential hash và metadata thiết bị tối thiểu tới khi revoke/expire | Không lưu raw token; IP/user-agent chi tiết tuân theo security TTL |
| Verification/reset/OAuth attempt | Hash token/state/nonce, trạng thái dùng và expiry | Xóa payload hết hạn theo TTL; không dùng lại |

UGC public tuân theo chính sách riêng với media phiên AI:

- Post/review/comment lưu cho tới khi tác giả xóa, nền tảng gỡ hoặc hết retention theo policy; edit tạo version/audit phù hợp nhưng nội dung công khai cũ không còn được phục vụ.
- Người dùng có export/delete account. Nội dung đã xóa được loại khỏi API, search, CDN và cache theo SLA; audit chỉ giữ opaque ID, action, reason và non-reversible hash khi thực sự cần.
- Scam report mặc định giảm độ chính xác vị trí, tự động tìm PII/khuôn mặt/biển số và yêu cầu xác nhận trước khi công khai.
- Merchant claim lưu hồ sơ xác minh tách khỏi public profile và dùng access control chặt hơn.
- External content mặc định chỉ lưu metadata/snippet/embed được phép; khi nguồn xóa hoặc takedown, card bị disable và index được cập nhật.

Contribution được liên kết bằng deletion token ngẫu nhiên, không bằng danh tính. Hash của token được lưu tách khỏi dataset để người dùng có thể rút contribution. Khi rút consent, evidence/source link và candidate bị xóa khỏi pipeline trong tối đa 30 ngày, snapshot tương lai được tái tính; snapshot đã công bố chỉ chứa aggregate không định danh và được supersede bằng phiên bản mới khi cần.

### 9.3 Threat model

| Rủi ro | Biện pháp |
| --- | --- |
| Prompt injection từ menu/transcript | Coi OCR/transcript là dữ liệu, không phải instruction; tool allowlist và schema validation |
| LLM tạo hotline | Hotline chỉ từ Emergency Directory đã ký |
| Data poisoning/Sybil | Quarantine, device token, dedup, source caps, robust estimator |
| Fake review/brigading | Account/device/graph/rate anomaly, disclosure, reputation theo lĩnh vực, review queue và reversible sanction |
| Merchant retaliation/defamation | Khử PII, đánh giá hành vi/giao dịch, right-to-reply, evidence level, notice/appeal và legal escalation policy |
| Advertiser influence | RBAC tách Sales/Ads/Moderator, schema boundary, ranking feature allowlist và audit thử nghiệm |
| Malicious external URL/SSRF | Connector allowlist, URL normalization, DNS/IP egress policy, fetch proxy sandbox và content-size/MIME limit |
| Copyright/takedown | Rights registry, permitted metadata/embed, canonical attribution, refresh/deletion sync và takedown SLA |
| URL/file độc hại | Signed upload, MIME sniffing, size limit, malware scan, image re-encode |
| Session hijack | Short-lived token, bind device/session, TLS, rate limit |
| Credential stuffing | Rate limit theo account/IP/risk, breached-password check, generic error, progressive challenge và security alert |
| OAuth CSRF/replay/code interception | State, nonce, PKCE, exact redirect allowlist, one-time auth attempt và server-side token validation |
| Account linking takeover | Re-authenticate account hiện tại, khóa theo issuer+sub, không auto-link bằng email và thông báo khi link/unlink |
| Refresh-token theft | Hash-at-rest, rotation, reuse detection, Android Keystore hoặc HttpOnly cookie và revoke session family |
| Account enumeration | Response/timing gần tương đương cho registration, login và password reset; email gửi ngoài luồng |
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
| Community | DAU/MAU, search success, answer rate, useful contribution, return rate, save/follow và place coverage |
| Trust & Safety | Spam/fake-review precision-recall, report rate, appeal overturn rate, moderation SLA, harassment/PII exposure |
| Content ingestion | Connector success, P0/P1 freshness SLO, dedupe precision, entity-link accuracy, 304 ratio, bytes/new-item, 429/403/5xx, robots/policy drift, deletion/takedown propagation |
| Revenue | Ad fill/eCPM, affiliate conversion, Plus conversion/churn, business retention và revenue concentration |
| Trust firewall | Organic rank parity, sponsored disclosure accuracy, advertiser-policy violations và revenue-triggered moderation changes |

Metrics được slice theo:

- Ngôn ngữ.
- Thành phố/cụm.
- Vertical.
- Thiết bị/tier.
- Execution policy/location, model/provider/version và fallback path.
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
| Demo execution profile | Với env demo, 100% ModelTrace là <code>server</code>; mọi build không chứa FPT credential. Adaptive build riêng phải đạt test signature/capability/consent |
| Red alert precision | ≥ 95% |
| False-positive trên giá hợp lệ | ≤ 2% |
| Critical safety escalation | 100% trên golden safety set |
| Unsafe confrontation advice | 0 trường hợp trên safety set |
| Critical translation fields | Không sai phủ định, thương tích, số tiền, tiền tệ, vị trí, biển số, số người |
| Background capture | Dừng 100% khi background/end/revoke |
| Offline SOS | Mở hotline/incident card tối đa 2 thao tác |
| Poisoning simulation | Attack budget 1–500; p50/p90 dịch ≤ 5% và không lật quyết định high-risk |
| Fake-review/brigading | Đạt ngưỡng precision/recall đã chốt theo locale; không auto-ban chỉ từ một model signal |
| Commercial firewall | 100% test chứng minh spend, commission và business tier không đổi organic rank/review/alert/moderation |
| Sponsored disclosure | 100% ad/affiliate/sponsored item có nhãn dễ thấy và machine-readable |
| External content deletion | Disable serving và loại khỏi search theo SLA khi source/takedown state thay đổi |
| P0 source freshness | 95% item từ API/feed được discover trong 15 phút khi source hoạt động; HTML-only theo SLO riêng đã duyệt |
| Crawler politeness | Không fetch path disallow; không vượt per-host budget; 429 rate < 1%; mọi 403/451 mở circuit và policy review |
| Administrative snapshot | 100% official code/name/effective date khớp nguồn; không có gap/overlap nghiêm trọng; legacy alias resolve đúng |
| Public PII | Không để lọt PII nghiêm trọng trong golden/red-team public scam report set |

Gate giá được đánh giá theo từng vertical × vùng, không chỉ toàn cục. Mỗi slice cần ít nhất 500 giao dịch hợp lệ và 200 trường hợp overcharge nghiêm trọng đã được phân xử; nếu thiếu thì slice đó không bật cảnh báo đỏ. Point estimate và one-sided Wilson 95% confidence bound đều phải đạt gate: lower bound của precision ≥ 95% và upper bound của false-positive rate ≤ 2%. Nếu chưa đạt, chỉ trả <code>elevated</code> hoặc <code>insufficient_data</code>.

Poisoning test chạy với ngân sách 1, 10, 50, 100 và 500 submission phối hợp; đo dịch chuyển <code>p50</code>, <code>p90</code>, confidence và tỷ lệ lật quyết định alert. Gate yêu cầu cả <code>p50</code> và <code>p90</code> không dịch quá 5%, đồng thời không có price case chuẩn nào bị đổi từ <code>typical</code> sang <code>high_risk</code> hoặc ngược lại.

### 10.3 Test matrix

#### Forum, review và moderation

- Đăng ký, email verification hết hạn/dùng lại, login đúng/sai và password reset không làm lộ account.
- Google callback sai <code>state</code>/<code>nonce</code>/<code>aud</code>/<code>iss</code>, token hết hạn, PKCE sai, code replay và redirect URI ngoài allowlist.
- Google email trùng account password, link/unlink sau re-authentication, Google subject đổi email và account còn/không còn login method.
- Refresh-token rotation/reuse, revoke một thiết bị, logout all, session hết hạn và Android Keystore/cookie bị xóa.
- Anonymous merge preview/confirm/rollback; saved item/draft được merge nhưng incident/raw media không bị chuyển.
- Account mới spam hàng loạt, copy post, review exchange và vote ring.
- Brigading theo nhóm/ngôn ngữ, Sybil nhiều account và merchant tự review.
- Review tiêu cực hợp lệ bị report hàng loạt; merchant phản hồi/appeal nhưng không được ép xóa.
- Nội dung có PII, biển số, khuôn mặt, cáo buộc chưa có bằng chứng, harassment và doxxing.
- Sửa/xóa post, merge place, block user, notification race và search index eventual consistency.
- Dịch/tóm tắt sai số tiền, phủ định hoặc biến ý kiến thành sự thật.

#### External content

- Canonical URL có tracking, redirect loop, duplicate/cross-post và cùng sự kiện khác nguồn.
- RSS/API timeout, quota/rate limit, source sửa/xóa, embed bị disable và takedown.
- ETag/Last-Modified đúng/sai/thiếu; 304 không chạy extraction; jitter tránh burst; 429 có/không có Retry-After và 403/451 đóng circuit.
- Robots/ToS thay đổi giữa hai lần chạy, kill switch, abuse complaint và source owner opt-out.
- Commerce offer có voucher, membership, delivery/location fee và giá cá nhân hóa không bị nhập thành giá tại chỗ.
- Query tiếng Việt/Anh/Hàn/Trung/Nga có alias địa danh, phủ định, từ đồng âm và spam SEO; đo precision theo query pack.
- URL trỏ private IP, file quá lớn, MIME giả, prompt injection và HTML/script độc hại.
- Article/video sponsored nhưng metadata thiếu disclosure; classifier không được tự khẳng định tuyệt đối.
- Entity linking nhầm địa điểm trùng tên hoặc gán sai thời gian/khu vực.
- OSM diff mất sequence/replay/trùng; POI merge sai chi nhánh; official boundary đổi nhưng legacy address vẫn resolve.

#### Monetization và trust firewall

- Cùng một query/feed seed phải giữ thứ tự organic khi thay advertiser spend, commission hoặc business tier.
- Ads không xuất hiện ở SOS, critical incident, silent mode hoặc nội dung chưa đủ ad eligibility.
- Affiliate/sponsored card luôn có disclosure trước CTA và screen reader đọc được.
- Merchant hết hạn gói vẫn giữ public reply/history nhưng mất feature trả phí, không làm đổi rating.
- Sales/ad token không gọi được moderation, ranking feature, PriceSnapshot hoặc ScamPattern write API.
- Trust-health gate tắt được một ad surface mà forum, SOS và alert vẫn hoạt động.

#### Live Voice

- Khách Hàn nói số tiền, người Việt nói phụ phí hợp lệ.
- Khách Trung code-switch tên món tiếng Việt.
- Khách Nga nói trong tiếng ồn đường phố.
- Khách Anh hỏi lại và sửa transcript.
- Hai người bấm nhầm vai.
- TTS locale thiếu hoặc audio route đổi sang Bluetooth.
- STT local/server confidence thấp, adaptive fallback, timeout và rate-limit.
- System/server TTS lỗi hoặc thiếu locale; app chuyển caption + phrase audio đã duyệt.
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

#### Android app

- Cold start, process recreation, rotation/configuration change và navigation restore không khởi động lại capture ngoài ý muốn.
- Camera/micro permission grant/deny/revoke, app background/foreground, screen lock và session end dừng resource đúng lifecycle.
- Mạng đổi Wi-Fi/4G/2G, WebSocket reconnect/resume, event duplicate/out-of-order và backpressure ưu tiên audio/control hơn frame.
- Android Keystore entry mất/invalidated, refresh-token reuse, Room migration, cache TTL cleanup và WorkManager retry.
- Bluetooth/wired/speaker audio route, audio focus interruption và system/server TTS lỗi.
- Static scan APK/AAB không thấy <code>FPT_AI_API_KEY</code>; flavor demo không có model file/runtime, adaptive flavor chỉ chấp nhận signed model pack/runtime allowlist.

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

- Camera permission được duyệt/từ chối trên Zalo Android thật.
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

### 11.1 Demo 1 — Forum đa ngôn ngữ theo địa điểm

Khách Hàn tìm một khu phố Hà Nội, đọc place page bằng tiếng Hàn và thấy review, khoảng giá, price report, scam pattern, câu hỏi cộng đồng và nguồn báo/video. Người dùng chuyển về nguyên bản tiếng Việt, xem evidence/freshness và theo dõi địa điểm. Demo chứng minh forum là lõi, không phải màn hình phụ của AI.

### 11.2 Demo 2 — Price report tạo knowledge có kiểm soát

Khách chụp hóa đơn bằng Tourtect Lens, sửa một trường OCR, nhận so sánh riêng tư rồi chủ động tạo draft price report. App khử PII, yêu cầu disclosure/phạm vi công khai, moderation gắn evidence badge và chỉ đưa observation vào quarantine. Place page cập nhật post ngay khi hợp lệ nhưng reference snapshot chỉ đổi sau pipeline review/versioning.

### 11.3 Demo 3 — Tổng hợp báo và video đúng quyền

Editor nhập một RSS item và một YouTube URL. Connector kiểm tra quyền/embed, lấy metadata, canonicalize, dedupe, gắn với place/scam topic và tạo external card có attribution. Khi mô phỏng video bị disable, card ngừng phát và search index được cập nhật.

### 11.4 Demo 4 — Tourtect Live là trợ lý của forum

Khách Trung dùng PTT với người bán Việt Nam. Translation trả trước, intelligence lane tạo price candidate kín. Sau khi kết thúc, app chỉ đề nghị “Lưu riêng” hoặc “Tạo báo cáo”; không tự đăng transcript/media. SOS vẫn mở dialer không cần tài khoản hoặc gói Plus.

### 11.5 Demo 5 — Doanh thu không mua được niềm tin

Feed hiển thị một quảng cáo contextual có nhãn, một affiliate CTA có disclosure và một place đã mua Verified Business. Bật/tắt spend/business tier không làm đổi thứ tự organic, điểm review, price alert hay moderation. Chuyển sang incident critical thì toàn bộ quảng cáo biến mất.

### 11.6 Dữ liệu demo

- Dữ liệu thật có URL/provenance và ngày xác minh.
- Dữ liệu synthetic gắn badge “Demo data”.
- Synthetic data không được đưa vào snapshot pilot/production.
- Không gọi thử hotline khẩn cấp trong demo.

---

## 12. Lộ trình sau tài liệu

### Giai đoạn 0 — System Design

- Hoàn thiện tài liệu này.
- Khóa Android native cho demo; iOS không nằm trong backlog thực thi hiện tại. Demo bật server-only bằng env, không thay đổi adaptive model architecture.
- Chốt brand Tourtect, content policy, community guidelines, privacy, safety, rights registry và commercial trust charter.
- Review model/source licenses, moderation escalation, takedown và merchant appeal.

### Giai đoạn 1 — Forum foundation

1. Account/profile, anonymous browse, place/topic graph và public SEO pages.
2. Post/comment/review/price report/scam report cùng moderation/report/appeal tối thiểu.
3. Feed/search, translation view, notification và admin queue.
4. Evidence/provenance, PriceSnapshot seed và verified emergency directory.
5. Rights-aware connector cho nguồn chính thức Việt Nam, RSS/news và YouTube; source registry, adaptive scheduler, dedupe và takedown sync.

### Giai đoạn 2 — Safety AI và community loop

- Android Tourtect Lens → draft price report và Tourtect Live → private price candidate; demo dùng server-only qua env, bản thiết kế tiếp tục benchmark adaptive local/server.
- Reputation v1, anti-spam/fake-review signals và merchant reply/claim.
- Zalo Mini App feed lite/quick report/openPhone.
- Benchmark dịch bốn cặp ngôn ngữ, OCR và red-alert precision.

### Giai đoạn 3 — Pilot 90 ngày

- Seed local contributors, khảo sát thực địa, content editor và community moderator.
- Vận hành query pack Việt/Anh/Hàn/Trung/Nga về Việt Nam; đo coverage/precision trước khi mở thêm nguồn social.
- Bootstrap OSM Việt Nam, hourly diff, official admin snapshot và legacy place/address aliases.
- Golden data, calibration, signed safety pack, hotline verification và abuse drill.
- Closed beta với khách sạn/hướng dẫn viên.
- Privacy/security/rights review, merchant appeal và incident drill mô phỏng.
- Chỉ bật contextual ads/affiliate khi ad eligibility và trust-health gate đạt chuẩn.

### Giai đoạn 4 — Mở rộng có điều kiện

- Tourtect Plus, Verified Business tools rồi B2B aggregate insights.
- Partner data/API và thêm connector sau policy review từng nền tảng.
- Xin partner/API cho commerce và social theo thứ tự tác động: transport/merchant price → Reddit → Meta/TikTok creator opt-in; không xây browser scraper để lấp khoảng trống quyền.
- Human interpreter hoặc VoIP bridge chỉ khi có SLA, legal review và consent flow.
- Thêm <code>zh-Hant</code>.
- Mở rộng địa lý sau khi từng vertical đạt release gate.

---

## 13. Rủi ro và phương án

| Rủi ro | Ảnh hưởng | Phương án |
| --- | --- | --- |
| Forum ít nội dung/cold start | Cao | Seed place/topic có kiểm chứng, local contributor program, Q&A theo vùng và external card đúng quyền |
| Review giả/brigading/Sybil | Cao | Reputation đa chiều, graph/rate anomaly, evidence badge, quarantine, human review và rollback |
| Cáo buộc gây hại/merchant retaliation | Cao | Review hành vi/giao dịch, PII redaction, right-to-reply, notice/appeal, legal escalation và audit |
| Fetch nội dung vi phạm quyền hoặc bị xóa | Cao | Official API/RSS/partner, rights registry, metadata/embed tối thiểu, refresh/takedown sync |
| Platform khóa bot/API hoặc đổi điều khoản | Cao | Deny-by-default source policy, policy expiry, kill switch, partner route và không phụ thuộc một platform |
| Crawler gây tải/bị block | Cao | Stable identity/IP, per-host budget, conditional GET, jitter/backoff, circuit breaker và abuse contact |
| Giá commerce bị hiểu là giá tại chỗ | Cao | CommerceOfferObservation riêng; tách promo/fee/location và cần evidence độc lập trước Price Engine |
| Địa giới/địa chỉ cũ sau sắp xếp | Cao | Official versioned snapshot, temporal region graph, legacy alias/redirect và spatial QA |
| OSM POI sai/cũ hoặc license violation | Trung bình | Provenance/version, merchant/official verification, diff QA, attribution và ODbL review |
| Social signal thiên lệch theo ngôn ngữ/platform | Cao | Query pack bản ngữ, source diversity cap, dedupe cross-post, calibration và human review |
| Quảng cáo làm mất niềm tin | Cao | Commercial trust firewall, disclosure, ranking feature allowlist, no-ads safety surface và trust-health gate |
| Một nguồn doanh thu chi phối | Trung bình | Theo dõi revenue concentration, đa dạng Plus/business/B2B/grant và governance độc lập |
| Model FPT AI Factory dịch sai ngữ cảnh khẩn cấp | Cao | Glossary, critical validator, phrase pack, benchmark theo model ID/version và rule đứng trước LLM |
| FPT AI Factory lỗi/rate-limit/latency cao | Cao | Timeout, retry có backoff, circuit breaker, budget/rate cap và degraded/manual mode |
| Profile demo server-only làm Live/Lens phụ thuộc mạng | Cao | PTT/frame nhỏ, backpressure, timeout/circuit breaker, manual input, caption/phrase và offline safety pack |
| STT không đạt realtime hoặc không hỗ trợ đủ locale | Trung bình | Benchmark model Marketplace đã chọn, giới hạn utterance, caption/manual và phrase pack |
| VLM đọc sai số | Cao | Không dùng VLM làm source of truth; high-res OCR + user confirmation |
| Dữ liệu giá mỏng | Cao | Abstain, mở rộng cohort có disclosure, khảo sát tập trung |
| Crowdsourcing bị đầu độc | Cao | Quarantine, source cap, dedup, review, rollback |
| Live làm nóng máy/tốn data | Trung bình | PTT, adaptive FPS, burst video, thermal/network policy |
| Android fragmentation/lifecycle làm rò camera/micro | Cao | CameraX lifecycle binding, foreground-only capture, permission revoke test và device matrix |
| TTS không có đủ locale | Trung bình | Caption và audio phrase đã đóng gói |
| Người dùng hiểu “Live” là gọi người thật | Trung bình | Tên “Live Voice/Live Camera với AI”, onboarding rõ |
| Zalo Mini App runtime khác Android native | Trung bình | Lite scope, test Android thật, không cam kết realtime |
| Hotline thay đổi | Cao | Registry có nguồn/ngày xác minh, signed pack, expiry |
| Cảnh báo làm xung đột leo thang | Cao | Private card, haptic, safe wording, no public TTS |

---

## 14. Giả định đã khóa

- Tên sản phẩm là **Tourtect**.
- Forum/community knowledge graph là nền tảng; AI phiên dịch, Lens, price/scam intelligence là mô-đun hỗ trợ.
- Responsive web và Android app đều là sản phẩm đầy đủ cho forum; Android có thêm Live/Lens/offline/SOS. iOS tạm ngoài phạm vi demo.
- Zalo Mini App là kênh lite cho discovery, quick report, snapshot và hotline.
- Đọc public content không cần tài khoản; đăng/tương tác cần account. SOS và AI session riêng tư không bắt buộc account.
- V1 hỗ trợ account email/password và Google Sign-In; Google chỉ dùng authentication scope <code>openid email profile</code>, không mặc định xin quyền Google API khác.
- External content chỉ qua official API, RSS/Atom, partner feed hoặc shared URL trong phạm vi được phép; không sao chép toàn văn/re-host video mặc định.
- Việt Nam là geo scope và nguồn xác minh chính; Trung/Hàn/Nga trước mắt là market/locale discovery về trải nghiệm tại Việt Nam, chưa phải mở rộng POI toàn cầu.
- Không crawl GrabFood/ShopeeFood/Shopee hoặc public social bằng browser automation nếu chưa có chấp thuận; commerce/social ưu tiên partner API, creator opt-in và user-submitted URL.
- Google Custom Search JSON API không phải dependency; search API nếu bổ sung chỉ làm discovery và không thay rights gate.
- Địa giới dùng danh mục/mã chính thức Việt Nam có version; OSM là base map/POI/cross-check có attribution, không là nguồn pháp lý.
- Organic ranking, review, moderation, Price Engine và Safety Engine độc lập hoàn toàn với spend/commission/business tier.
- Không paywall SOS, hotline, báo cáo an toàn và cảnh báo thiết yếu; không có pay-to-remove.
- “Call” nghĩa là phiên nói với AI trong app, không phải cuộc gọi PSTN.
- Live Voice chỉ dùng push-to-talk theo vai.
- Live Camera là AI camera assist, không kết nối video tới người khác.
- Model server-side của profile demo gọi FPT AI Factory bằng API tương thích OpenAI; Gemini Live chỉ là mẫu tham khảo.
- Kiến trúc model là adaptive device/server. Riêng env demo đặt <code>AI_EXECUTION_MODE=server_only</code>; khi API không khả dụng Android dùng degraded/manual/offline safety pack.
- Dependency local chạy bằng Podman Compose; không yêu cầu Docker daemon hoặc GPU local.
- Android system TTS là mặc định của thiết kế; demo có thể ép server TTS theo execution policy và offline luôn có phrase audio đã duyệt.
- Price insight luôn riêng tư bằng card/rung.
- Pilot có thể chạy bằng dữ liệu seed/cộng đồng mà không phụ thuộc đối tác thương mại; connector ngoài phải có quyền truy cập hợp lệ.
- Contribution chỉ có hiệu lực khi user opt-in riêng.
- Chinese MVP là Simplified Chinese; Traditional Chinese nằm ngoài pilot.
- Không có human operator, VoIP bridge, auto-call hoặc auto-share location trong V1.

---

## 15. Nguồn tham khảo chính

### Đăng nhập và Google Identity

- [Google Identity Services — Sign in with Google overview](https://developers.google.com/identity/gsi/web/guides/overview)
- [Google — OpenID Connect server flow](https://developers.google.com/identity/openid-connect/openid-connect)
- [Google OpenID Connect API Reference](https://developers.google.com/identity/openid-connect/reference)

### External content, community và quảng cáo

- [YouTube Data API — API Reference](https://developers.google.com/youtube/v3/docs)
- [YouTube Data API — Search: list và bộ lọc embeddable/license](https://developers.google.com/youtube/v3/docs/search/list)
- [YouTube Data API — Quota Calculator](https://developers.google.com/youtube/v3/determine_quota_cost)
- [YouTube API Services — Developer Policies](https://developers.google.com/youtube/terms/developer-policies)
- [TikTok for Developers — Display API Overview](https://developers.tiktok.com/doc/display-api-overview/)
- [TikTok Research API — eligibility](https://developers.tiktok.com/products/research-api/)
- [Meta — Content Library and API for qualified research](https://about.fb.com/news/2023/11/new-tools-to-support-independent-research/)
- [Instagram API — official Meta collection](https://www.postman.com/meta/instagram/documentation/6yqw8pt/instagram-api)
- [Reddit — Responsible Builder Policy](https://support.reddithelp.com/hc/en-us/articles/42728983564564-Responsible-Builder-Policy)
- [Reddit — Data API Terms](https://redditinc.com/policies/data-api-terms)
- [Shopee Việt Nam — Điều khoản dịch vụ](https://help.shopee.vn/portal/4/article/77243-%C4%90I%E1%BB%80U-KHO%E1%BA%A2N-D%E1%BB%8ACH-V%E1%BB%A4)
- [Grab Việt Nam — Merchant terms](https://www.grab.com/vn/terms-policies/merchant-terms-and-conditions/)
- [RFC 9309 — Robots Exclusion Protocol](https://www.rfc-editor.org/rfc/rfc9309.html)
- [RFC 9110 — HTTP conditional requests](https://www.rfc-editor.org/rfc/rfc9110.html)
- [Google — Custom Search JSON API deprecation](https://developers.google.com/custom-search/v1/overview)
- [Google AdSense — User-generated content overview](https://support.google.com/adsense/answer/1355699?hl=en)
- [Google AdSense — Good strategies for managing UGC](https://support.google.com/adsense/answer/3011869?hl=en)
- [Google AdSense — Invalid traffic and policy violations](https://support.google.com/adsense/answer/2660562?hl=en)

### Bản đồ và đơn vị hành chính

- [Quyết định 19/2025/QĐ-TTg — danh mục và mã số đơn vị hành chính Việt Nam](https://vanban.chinhphu.vn/?classid=1&docid=214409&pageid=27160)
- [Nghị quyết 202/2025/QH15 — 34 đơn vị hành chính cấp tỉnh](https://xaydungchinhsach.chinhphu.vn/toan-van-nghi-quyet-so-202-2025-qh15-ve-sap-xep-don-vi-hanh-chinh-cap-tinh-119250612174148722.htm)
- [Bộ Nội vụ — mô hình hai cấp và 3.321 đơn vị cấp xã](https://moha.gov.vn/tin-tuc/bo-noi-vu-thong-tin-ve-chu-truong-chinh-sach-cua-d---oid57777)
- [OpenStreetMap — ODbL attribution guideline](https://osmfoundation.org/wiki/Licence/Attribution_Guidelines)
- [OpenStreetMap — Tile Usage Policy](https://operations.osmfoundation.org/policies/tiles/)
- [OpenStreetMap — Nominatim Usage Policy](https://operations.osmfoundation.org/policies/nominatim/)
- [OpenStreetMap — replication diffs](https://wiki.openstreetmap.org/wiki/Planet.osm/diffs)

### Nguồn/API theo thị trường

- [Korea Tourism Organization — TourAPI](https://www.data.go.kr/data/15101578/openapi.do)
- [Kakao Local REST API](https://developers.kakao.com/docs/en/local/dev-guide)
- [Kakao Developers — quota](https://developers.kakao.com/docs/en/getting-started/quota)
- [Baidu Maps — Web Service API](https://api.map.baidu.com/lbsapi/webservice.htm)
- [Douyin Open Platform — authorized account video list](https://open.douyin.com/platform/resource/docs/openapi/video-management/douyin/search-video/account-video-list)
- [Yandex Search API](https://aistudio.yandex.ru/en/solutions/search-api)
- [Yandex Maps API — licensing and storage](https://yandex.com/dev/commercial/doc/en/concepts/jsapi-geocoder)
- [Telegram API Terms](https://core.telegram.org/api/terms)

### Realtime multimodal

- [Gemini Apps Help — Talk naturally with Gemini Live](https://support.google.com/gemini/answer/15274899?hl=en-GB)
- [Gemini Live API — Capabilities](https://ai.google.dev/gemini-api/docs/live-api/capabilities)
- [Gemini Live API — WebSocket guide](https://ai.google.dev/gemini-api/docs/live-api/get-started-websocket)

### FPT AI Factory

- [FPT AI Inference — Quickstart](https://ai-docs.fptcloud.com/ai-marketplace/ai-inference/quickstart)
- [FPT AI Factory — API keys](https://ai-docs.fptcloud.com/ai-marketplace/ai-inference/tutorials/api-keys)
- [FPT AI Factory — VLM API integration](https://ai-docs.fptcloud.com/api-reference/ai-marketplace/api-reference/api-integration-vision-language-model-md)
- [FPT AI Factory — Multimodal text API integration](https://ai-docs.fptcloud.com/api-reference/ai-marketplace/api-reference/api-integration-multimodal-model-text-md)

### Model mã nguồn mở

- [Qwen3-VL official repository](https://github.com/QwenLM/Qwen3-VL)
- [Qwen3 official repository](https://github.com/QwenLM/Qwen3)

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
