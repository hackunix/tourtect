# Community Knowledge Graph

> Tách từ `system-design.md` — mục 6.1.

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

> [!IMPORTANT]
> **Không sử dụng dữ liệu cũ trước thời điểm sáp nhập (trước 01/07/2025):** 
> Toàn bộ dữ liệu giá giao dịch, hóa đơn, price report hay khảo sát thực địa phát sinh trước ngày 01/07/2025 (trước khi áp dụng chính quyền 2 cấp) không được đưa vào tính toán giá tham chiếu (Reference Price) hiện tại. Ranh giới địa giới hành chính thay đổi lớn và độ trượt giá cũ làm dữ liệu này không còn tương thích hoặc không thể ánh xạ chính xác sang các cohort giá mới, có thể gây ra sai lệch nghiêm trọng hoặc cảnh báo giả. Dữ liệu này chỉ được lưu trữ dạng historical archive và bị loại hoàn toàn khỏi các pipeline hoạt động của Price Engine.

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

Đơn vị hành chính chưa đủ để mô hình hóa giá. Một phường có thể đồng thời chứa phố du lịch, chợ dân sinh, sân bay/bến xe, khu văn phòng và khu dân cư với mặt bằng giá khác nhau. Vì vậy <code>PricingZone</code> là lớp nghiệp vụ tách khỏi <code>Region</code>:

- <code>pricing_zone_id</code>, geometry, tên/loại zone, <code>parent_region_id</code>, nguồn và version.
- Loại zone ban đầu: residential, local_market, tourist_core, attraction, transport_hub, airport, nightlife, office/commercial hoặc event_temporary.
- Zone có <code>valid_from/valid_to</code>; zone sự kiện hoặc sân bay có thể có rule/phụ phí riêng.
- Một observation luôn gắn đơn vị hành chính theo snapshot tại <code>observed_at</code>; <code>pricing_zone_id</code> chỉ bổ sung bối cảnh và không thay mã hành chính.
- Geometry/mapping từ OSM chỉ là candidate. Zone ảnh hưởng reference price phải được data steward duyệt, có provenance và rollback.

Phân khúc giá cũng tách thành các chiều khách quan, không suy từ diện mạo người bán:

| Chiều | Giá trị gợi ý | Mục đích |
| --- | --- | --- |
| <code>service_segment</code> | budget, standard, premium, luxury, regulated | Không so quán bình dân với nhà hàng cao cấp hoặc taxi thường với xe premium |
| <code>venue_type</code> | fixed_shop, casual_eatery, street_stall, mobile_vendor, market_stall, attraction_concession, transport_vendor, peer_to_peer | So đúng mô hình phục vụ và cost structure; không phải risk score |
| <code>transaction_context</code> | posted_price, verbal_quote, negotiated, metered, platform_booked, unsolicited_goods_service, donation_solicitation | Phân biệt giá niêm yết, mặc cả, đồng hồ, giao dịch nền tảng và tình huống không phải mua bán bình thường |
| <code>fulfilment_context</code> | dine_in, takeaway, delivery, curbside, in_vehicle, guided_service | Tách phí phục vụ/giao hàng/bối cảnh thực hiện |

Các giá trị segment cần evidence từ menu/biển hiệu, loại dịch vụ, thuộc tính place hoặc xác nhận của user. Thiếu segment quan trọng thì hỏi lại hoặc trả <code>insufficient_data</code>; không để model tự suy đoán từ khuôn mặt, trang phục hoặc điều kiện kinh tế.

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
