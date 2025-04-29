package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"strconv"
)

type YdbUsersRepository struct {
	db *ydb.Driver
}

func NewYdbUsersRepository(db *ydb.Driver) *YdbUsersRepository {
	return &YdbUsersRepository{db: db}
}

func (y *YdbUsersRepository) GetUserByTelegramId(ctx context.Context, telegramId string) (*entities.User, error) {
	row, err := y.db.Query().QueryRow(
		ctx,
		`
		DECLARE $telegram_id AS Utf8;

		SELECT u.id, u.telegram_id 
		FROM users u
		WHERE u.telegram_id = $telegram_id;
		`,
		query.WithParameters(ydb.ParamsBuilder().Param("$telegram_id").Text(telegramId).Build()),
		query.WithTxControl(query.NoTx()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	var user entities.User

	if err := row.ScanStruct(&user); err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &user, nil
}

func (y *YdbUsersRepository) AddUser(ctx context.Context, user *entities.User) error {
	userUUID := uuid.New()
	telegramID := strconv.FormatInt(user.TelegramId, 10)

	err := y.db.Query().Exec(
		ctx,
		``,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$id").Uuid(userUUID).
			Param("$telegram_id").Text(telegramID).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}
