export type Locale = "vi-VN" | "en";
export type FeedMode = "following" | "nearby" | "latest" | "trending" | "safety";

export interface Coordinates { latitude: number; longitude: number }
export interface PlaceSummary {
  place_id: string; name: string; category: string; region_id?: string; address?: string;
  coordinates?: Coordinates; aliases?: string[]; post_count?: number; average_rating?: number;
  freshness?: string; distance_m?: number; created_at?: string;
}
export interface PlaceDetail extends PlaceSummary {
  description?: string; phone?: string; website?: string; opening_hours?: string;
  review_count?: number; updated_at?: string;
}
export interface AuthorSummary { principal_id: string; display_name: string }
export interface PlaceAttachmentData { place_id: string; name: string; category: string; region_id: string; freshness: string }
export interface Post {
  post_id: string; author_id: string; author?: AuthorSummary; post_type: string; original_locale: string;
  title: string; body: string; region_id?: string; place_ids?: string[]; places?: PlaceAttachmentData[];
  evidence_level: string; commercial_disclosure?: string; moderation_status: string;
  structured_data?: Record<string, unknown>; useful_count?: number; comment_count?: number;
  viewer_useful?: boolean; viewer_saved?: boolean; reason_codes?: string[];
  created_at: string; updated_at?: string;
}
export interface CursorInfo { next_cursor?: string; has_more?: boolean }
export interface FeedResponse { items: Post[]; pagination: CursorInfo }
export interface PlaceListResponse { items: PlaceSummary[]; pagination: CursorInfo }
export interface SearchResponse { query: string; tab: string; places: PlaceSummary[]; posts: Post[] }
export interface Comment {
  comment_id: string; post_id: string; author?: AuthorSummary; author_id?: string;
  parent_comment_id?: string; body: string; moderation_status: string; created_at: string; updated_at: string;
}
export interface CommentListResponse { items: Comment[]; pagination: CursorInfo }
export interface Notification { notification_id: string; kind: string; message: string; post_id?: string; read_at?: string; created_at: string }

export interface Money { amount_minor: string; currency: string; exponent: number }
export interface PriceCheckRequest {
  vertical: "taxi" | "exchange" | "food" | "tour" | "street_retail"; raw_item: string; money: Money;
  unit: string; region_id: string; pricing_zone_id?: string; service_segment: "budget" | "standard" | "premium" | "luxury" | "regulated";
  venue_type: string; transaction_context: string; observed_at: string; extraction_confidence?: number; user_confirmed?: boolean;
}
export interface PriceReference { p10_minor?: string; p50_minor?: string; p90_minor?: string; currency?: string; exponent?: number; unit?: string; effective_sample_size?: number }
export interface PriceInsight {
  alert_level: "typical" | "elevated" | "high_risk" | "insufficient_data"; observed: Money; reference?: PriceReference;
  deviation_ratio?: number; confidence: number; comparison_scope?: string; freshness: string; reasons: string[];
  possible_benign_explanations?: string[]; dataset_version: string; trace_id: string;
}
export interface SafetyAssessmentRequest { observed_facts: string[]; threat_indicators?: string[]; injury_indicators?: string[]; confinement_indicators?: string[]; coercion_indicators?: string[]; ability_to_leave?: boolean }
export interface SafetyAssessment {
  urgency: "critical" | "urgent" | "non_emergency" | "information"; safe_actions: string[]; approved_action_codes?: string[];
  explanation_codes?: string[]; silent_mode_recommended?: boolean; surface_emergency_options?: boolean;
  emergency_service_ids?: string[]; safety_directory_version: string; confidence: number; trace_id: string;
}

export class ProblemDetailError extends Error {
  constructor(public status: number, message: string, public requestId?: string, public problem?: Record<string, unknown>) { super(message); this.name = "ProblemDetailError" }
}

const BASE_URL = typeof window === "undefined"
  ? (process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080")
  : (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080");

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const requestId = crypto.randomUUID();
  const response = await fetch(`${BASE_URL}${path}`, {
    cache: "no-store", ...options,
    headers: { "Content-Type": "application/json", "X-Request-ID": requestId, ...options?.headers },
  });
  if (!response.ok) {
    let problem: Record<string, unknown> = {};
    try { problem = await response.json() as Record<string, unknown>; } catch { /* non-JSON upstream */ }
    const detail = typeof problem.detail === "string" ? problem.detail : typeof problem.title === "string" ? problem.title : "API request failed";
    throw new ProblemDetailError(response.status, detail, String(problem.request_id || response.headers.get("X-Request-ID") || requestId), problem);
  }
  return response.json() as Promise<T>;
}

function query(params: Record<string, string | number | undefined>) {
  const values = new URLSearchParams(); Object.entries(params).forEach(([key, value]) => { if (value !== undefined && value !== "") values.set(key, String(value)); });
  return values.toString();
}

export const api = {
  getFeed: (params: { mode: FeedMode; region_id?: string; cursor?: string; limit?: number }) => apiFetch<FeedResponse>(`/v1/feed?${query(params)}`),
  search: (q: string, tab = "top") => apiFetch<SearchResponse>(`/v1/search?${query({ q, tab })}`),
  getPlaces: (params: { q?: string; region_id?: string; limit?: number } = {}) => apiFetch<PlaceListResponse>(`/v1/places?${query(params)}`),
  getPlace: (id: string) => apiFetch<PlaceDetail>(`/v1/places/${id}`),
  createDraft: (body: { post_type: string; original_locale: string; title: string; body: string; region_id?: string; place_ids?: string[]; structured_data?: Record<string, unknown> }) => apiFetch<Post>("/v1/posts/drafts", { method: "POST", body: JSON.stringify(body) }),
  publishPost: (id: string) => apiFetch<Post>(`/v1/posts/${id}/publish`, { method: "POST" }),
  getComments: (id: string) => apiFetch<CommentListResponse>(`/v1/posts/${id}/comments`),
  addComment: (id: string, body: string, parent_comment_id?: string) => apiFetch<Comment>(`/v1/posts/${id}/comments`, { method: "POST", body: JSON.stringify({ body, parent_comment_id }) }),
  setUseful: (id: string, active: boolean) => apiFetch<{ post_id: string; vote: boolean }>(`/v1/posts/${id}/votes/useful`, { method: active ? "PUT" : "DELETE" }),
  setSaved: (id: string, active: boolean) => apiFetch<{ post_id: string; saved: boolean }>(`/v1/saved/posts/${id}`, { method: active ? "PUT" : "DELETE" }),
  getSaved: () => apiFetch<FeedResponse>("/v1/saved"),
  getNotifications: () => apiFetch<{ items: Notification[]; pagination: CursorInfo }>("/v1/notifications"),
  checkPrice: (body: PriceCheckRequest) => apiFetch<PriceInsight>("/v1/price-checks", { method: "POST", body: JSON.stringify(body) }),
  assessSafety: (body: SafetyAssessmentRequest) => apiFetch<SafetyAssessment>("/v1/safety/assessments", { method: "POST", body: JSON.stringify(body) }),
};
