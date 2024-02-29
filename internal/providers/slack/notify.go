package slack

import (
	"github.com/slack-go/slack"
)

func notify(api *slack.Client, message, channelID string) error {
	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		return err
	}

	return nil
}
