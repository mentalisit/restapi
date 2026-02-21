package models

import (
	"fmt"
	"slices"
	"strconv"
	"time"
)

type InMessageV2 struct {
	Text        string
	Tip         string
	NameNick    string
	Username    string
	UserId      string
	NameMention string
	Messenger   Info
	Config      CorporationConfigV2
	Options     Options
}

type Options []string

func (o *Options) Contains(s string) bool {
	return slices.Contains(*o, s)
}
func (o *Options) Remove(s string) {
	for i, item := range *o {
		if item == s {
			*o = append((*o)[:i], (*o)[i+1:]...)
			return
		}
	}
}
func (o *Options) Add(s string) {
	if o.Contains(s) {
		panic(fmt.Sprintf("option %s already add to opt %+v\n", s, *o))
	}
	*o = append(*o, s)
}

func (i Info) ToMap() map[string]string {
	m := make(map[string]string)

	if i.TypeMessenger != "" {
		m[MType] = i.TypeMessenger
	}
	if i.MessageId != "" {
		m[MMId] = i.MessageId
	}
	if i.ChannelId != "" {
		m[MChId] = i.ChannelId
	}
	if i.ChannelName != "" {
		m[MChName] = i.ChannelName
	}
	if i.GuildId != "" {
		m[MGuId] = i.GuildId
	}
	if i.GuildName != "" {
		m[MGuName] = i.GuildName
	}
	if i.GuildAvatarUrl != "" {
		m[MGuAvatarUrl] = i.GuildAvatarUrl
	}
	if i.Alias != "" {
		m[MAlias] = i.Alias
	}
	if i.UserAvatarUrl != "" {
		m[MUsAvatarUrl] = i.UserAvatarUrl
	}
	if i.GameCorporation != "" {
		m[MGC] = i.GameCorporation
	}
	if i.GameCorporationId != "" {
		m[MGCId] = i.GameCorporationId
	}
	if i.Language != "" {
		m[MLang] = i.Language
	}
	if i.ConfigParamId != 0 {
		m[MConfPar] = strconv.Itoa(i.ConfigParamId)
	}
	if !i.CreatedAt.IsZero() {
		m[MCreAt] = i.CreatedAt.Format(time.RFC3339)
	}

	// Handle Options field - convert map[string]bool to comma-separated string
	if len(i.Options) > 0 {
		optionsStr := ""
		for option, enabled := range i.Options {
			if enabled {
				if optionsStr != "" {
					optionsStr += ","
				}
				optionsStr += option
			}
		}
		if optionsStr != "" {
			m[MOptions] = optionsStr
		}
	}

	return m
}

const (
	MType        = "TypeMessenger"
	MMId         = "MessageId"
	MChId        = "ChannelId"
	MChName      = "ChannelName"
	MGuId        = "GuildId"
	MGuName      = "GuildName"
	MGuAvatarUrl = "GuildAvatarUrl"
	MAlias       = "Alias"
	MUsAvatarUrl = "UserAvatarUrl"
	MGC          = "GameCorporation"
	MGCId        = "GameCorporationId"
	MLang        = "Language"
	MConfPar     = "ConfigParamId"
	MCreAt       = "CreatedAt"
	MOptions     = "Options"
)
const (
	OptionReaction        = "Reaction"
	OptionInClient        = "InClient"
	OptionQueue           = "Queue"
	OptionPl30            = "Pl30"
	OptionMinusMin        = "MinusMin"
	OptionMinusMinNext    = "MinusMinNext"
	OptionEdit            = "Edit"
	OptionUpdate          = "Update"
	OptionUpdateAutoHelp  = "UpdateAutoHelp"
	OptionMessageUpdateDS = "MessageUpdateDS"
	OptionMessageUpdateTG = "MessageUpdateTG"
	OptionElseTrue        = "ElseTrue"
	OptionQueueAll        = "QueueAll"
	OptionPlus            = "Plus"
	OptionAutoHelp        = "AutoHelp"
	OptionCleanChat       = "CleanChat"
)

type CorporationConfigV2 struct {
	Uid         string
	Channels    ChannelsMap
	HelpMessage HelpMessage
}
type HelpMessage map[string]*Info

type ChannelsMap map[string]*Info

type Info struct {
	TypeMessenger     string          `json:"TypeMessenger,omitempty"`
	MessageId         string          `json:"MessageId,omitempty"`
	ChannelId         string          `json:"ChannelId,omitempty"`
	ChannelName       string          `json:"ChannelName,omitempty"`
	GuildId           string          `json:"GuildId,omitempty"`
	GuildName         string          `json:"GuildName,omitempty"`
	GuildAvatarUrl    string          `json:"GuildAvatarUrl,omitempty"`
	Alias             string          `json:"Alias,omitempty"`
	UserAvatarUrl     string          `json:"UserAvatarUrl,omitempty"`
	GameCorporation   string          `json:"GameCorporation,omitempty"`
	GameCorporationId string          `json:"GameCorporationId,omitempty"`
	Language          string          `json:"Language,omitempty"`
	ConfigParamId     int             `json:"ConfigParamId,omitempty"`
	CreatedAt         time.Time       `json:"CreatedAt,omitempty"`
	Options           map[string]bool `json:"Options,omitempty"`
}

type CorpInfo struct {
	ID         int64     `json:"id" db:"id"`
	ConfigName string    `json:"config_name" db:"config_name"`
	CorpName   string    `json:"corp_name" db:"corp_name"`
	CorpID     string    `json:"corp_id" db:"corp_id"`
	Level      int       `json:"level" db:"level"`
	Percent    int       `json:"percent" db:"percent"`
	XP         int       `json:"xp" db:"xp"`
	Webhook    bool      `json:"webhook" db:"webhook"`
	DateEnded  time.Time `json:"date_ended" db:"date_ended"`
	LastUpdate time.Time `json:"last_update" db:"last_update"`
}
