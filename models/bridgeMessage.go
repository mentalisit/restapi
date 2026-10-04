package models

import (
	"encoding/json"
	"strings"
)

// Режимы канала моста: канал служит источником сообщений, приёмником или тем и
// другим.
const (
	ModeBoth  = "both"
	ModeRead  = "read"
	ModeWrite = "write"
)

type ToBridgeMessage struct {
	Text          string              `json:"text"`
	Sender        string              `json:"sender"`
	SenderId      string              `json:"senderId"`
	Tip           string              `json:"tip"`
	ChatId        string              `json:"chatId"`
	MesId         string              `json:"mesId"`
	GuildId       string              `json:"guildId"`
	TimestampUnix int64               `json:"timestampUnix"`
	Extra         []FileInfo          `json:"extra"`
	Avatar        string              `json:"avatar"`
	Reply         *BridgeMessageReply `json:"reply"`
	ReplyMap      map[string]string   `json:"replyMap"`
	ConfigID      string              `json:"configId"` // UUIDv7
	Config        *BridgeConfig       `json:"config"`   // Resolved from in-memory cache
}
type FileInfo struct {
	Name   string `json:"name"`
	Data   []byte `json:"data"`
	URL    string `json:"URL"`
	Size   int64  `json:"size"`
	FileID string `json:"fileID"`
}
type BridgeMessageReply struct {
	TimeMessage int64  `json:"time_message"`
	Text        string `json:"text"`
	Avatar      string `json:"avatar"`
	UserName    string `json:"userName"`
}

// BridgeConfig — конфигурация реле.
//
// В bot.bridge_config настройки лежат в колонке options отдельным JSON, а здесь
// они разложены по полям структуры: структуру из модуля restapi используют ещё
// и мессенджеры, поэтому менять её форму нельзя. Options заполняется из
// колонки методом ParseOptions, остальные поля приходят из самой строки.
type BridgeConfig struct {
	Id      int                        `json:"id"`
	UUID    string                     `json:"uuid"` // uuidV7
	Host    string                     `json:"host"`
	Channel map[string][]BridgeChannel `json:"channel"`
	Options BridgeOptions              `json:"options"`
}

// BridgeOptions — настройки реле в том виде, в котором они лежат в колонке
// options bot.bridge_config.
//
// Отдельная структура нужна, чтобы не плодить колонки: настройка — это одно
// поле JSON, и её добавление не требует ни ALTER TABLE, ни правок во всех
// сервисах сразу. Неизвестные ключи игнорируются, поэтому настройку может
// прочитать только мессенджер или только будущая версия моста.
type BridgeOptions struct {
	NameRelay         string   `json:"name_relay"`
	Role              []string `json:"role"`
	ForbiddenPrefixes []string `json:"forbidden_prefixes"`
	EnablePolls       bool     `json:"enable_polls"`
	ShowAvatars       bool     `json:"show_avatars"`
	ShowBadges        bool     `json:"show_badges"`
	NativeReply       bool     `json:"native_reply"`
	// ReplyPrefix — nil означает значение по умолчанию, пустая строка — цитату
	// без подписи.
	ReplyPrefix       *string `json:"reply_prefix"`
	SplitLongMessages bool    `json:"split_long_messages"`
	ForwardedPrefix   string  `json:"forwarded_prefix"`
	IgnoreBots        bool    `json:"ignore_bots"`
}

// DefaultForwardedPrefix — подпись пересылки, пока она не задана в реле.
const DefaultForwardedPrefix = "Forwarded from"

// DefaultBridgeOptions — значения по умолчанию.
//
// Разбор идёт поверх них, поэтому отсутствующие в options ключи получают эти
// значения, а не нулевые: show_avatars и show_badges выключили бы аватары и
// значки молча, без ошибки.
func DefaultBridgeOptions() BridgeOptions {
	return BridgeOptions{
		Role:              []string{},
		ForbiddenPrefixes: []string{},
		ShowAvatars:       true,
		ShowBadges:        true,
		ForwardedPrefix:   DefaultForwardedPrefix,
	}
}

