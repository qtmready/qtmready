package activities

import (
	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/events"
	"go.breu.io/quantm/internal/hooks/slack/attach"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func fields_lines_exceeded(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		attach.Repo(event),
		attach.Branch(event),
		attach.Threshold(),
		attach.TotalLinesCount(event),
		attach.LinesAdded(event),
		attach.LinesDeleted(event),
		attach.AddedFiles(event),
		attach.DeletedFiles(event),
		attach.ModifiedFiles(event),
		attach.RenameFiles(event),
	}

	return fields
}

func fields_merge_conflict(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		attach.Repo(event),
		attach.BranchMerge(event),
		attach.CurrentHead(event),
		attach.ConflictHead(),
		attach.AffectedFiles(),
	}

	return fields
}
