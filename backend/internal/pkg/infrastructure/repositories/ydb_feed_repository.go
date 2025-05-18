package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"strings"
	"time"
)

type YdbFeedRepository struct {
	db *ydb.Driver
}

func NewYdbFeedRepository(driver *ydb.Driver) *YdbFeedRepository {
	return &YdbFeedRepository{
		db: driver,
	}
}

func (y *YdbFeedRepository) AddArticleToFeed(ctx context.Context, article *entities.Article) error {
	articleUUID := uuid.New()

	return y.db.Table().DoTx(ctx, func(ctx context.Context, tx table.TransactionActor) error {
		result, err := tx.Execute(
			ctx, `
			DECLARE $id AS Uuid;
			DECLARE $added_at AS Datetime;	
			DECLARE $published_at AS Datetime;	
			DECLARE $source_id AS Uuid;	
			DECLARE $title AS String;	
			DECLARE $text AS String;	
			DECLARE $url AS String;	
			DECLARE $preview_url AS String;	

			INSERT INTO 
				articles(id, added_at, published_at, source_id, title, text, url, preview_url)
			VALUES 
				($id, $added_at, $published_at, $source_id, $title, $text, $url, $preview_url)
			`, table.NewQueryParameters(
				table.ValueParam("$id", types.UuidValue(articleUUID)),
				table.ValueParam("$added_at", types.DatetimeValueFromTime(article.AddedAt)),
				table.ValueParam("$published_at", types.DatetimeValueFromTime(article.PublishedAt)),
				table.ValueParam("$source_id", types.UuidValue(article.SourceId)),
				table.ValueParam("$title", types.StringValueFromString(article.Title)),
				table.ValueParam("$text", types.StringValueFromString(article.Text)),
				table.ValueParam("$url", types.StringValueFromString(article.Url)),
				table.ValueParam("$preview_url", types.StringValueFromString(article.PreviewUrl)),
			))

		if err != nil {
			return err
		}

		_ = result.Close()

		result, err = tx.Execute(ctx,
			`
			DECLARE $article_id AS Uuid;
			DECLARE $source_id AS Uuid;

			UPSERT INTO article_user_infos(article_id, user_id, starred, read)
			SELECT 
				$article_id, user_id, False, False
			FROM feed_source_user_infos
			WHERE source_id = $source_id AND disabled = False;
			`, table.NewQueryParameters(
				table.ValueParam("$article_id", types.UuidValue(articleUUID)),
				table.ValueParam("$source_id", types.UuidValue(article.SourceId)),
			))

		if err != nil {
			return err
		}

		_ = result.Close()

		return nil
	})
}

func (y *YdbFeedRepository) GetFeedByUser(ctx context.Context, userId entities.UserId) ([]*entities.Article, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user uuid: %w", err)
	}

	result, err := y.db.Query().Query(ctx, selectFeedByUserIdSql,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).Build(),
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
		FROM article_user_infos i
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

	row, err := y.db.Query().QueryRow(ctx, selectArticleFromFeedByUserIdAndArticleIdSql,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$article_id").Uuid(articleId).
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
		Param("article_id").Uuid(articleId)

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

func (y *YdbFeedRepository) DeleteOrphanedArticles(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (y *YdbFeedRepository) DeleteArticlesAddedBefore(ctx context.Context, datetime time.Time) error {
	//TODO implement me
	panic("implement me")
}
