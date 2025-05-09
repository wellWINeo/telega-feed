package entities

import "github.com/google/uuid"

type UserId = string

type User struct {
	Id         uuid.UUID `json:"id,omitempty" sql:"id"`
	TelegramId string    `json:"telegram_id,omitempty" sql:"telegram_id"`
}
