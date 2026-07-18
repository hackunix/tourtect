# System Context

> Tách từ `system-design.md` — mục 5.1.

### 5.1 System context

~~~mermaid
flowchart LR
    Tourist["Du khách"]
    Community["Cộng đồng địa phương / expat"]
    Merchant["Chủ cơ sở"]
    Reviewer["Moderator / Editor / Data steward"]
    Sources["Official API / RSS / partner / khảo sát"]

    Android["Tourtect Android\nForum + Live + Lens"]
    Web["Tourtect Responsive Web\nForum + Search + Places"]
    Zalo["Zalo Mini App Lite"]
    Admin["Admin / Moderation Web"]
    Platform["Tourtect Platform"]

    Tourist --> Android
    Tourist --> Web
    Community --> Android
    Community --> Web
    Merchant --> Web
    Reviewer --> Admin

    Android --> Platform
    Web --> Platform
    Zalo --> Platform
    Admin --> Platform
    Sources --> Platform
~~~
