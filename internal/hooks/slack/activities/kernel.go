package activities

import (
	"context"
	"encoding/json"
	"fmt"
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
	Kernel struct{}
)

const (
	footer = "Powered by quantm.io"
)

var (
	ts = json.Number(strconv.FormatInt(time.Now().Unix(), 10))
)

func (k *Kernel) NotifyLinesExceed(
	ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Diff],
) error {
	var err error

	token := ""
	target := ""

	if event.Subject.UserID != uuid.Nil {
		token, target, err = k.to_user(ctx, event.Subject.UserID)
		if err != nil {
			return err
		}
	} else {
		token, target, err = k.to_repo(ctx, event.Subject.ID)
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
		Ts:         ts,
	}

	client, err := config.GetSlackClient(token)
	if err != nil {
		return err
	}

	return fns.SendMessage(client, target, attachment)
}

func (k *Kernel) NotifyMergeConflict(
	ctx context.Context, event *events.Event[eventsv1.ChatHook, eventsv1.Merge],
) error {
	var err error

	token := ""
	target := ""

	if event.Subject.UserID != uuid.Nil {
		token, target, err = k.to_user(ctx, event.Subject.UserID)
		if err != nil {
			return err
		}
	} else {
		token, target, err = k.to_repo(ctx, event.Subject.ID)
		if err != nil {
			return err
		}
	}

	client, err := config.GetSlackClient(token)
	if err != nil {
		return err
	}

	attachment := slack.Attachment{
		Color: "warning",
		Pretext: fmt.Sprintf(`We've detected a merge conflict in your feature branch, <%s/tree/%s|%s>.
    This means there are changes in your branch that clash with recent updates on the main branch (trunk).`,
			event.Context.Source, event.Payload.HeadBranch, event.Payload.HeadBranch),
		Fallback:   "Merge Conflict Detected",
		MarkdownIn: []string{"fields"},
		Footer:     footer,
		Fields:     fns.MergeConflictFields(event),
		Ts:         ts,
	}

	return fns.SendMessage(client, target, attachment)
}

func (k *Kernel) to_user(ctx context.Context, link_to uuid.UUID) (string, string, error) {
	msg, err := db.Queries().GetChatLink(ctx, link_to)
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

func (k *Kernel) to_repo(ctx context.Context, link_to uuid.UUID) (string, string, error) {
	msg, err := db.Queries().GetChatLink(ctx, link_to)
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
