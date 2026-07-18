# Tourtect Frontend Design System

**Requirement IDs:** FO-01–FO-08, PC-05–PC-07, SE-02–SE-06  
**Status:** Web source of truth; Android should reuse the semantics, not the implementation.

## Visual principles

Tourtect is a calm, conversation-first travel community. The primary surface is a centered feed rather than a dashboard. Posts form a continuous reading stream separated by thin rules; cards, rounded containers and elevation are reserved for inputs, dialogs, contextual attachments and urgent safety content.

The visual direction may use familiar social-product proportions but must not copy Meta/Threads marks, proprietary icons, illustrations, wording or motion. Safety and privacy take precedence over visual minimalism, engagement and monetization.

## Tokens

### Color

| Token | Meaning |
| --- | --- |
| `background`, `surface`, `surface-subtle` | Page, primary reading surface and quiet attachments |
| `text-primary`, `text-secondary`, `border` | Main copy, metadata and separators |
| `interactive`, `focus` | Ordinary actions and visible keyboard focus |
| `price-typical`, `price-elevated`, `price-high` | Deterministic Price Engine states |
| `status-uncertain` | `insufficient_data`, stale or unavailable data |
| `safety-urgent`, `safety-critical` | Backend-approved urgent states only |
| `evidence-verified` | Verified receipt/source indicators |

Every token has light and dark values. Components must not hard-code status colors or communicate state by color alone. Red is not an ordinary action color and does not prove fraud.

### Typography and spacing

- UI font: Noto Sans, then Inter/system UI and installed Noto CJK/Korean fallbacks.
- Body copy uses comfortable 1.5–1.6 line height; author names use medium/semibold weight; technical IDs alone may use monospace.
- Spacing follows a 4px base: `4, 8, 12, 16, 20, 24, 32`.
- Main feed is at most 680px. Desktop navigation and context panels never enlarge the reading column.

### Shape, elevation and motion

- Small/medium/large radii are `8/12/16px`; posts themselves remain flat.
- Shadows are limited to modal dialogs and floating sheets.
- Motion is short and functional: tab state, composer expansion, skeletons and bottom sheets.
- `prefers-reduced-motion` removes non-essential transitions and animation never delays SOS.

## Component anatomy

### Post

1. Avatar, author and timestamp.
2. Title/body in the original locale; translated content must be labelled and retain an original-content control.
3. Compact post-type and evidence indicators.
4. Optional place attachment or structured review/price/safety context.
5. Reply, useful, save and share actions.

Useful votes are community feedback, never evidence. Scam reports use neutral “reported risk pattern” language and enter moderation before publication. Thread connectors are shown only for related comments; Web visually indents at most two levels.

### Price status

Price insight shows entered amount, reference range, confidence, freshness, sample size, explanation and dataset version. `elevated` uses caution, `high_risk` uses strong caution, and `insufficient_data` uses neutral uncertainty. No result labels a merchant or acts as legal proof.

### Safety status

Urgent content includes an icon, heading and explicit action text. It renders only backend-approved actions and versioned directory data. Tourtect never invents a hotline, automatically calls, or shares location/incident data. Critical alerts have no social reaction count or celebratory motion and remain usable in silent mode.

## Responsive navigation

- Desktop: compact left rail, 600–680px feed, optional right context panel.
- Tablet: compact rail and centered feed; context panel is removed before the feed is compressed.
- Mobile: single-column feed and bottom navigation for Home, Search, Create, Saved and Notifications. SOS remains a separate safety action and is never hidden in Profile.
- Live and Lens are shown only when their routes and APIs are functional; disabled visual demos are not permitted.

## Accessibility and content integrity

- All actions have accessible names; icons never carry meaning alone.
- Tabs, composer, threads and dialogs support keyboard order and visible focus.
- Status changes use `status`/`alert` semantics where appropriate and touch targets are at least 44px on mobile.
- Light/dark contrast, zoom, Vietnamese diacritics, Korean, Simplified Chinese and Cyrillic content must remain readable.
- Initial loading uses feed-shaped skeletons; mutation progress is inline.
- API failures show an unavailable state, retry and request ID where available. Cached content must be marked stale; fake fallback posts are forbidden.
- Organic ranking uses relevance, freshness, evidence, source diversity, usefulness and safety priority only. Advertising spend, affiliate commission and business tier are prohibited inputs.
