package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

type YdbSummariesRepository struct {
	db *ydb.Driver
}

func NewYdbSummariesRepository(db *ydb.Driver) *YdbSummariesRepository {
	return &YdbSummariesRepository{db: db}
}

func (y YdbSummariesRepository) GetSummary(ctx context.Context, articleId entities.ArticleId) (*entities.Summary, error) {
	articleUUID, err := uuid.Parse(articleId)
	if err != nil {
		return nil, fmt.Errorf("invalid article uuid: %w", err)
	}

	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $article_id AS Uuid;

		SELECT 
			s.id,
			s.generated_at,
			s.text
		FROM summaries s
		WHERE s.id = $article_id
		ORDER BY s.generated_at DESC
		LIMIT 1;	
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$article_id").Uuid(articleUUID).Build()),
		query.WithTxControl(query.NoTx()),
	)

	if err != nil {
		return nil, err
	}

	var summary entities.Summary

	if err := row.ScanStruct(&summary); err != nil {
		return nil, err
	}

	return &summary, nil
}

func (y YdbSummariesRepository) AddSummary(ctx context.Context, articleId entities.ArticleId, summary *entities.Summary) error {
	id := uuid.New()

	articleUUID, err := uuid.Parse(articleId)
	if err != nil {
		return fmt.Errorf("invalid article id: %s", articleId)
	}

	err = y.db.Query().Exec(
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
			Param("$article_id").Uuid(articleUUID).
			Param("$text").Text(summary.Text).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}
