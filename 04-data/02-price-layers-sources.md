# Các lớp dữ liệu giá và phân tầng nguồn

> Tách từ `system-design.md` — mục 6.2–6.3.

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
