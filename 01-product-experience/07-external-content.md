# External Content và chiến lược cập nhật dữ liệu

> Tách từ `system-design.md` — mục 3.10.

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
