# Điều hướng và wireflow đa kênh

> Tách từ `system-design.md` — mục 3.1–3.2.

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
