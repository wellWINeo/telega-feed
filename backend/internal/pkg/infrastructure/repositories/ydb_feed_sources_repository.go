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
)

type YdbFeedSourceRepository struct {
	db *ydb.Driver
}

func NewYdbFeedSourceRepository(db *ydb.Driver) *YdbFeedSourceRepository {
	return &YdbFeedSourceRepository{db: db}
}

func (y *YdbFeedSourceRepository) AddSource(ctx context.Context, userId entities.UserId, source *entities.FeedSource) (entities.FeedSourceId, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	var sourceId entities.FeedSourceId

	err = y.db.Table().DoTx(ctx, func(ctx context.Context, tx table.TransactionActor) error {
		result, err := tx.Execute(ctx, selectFeedSourceByFeedUrlSql, table.NewQueryParameters(
			table.ValueParam("$feed_url", types.StringValueFromString(source.FeedUrl)),
		))

		defer func() {
			_ = result.Close()
		}()

		if err != nil {
			return err
		}

		if !result.NextResultSet(ctx) || result.CurrentResultSet().RowCount() == 0 {
			sourceId = uuid.New()

			result, err := tx.Execute(ctx, insertFeedSourceSql, table.NewQueryParameters(
				table.ValueParam("$id", types.UuidValue(sourceId)),
				table.ValueParam("$feed_url", types.StringValueFromString(source.FeedUrl)),
				table.ValueParam("$type", types.StringValueFromString(source.Type)),
			))

			if err != nil {
				return err
			}

			_ = result.Close()
		} else {
			if !result.NextRow() {
				return fmt.Errorf("no row found")
			}

			if err := result.ScanWithDefaults(&sourceId); err != nil {
				return err
			}
		}

		result, err = tx.Execute(ctx, insertFeedSourceUserInfoSql, table.NewQueryParameters(
			table.ValueParam("$source_id", types.UuidValue(sourceId)),
			table.ValueParam("$user_id", types.UuidValue(userUUID)),
			table.ValueParam("$name", types.StringValueFromString(source.Name)),
		))

		if err != nil {
			return err
		}

		_ = result.Close()

		return nil

	}, table.WithTxSettings(table.TxSettings(table.WithSerializableReadWrite())))

	if err != nil {
		return uuid.Nil, err
	}

	return sourceId, nil
}

func (y *YdbFeedSourceRepository) GetSources(ctx context.Context, userId entities.UserId) ([]*entities.FeedSource, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	// query
	result, err := y.db.Query().Query(ctx, selectFeedSourcesByUserIdSql,
		query.WithParameters(ydb.ParamsBuilder().Param("$user_id").Uuid(userUUID).Build()),
		query.WithTxControl(query.SnapshotReadOnlyTxControl()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query feed_sources: %w", err)
	}

	return mapResultSet[entities.FeedSource](result, ctx)
}

func (y *YdbFeedSourceRepository) GetSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId) (*entities.FeedSource, error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user UUID: %w", err)
	}

	row, err := y.db.Query().QueryRow(ctx, selectFeedSourceByUserIdAndSourceIdSql,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$source_id").Uuid(sourceId).
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

func (y *YdbFeedSourceRepository) GetSourcesForFeedUpdate(ctx context.Context) ([]*entities.FeedSource, error) {
	result, err := y.db.Query().
		Query(ctx, selectNotDisabledFeedSourcesSql, query.WithTxControl(query.SnapshotReadOnlyTxControl()))

	if err != nil {
		return nil, err
	}

	defer func() { _ = result.Close(ctx) }()

	return mapResultSet[entities.FeedSource](result, ctx)
}

func (y *YdbFeedSourceRepository) UpdateSource(ctx context.Context, userId entities.UserId, sourceId entities.FeedSourceId, patch *entities.FeedSourcePatch) error {
	if !patch.Name.HasValue() && !patch.Disabled.HasValue() {
		return fmt.Errorf("empty patch for update")
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("failed to parse user UUID: %w", err)
	}

	sqlBuilder := strings.Builder{}
	paramsBuilder := ydb.ParamsBuilder().
		Param("$user_id").Uuid(userUUID).
		Param("$source_id").Uuid(sourceId)

	sqlBuilder.WriteString("DECLARE $user_id AS Uuid;\nDECLARE $source_id AS Uuid;\n")

	if patch.Name.HasValue() {
		sqlBuilder.WriteString("DECLARE $name AS String;\n")

		value, _ := patch.Name.Value()
		paramsBuilder = paramsBuilder.Param("$name").Any(types.StringValueFromString(value))
	}

	if patch.Disabled.HasValue() {
		sqlBuilder.WriteString("DECLARE $disabled AS Bool;\n")

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

	sqlBuilder.WriteString("\nWHERE user_id = $user_id AND source_id = $source_id;")

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

	err = y.db.Query().Exec(
		ctx,
		`
		DECLARE $user_id AS Uuid;
		DECLARE $source_id AS Uuid;
		
		DELETE FROM feed_source_user_infos 
		WHERE user_id = $user_id AND source_id = $source_id;
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$user_id").Uuid(userUUID).
			Param("$source_id").Uuid(sourceId).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}

func (y *YdbFeedSourceRepository) DeleteOrphanedSources(ctx context.Context) error {
	return y.db.Table().DoTx(ctx, func(ctx context.Context, tx table.TransactionActor) error {
		result, err := tx.Execute(ctx, `
			DELETE FROM feed_sources ON
			SELECT *
			FROM feed_sources s
			LEFT ONLY feed_source_user_infos i
				ON s.id = i.source_id
			`, table.NewQueryParameters())

		defer func() { _ = result.Close() }()

		return err
	}, table.WithTxSettings(table.TxSettings(table.WithSnapshotReadOnly())))
}
