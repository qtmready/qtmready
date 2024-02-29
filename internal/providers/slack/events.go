package slack

import (
	"context"
	"log/slog"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type (
	SocketEventHandler  func(context.Context, *socketmode.Event, *slack.Client, *socketmode.Client) error
	SockerEventHandlers map[socketmode.EventType]SocketEventHandler
)

var (
	helloMessageSent bool
)

func RunSocket() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go handleSocket(ctx, Instance().Client, Instance().Socket)

	return Instance().Socket.Run()
}

// handleSocket handles the socketmode events.
func handleSocket(ctx context.Context, client *slack.Client, socket *socketmode.Client) {
	handlers := SockerEventHandlers{
		socketmode.EventTypeEventsAPI:    handleEventTypeEventsAPI,
		socketmode.EventTypeSlashCommand: handleEventTypeSlashCommand,
		socketmode.EventTypeInteractive:  handleEventTypeInteractive,
		socketmode.EventTypeHello:        handleEventTypeHello,
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down socketmode listener...")
			return

		case event := <-socket.Events:
			if handler, ok := handlers[event.Type]; ok {
				// Print event and arguments for debugging
				slog.Info("handling event", slog.Any("event_type", event.Type), slog.Any("event", event))

				if err := handler(ctx, &event, client, socket); err != nil {
					slog.Error("unable to handle socket event", slog.Any("event", event.Type), slog.Any("error", err.Error()))
				} else {
					socket.Ack(*event.Request)
				}
			} else {
				slog.Warn("unknown event, ignoring ...", slog.Any("event", event.Type))
			}
		}
	}
}

func handleEventTypeHello(ctx context.Context, event *socketmode.Event, client *slack.Client, socket *socketmode.Client) error {
	if !helloMessageSent {
		// Extract necessary information from the event payload
		channelID := "C05J9NXGM1P"
		message := "Message to post in the channel"

		// Call Slack API to post the message
		if err := notify(client, message, channelID); err != nil {
			slog.Info("Failed to post message to channel", slog.Any("e", err))
			return err
		}

		helloMessageSent = true
	}

	return nil
}

// handleEventTypeEventsAPI handles the events API event.
func handleEventTypeEventsAPI(ctx context.Context, event *socketmode.Event, client *slack.Client, socket *socketmode.Client) error {
	_, ok := event.Data.(slackevents.EventsAPIEvent)

	if !ok {
		return NewSocketEventPayloadError(event.Type)
	}

	return nil
}

// handleEventTypeSlashCommand handles the slash command event.
func handleEventTypeSlashCommand(ctx context.Context, event *socketmode.Event, client *slack.Client, socket *socketmode.Client) error {
	_, ok := event.Data.(slack.SlashCommand)
	if !ok {
		return NewSocketEventPayloadError(event.Type)
	}

	return nil
}

func handleEventTypeInteractive(ctx context.Context, event *socketmode.Event, client *slack.Client, socket *socketmode.Client) error {
	// Extract the callback from the event
	_, ok := event.Data.(slack.InteractionCallback)

	if !ok {
		return NewSocketEventPayloadError(event.Type)
	}

	return nil
}

// handleEventsAPIPayload handles the events API payload.
func handleEventsAPIPayload(event slackevents.EventsAPIEvent) error {
	if event.Type == slackevents.CallbackEvent {
		if e, ok := event.InnerEvent.Data.(*slackevents.AppMentionEvent); ok {
			slog.Info("handleEventMessage event", slog.Any("event", e))
		} else {
			return ErrInvalidEventPayload
		}
	} else {
		return ErrInvalidEvent
	}

	return nil
}
