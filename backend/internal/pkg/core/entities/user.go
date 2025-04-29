package entities

type UserId = string

type User struct {
	Id         UserId `json:"id,omitempty"`
	TelegramId int64  `json:"telegram_id,omitempty"`
}
