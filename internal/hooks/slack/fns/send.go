package fns

import (
	"log/slog"

	"github.com/slack-go/slack"
)

func SendMessage(client *slack.Client, channelID string, attachment slack.Attachment) error {
	// Send message
	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		slog.Error("Error sending message to channel ", ": ", slog.Any("e", err))
		return err
	}

	return nil
}
