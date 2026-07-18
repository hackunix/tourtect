# Public Types — Price Check

> Tách từ `system-design.md` — mục 7.3, dòng 1677–1822 và ghi chú 1844–1846.
> Được tách lại thành code fence độc lập để chỉ nạp nhóm type cần dùng.

## Nhóm kiểu Price Check

~~~typescript
interface Money {
  amount_minor: string;
  currency: string;
  exponent: number;
}

interface PriceCandidateBase {
  canonical_item_id?: string;
  raw_item: string;
  money: Money;
  unit: string;
  region_id: string;
  admin_snapshot_version: string;
  pricing_zone_id?: string;
  service_segment: ServiceSegment;
  venue_type: VenueType;
  transaction_context: TransactionContext;
  fulfilment_context?: string;
  observed_at: string;
  extraction_confidence: number;
  user_confirmed: boolean;
  comparison_readiness: "ready" | "needs_confirmation" | "insufficient";
  missing_required_fields: string[];
}

interface TaxiAttributes {
  distance_m?: number;
  duration_s?: number;
  vehicle_class?: string;
  airport_trip?: boolean;
  tolls_included?: boolean;
  waiting_fee_included?: boolean;
}

interface ExchangeAttributes {
  source_currency: string;
  target_currency: string;
  quoted_rate?: string;
  fee?: Money;
  rate_timestamp: string;
}

interface FoodAttributes {
  portion_size?: string;
  quantity: number;
  tax_included?: boolean;
  service_charge_included?: boolean;
}

interface TourAttributes {
  duration_minutes: number;
  group_type: "private" | "group";
  guide_language?: Locale;
  inclusions: string[];
  season?: string;
}

interface StreetRetailAttributes {
  item_category: string;
  quantity: number;
  weight_or_size?: string;
  item_condition?: "new" | "used" | "unknown";
  authenticity_claim?: "none" | "official" | "handmade" | "unknown";
}

type PriceCandidate =
  | (PriceCandidateBase & {
      vertical: "taxi";
      attributes: TaxiAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "exchange";
      attributes: ExchangeAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "food";
      attributes: FoodAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "tour";
      attributes: TourAttributes;
    })
  | (PriceCandidateBase & {
      vertical: "street_retail";
      attributes: StreetRetailAttributes;
    });

interface CaptureCreateResponse {
  capture_id: string;
  upload_url: string;
  required_headers: Record<string, string>;
  expires_at: string;
}

interface CaptureFinalizeRequest {
  object_etag: string;
  sha256: string;
  media_type: "image/jpeg" | "image/png";
  redaction_applied: boolean;
  client_ocr_model?: string;
}

type PriceCheckRequest =
  | {
      source: "capture";
      capture_id: string;
      candidate: PriceCandidate;
    }
  | {
      source: "manual";
      candidate: PriceCandidate;
    };

interface PriceInsight {
  alert_level: AlertLevel;
  observed: Money;
  reference?: {
    p10_minor: string;
    p50_minor: string;
    p90_minor: string;
    currency: string;
    exponent: number;
    unit: string;
    region_id: string;
    pricing_zone_id?: string;
    service_segment: ServiceSegment;
    venue_type: VenueType;
    geo_fallback_level: "exact_zone" | "commune" | "province" | "province_cluster" | "national_vertical";
    effective_sample_size: number;
    independent_source_count: number;
    adjustment_profile_version?: string;
    adjustments?: Array<{
      kind: "geo" | "segment" | "venue" | "context" | "temporal";
      factor: number;
      confidence: number;
    }>;
  };
  confidence: number;
  comparison_scope: string;
  freshness: string;
  reasons: string[];
  possible_benign_explanations: string[];
  dataset_version: string;
  trace_id: string;
}
~~~

Chỉ candidate có <code>comparison_readiness = "ready"</code> mới được Price Engine đánh giá. Candidate thiếu thuộc tính vẫn được nhận để UI hỏi lại, nhưng kết quả phải là <code>insufficient_data</code> cùng <code>missing_required_fields</code>; server luôn validate schema theo vertical, không tin type phía client.
