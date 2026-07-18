-- Tourtect seed data — synthetic Hanoi data for development
-- No real personal data. All content is fictional.
-- This file is idempotent: uses ON CONFLICT DO NOTHING or checks.

BEGIN;

-- ============================================================
-- Seed principal (system/anonymous author for seed posts)
-- ============================================================
INSERT INTO principals (principal_id, display_name, primary_email, email_verified, status, locale)
VALUES
    ('019078a0-0001-7000-8000-000000000001', 'Seed User Alpha', 'seed-alpha@example.test', true, 'active', 'vi-VN'),
    ('019078a0-0001-7000-8000-000000000002', 'Seed User Beta', 'seed-beta@example.test', true, 'active', 'en')
ON CONFLICT (principal_id) DO NOTHING;

-- ============================================================
-- 5 Hanoi places with PostGIS coordinates
-- ============================================================
INSERT INTO places (place_id, name, category, region_id, address, description, coordinates, phone, website, opening_hours, post_count, review_count, average_rating, freshness)
VALUES
    (
        '019078a0-1001-7000-8000-000000000001',
        'Phở Thìn Bờ Hồ',
        'food',
        'hanoi-hoan-kiem',
        '13 Lò Đúc, Hai Bà Trưng, Hà Nội',
        'Quán phở bò nổi tiếng với nước dùng đậm đà và thịt bò tái chín.',
        ST_SetSRID(ST_MakePoint(105.8575, 21.0115), 4326)::geography,
        NULL, NULL, '06:00-10:00, 18:00-22:00',
        3, 2, 4.50, now() - interval '2 days'
    ),
    (
        '019078a0-1001-7000-8000-000000000002',
        'Bún Chả Hương Liên',
        'food',
        'hanoi-dong-da',
        '24 Lê Văn Hưu, Hai Bà Trưng, Hà Nội',
        'Quán bún chả nổi tiếng sau chuyến thăm của cựu Tổng thống Obama.',
        ST_SetSRID(ST_MakePoint(105.8505, 21.0135), 4326)::geography,
        NULL, NULL, '10:00-14:30',
        2, 1, 4.20, now() - interval '5 days'
    ),
    (
        '019078a0-1001-7000-8000-000000000003',
        'Chợ Đồng Xuân',
        'street_retail',
        'hanoi-hoan-kiem',
        'Đồng Xuân, Hoàn Kiếm, Hà Nội',
        'Chợ lớn nhất khu phố cổ, bán đủ loại hàng hóa từ thực phẩm đến quần áo.',
        ST_SetSRID(ST_MakePoint(105.8500, 21.0380), 4326)::geography,
        NULL, NULL, '06:00-18:00',
        2, 0, 0, now() - interval '10 days'
    ),
    (
        '019078a0-1001-7000-8000-000000000004',
        'Hồ Hoàn Kiếm',
        'tour',
        'hanoi-hoan-kiem',
        'Hoàn Kiếm, Hà Nội',
        'Hồ thiêng giữa trung tâm Hà Nội với Tháp Rùa và Đền Ngọc Sơn.',
        ST_SetSRID(ST_MakePoint(105.8523, 21.0288), 4326)::geography,
        NULL, 'https://example.test/hoan-kiem', '24/7',
        1, 0, 0, now() - interval '1 day'
    ),
    (
        '019078a0-1001-7000-8000-000000000005',
        'Nội Bài Taxi Stand',
        'taxi',
        'hanoi-soc-son',
        'Sân bay Quốc tế Nội Bài, Sóc Sơn, Hà Nội',
        'Bãi đón taxi chính thức tại sân bay Nội Bài.',
        ST_SetSRID(ST_MakePoint(105.8038, 21.2187), 4326)::geography,
        NULL, NULL, '24/7',
        2, 0, 0, now() - interval '3 days'
    )
ON CONFLICT (place_id) DO NOTHING;

-- ============================================================
-- Place aliases
-- ============================================================
INSERT INTO place_aliases (place_id, alias, locale)
VALUES
    ('019078a0-1001-7000-8000-000000000001', 'Pho Thin Bo Ho', 'en'),
    ('019078a0-1001-7000-8000-000000000001', 'Phở Thìn', 'vi-VN'),
    ('019078a0-1001-7000-8000-000000000002', 'Bun Cha Obama', 'en'),
    ('019078a0-1001-7000-8000-000000000002', 'Bún Chả Obama', 'vi-VN'),
    ('019078a0-1001-7000-8000-000000000003', 'Dong Xuan Market', 'en'),
    ('019078a0-1001-7000-8000-000000000004', 'Hoan Kiem Lake', 'en'),
    ('019078a0-1001-7000-8000-000000000004', 'Sword Lake', 'en'),
    ('019078a0-1001-7000-8000-000000000005', 'Noi Bai Airport Taxi', 'en')
ON CONFLICT DO NOTHING;

