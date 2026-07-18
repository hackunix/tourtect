# Cold Start và chống Data Poisoning

> Tách từ `system-design.md` — mục 6.7–6.8.

### 6.7 Cold start không phụ thuộc đối tác

- Seed dữ liệu có provenance từ nguồn công khai chính thức.
- Khảo sát nhỏ phân tầng tại ba cụm pilot.
- Sampling frame là <code>admin unit × pricing zone × service segment × venue type × daypart</code>; ưu tiên item chuẩn hóa và đặt minimum independent-source count trước khi publish snapshot.
- Seed riêng cho fixed shop, casual eatery, street stall/mobile vendor và market stall để không lấy giá nhà hàng làm chuẩn cho giao dịch đường phố hoặc ngược lại.
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
