-- Seed users
INSERT INTO users (id, telegram_id) VALUES
    (Uuid("11111111-1111-1111-1111-111111111111"), "user1"),
    (Uuid("22222222-2222-2222-2222-222222222222"), "user2"),
    (Uuid("33333333-3333-3333-3333-333333333333"), "user3");

-- Seed feed sources
INSERT INTO feed_sources (id, feed_url, type) VALUES
    (Uuid("44444444-4444-4444-4444-444444444444"), "http://nginx/fsf_news.xml", "rdf"),
    (Uuid("55555555-5555-5555-5555-555555555555"), "http://nginx/opennews_all.rss", "rss");

-- Seed user subscriptions to feed sources
INSERT INTO feed_source_user_infos (user_id, source_id, name, disabled) VALUES
    (Uuid("22222222-2222-2222-2222-222222222222"), Uuid("44444444-4444-4444-4444-444444444444"), "FSF feed", False),
    (Uuid("33333333-3333-3333-3333-333333333333"), Uuid("55555555-5555-5555-5555-555555555555"), "opennet", True);

-- Seed articles
-- INSERT INTO articles (id, source_id, added_at, published_at, title, text, url, preview_url) VALUES
--     (Uuid("77777777-7777-7777-7777-777777777777"), Uuid("44444444-4444-4444-4444-444444444444"), DateTime("2023-01-01T10:00:00Z"), DateTime("2023-01-01T09:00:00Z"), "New AI Breakthrough", "Researchers have made a significant breakthrough...", "https://tech-blog.example.com/ai-breakthrough", "https://tech-blog.example.com/preview1.jpg"),
--     (Uuid("88888888-8888-8888-8888-888888888888"), Uuid("44444444-4444-4444-4444-444444444444"), DateTime("2023-01-02T11:00:00Z"), DateTime("2023-01-02T10:00:00Z"), "Quantum Computing Advances", "New quantum processor achieves...", "https://tech-blog.example.com/quantum", "https://tech-blog.example.com/preview2.jpg"),
--     (Uuid("99999999-9999-9999-9999-999999999999"), Uuid("55555555-5555-5555-5555-555555555555"), DateTime("2023-01-03T12:00:00Z"), DateTime("2023-01-03T11:00:00Z"), "Global Summit Concludes", "World leaders agreed on new climate...", "https://news.example.com/summit", NULL),
--     (Uuid("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), Uuid("66666666-6666-6666-6666-666666666666"), DateTime("2023-01-04T13:00:00Z"), DateTime("2023-01-04T12:00:00Z"), "Championship Finals", "Team A wins against Team B in a thrilling...", "https://sports.example.com/finals", "https://sports.example.com/preview1.jpg");

-- Seed user article interactions
-- INSERT INTO article_user_infos (article_id, user_id, starred, read) VALUES
--     (Uuid("77777777-7777-7777-7777-777777777777"), Uuid("11111111-1111-1111-1111-111111111111"), True, True),
--     (Uuid("88888888-8888-8888-8888-888888888888"), Uuid("11111111-1111-1111-1111-111111111111"), False, True),
--     (Uuid("77777777-7777-7777-7777-777777777777"), Uuid("22222222-2222-2222-2222-222222222222"), False, False),
--     (Uuid("99999999-9999-9999-9999-999999999999"), Uuid("11111111-1111-1111-1111-111111111111"), False, False),
--     (Uuid("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), Uuid("33333333-3333-3333-3333-333333333333"), True, True);

-- Seed summaries
-- INSERT INTO summaries (id, generated_at, article_id, text) VALUES
--     (Uuid("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), Timestamp("2023-01-01T10:30:00Z"), Uuid("77777777-7777-7777-7777-777777777777"), "Researchers achieved a breakthrough in AI by developing a new algorithm that significantly improves processing efficiency."),
--     (Uuid("cccccccc-cccc-cccc-cccc-cccccccccccc"), Timestamp("2023-01-04T13:30:00Z"), Uuid("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), "Team A won the championship finals against Team B with a score of 3-2 in an overtime match that kept fans on edge.");

-- Seed digests
-- INSERT INTO digests (id, user_id, generated_at, text) VALUES
--     (Uuid("dddddddd-dddd-dddd-dddd-dddddddddddd"), Uuid("11111111-1111-1111-1111-111111111111"), Timestamp("2023-01-05T08:00:00Z"), "Your daily digest:\n1. New AI Breakthrough (read)\n2. Quantum Computing Advances (read)\n3. Global Summit Concludes (new)"),
--     (Uuid("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"), Uuid("33333333-3333-3333-3333-333333333333"), Timestamp("2023-01-05T09:00:00Z"), "Your sports update:\nChampionship Finals: Team A wins against Team B (read)");