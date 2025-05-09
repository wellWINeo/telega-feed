package repositories

// SQL
const (
	selectFeedSourceByFeedUrlSql = `
DECLARE $feed_url AS String;

SELECT s.id AS id,
FROM feed_sources s
WHERE s.feed_url = $feed_url
LIMIT 1;
`

	insertFeedSourceSql = `
DECLARE $id AS Uuid;
DECLARE $feed_url AS String;
DECLARE $type AS String;

INSERT INTO feed_sources (id, feed_url, type)
VALUES ($id, $feed_url, $type);
`

	insertFeedSourceUserInfoSql = `
DECLARE $source_id AS Uuid;
DECLARE $user_id AS Uuid;
DECLARE $name AS String;

INSERT INTO feed_source_user_infos (user_id, source_id, name)
VALUES ($user_id, $source_id, $name);
`

	selectFeedSourcesByUserIdSql = `
DECLARE $user_id AS Uuid;

SELECT
	s.id as id,
	i.name as name,
	s.feed_url as feed_url,
	s.type as type,
	i.disabled as disabled
FROM feed_source_user_infos i
INNER JOIN feed_sources s
	ON i.source_id = s.id
WHERE i.user_id = $user_id;
`

	selectFeedSourceByUserIdAndSourceIdSql = `
DECLARE $user_id AS Uuid;
DECLARE $source_id AS Uuid;

SELECT
	s.id AS id,
	i.name AS name,
	s.feed_url AS feed_url,
	s.type AS type,
	i.disabled AS disabled
FROM feed_sources s
INNER JOIN feed_source_user_infos i
	ON i.source_id = s.id
WHERE 
	i.source_id = $source_id 
	AND i.user_id = $user_id
LIMIT 1;
`

	selectNotDisabledFeedSourcesSql = `
SELECT DISTINCT
	s.id AS id,
	'' AS name,
	s.feed_url AS feed_url,
	s.type AS type,
	False AS disabled
FROM feed_sources s
INNER JOIN feed_source_user_infos i
	ON i.source_id = s.id
WHERE i.disabled = False;
`
)
