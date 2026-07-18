# Metrics và Release Gates

> Tách từ `system-design.md` — mục 10.1–10.2.

### 10.1 Metrics

| Nhóm | Metrics |
| --- | --- |
| OCR/Vision | CER, field extraction F1, price/currency exact match, scene classification F1 |
| ASR | WER theo locale/accent/noise, critical entity exact match |
| Translation | Human adequacy, critical token preservation, latency, semantic parity |
| Price | Red precision, false-positive rate, recall severe overcharge, calibration, abstention |
| Data | Coverage, freshness, source diversity, effective sample size, drift |
| Safety | Critical escalation recall, unsafe-advice rate, playbook parity |
| Realtime | P50/P95 translation latency, disconnect, resume success, dropped frame/audio |
| Device | Battery drain, peak memory, thermal throttling, crash-free sessions |
| Privacy | TTL deletion success, consent errors, raw-content log violations |
| Emergency | Dialer success, hotline lookup availability, offline pack validity |
| Community | DAU/MAU, search success, answer rate, useful contribution, return rate, save/follow và place coverage |
| Trust & Safety | Spam/fake-review precision-recall, report rate, appeal overturn rate, moderation SLA, harassment/PII exposure |
| Content ingestion | Connector success, P0/P1 freshness SLO, dedupe precision, entity-link accuracy, 304 ratio, bytes/new-item, 429/403/5xx, robots/policy drift, deletion/takedown propagation |
| Revenue | Ad fill/eCPM, affiliate conversion, Plus conversion/churn, business retention và revenue concentration |
| Trust firewall | Organic rank parity, sponsored disclosure accuracy, advertiser-policy violations và revenue-triggered moderation changes |

Metrics được slice theo:

- Ngôn ngữ.
- Admin snapshot/đơn vị cấp tỉnh-xã, pricing zone và geo fallback level.
- Vertical.
- Service segment, venue type và transaction context; audit fairness để bảo đảm context không biến thành proxy kết tội nhóm người.
- Thiết bị/tier.
- Execution policy/location, model/provider/version và fallback path.
- Network condition.
- Dữ liệu mới/cũ.

### 10.2 Release gates

| Gate | Mục tiêu |
| --- | --- |
| PTT → audio dịch đầu tiên | P95 ≤ 2 giây trên 4G ổn định |
| Price insight từ lời nói | P95 ≤ 3 giây |
| Live Camera observation | P95 ≤ 3 giây |
| OCR + confirmed price result | P95 ≤ 5 giây |
| Camera sampling | Không vượt 1 FPS |
| Demo execution profile | Với env demo, 100% ModelTrace là <code>server</code>; mọi build không chứa FPT credential. Adaptive build riêng phải đạt test signature/capability/consent |
| Red alert precision | ≥ 95% |
| False-positive trên giá hợp lệ | ≤ 2% |
| Critical safety escalation | 100% trên golden safety set |
| Unsafe confrontation advice | 0 trường hợp trên safety set |
| Critical translation fields | Không sai phủ định, thương tích, số tiền, tiền tệ, vị trí, biển số, số người |
| Background capture | Dừng 100% khi background/end/revoke |
| Offline SOS | Mở hotline/incident card tối đa 2 thao tác |
| Poisoning simulation | Attack budget 1–500; p50/p90 dịch ≤ 5% và không lật quyết định high-risk |
| Fake-review/brigading | Đạt ngưỡng precision/recall đã chốt theo locale; không auto-ban chỉ từ một model signal |
| Commercial firewall | 100% test chứng minh spend, commission và business tier không đổi organic rank/review/alert/moderation |
| Sponsored disclosure | 100% ad/affiliate/sponsored item có nhãn dễ thấy và machine-readable |
| External content deletion | Disable serving và loại khỏi search theo SLA khi source/takedown state thay đổi |
| P0 source freshness | 95% item từ API/feed được discover trong 15 phút khi source hoạt động; HTML-only theo SLO riêng đã duyệt |
| Crawler politeness | Không fetch path disallow; không vượt per-host budget; 429 rate < 1%; mọi 403/451 mở circuit và policy review |
| Administrative snapshot | 100% official code/name/effective date khớp nguồn; không có gap/overlap nghiêm trọng; legacy alias resolve đúng |
| Public PII | Không để lọt PII nghiêm trọng trong golden/red-team public scam report set |

Gate giá được đánh giá theo từng vertical × vùng, không chỉ toàn cục. Mỗi slice cần ít nhất 500 giao dịch hợp lệ và 200 trường hợp overcharge nghiêm trọng đã được phân xử; nếu thiếu thì slice đó không bật cảnh báo đỏ. Point estimate và one-sided Wilson 95% confidence bound đều phải đạt gate: lower bound của precision ≥ 95% và upper bound của false-positive rate ≤ 2%. Nếu chưa đạt, chỉ trả <code>elevated</code> hoặc <code>insufficient_data</code>.

Poisoning test chạy với ngân sách 1, 10, 50, 100 và 500 submission phối hợp; đo dịch chuyển <code>p50</code>, <code>p90</code>, confidence và tỷ lệ lật quyết định alert. Gate yêu cầu cả <code>p50</code> và <code>p90</code> không dịch quá 5%, đồng thời không có price case chuẩn nào bị đổi từ <code>typical</code> sang <code>high_risk</code> hoặc ngược lại.
