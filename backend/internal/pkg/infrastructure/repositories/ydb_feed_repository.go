package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"strings"
)

type YdbFeedRepository struct {
	db *ydb.Driver
}

func NewYdbFeedRepository(driver *ydb.Driver) *YdbFeedRepository {
	return &YdbFeedRepository{
		db: driver,
	}
}

func (y *YdbFeedRepository) AddArticlesToFeed(ctx context.Context, articles []*entities.Article) error {
	//TODO implement me
	panic("implement me")
}

func (y *YdbFeedRepository) GetFeedByUser(ctx context.Context, userId entities.UserId) ([]*entities.Article, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user uuid: %w", err)
	}

	result, err := y.db.Query().Query(
		ctx,
		`
		DECLARE $user_id AS Uuid;

		SELECT 
			a.id,
			a.added_at,
			a.published_at,
			a.title,
			a.text,
			a.url,
			a.preview_url,
			i.starred,
			i.read
		FROM article_user_infos i
		INNER JOIN articles a
			ON i.article_id = a.id
		ORDER BY
			i.starred DESC,	
			i.read ASC,
			a.published_at DESC
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Build(),
		),
		query.WithTxControl(query.SnapshotReadOnlyTxControl()),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = result.Close(ctx)
	}()

	return mapResultSet[entities.Article](result, ctx)
}

func (y *YdbFeedRepository) GetTodayArticles(ctx context.Context, userId entities.UserId) ([]*entities.Article, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %w", err)
	}

	// query
	result, err := y.db.Query().Query(
		ctx,
		`
		DECLARE $user_id as Uuid;

		SELECT *
		FROM article_user_info i
		INNER JOIN articles a
		ON a.id = i.article_id
		WHERE i.user_id = $user_id AND a.added_at = CurrentUtcDate();
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Build(),
		),
		query.WithTxControl(query.SnapshotReadOnlyTxControl()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get articles for today: %w", err)
	}

	defer func() {
		_ = result.Close(ctx)
	}()

	// map
	return mapResultSet[entities.Article](result, ctx)
}

func (y *YdbFeedRepository) GetArticleById(ctx context.Context, userId entities.UserId, articleId entities.ArticleId) (*entities.Article, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %w", err)
	}

	articleUUID, err := uuid.Parse(articleId)
	if err != nil {
		return nil, fmt.Errorf("invalid article id %w", err)
	}

	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $user_id as Uuid;
		DECLARE $article_id as Uuid;

		SELECT 
			a.id,
			a.added_at,
			a.published_at,
			a.title,
			a.text,
			a.url,
			a.preview_url,
			i.starred,
			i.read
		FROM article_user_infos i
		INNER JOIN articles a
		ON a.id = i.article_id
		WHERE i.user_id = $user_id AND i.article_id = $article_id;
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$article_id").Uuid(articleUUID).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get article by id %w", err)
	}

	var article entities.Article

	if err := row.ScanStruct(&article); err != nil {
		return nil, fmt.Errorf("failed to map result to article entity: %w", err)
	}

	return &article, nil
}

func (y *YdbFeedRepository) UpdateArticle(
	ctx context.Context,
	userId entities.UserId,
	articleId entities.ArticleId,
	patch *entities.ArticlePatch,
) (*entities.Article, error) {
	if !patch.Starred.HasValue() && !patch.Read.HasValue() {
		return nil, fmt.Errorf("patch is empty")
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id %w", err)
	}

	articleUUID, err := uuid.Parse(articleId)
	if err != nil {
		return nil, fmt.Errorf("invalid article id %w", err)
	}

	sqlBuilder := strings.Builder{}

	sqlBuilder.WriteString("DECLARE $user_id as Uuid;\nDECLARE $article_id as Uuid;\n")

	if patch.Starred.HasValue() {
		sqlBuilder.WriteString("DECLARE $starred AS Boolean;\n")
	}

	if patch.Read.HasValue() {
		sqlBuilder.WriteString("DECLARE $read AS Boolean;\n")
	}

	sqlBuilder.WriteString("UPDATE article_user_infos i\n SET ")

	if patch.Starred.HasValue() {
		sqlBuilder.WriteString("starred = $starred")
	}

	if patch.Read.HasValue() {
		if patch.Starred.HasValue() {
			sqlBuilder.WriteString(", ")
		}

		sqlBuilder.WriteString("read = $read")
	}

	ydbParams := ydb.ParamsBuilder().
		Param("user_id").Uuid(userUUID).
		Param("article_id").Uuid(articleUUID)

	if patch.Starred.HasValue() {
		value, _ := patch.Starred.Value()
		ydbParams = ydbParams.Param("starred").Bool(value)
	}

	if patch.Read.HasValue() {
		value, _ := patch.Read.Value()
		ydbParams = ydbParams.Param("read").Bool(value)
	}

	err = y.db.Query().Exec(
		ctx,
		sqlBuilder.String(),
		query.WithParameters(ydbParams.Build()),
		query.WithTxControl(query.NoTx()),
	)

	return nil, err
}
