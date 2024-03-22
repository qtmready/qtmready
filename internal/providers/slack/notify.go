package slack

import (
	"log/slog"

	"github.com/slack-go/slack"
)

func notify(client *slack.Client, channelID, message string) error {
	attachment := slack.Attachment{
		Color: "danger",
		Title: "Message from quantm",
		Text:  message,
	}

	// Send message
	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		slog.Error("Error sending message to channel ", channelID, ": ", err)
		return err
	}

	return nil
}
