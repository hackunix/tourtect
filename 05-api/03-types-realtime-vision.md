# Public Types — Realtime và Vision

> Tách từ `system-design.md` — mục 7.3, dòng 1638–1676.
> Được tách lại thành code fence độc lập để chỉ nạp nhóm type cần dùng.

## Nhóm kiểu Realtime và Vision

~~~typescript
interface RealtimeSessionConfig {
  channel: "android_app";
  mode: LiveMode;
  tourist_locale: Locale;
  local_locale: "vi-VN";
  region_id: string;
  execution_policy: ExecutionPolicy;
  enabled_scopes: ConsentScope[];
}

interface ModelTrace {
  provider: string;
  model: string;
  model_version: string;
  execution_location: ExecutionLocation;
  latency_ms: number;
  confidence?: number;
  fallback_reason?: string;
}

interface VisionObservation {
  observation_id: string;
  scene_type:
    | "menu"
    | "food"
    | "taxi_meter"
    | "receipt"
    | "price_board"
    | "exchange_counter"
    | "unknown";
  regions_of_interest: Array<{
    kind: "text" | "price" | "object";
    box: [number, number, number, number];
    confidence: number;
  }>;
  requires_capture: boolean;
  model_trace: ModelTrace;
}
~~~
