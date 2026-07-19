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

export type AssistantIntent =
  | "general_travel_question" | "place_discovery" | "place_information" | "price_check"
  | "price_explanation" | "translation" | "live_translation" | "menu_or_receipt_analysis"
  | "scam_pattern_assessment" | "safety_assessment" | "emergency_help" | "community_search"
  | "create_report_draft" | "unknown";
export type AssistantInputType = "text" | "voice_transcript" | "image_capture" | "structured_price_candidate" | "structured_safety_facts";
export type AssistantSafetyState = "critical" | "urgent" | "non_emergency" | "information" | "unknown";

export interface CreateAssistantSessionRequest {
  locale: string; target_locale?: string; place_id?: string; approximate_region?: string;
  interaction_mode?: "text" | "voice" | "lens" | "mixed"; processing_consent?: boolean;
}
export interface AssistantSessionContext {
  locale: string; target_locale?: string; place_id?: string; approximate_region?: string;
  interaction_mode: string; current_safety_state?: AssistantSafetyState;
  user_confirmed_facts?: string[]; active_capture_ids?: string[];
  consent_state: { processing: boolean; contribution: boolean; publish: boolean };
}
export interface AssistantSession {
  session_id: string; version: number; created_at: string; updated_at: string; expires_at: string;
  context: AssistantSessionContext; recent_responses?: AssistantResponse[];
}
export interface AssistantMessageRequest {
  message_id: string; input_type: AssistantInputType; text?: string; locale?: string;
  place_id?: string; capture_id?: string; user_confirmed?: boolean; structured_data?: Record<string, unknown>;
}
export interface AssistantEvidence {
  evidence_id: string;
  source_type: "official_source" | "community_post" | "verified_price_observation" | "price_snapshot" | "safety_directory" | "scam_pattern" | "place_record" | "session_fact";
  source_id: string; title: string; summary: string; observed_at?: string;
  freshness: "fresh" | "aging" | "stale" | "unknown";
  evidence_level: "official" | "verified" | "community" | "session_confirmed"; source_url?: string;
}
export interface AssistantToolResult {
  tool_result_id: string; tool_name: string; status: "succeeded" | "insufficient_data" | "degraded" | "failed";
  duration_ms: number; output: Record<string, unknown>; error_category?: string;
}
export interface AssistantConfirmation {
  confirmation_id: string;
  action: "create_report_draft" | "publish_report" | "upload_evidence" | "save_transcript" | "share_location" | "open_dialer" | "contact_trusted_person" | "contribute_observation";
  title: string; description: string; expires_at: string;
}
export interface AssistantSuggestedAction {
  action_id: string; label: string;
  action_type: "clarify" | "deep_link" | "confirmation" | "manual_price_check" | "manual_safety_assessment" | "offline_directory" | "save_private_draft";
  target?: string; requires_confirmation: boolean;
}
export interface AssistantResponse {
  assistant_message_id: string; intent: AssistantIntent; message: string; confidence: number;
  evidence: AssistantEvidence[]; tool_results: AssistantToolResult[];
  requested_confirmation?: AssistantConfirmation; suggested_actions: AssistantSuggestedAction[];
  safety_state: AssistantSafetyState; factors_considered: string[]; missing_information: string[];
  freshness?: string; dataset_version?: string; fallback_used: boolean; trace_id: string;
}
export interface AssistantConfirmationRequest { confirmation_id: string; decision: "confirmed" | "rejected" }
export interface AssistantConfirmationResult { confirmation_id: string; action: string; status: "confirmed" | "rejected"; executed_at: string; result_id?: string; target?: string }
export interface AssistantFeedbackRequest {
  assistant_message_id: string;
  feedback_type: "helpful" | "not_helpful" | "correction" | "translation_incorrect" | "false_positive" | "confirm_extraction" | "contribute_redacted_observation";
  field?: string; original_value?: string; corrected_value?: string; consent_to_contribute?: boolean;
}
export interface AssistantFeedbackReceipt { feedback_id: string; status: "quarantined"; created_at: string }

export class ProblemDetailError extends Error {
  constructor(public status: number, message: string, public requestId?: string, public problem?: Record<string, unknown>) { super(message); this.name = "ProblemDetailError" }
}

const BASE_URL = typeof window === "undefined"
  ? (process.env.API_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080")
  : "";

export function safeUUID(): string {
  if (typeof crypto !== "undefined" && typeof crypto.randomUUID === "function") {
    return crypto.randomUUID();
  }
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const requestId = safeUUID();
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
  if (response.status === 204) return undefined as T;
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
  createAssistantSession: (body: CreateAssistantSessionRequest, signal?: AbortSignal) => apiFetch<AssistantSession>("/v1/assistant/sessions", { method: "POST", body: JSON.stringify(body), signal }),
  getAssistantSession: (sessionId: string, signal?: AbortSignal) => apiFetch<AssistantSession>(`/v1/assistant/sessions/${sessionId}`, { signal }),
  deleteAssistantSession: (sessionId: string, signal?: AbortSignal) => apiFetch<void>(`/v1/assistant/sessions/${sessionId}`, { method: "DELETE", signal }),
  createAssistantMessage: (sessionId: string, body: AssistantMessageRequest, signal?: AbortSignal) => apiFetch<AssistantResponse>(`/v1/assistant/sessions/${sessionId}/messages`, { method: "POST", body: JSON.stringify(body), signal }),
  confirmAssistantAction: (sessionId: string, body: AssistantConfirmationRequest, signal?: AbortSignal) => apiFetch<AssistantConfirmationResult>(`/v1/assistant/sessions/${sessionId}/confirmations`, { method: "POST", body: JSON.stringify(body), signal }),
  createAssistantFeedback: (sessionId: string, body: AssistantFeedbackRequest, signal?: AbortSignal) => apiFetch<AssistantFeedbackReceipt>(`/v1/assistant/sessions/${sessionId}/feedback`, { method: "POST", body: JSON.stringify(body), signal }),
};
