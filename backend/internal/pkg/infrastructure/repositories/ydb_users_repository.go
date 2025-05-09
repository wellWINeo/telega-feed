package repositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
	"io"
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

func (y *YdbUsersRepository) GetUserById(ctx context.Context, userID entities.UserId) (*entities.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	row, err := y.db.Query().QueryRow(ctx, `
		DECLARE $user_id AS Uuid;
		
		SELECT u.id, u.telegram_id FROM users u WHERE u.id = $user_id;
		`, query.WithParameters(ydb.ParamsBuilder().Param("$user_id").Uuid(userUUID).Build()),
		query.WithTxControl(query.NoTx()),
	)

	if errors.Is(err, io.EOF) {
		return nil, nil
	}

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

	err := y.db.Query().Exec(
		ctx,
		`
		DECLARE $id AS Uuid;
		DECLARE $telegram_id AS String;

		INSERT INTO users (id, telegram_id)
		VALUES ($id, $telegram_id);
		`,
		query.WithParameters(ydb.ParamsBuilder().
			Param("$id").Uuid(userUUID).
			Param("$telegram_id").Text(user.TelegramId).
			Build(),
		),
		query.WithTxControl(query.NoTx()),
	)

	return err
}