-- ============================================================
-- 6 public posts + 2 drafts
-- ============================================================
INSERT INTO posts (post_id, author_id, post_type, original_locale, title, body, moderation_status, created_at, updated_at)
VALUES
    -- Public posts
    (
        '019078a0-2001-7000-8000-000000000001',
        '019078a0-0001-7000-8000-000000000001',
        'review', 'vi-VN',
        'Phở Thìn Bờ Hồ — phở ngon nhất Hà Nội?',
        'Phở ở đây nước dùng rất đậm, thịt bò tái chín mềm. Giá 55.000đ/bát, hợp lý so với khu Hoàn Kiếm. Quán nhỏ nhưng phục vụ nhanh.',
        'published', now() - interval '30 days', now() - interval '30 days'
    ),
    (
        '019078a0-2001-7000-8000-000000000002',
        '019078a0-0001-7000-8000-000000000002',
        'tip', 'en',
        'How to get a fair taxi from Noi Bai Airport',
        'Use the official taxi stand on the right side after exiting arrivals. Grab and Be are available too. Typical fare to Old Quarter is 350,000-400,000 VND by meter.',
        'published', now() - interval '20 days', now() - interval '20 days'
    ),
    (
        '019078a0-2001-7000-8000-000000000003',
        '019078a0-0001-7000-8000-000000000001',
        'price_report', 'vi-VN',
        'Giá bún chả tại Hương Liên tháng 7/2026',
        'Bún chả set đầy đủ: 50.000đ. Nem rán: 10.000đ/cái. Nước ngọt: 15.000đ. Giá ổn định so với lần trước.',
        'published', now() - interval '15 days', now() - interval '15 days'
    ),
    (
        '019078a0-2001-7000-8000-000000000004',
        '019078a0-0001-7000-8000-000000000002',
        'discussion', 'en',
        'Dong Xuan Market — worth visiting or tourist trap?',
        'The ground floor has mostly wholesale goods. Upper floors have more tourist items. Bargaining is expected — start at 50% of the asking price.',
        'published', now() - interval '10 days', now() - interval '10 days'
    ),
    (
        '019078a0-2001-7000-8000-000000000005',
        '019078a0-0001-7000-8000-000000000001',
        'scam_report', 'vi-VN',
        'Cảnh báo: Taxi giả mạo ở Nội Bài',
        'Có nhiều xe taxi giả mạo thương hiệu lớn, đón khách ở khu vực sân bay. Luôn kiểm tra bảng hiệu, đồng hồ tính tiền và biển số xe.',
        'published', now() - interval '5 days', now() - interval '5 days'
    ),
    (
        '019078a0-2001-7000-8000-000000000006',
        '019078a0-0001-7000-8000-000000000002',
        'tip', 'en',
        'Walking around Hoan Kiem Lake — best times',
        'Early morning (6-7 AM) is great for seeing locals exercise. Weekend evenings are lively with street performers. The pedestrian zone is active Fri-Sun evenings.',
        'published', now() - interval '2 days', now() - interval '2 days'
    ),
    -- Draft posts
    (
        '019078a0-2001-7000-8000-000000000007',
        '019078a0-0001-7000-8000-000000000001',
        'review', 'vi-VN',
        'Draft: Review Chợ Đồng Xuân [chưa hoàn thành]',
        'Bản nháp — cần bổ sung thêm thông tin về giá và trải nghiệm.',
        'draft', now() - interval '1 day', now() - interval '1 day'
    ),
    (
        '019078a0-2001-7000-8000-000000000008',
        '019078a0-0001-7000-8000-000000000002',
        'question', 'en',
        'Draft: Is it safe to walk in Old Quarter at night?',
        'Want to ask about safety in the Old Quarter area after 10 PM. Need to add more context.',
        'draft', now() - interval '12 hours', now() - interval '12 hours'
    )
ON CONFLICT (post_id) DO NOTHING;

-- Post-place links
INSERT INTO post_place_links (post_id, place_id)
VALUES
    ('019078a0-2001-7000-8000-000000000001', '019078a0-1001-7000-8000-000000000001'),
    ('019078a0-2001-7000-8000-000000000002', '019078a0-1001-7000-8000-000000000005'),
    ('019078a0-2001-7000-8000-000000000003', '019078a0-1001-7000-8000-000000000002'),
    ('019078a0-2001-7000-8000-000000000004', '019078a0-1001-7000-8000-000000000003'),
    ('019078a0-2001-7000-8000-000000000005', '019078a0-1001-7000-8000-000000000005'),
    ('019078a0-2001-7000-8000-000000000006', '019078a0-1001-7000-8000-000000000004'),
    ('019078a0-2001-7000-8000-000000000007', '019078a0-1001-7000-8000-000000000003')
ON CONFLICT DO NOTHING;

