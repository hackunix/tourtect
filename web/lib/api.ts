// Typed Tourtect API Client for Next.js

export interface Coordinates {
  latitude: number;
  longitude: number;
}

export interface PlaceSummary {
  place_id: string;
  name: string;
  category: string;
  region_id?: string;
  address?: string;
  coordinates: Coordinates;
  aliases?: string[];
  post_count?: number;
  average_rating?: number;
  freshness?: string;
  distance_m?: number;
  created_at: string;
}

export interface PlaceDetail extends PlaceSummary {
  description?: string;
  phone?: string;
  website?: string;
  opening_hours?: string;
  review_count?: number;
  updated_at?: string;
}

export interface PlaceListResponse {
  items: PlaceSummary[];
  pagination: {
    next_cursor?: string;
    has_more?: boolean;
  };
}

export interface Post {
  post_id: string;
  author_id: string;
  post_type: string;
  original_locale: string;
  title: string;
  body: string;
  place_ids?: string[];
  evidence_level: string;
  commercial_disclosure: string;
  moderation_status: string;
  created_at: string;
  updated_at: string;
}

export interface PostListResponse {
  items: Post[];
  pagination: {
    next_cursor?: string;
    has_more?: boolean;
  };
}

export interface Money {
  amount_minor: string;
  currency: string;
  exponent: number;
}

export interface PriceCheckRequest {
  vertical: "taxi" | "exchange" | "food" | "tour" | "street_retail";
  raw_item: string;
  money: Money;
  unit: string;
  region_id: string;
  pricing_zone_id?: string;
  service_segment: "budget" | "standard" | "premium" | "luxury" | "regulated";
  venue_type: string;
  transaction_context: string;
  observed_at: string;
  extraction_confidence?: number;
  user_confirmed?: boolean;
}

export interface PriceReference {
  p10_minor?: string;
  p50_minor?: string;
  p90_minor?: string;
  currency?: string;
  exponent?: number;
  unit?: string;
  region_id?: string;
  pricing_zone_id?: string;
  service_segment?: string;
  venue_type?: string;
  geo_fallback_level?: string;
  effective_sample_size?: number;
}

export interface PriceInsight {
  alert_level: "typical" | "elevated" | "high_risk" | "insufficient_data";
  observed: Money;
  reference?: PriceReference;
  deviation_ratio?: number;
  confidence: number;
  comparison_scope: string;
  freshness: string;
  reasons: string[];
  possible_benign_explanations: string[];
  dataset_version: string;
  trace_id: string;
}

export interface SafetyAssessmentRequest {
  observed_facts: string[];
  user_reported_state?: string;
  threat_indicators?: string[];
  injury_indicators?: string[];
  confinement_indicators?: string[];
  coercion_indicators?: string[];
  ability_to_leave?: boolean;
  user_confirmed_facts?: string[];
}

export interface SafetyAssessment {
  urgency: "critical" | "urgent" | "non_emergency" | "information";
  safe_actions: string[];
  approved_action_codes?: string[];
  explanation_codes?: string[];
  silent_mode_recommended?: boolean;
  surface_emergency_options?: boolean;
  emergency_service_ids?: string[];
  safety_directory_version: string;
  confidence: number;
  trace_id: string;
}

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const url = `${BASE_URL}${path}`;
  const response = await fetch(url, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      "X-Request-ID": crypto.randomUUID(),
      ...options?.headers,
    },
  });

  if (!response.ok) {
    let errorDetail = "API request failed";
    try {
      const prob = await response.json();
      errorDetail = prob.detail || prob.title || errorDetail;
    } catch {
      // Ignore parse errors on non-json status bodies
    }
    throw new Error(errorDetail);
  }

  return response.json() as Promise<T>;
}

export const api = {
  getPlaces: (params?: { q?: string; category?: string; region_id?: string; lat?: number; lon?: number; limit?: number; cursor?: string }) => {
    const qps = new URLSearchParams();
    if (params) {
      if (params.q) qps.set("q", params.q);
      if (params.category) qps.set("category", params.category);
      if (params.region_id) qps.set("region_id", params.region_id);
      if (params.lat !== undefined) qps.set("lat", String(params.lat));
      if (params.lon !== undefined) qps.set("lon", String(params.lon));
      if (params.limit !== undefined) qps.set("limit", String(params.limit));
      if (params.cursor) qps.set("cursor", params.cursor);
    }
    return apiFetch<PlaceListResponse>(`/v1/places?${qps.toString()}`, { cache: "no-store" });
  },

  getPlace: (placeId: string) => {
    return apiFetch<PlaceDetail>(`/v1/places/${placeId}`, { cache: "no-store" });
  },

  getPosts: (params?: { place_id?: string; post_type?: string; limit?: number; cursor?: string }) => {
    const qps = new URLSearchParams();
    if (params) {
      if (params.place_id) qps.set("place_id", params.place_id);
      if (params.post_type) qps.set("post_type", params.post_type);
      if (params.limit !== undefined) qps.set("limit", String(params.limit));
      if (params.cursor) qps.set("cursor", params.cursor);
    }
    return apiFetch<PostListResponse>(`/v1/posts?${qps.toString()}`, { cache: "no-store" });
  },

  createDraft: (body: { post_type: string; original_locale: string; title: string; body: string; place_ids?: string[] }) => {
    return apiFetch<Post>("/v1/posts/drafts", {
      method: "POST",
      body: JSON.stringify(body),
    });
  },

  publishPost: (postId: string) => {
    return apiFetch<Post>(`/v1/posts/${postId}/publish`, {
      method: "POST",
    });
  },

  checkPrice: (body: PriceCheckRequest) => {
    return apiFetch<PriceInsight>("/v1/price-checks", {
      method: "POST",
      body: JSON.stringify(body),
    });
  },

  assessSafety: (body: SafetyAssessmentRequest) => {
    return apiFetch<SafetyAssessment>("/v1/safety/assessments", {
      method: "POST",
      body: JSON.stringify(body),
    });
  },
};
