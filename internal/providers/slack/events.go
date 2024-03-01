package slack

import (
	"log/slog"

	"github.com/slack-go/slack"
)

func RunSlack() error {
	return handleMessageEvent(Instance().Client)
}

func handleMessageEvent(client *slack.Client) error {
	channelID := "C05J9NXGM1P" // get the channel_id from database

	if err := notify(client, channelID); err != nil {
		slog.Info("Failed to post message to channel", slog.Any("e", err))
		return err
	}

	return nil
}
