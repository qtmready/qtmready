package activities

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/db"
	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/slack/cast"
	"go.breu.io/quantm/internal/hooks/slack/config"
	"go.breu.io/quantm/internal/hooks/slack/fns"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

type (
	// Activities groups all the activities for the slack provider.
	Activities struct{}
)

const (
	footer = "Powered by quantm.io"
)

func (a *Activities) NotifyLinesExceed(
	ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Diff],
) error {
	var err error

	token := ""
	target := ""

	if event.Subject.UserID != uuid.Nil {
		token, target, err = a.to_user(ctx, event.Subject.UserID)
		if err != nil {
			return err
		}
	} else {
		token, target, err = a.to_repo(ctx, event.Subject.ID)
		if err != nil {
			return err
		}
	}

	attachment := slack.Attachment{
		Color:      "warning",
		Pretext:    "The number of lines in this pull request exceeds the allowed threshold. Please review and adjust accordingly.",
		Fallback:   "Line Exceed Detected",
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Fields:     fns.LineExceedFields(event),
		Ts:         json.Number(strconv.FormatInt(time.Now().Unix(), 10)),
	}

	client, err := config.GetSlackClient(token)
	if err != nil {
		return err
	}

	return fns.SendMessage(client, target, attachment)
}

func (a *Activities) to_user(ctx context.Context, link_to uuid.UUID) (string, string, error) {
	msg, err := db.Queries().GetMessagesByLinkTo(ctx, link_to)
	if err != nil {
		return "", "", err
	}

	d, err := cast.ByteToMessageProviderSlackUserInfo(msg.Data)
	if err != nil {
		return "", "", err
	}

	token, err := fns.Reveal(d.BotToken, d.ProviderTeamID)
	if err != nil {
		return "", "", err
	}

	return token, d.ProviderUserID, nil
}

func (a *Activities) to_repo(ctx context.Context, link_to uuid.UUID) (string, string, error) {
	msg, err := db.Queries().GetMessagesByLinkTo(ctx, link_to)
	if err != nil {
		return "", "", err
	}

	d, err := cast.ByteToMessageProviderSlackData(msg.Data)
	if err != nil {
		return "", "", err
	}

	token, err := fns.Reveal(d.BotToken, d.WorkspaceID)
	if err != nil {
		return "", "", err
	}

	return token, d.ChannelID, nil
}
