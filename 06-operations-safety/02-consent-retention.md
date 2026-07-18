# Consent và Retention

> Tách từ `system-design.md` — mục 9.1–9.2.

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
