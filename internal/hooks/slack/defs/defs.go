package defs

import (
	"encoding/json"
)

// Kind constants.
const (
	KindBot  = "bot"
	KindUser = "user"
)

type (
	MessageProviderSlackData struct {
		BotToken      string `json:"bot_token"`
		ChannelID     string `json:"channel_id"`
		ChannelName   string `json:"channel_name"`
		WorkspaceID   string `json:"workspace_id"`
		WorkspaceName string `json:"workspace_name"`
	}

	MessageProviderSlackUserInfo struct {
		BotToken       string `json:"bot_token"`
		UserToken      string `json:"user_token"`
		ProviderUserID string `json:"provider_user_id"`
		ProviderTeamID string `json:"provider_team_id"`
	}

	MessageProviderData interface {
		Marshal() ([]byte, error)
		Unmarshal(data []byte) error
	}
)

// Implement Marshal and Unmarshal for MessageProviderSlackData.
func (m *MessageProviderSlackData) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageProviderSlackData) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}

// Implement Marshal and Unmarshal for MessageProviderSlackUserInfo.
func (m *MessageProviderSlackUserInfo) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *MessageProviderSlackUserInfo) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