// ParseOptions разбирает содержимое колонки options и заполняет настройки
// конфигурации. Отсутствующие ключи остаются со значениями по умолчанию.
func (c *BridgeConfig) ParseOptions(data []byte) error {
	options := DefaultBridgeOptions()
	if len(data) > 0 {
		if err := json.Unmarshal(data, &options); err != nil {
			return err
		}
	}
	c.Options = options
	return nil
}

// BridgeChannel — канал реле: один получатель или источник сообщений.
type BridgeChannel struct {
	ChannelId       string            `json:"channel_id"`
	GuildId         string            `json:"guild_id"`
	CorpChannelName string            `json:"corp_channel_name"`
	AliasName       string            `json:"alias_name"`
	MappingRoles    map[string]string `json:"mapping_roles"`

	EnablePolls bool              `json:"enable_polls"`
	AllowMedia  *bool             `json:"allow_media"`
	Mode        string            `json:"mode"` // "both" | "read" | "write"
	Settings    map[string]string `json:"settings"`
}

// MediaAllowed сообщает, можно ли пересылать вложения канала.
//
// Указатель, а не bool: у старого моста такого ограничения не было вовсе, и
// конфигурации, перенесённые из rs_bot.bridge_config, ключа allow_media не
// содержат. С bool такие каналы молча теряли картинки — ключ нужно было
// проставить в каждом реле вручную.
func (c BridgeChannel) MediaAllowed() bool {
	return c.AllowMedia == nil || *c.AllowMedia
}

// RelayMode — режим канала. Пустое и неизвестное значение равны "both": каналы,
// настроенные до появления режимов, работают как раньше.
func (c BridgeChannel) RelayMode() string {
	switch strings.ToLower(strings.TrimSpace(c.Mode)) {
	case ModeRead:
		return ModeRead
	case ModeWrite:
		return ModeWrite
	default:
		return ModeBoth
	}
}

// Readable сообщает, что канал служит источником сообщений. В "write" бот
// только отправляет, поэтому читать из такого канала нечего.
func (c BridgeChannel) Readable() bool {
	return c.RelayMode() != ModeWrite
}

// Writable сообщает, что в канал можно отправлять. В "read" канал только
// читается, запись в него запрещена.
func (c BridgeChannel) Writable() bool {
	return c.RelayMode() != ModeRead
}

// DefaultReplyPrefix — подпись цитаты, когда она указана в конфигурации.
const DefaultReplyPrefix = "Reply:"

// QuotePrefix возвращает подпись цитаты для реле: nil — значение по умолчанию,
// явно пустая строка — цитата без подписи.
func (c BridgeConfig) QuotePrefix() string {
	if c.Options.ReplyPrefix == nil {
		return DefaultReplyPrefix
	}

	return strings.TrimSpace(*c.Options.ReplyPrefix)
}

// ReplyQuote собирает строку цитаты для текста: подпись из конфигурации и сам
// текст. Пустая подпись отправляет цитату без неё.
func ReplyQuote(prefix, text string) string {
	quote := strings.TrimSpace(text)
	if prefix == "" {
		return quote
	}

	return prefix + " " + quote
}

type BridgeSendToMessenger struct {
	Text      string              `json:"text"`
	Sender    string              `json:"sender"`
	ChannelId []string            `json:"channelId"`
	Avatar    string              `json:"avatar"`
	Extra     []FileInfo          `json:"extra"`
	Reply     *BridgeMessageReply `json:"reply"`
	ReplyMap  map[string]string   `json:"replyMap"`
}
type MessageIds struct {
	MessageId string
	ChatId    string
}

type Request struct {
	Data    map[string]string `json:"data"`
	Options []string          `json:"options"`
}
