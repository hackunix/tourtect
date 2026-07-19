# Source Access and Crawler Policy

This document defines the strict security and compliance boundaries for Tourtect's external crawling and sandboxed fetching system.

## Prohibited Ingestion Practices

To ensure legal compliance and avoid being marked as a malicious agent, Tourtect **explicitly prohibits** the following techniques:
1. **Private Endpoint Scraping**: Crawling endpoints not exposed publicly (e.g. reverse-engineering private mobile app APIs).
2. **Access Control / Login Bypass**: Using authenticated user sessions to fetch data behind membership walls without license.
3. **Paywall Bypassing**: Fetching premium copyrighted materials that require subscriptions.
4. **CAPTCHA Circumvention**: Integrating automated CAPTCHA solvers (e.g. 2Captcha) to force access.
5. **Evasive Fingerprinting**: Faking browser footprints, rotating TLS signatures, or lying about client identity to bypass rate limits.
6. **Aggressive Proxy Rotation**: Utilizing residential proxy networks for the sole purpose of circumventing host blocks.

---

## Allowed Access Modes

1. **Official APIs**: Platform APIs (e.g. YouTube Data API, Kakao Local) with approved API keys and quotas.
2. **Partner Integrations**: Direct feeds, API push endpoints, or shared database replication schemas.
3. **RSS / Atom Feeds**: Standard news or warning broadcast channels.
4. **Sitemaps**: Discovering site URLs via published `sitemap.xml` entries.
5. **Allowlisted HTML**: Crawling public pages on approved domains, provided `robots.txt` does not prohibit it.
6. **User Submissions**: Fetching individual URLs shared directly by Tourtect users for validation.

---

## Crawler Identity and Header Policy

All crawlers must send a unified, transparent user-agent header containing contact details:
```
User-Agent: TourtectBot/0.1 (+https://tourtect.example/crawler; data@tourtect.example)
```
If a site operator requests opt-out via `robots.txt` or an email to `data@tourtect.example`, the domain is immediately added to the blocked list in `source_policies`.

---

## Scheduler Policy & Host Concurrency

To minimize load on public servers:
- **Default Delay**: 5 seconds between consecutive requests on the same host (RPS $\le$ 0.2).
- **Concurrency Limit**: Maximum of 1 concurrent fetch request per host.
- **Circuit Breaker**: Auto-disable a connector if error rates exceed 5% (5xx status) or if 429 (Too Many Requests) is returned.
- **Conditional GET**: Always include `If-None-Match` (ETag) or `If-Modified-Since` (Last-Modified) headers to utilize 304 Not Modified responses, preventing extraction runs on unchanged pages.

---

## SSRF Security Gate (Fetch Proxy)

The fetch layer must implement strict validation on target URLs to prevent Server-Side Request Forgery (SSRF). The fetch proxy rejects requests pointing to:
- **Loopback addresses**: `localhost`, `127.0.0.1`, `::1`
- **Private subnets**: `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`
- **Link-Local**: `169.254.0.0/16` (preventing cloud metadata exposure, e.g. AWS/GCP instance details).
- **Unsupported protocols**: Rejects `file://`, `ftp://`, `data:`, `gopher://`, `javascript:`. Only `https://` (and `http://` under explicit allowed rules) are accepted.
- **DNS Re-binding Protection**: Re-resolve host IP and check safety on *every* redirect in the chain.
