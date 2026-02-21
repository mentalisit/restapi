package models

type IncomingMessage struct {
	Text         string
	DmChat       string
	Name         string
	MentionName  string
	NameId       string
	NickName     string
	Avatar       string
	AvatarF      string
	ChannelId    string
	GuildId      string
	GuildName    string
	GuildAvatar  string
	GuildAvatarF string
	Type         string
	Language     string
}

type DsMembersRoles struct {
	Userid  string   `json:"userid"`
	RolesId []string `json:"rolesId"`
}

// CorpRole represents a corporation role data structure
type CorpRole struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
