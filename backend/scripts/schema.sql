CREATE TABLE IF NOT EXISTS users (
    id Uuid,
    telegram_id String NOT NULL,

    PRIMARY KEY (id),
    INDEX ix_users_telegram_id GLOBAL UNIQUE SYNC ON (telegram_id)
);

CREATE TABLE IF NOT EXISTS feed_sources (
    id Uuid,
    feed_url String NOT NULL,
    type String,

    PRIMARY KEY (id),
    INDEX ix_sources_feed_url_unique GLOBAL UNIQUE SYNC ON (feed_url)
);

CREATE TABLE IF NOT EXISTS feed_source_user_infos (
    user_id Uuid,
    source_id Uuid,

    name String,
    disabled Bool DEFAULT False,

    PRIMARY KEY (user_id, source_id)
);

CREATE TABLE IF NOT EXISTS articles (
    id Uuid,
    source_id Uuid NOT NULL,
    added_at Datetime,
    published_at Datetime,
    title Utf8 NOT NULL,
    text Utf8,
    url Utf8 NOT NULL,
    preview_url Utf8,

    PRIMARY KEY (id),
    INDEX ix_articles_source_id GLOBAL SYNC ON (source_id)
);

CREATE TABLE IF NOT EXISTS article_user_infos (
    article_id Uuid NOT NULL,
    user_id Uuid NOT NULL,
    starred Bool DEFAULT FALSE,
    read Bool DEFAULT FALSE,

    PRIMARY KEY (article_id, user_id),
    INDEX ix_article_user_infos_user_id GLOBAL SYNC ON (user_id)
);

CREATE TABLE IF NOT EXISTS summaries (
    id Uuid,
    generated_at Timestamp,
    article_id Uuid NOT NULL,
    text Utf8 NOT NULL,

    PRIMARY KEY (id),
    INDEX ix_summaries_article_id GLOBAL SYNC ON (article_id)
);

CREATE TABLE IF NOT EXISTS digests (
    id Uuid,
    user_id Uuid,
    generated_at Timestamp,
    text Utf8 NOT NULL,

    PRIMARY KEY (id),
    INDEX ix_digests_user_id_generated_at GLOBAL SYNC ON (user_id, generated_at)
);
