package slack

import (
	"log/slog"

	"github.com/slack-go/slack"
)

func NotifyOnSlack(message string) error {
	return handleMessageEvent(SlackClient(), message)
}

func handleMessageEvent(client *slack.Client, message string) error {
	channelID := "C05HQE6M41L" // TODO: get the channel_id from database

	if err := notify(client, channelID, message); err != nil {
		slog.Info("Failed to post message to channel", slog.Any("e", err))
		return err
	}

	return nil
}
