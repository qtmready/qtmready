package slack

import (
	"log/slog"

	"github.com/slack-go/slack"
)

func notify(client *slack.Client, channelID string) error {
	attachment := slack.Attachment{
		Pretext: "Github",
		Text:    "Hello from quantm!",
	}

	// Send message
	_, _, err := client.PostMessage(
		channelID,
		slack.MsgOptionText("Hello", false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)

	if err != nil {
		slog.Error("Error sending message to channel ", channelID, ": ", err)
		return err
	}

	return nil
}
