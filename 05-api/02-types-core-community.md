# Public Types — Core, Account và Community

> Tách từ `system-design.md` — mục 7.3, dòng 1477–1637.
> Được tách lại thành code fence độc lập để chỉ nạp nhóm type cần dùng.

## Nhóm kiểu Core, Account và Community

~~~typescript
type Locale = "vi-VN" | "ko-KR" | "zh-Hans" | "en" | "ru-RU";

type Channel = "android_app" | "responsive_web" | "zalo_mini_app" | "admin_web";
type PostType =
  | "discussion"
  | "question"
  | "review"
  | "price_report"
  | "scam_report"
  | "tip"
  | "official_alert"
  | "external_link";
type EvidenceLevel = "none" | "metadata" | "verified_receipt" | "verified_source";
type ModerationStatus = "draft" | "pending" | "published" | "limited" | "removed" | "appealed";
type CommercialDisclosure = "none" | "invited" | "gifted" | "affiliate" | "employee" | "sponsored";
type RevenueSurface = "contextual_ad" | "affiliate" | "plus" | "business" | "sponsored" | "b2b";
type LiveMode = "voice" | "camera";
type SpeakerRole = "tourist" | "local";
type ExecutionLocation = "device" | "server";
type ExecutionPolicy = "adaptive" | "server_only" | "local_only";
type Vertical = "taxi" | "exchange" | "food" | "tour" | "street_retail";
type ServiceSegment = "budget" | "standard" | "premium" | "luxury" | "regulated";
type VenueType =
  | "fixed_shop"
  | "casual_eatery"
  | "street_stall"
  | "mobile_vendor"
  | "market_stall"
  | "attraction_concession"
  | "transport_vendor"
  | "peer_to_peer";
type TransactionContext =
  | "posted_price"
  | "verbal_quote"
  | "negotiated"
  | "metered"
  | "platform_booked"
  | "unsolicited_goods_service"
  | "donation_solicitation";
type ScamBehaviorType =
  | "price_switch"
  | "unit_or_currency_ambiguity"
  | "false_official_fee"
  | "unsolicited_goods_service"
  | "bait_and_switch"
  | "withheld_change_or_document"
  | "staged_damage_claim"
  | "coercive_solicitation"
  | "confinement_or_threat";

type AlertLevel =
  | "typical"
  | "elevated"
  | "high_risk"
  | "insufficient_data";

type ConsentScope =
  | "process_microphone"
  | "process_camera"
  | "precise_location"
  | "share_incident"
  | "contribute_redacted_data";

type AccountStatus = "pending_email_verification" | "active" | "suspended" | "scheduled_for_deletion";
type IdentityProvider = "password" | "google";

interface Account {
  account_id: string;
  display_name: string;
  primary_email_masked: string;
  email_verified: boolean;
  status: AccountStatus;
  locale: Locale;
  created_at: string;
}

interface FederatedIdentity {
  identity_id: string;
  account_id: string;
  provider: "google";
  issuer: string;
  subject: string;
  email_at_link_time_masked?: string;
  linked_at: string;
}

interface AccountSession {
  session_id: string;
  account_id: string;
  device_label?: string;
  created_at: string;
  last_seen_at: string;
  expires_at: string;
  revoked_at?: string;
}

interface OAuthAttempt {
  attempt_id: string;
  provider: "google";
  state_hash: string;
  nonce_hash: string;
  pkce_challenge: string;
  redirect_uri_id: string;
  expires_at: string;
  consumed_at?: string;
}

interface Post {
  post_id: string;
  author_id: string;
  post_type: PostType;
  original_locale: Locale;
  title: string;
  body: string;
  place_ids: string[];
  topic_ids: string[];
  region_id?: string;
  evidence_level: EvidenceLevel;
  commercial_disclosure: CommercialDisclosure;
  moderation_status: ModerationStatus;
  created_at: string;
  updated_at: string;
}

interface Review {
  review_id: string;
  post_id: string;
  place_id: string;
  visited_at?: string;
  overall_rating: number;
  price_transparency_rating?: number;
  service_rating?: number;
  safety_rating?: number;
  value_rating?: number;
  evidence_level: EvidenceLevel;
  commercial_disclosure: CommercialDisclosure;
}

interface ExternalContent {
  external_content_id: string;
  platform: string;
  canonical_url: string;
  source_content_id?: string;
  original_locale?: string;
  published_at?: string;
  last_checked_at: string;
  rights_status: "embed_allowed" | "metadata_only" | "partner_licensed" | "blocked";
  source_state: "active" | "changed" | "deleted" | "embed_disabled" | "takedown" | "expired";
  place_ids: string[];
  topic_ids: string[];
}

interface ReputationProfile {
  user_id: string;
  local_knowledge: number;
  price_evidence: number;
  translation: number;
  safety: number;
  last_recalculated_at: string;
}
~~~
