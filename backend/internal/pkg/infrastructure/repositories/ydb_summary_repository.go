package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"io"
)

type YdbSummariesRepository struct {
	db *ydb.Driver
}

func NewYdbSummariesRepository(db *ydb.Driver) *YdbSummariesRepository {
	return &YdbSummariesRepository{db: db}
}

func (y *YdbSummariesRepository) GetSummary(ctx context.Context, articleId entities.ArticleId) (*entities.Summary, error) {
	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $article_id AS Uuid;

		SELECT 
			s.id,
			s.generated_at,
			s.text
		FROM summaries s
		WHERE s.article_id= $article_id
		ORDER BY s.generated_at DESC
		LIMIT 1;	
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$article_id").Uuid(articleId).Build()),
		query.WithTxControl(query.NoTx()),
	)

	if errors.Is(err, io.EOF) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var summary entities.Summary

	if err := row.ScanStruct(&summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func (y *YdbSummariesRepository) AddSummary(ctx context.Context, articleId entities.ArticleId, summary *entities.Summary) error {
	id := uuid.New()

	err := y.db.Query().Exec(
		ctx,
		`
		DECLARE $id AS Uuid;
		DECLARE $article_id AS Uuid;
		DECLARE $text AS Utf8;

		INSERT INTO summaries(id, article_id, text) 
		VALUES ($id, $article_id, $text);
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$id").Uuid(id).
			Param("$article_id").Uuid(articleId).
			Param("$text").Text(summary.Text).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}

func (y *YdbSummariesRepository) DeleteOrphanedSummaries(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
