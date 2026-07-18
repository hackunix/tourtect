# Thuật toán cảnh báo giá

> Tách từ `system-design.md` — mục 6.6.

### 6.6 Thuật toán cảnh báo

Đầu vào bắt buộc:

- Item/dịch vụ đã chuẩn hóa.
- Giá, tiền tệ và đơn vị đã được xác nhận.
- Đơn vị hành chính/version, pricing zone nếu có và thời điểm.
- Service segment, venue type và transaction/fulfilment context đã xác nhận hoặc có nguồn đáng tin.
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

1. Nếu context là <code>donation_solicitation</code> hoặc không có item/unit rõ, Price Engine abstain; nếu có gây áp lực thì chuyển Scam/Safety Engine.
2. Nếu vi phạm một biểu phí/trần giá chính thức đã được xác minh, trả <code>high_risk</code> với rule provenance.
3. Nếu không đủ dữ liệu, không đủ thuộc tính, cohort tương đương hoặc confidence thấp, trả <code>insufficient_data</code>.
4. Nếu giá nhỏ hơn <code>p10</code>, trả <code>typical</code> kèm lý do “thấp hơn khoảng thường gặp”; hệ thống chỉ phát hiện overcharge, không coi giá thấp là scam.
5. Nếu giá nằm từ <code>p10</code> đến <code>p90</code>, trả <code>typical</code>.
6. Nếu giá vượt <code>p90</code> nhưng chưa thỏa red gate, trả <code>elevated</code>; UI phân biệt “cao nhẹ” với chênh lệch đã vượt materiality threshold.
7. Chỉ trả <code>high_risk</code> khi đồng thời:
   - Giá vượt upper reference bound.
   - Chênh lệch vượt materiality threshold theo ngành.
   - OCR/đơn vị/thuộc tính đã xác nhận.
   - Confidence ≥ ngưỡng được calibration để đạt precision gate.
   - Không phải trường hợp quán ăn có niêm yết giá rõ ràng (vertical là <code>food</code> và <code>transaction_context</code> là <code>posted_price</code>). Nếu quán ăn có niêm yết giá rõ ràng, dù giá cao hơn bình thường cũng không coi là chặt chém, hệ thống chỉ trả kết quả tối đa là <code>elevated</code> để phản ánh tính minh bạch.

Ngưỡng không hard-code chung cho năm vertical. Chúng nằm trong versioned configuration, được chọn trên validation set và phát hành cùng model card.
