# Public Types — Scam Assessment

> Tách từ `system-design.md` — mục 7.3, dòng 1823–1841.
> Được tách lại thành code fence độc lập để chỉ nạp nhóm type cần dùng.

## Nhóm kiểu Scam Assessment

~~~typescript
interface ScamAssessment {
  urgency: "critical" | "urgent" | "non_emergency" | "information";
  matched_pattern_ids: string[];
  confirmed_behavior_types: ScamBehaviorType[];
  reported_behavior_types: ScamBehaviorType[];
  seller_context?: VenueType;
  confirmed_facts: string[];
  ai_inferences: string[];
  safe_actions: string[];
  do_not: string[];
  follow_up_questions: string[];
  escalation?: {
    incident_type: string;
    emergency_service_ids: string[];
  };
  confidence: number;
  playbook_version: string;
  trace_id: string;
}
~~~
