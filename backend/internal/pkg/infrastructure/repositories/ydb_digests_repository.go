package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	"io"
)

type YdbDigestsRepository struct {
	db *ydb.Driver
}

func NewYdbDigestsRepository(db *ydb.Driver) *YdbDigestsRepository {
	return &YdbDigestsRepository{db: db}
}

func (y *YdbDigestsRepository) FindLatestDigestForToday(ctx context.Context, userId entities.UserId) (bool, string, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return false, "", fmt.Errorf("failed to parse user UUID: %w", err)
	}

	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $user_id AS Uuid;

		SELECT text 
		FROM digests 
		WHERE user_id = $user_id AND generated_at = CurrentUtcDate()
		LIMIT 1;
		`,
		query.WithParameters(ydb.ParamsBuilder().Param("$user_id").Uuid(userUUID).Build()),
		query.WithTxControl(query.NoTx()),
	)

	if errors.Is(err, io.EOF) {
		return false, "", nil
	}

	if err != nil {
		return false, "", fmt.Errorf("failed to find latest digest: %w", err)
	}

	var text string

	if err = row.Scan(&text); err != nil {
		return false, "", fmt.Errorf("failed to find latest digest: %w", err)
	}

	return true, text, nil
}

func (y *YdbDigestsRepository) AddDigest(ctx context.Context, userId entities.UserId, digest string) error {
	id := uuid.New()

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("invalid user uuid: %w", err)
	}

	err = y.db.Query().Exec(
		ctx,
		`
		DECLARE $uuid AS Uuid;
		DECLARE $user_id AS Uuid;
		DECLARE $digest AS String;
		
		INSERT INTO digests(id, user_id, text) VALUES ($uuid, $user_id, $digest);
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$uuid").Uuid(id).
			Param("$digest").Any(types.StringValueFromString(digest)).
			Param("$user_id").Uuid(userUUID).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}
