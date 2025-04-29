package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"strings"
)

type YdbFeedSourceRepository struct {
	db *ydb.Driver
}

func NewYdbFeedSourceRepository(db *ydb.Driver) *YdbFeedSourceRepository {
	return &YdbFeedSourceRepository{db: db}
}

func (y *YdbFeedSourceRepository) AddSource(ctx context.Context, userId entities.UserId, source *entities.FeedSource) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("failed to parse user UUID: %w", err)
	}

	tx := query.TxControl(
		query.BeginTx(
			query.WithSerializableReadWrite(),
		),
	)

	// checks is source with such url already exists
	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $feed_url AS Utf8;

		SELECT s.id,
		FROM feed_sources s
		WHERE s.feed_url = $feed_url
		LIMIT 1;
		`,
		query.WithParameters(ydb.ParamsBuilder().Param("$feed_url").Text(source.FeedUrl).Build()),
		query.WithTxControl(tx),
	)

	var sourceUUID uuid.UUID

	if errors.Is(err, sql.ErrNoRows) {
		sourceUUID = uuid.New()

		// create feed source
		err := y.db.Query().Exec(
			ctx,
			`
			DECLARE $id AS Uuid;
			DECLARE $feed_url AS Utf8;
			DECLARE $type AS Utf8;

			INSERT INTO feed_sources (id, feed_url, type)
			VALUES ($id, $feed_url, $type);
			`,
			query.WithParameters(ydb.ParamsBuilder().
				Param("$id").Uuid(sourceUUID).
				Param("$feed_url").Text(source.FeedUrl).
				Param("$type").Text("default").
				Build(),
			),
			query.WithTxControl(tx),
		)

		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		if err := row.Scan(&sourceUUID); err != nil {
			return err
		}
	}

	// add user's feed source metadata
	err = y.db.Query().Exec(
		ctx,
		`
		DECLARE $user_id AS Uuid;
		DECLARE $source_id AS Uuid;

		INSERT INTO user_feed_sources (user_id, source_id)
		VALUES ($user_id, $source_id);`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$source_id").Uuid(sourceUUID).
			Build(),
		),
		query.WithTxControl(tx),
	)

	return err
}

func (y *YdbFeedSourceRepository) GetSources(ctx context.Context, userId entities.UserId) ([]*entities.FeedSource, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	// query
	result, err := y.db.Query().Query(
		ctx,
		`
		DECLARE $user_id AS Uuid;
		
		SELECT
			s.id,
			i.name,
			s.feed_url,
			i.disabled
		FROM source_user_infos i
		INNER JOIN feed_sources s
			ON i.source_id = s.id
		WHERE i.user_id = $user_id;
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = result.Close(ctx)
	}()

	// map
	return mapResultSet[entities.FeedSource](result, ctx)
}

func (y *YdbFeedSourceRepository) GetSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId) (*entities.FeedSource, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	sourceUUID, err := uuid.Parse(sourceId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source UUID: %w", err)
	}

	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $user_id AS Uuid;
		DECLARE $source_id AS Uuid;
	
		SELECT
			s.id,
			i.name,
			s.feed_url,
			s.type,
			i.disabled
		FROM feed_sources s
		INNER JOIN source_user_infos i
			ON i.source_id = $source_id AND i.user_id = $user_id
		LIMIT 1;
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$source_id").Uuid(sourceUUID).
			Build(),
		),
		query.WithTxControl(query.SnapshotReadOnlyTxControl()),
	)

	if err != nil {
		return nil, err
	}

	var feedSource entities.FeedSource
	if err := row.ScanStruct(&feedSource); err != nil {
		return nil, err
	}

	return &feedSource, nil
}

func (y *YdbFeedSourceRepository) UpdateSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId, patch *entities.FeedSourcePatch) error {
	if !patch.Name.HasValue() && !patch.Disabled.HasValue() {
		return fmt.Errorf("empty patch for update")
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("failed to parse user UUID: %w", err)
	}

	sourceUUID, err := uuid.Parse(sourceId)
	if err != nil {
		return fmt.Errorf("failed to parse source UUID: %w", err)
	}

	sqlBuilder := strings.Builder{}
	paramsBuilder := ydb.ParamsBuilder().
		Param("$user_id").Uuid(userUUID).
		Param("$source_id").Uuid(sourceUUID)

	sqlBuilder.WriteString("DECLARE $user_id AS Uuid;\nDECLARE $source_id AS Uuid;\n")

	if patch.Name.HasValue() {
		sqlBuilder.WriteString("DECLARE $name AS Utf8;\n")

		value, _ := patch.Name.Value()
		paramsBuilder = paramsBuilder.Param("$name").Text(value)
	}

	if patch.Disabled.HasValue() {
		sqlBuilder.WriteString("DECLARE $disabled AS Boolean;\n")

		value, _ := patch.Disabled.Value()
		paramsBuilder = paramsBuilder.Param("$disabled").Bool(value)
	}

	sqlBuilder.WriteString("UPDATE feed_source_user_infos\nSET ")

	if patch.Name.HasValue() {
		sqlBuilder.WriteString("name = $name")
	}

	if patch.Disabled.HasValue() {
		if patch.Name.HasValue() {
			sqlBuilder.WriteString(", ")
		}

		sqlBuilder.WriteString("disabled = $disabled")
	}

	err = y.db.Query().Exec(
		ctx,
		sqlBuilder.String(),
		query.WithParameters(paramsBuilder.Build()),
		query.WithTxControl(query.NoTx()),
	)

	return err
}

func (y *YdbFeedSourceRepository) DeleteSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("failed to parse user id: %w", err)
	}

	sourceUUID, err := uuid.Parse(sourceId)
	if err != nil {
		return fmt.Errorf("failed to parse source id: %w", err)
	}

	err = y.db.Query().Exec(
		ctx,
		`
		DECLARE $user_id AS Uuid;
		DECLARE $source_id AS Uuid;
		
		DELETE FROM user_feed_sources 
		WHERE user_id = $user_id AND source_id = $source_id;
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$source_id").Uuid(sourceUUID).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}