-- ============================================================
-- 2 price snapshots + 6 price observations
-- ============================================================
INSERT INTO price_snapshots (snapshot_id, vertical, region_id, service_segment, venue_type, unit, currency, p10_minor, p50_minor, p90_minor, exponent, sample_size, independent_source_count, version, effective_from)
VALUES
    (
        '019078a0-3001-7000-8000-000000000001',
        'food', 'hanoi-hoan-kiem', 'budget', 'casual_eatery',
        'bowl', 'VND', 40000, 55000, 75000, 0,
        25, 8, 'snap-food-202607-v1', now() - interval '60 days'
    ),
    (
        '019078a0-3001-7000-8000-000000000002',
        'taxi', 'hanoi-soc-son', 'standard', 'transport_vendor',
        'trip', 'VND', 300000, 380000, 450000, 0,
        40, 15, 'snap-taxi-202607-v1', now() - interval '45 days'
    )
ON CONFLICT (snapshot_id) DO NOTHING;

INSERT INTO price_observations (observation_id, snapshot_id, vertical, region_id, raw_item, amount_minor, currency, exponent, unit, service_segment, venue_type, transaction_context, extraction_confidence, user_confirmed, source, observed_at)
VALUES
    ('019078a0-4001-7000-8000-000000000001', '019078a0-3001-7000-8000-000000000001', 'food', 'hanoi-hoan-kiem', 'Phở bò tái chín', 55000, 'VND', 0, 'bowl', 'budget', 'casual_eatery', 'posted_price', 0.95, true, 'user_report', now() - interval '30 days'),
    ('019078a0-4001-7000-8000-000000000002', '019078a0-3001-7000-8000-000000000001', 'food', 'hanoi-hoan-kiem', 'Phở bò tái', 50000, 'VND', 0, 'bowl', 'budget', 'casual_eatery', 'posted_price', 0.90, true, 'user_report', now() - interval '25 days'),
    ('019078a0-4001-7000-8000-000000000003', '019078a0-3001-7000-8000-000000000001', 'food', 'hanoi-hoan-kiem', 'Bún chả set', 50000, 'VND', 0, 'bowl', 'budget', 'casual_eatery', 'posted_price', 0.92, true, 'user_report', now() - interval '15 days'),
    ('019078a0-4001-7000-8000-000000000004', '019078a0-3001-7000-8000-000000000002', 'taxi', 'hanoi-soc-son', 'Airport taxi to Old Quarter', 380000, 'VND', 0, 'trip', 'standard', 'transport_vendor', 'metered', 0.85, true, 'user_report', now() - interval '20 days'),
    ('019078a0-4001-7000-8000-000000000005', '019078a0-3001-7000-8000-000000000002', 'taxi', 'hanoi-soc-son', 'Noi Bai to Hoan Kiem', 350000, 'VND', 0, 'trip', 'standard', 'transport_vendor', 'metered', 0.80, false, 'user_report', now() - interval '10 days'),
    ('019078a0-4001-7000-8000-000000000006', '019078a0-3001-7000-8000-000000000002', 'taxi', 'hanoi-soc-son', 'Airport taxi', 420000, 'VND', 0, 'trip', 'standard', 'transport_vendor', 'verbal_quote', 0.70, false, 'user_report', now() - interval '5 days')
ON CONFLICT (observation_id) DO NOTHING;

-- ============================================================
-- Safety directory — 1 version with approved emergency numbers
-- ============================================================
INSERT INTO safety_directory_versions (version_id, version, description, published_at)
VALUES
    ('019078a0-5001-7000-8000-000000000001', 'safety-v2026.07', 'Initial safety directory for Hanoi', now() - interval '30 days')
ON CONFLICT (version_id) DO NOTHING;

INSERT INTO safety_directory_entries (version_id, region_id, service_name, service_type, phone_number, description, locale)
VALUES
    ('019078a0-5001-7000-8000-000000000001', 'hanoi', 'Công an (Police)', 'police', '113', 'Số điện thoại khẩn cấp công an Việt Nam', 'vi-VN'),
    ('019078a0-5001-7000-8000-000000000001', 'hanoi', 'Cấp cứu (Ambulance)', 'ambulance', '115', 'Số điện thoại cấp cứu y tế', 'vi-VN'),
    ('019078a0-5001-7000-8000-000000000001', 'hanoi', 'Cứu hỏa (Fire)', 'fire', '114', 'Số điện thoại cứu hỏa', 'vi-VN'),
    ('019078a0-5001-7000-8000-000000000001', 'hanoi', 'Tourist Police Hanoi', 'tourist_police', '069-942-0626', 'Công an du lịch Hà Nội', 'vi-VN'),
    ('019078a0-5001-7000-8000-000000000001', 'hanoi', 'SOS Vietnam Hotline', 'hotline', '1900-599-920', 'Đường dây hỗ trợ du khách 24/7', 'vi-VN')
ON CONFLICT DO NOTHING;

COMMIT;
