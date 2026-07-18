# Mô hình dữ liệu giá và an toàn

> Tách từ `system-design.md` — mục 6.4.

### 6.4 Mô hình dữ liệu giá và an toàn

#### PriceObservation

| Trường | Ý nghĩa |
| --- | --- |
| <code>observation_id</code> | UUID |
| <code>canonical_item_id</code> | Item/dịch vụ chuẩn hóa |
| <code>vertical</code> | taxi, exchange, food, tour, street_retail |
| <code>amount</code>, <code>currency</code>, <code>unit</code> | Giá đã chuẩn hóa; amount lưu bằng database NUMERIC hoặc integer minor units, không dùng float |
| <code>attributes</code> | Thuộc tính so sánh theo vertical |
| <code>region_id</code>, <code>admin_snapshot_version</code> | Đơn vị hành chính có hiệu lực tại thời điểm observation và phiên bản mapping |
| <code>pricing_zone_id</code>, <code>geo_precision</code> | Khu vực giá tùy chọn và độ chính xác; không mặc định lưu GPS chính xác |
| <code>service_segment</code>, <code>venue_type</code> | Phân khúc và loại điểm bán dùng để chọn đúng cohort; không dùng làm scam score |
| <code>transaction_context</code>, <code>fulfilment_context</code> | Cách báo giá/thỏa thuận và cách cung cấp dịch vụ |
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
| <code>region_id</code>, <code>admin_snapshot_version</code> | Đơn vị hành chính và phiên bản địa giới có hiệu lực |
| <code>pricing_zone_id</code> | Zone vi mô nếu đủ mẫu; null khi snapshot ở cấp hành chính |
| <code>service_segment</code>, <code>venue_type</code> | Phân khúc và mô hình điểm bán của cohort |
| <code>geo_fallback_level</code> | exact_zone, commune, province, province_cluster hoặc national_vertical |
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
| <code>adjustment_profile_version</code> | Phiên bản factor/model điều chỉnh địa lý, phân khúc, venue và thời gian |
| <code>rule_provenance</code> | Luật/trần giá chính thức nếu có |
| <code>published_at</code>, <code>approved_by</code> | Governance |

#### ScamPattern

- <code>pattern_id</code>, vùng/jurisdiction.
- <code>behavior_type</code>: price_switch, unit_or_currency_ambiguity, false_official_fee, unsolicited_goods_service, bait_and_switch, withheld_change_or_document, staged_damage_claim, coercive_solicitation hoặc confinement_or_threat.
- <code>seller_context</code>/<code>venue_type</code> chỉ phục vụ retrieval pattern phù hợp, không tự tăng risk score.
- Tín hiệu đa ngôn ngữ.
- Dữ kiện bắt buộc, evidence tối thiểu và các tín hiệu chống phân biệt đối xử.
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
