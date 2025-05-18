package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

// SQL
const (
	selectFeedByUserIdSql = `
DECLARE $user_id AS Uuid;

SELECT 
	a.id AS id,
	a.source_id AS source_id,
	a.added_at AS added_at,
	a.published_at AS published_at,
	a.title AS title,
	a.text AS text,
	a.url AS url,
	a.preview_url AS preview_url,
	i.starred AS starred,
	i.read AS read
FROM article_user_infos i
INNER JOIN articles a
	ON i.article_id = a.id
WHERE i.user_id = $user_id
ORDER BY
	starred DESC,	
	read ASC,
	published_at DESC
`

	selectArticleFromFeedByUserIdAndArticleIdSql = `
DECLARE $user_id as Uuid;
DECLARE $article_id as Uuid;

SELECT 
	a.id AS id,
	a.source_id AS source_id,
	a.added_at AS added_at,
	a.published_at AS published_at,
	a.title AS title,
	a.text AS text,
	a.url AS url,
	a.preview_url AS preview_url,
	i.starred AS starred,
	i.read AS read
FROM article_user_infos i
INNER JOIN articles a
ON a.id = i.article_id
WHERE i.user_id = $user_id AND i.article_id = $article_id
LIMIT 1;
`
)

// helper functions

func articlesToListValue(articles []*entities.Article) types.Value {
	articlesList := make([]types.Value, 0, len(articles))
	for _, article := range articles {
		articlesList = append(articlesList, types.StructValue(
			types.StructFieldValue("id", types.UuidValue(uuid.Nil)),
			types.StructFieldValue("source_id", types.UuidValue(uuid.Nil)),
			types.StructFieldValue("added_at", types.DatetimeValueFromTime(article.AddedAt)),
			types.StructFieldValue("published_at", types.DatetimeValueFromTime(article.PublishedAt)),
			types.StructFieldValue("title", types.StringValueFromString(article.Title)),
			types.StructFieldValue("text", types.StringValueFromString(article.Text)),
			types.StructFieldValue("url", types.StringValueFromString(article.Url)),
			types.StructFieldValue("preview_url", types.StringValueFromString(article.PreviewUrl)),
		))
	}

	return types.ListValue(articlesList...)
}
