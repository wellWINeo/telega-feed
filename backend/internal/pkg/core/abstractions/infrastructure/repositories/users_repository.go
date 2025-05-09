package abstractrepositories

import (
	"TelegaFeed/internal/pkg/core/entities"
	"context"
)

type UsersRepository interface {
	GetUserByTelegramId(ctx context.Context, telegramId string) (*entities.User, error)
	GetUserById(ctx context.Context, userID entities.UserId) (*entities.User, error)
	AddUser(ctx context.Context, user *entities.User) error
}
