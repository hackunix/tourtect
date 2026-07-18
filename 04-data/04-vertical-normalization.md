# Chuẩn hóa theo Vertical

> Tách từ `system-design.md` — mục 6.5.

### 6.5 Chuẩn hóa theo vertical

| Vertical | Thuộc tính bắt buộc |
| --- | --- |
| Taxi/đưa đón | Khoảng cách, thời gian, loại xe, giờ cao điểm, sân bay, cầu đường, phí chờ |
| Đổi tiền | Cặp tiền, timestamp, số đưa/nhận, phí công khai, effective rate |
| Ăn uống | Món/SKU, khẩu phần, số lượng, thuế, service charge, loại địa điểm |
| Tour | Thời lượng, private/group, ngôn ngữ, lịch trình, vé/xe/bữa ăn, mùa |
| Bán lẻ thiết yếu/đường phố | SKU/category, số lượng, trọng lượng/kích cỡ, tình trạng, niêm yết hay mặc cả, claim chính hãng/thủ công và loại điểm bán; item không chuẩn hóa phải abstain |

#### 6.5.1 Chọn cohort theo địa giới và phân khúc

Price Engine chọn cohort từ hẹp đến rộng, nhưng chỉ mở rộng một chiều tại một bước và luôn giữ các thuộc tính material của item:

1. Cùng <code>pricing_zone_id + service_segment + venue_type + vertical attributes</code>.
2. Cùng đơn vị cấp xã/phường/đặc khu + segment + venue type.
3. Cùng đơn vị cấp tỉnh + segment + urban/tourism cluster tương đương.
4. Cụm tỉnh/thành tương đồng hoặc national vertical baseline, chỉ khi item đủ chuẩn hóa.
5. Nếu vẫn thiếu <code>effective_sample_size</code> hoặc independent source count, trả <code>insufficient_data</code>.

Không mở rộng từ budget sang premium, từ street stall sang attraction concession hoặc từ posted price sang unsolicited goods chỉ để có đủ mẫu. Response phải công khai scope thực tế, cấp fallback, sample size, freshness và confidence. Khi có đủ dữ liệu, mô hình production có thể dùng hierarchical partial pooling để làm mượt estimate giữa zone/xã/tỉnh nhưng snapshot phát hành vẫn phải giải thích được.

#### 6.5.2 Điều chỉnh reference price

Profile demo ưu tiên snapshot cohort đã materialize vì dễ kiểm tra và giải thích. Nếu phải điều chỉnh từ một base snapshot rộng hơn, dùng <code>ReferenceAdjustmentProfile</code> có version:

- <code>base_snapshot_id</code> ở cấp tỉnh/cụm tương đồng.
- <code>geo_factor</code> cho xã/phường hoặc pricing zone.
- <code>segment_factor</code> cho budget/standard/premium/luxury/regulated.
- <code>venue_factor</code> và <code>context_factor</code> khi có đủ dữ liệu độc lập.
- <code>temporal_factor</code> cho daypart, cuối tuần, mùa/sự kiện và ngày hiệu lực.
- Sample size, independent source count, confidence interval, cap/floor, version và reviewer cho từng factor.

Với mỗi quantile <code>q ∈ {p10,p50,p90}</code>, demo có thể tính <code>adjusted_q = base_q × geo_factor × segment_factor × venue_factor × temporal_factor</code>. Chỉ áp dụng factor đạt minimum evidence; factor thiếu được thay bằng <code>1.0</code> và hạ confidence. Mọi factor bị cap để một zone ít dữ liệu không làm reference nhảy cực đoan, đồng thời response trả breakdown “base + các điều chỉnh” cho người dùng/reviewer.

Khi dữ liệu đủ lớn, thay phép nhân độc lập bằng robust hierarchical model trên log price: item baseline + province/commune/zone effect + segment/venue effect + thuộc tính vertical + temporal effect. Partial pooling kéo zone ít mẫu về parent distribution; model không có feature về danh tính, diện mạo, dân tộc, giới hoặc hoàn cảnh kinh tế của người bán.
