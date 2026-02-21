package models

import (
	"time"

	"github.com/google/uuid"
)

type Button struct {
	Text string
	Data string
}

type Timer struct {
	//Id       string `bson:"_id"`
	Tip    string `bson:"tip"`
	ChatId string `bson:"chatId"`
	MesId  string `bson:"mesId"`
	Timed  int    `bson:"timed"`
}

type EmodjiUser struct {
	Id                                int
	Tip, Name, Em1, Em2, Em3, Em4     string
	Module1, Module2, Module3, Weapon string
}

// MultiAccount представляет мульти-аккаунт пользователя
type MultiAccount struct {
	UUID             uuid.UUID `json:"uuid"`
	Nickname         string    `json:"nickname"`
	TelegramID       string    `json:"telegram_id"`
	TelegramUsername string    `json:"telegram_username"`
	DiscordID        string    `json:"discord_id"`
	DiscordUsername  string    `json:"discord_username"`
	WhatsappID       string    `json:"whatsapp_id"`
	WhatsappUsername string    `json:"whatsapp_username"`
	CreatedAt        time.Time `json:"created_at"`
	AvatarURL        string    `json:"avatar_url"`
	Alts             []string  `json:"alts"`
}

type Tech struct {
	Uid      uuid.UUID
	Username string
	Tech     []byte
}

type MultiAccountCorpMember struct {
	Uid        uuid.UUID   `db:"uid"`
	GuildIds   []uuid.UUID `db:"guildids"`
	TimeZona   string      `db:"timezona"`
	ZonaOffset int         `db:"zonaoffset"`
	AfkFor     string      `db:"afkfor"`
}
