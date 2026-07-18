# Threat Model, Safety Wording và Compliance

> Tách từ `system-design.md` — mục 9.3–9.5.

### 9.3 Threat model

| Rủi ro | Biện pháp |
| --- | --- |
| Prompt injection từ menu/transcript | Coi OCR/transcript là dữ liệu, không phải instruction; tool allowlist và schema validation |
| LLM tạo hotline | Hotline chỉ từ Emergency Directory đã ký |
| Data poisoning/Sybil | Quarantine, device token, dedup, source caps, robust estimator |
| Fake review/brigading | Account/device/graph/rate anomaly, disclosure, reputation theo lĩnh vực, review queue và reversible sanction |
| Merchant retaliation/defamation | Khử PII, đánh giá hành vi/giao dịch, right-to-reply, evidence level, notice/appeal và legal escalation policy |
| Advertiser influence | RBAC tách Sales/Ads/Moderator, schema boundary, ranking feature allowlist và audit thử nghiệm |
| Malicious external URL/SSRF | Connector allowlist, URL normalization, DNS/IP egress policy, fetch proxy sandbox và content-size/MIME limit |
| Copyright/takedown | Rights registry, permitted metadata/embed, canonical attribution, refresh/deletion sync và takedown SLA |
| URL/file độc hại | Signed upload, MIME sniffing, size limit, malware scan, image re-encode |
| Session hijack | Short-lived token, bind device/session, TLS, rate limit |
| Credential stuffing | Rate limit theo account/IP/risk, breached-password check, generic error, progressive challenge và security alert |
| OAuth CSRF/replay/code interception | State, nonce, PKCE, exact redirect allowlist, one-time auth attempt và server-side token validation |
| Account linking takeover | Re-authenticate account hiện tại, khóa theo issuer+sub, không auto-link bằng email và thông báo khi link/unlink |
| Refresh-token theft | Hash-at-rest, rotation, reuse detection, Android Keystore hoặc HttpOnly cookie và revoke session family |
| Account enumeration | Response/timing gần tương đương cho registration, login và password reset; email gửi ngoài luồng |
| Rò raw media trong log | Structured logging allowlist; cấm payload content |
| Admin lạm dụng | RBAC, least privilege, audit, dual approval cho snapshot/playbook |
| Model hallucination | Rule engine, constrained schema, confidence, provenance, abstention |
| Người bán thấy cảnh báo | Private card, haptic, headphone-only explain |
| Camera thu người ngoài | Explicit indicator, foreground-only, no continuous recording, consent reminder |

### 9.4 Safety wording

- “Trong khoảng thường gặp.”
- “Cao hơn mặt bằng của nhóm so sánh; hãy kiểm tra phụ phí.”
- “Rất bất thường so với dữ liệu hiện có; đây không phải kết luận pháp lý.”
- “Chưa đủ dữ liệu để đánh giá đáng tin.”
- “Giá cao hơn trung bình của nhóm so sánh nhưng đã được niêm yết công khai rõ ràng.” (Dùng cho trường hợp quán ăn có niêm yết giá, tối đa chỉ cảnh báo mức `elevated`, không coi là chặt chém).

Không dùng:

- “Người bán này là lừa đảo.”
- “Tuyệt đối không trả tiền” khi người dùng đang bị giữ/đe dọa.
- “Hãy quay phim/đuổi theo/giằng lại tiền.”
- Các cảnh báo mang tính cáo buộc chặt chém đối với quán ăn đã thực hiện niêm yết giá công khai rõ ràng.

### 9.5 Compliance baseline

- Thiết kế mapping consent, purpose limitation, deletion, export và incident handling theo Luật Bảo vệ dữ liệu cá nhân Việt Nam hiện hành.
- Cần legal review trước pilot công khai, đặc biệt với dữ liệu trẻ vị thành niên, sinh trắc học, location và chia sẻ cho cơ quan chức năng.
- Không tuyên bố “zero retention” nếu hạ tầng/model server chưa chứng minh được.

---
