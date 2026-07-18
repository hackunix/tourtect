# Price Check Response mẫu

> Tách từ `system-design.md` — mục 7.5.

### 7.5 Price check response mẫu

~~~json
{
  "check_id": "pc_01...",
  "status": "completed",
  "candidate": {
    "vertical": "food",
    "raw_item": "pho bo",
    "money": {
      "amount_minor": "180000",
      "currency": "VND",
      "exponent": 0
    },
    "unit": "bowl",
    "region_id": "vn_hn_phuong_hoan_kiem",
    "admin_snapshot_version": "vn-admin-2025.07",
    "pricing_zone_id": "hn_old_quarter_core",
    "service_segment": "budget",
    "venue_type": "casual_eatery",
    "transaction_context": "posted_price",
    "fulfilment_context": "dine_in",
    "observed_at": "2026-07-18T03:15:00Z",
    "attributes": {
      "quantity": 1,
      "portion_size": "regular",
      "tax_included": true,
      "service_charge_included": false
    },
    "extraction_confidence": 0.94,
    "user_confirmed": true,
    "comparison_readiness": "ready",
    "missing_required_fields": []
  },
  "insight": {
    "alert_level": "elevated",
    "observed": {
      "amount_minor": "180000",
      "currency": "VND",
      "exponent": 0
    },
    "reference": {
      "p10_minor": "55000",
      "p50_minor": "75000",
      "p90_minor": "120000",
      "currency": "VND",
      "exponent": 0,
      "unit": "bowl",
      "region_id": "vn_hn_phuong_hoan_kiem",
      "pricing_zone_id": "hn_old_quarter_core",
      "service_segment": "budget",
      "venue_type": "casual_eatery",
      "geo_fallback_level": "commune",
      "effective_sample_size": 42,
      "independent_source_count": 11,
      "adjustment_profile_version": "price-adjust-v2026.07.1",
      "adjustments": [
        { "kind": "geo", "factor": 1.12, "confidence": 0.81 },
        { "kind": "segment", "factor": 0.86, "confidence": 0.84 }
      ]
    },
    "confidence": 0.78,
    "comparison_scope": "Phường Hoàn Kiếm, phân khúc bình dân; mở rộng từ lõi Phố cổ ra cấp phường",
    "freshness": "2026-07-01",
    "reasons": [
      "Giá cao hơn khoảng thường gặp của nhóm so sánh"
    ],
    "possible_benign_explanations": [
      "Khẩu phần hoặc loại thịt có thể khác",
      "Có thể đã gồm thuế hoặc phí phục vụ"
    ],
    "dataset_version": "price-v2026.07.1",
    "trace_id": "tr_01..."
  }
}
~~~

---
