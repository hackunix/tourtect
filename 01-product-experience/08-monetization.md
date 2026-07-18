# Mô hình doanh thu và Commercial Trust Firewall

> Tách từ `system-design.md` — mục 3.11.

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
