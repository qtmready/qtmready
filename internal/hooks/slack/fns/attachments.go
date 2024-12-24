package fns

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

func LineExceedFields(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		create_attachment_repository(event),
		create_attachment_branch(event),
		create_attachment_threshold(),
		create_attachment_total_lines_count(event),
		create_attachment_lines_added(event),
		create_attachment_lines_deleted(event),
		create_attachment_added_files(event),
		create_attachment_deleted_files(event),
		create_attachment_modified_files(event),
		create_attachment_rename_files(event),
	}

	return fields
}

func MergeConflictFields(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) []slack.AttachmentField {
	fields := []slack.AttachmentField{
		create_attachment_repository(event),
		create_attachment_branch_merge(event),
		create_attachment_current_head(event),
		create_attachment_conflict_head(),
		create_attachment_affected_files(),
	}

	return fields
}

func extract_repo(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return parts[len(parts)-1]
}

func format_files(files []string) string {
	result := ""
	for _, file := range files {
		result += "- " + file + "\n"
	}

	return result
}

func format_renamed_files(files []*eventsv1.RenamedFile) string {
	result := ""
	for _, file := range files {
		result += fmt.Sprintf("- %s -> %s\n", file.GetOld(), file.GetNew())
	}

	return result
}

func create_attachment_repository[E events.Payload](event *events.Event[eventsv1.ChatHook, E]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Repository*",
		Value: fmt.Sprintf("<%s|%s>", event.Context.Source, extract_repo(event.Context.Source)),
		Short: true,
	}
}

func create_attachment_branch(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, "", ""),
		Short: true,
	}
}

func create_attachment_branch_merge(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, event.Payload.BaseBranch, event.Payload.BaseBranch),
		Short: true,
	}
}

func create_attachment_threshold() slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Threshold*",
		Value: fmt.Sprintf("%d", 0),
		Short: true,
	}
}

func create_attachment_total_lines_count(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Total Lines Count*",
		Value: fmt.Sprintf("%d", event.Payload.GetLines().GetAdded()+event.Payload.GetLines().GetRemoved()),
		Short: true,
	}
}

func create_attachment_lines_added(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Lines Added*",
		Value: fmt.Sprintf("%d", event.Payload.GetLines().GetAdded()),
		Short: true,
	}
}

func create_attachment_lines_deleted(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Lines Deleted*",
		Value: fmt.Sprintf("%d", event.Payload.GetLines().GetRemoved()),
		Short: true,
	}
}

func create_attachment_added_files(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Added Files",
		Value: format_files(event.Payload.GetFiles().GetAdded()),
		Short: false,
	}
}

func create_attachment_deleted_files(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Deleted Files",
		Value: format_files(event.Payload.GetFiles().GetDeleted()),
		Short: false,
	}
}

func create_attachment_modified_files(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Modified Files",
		Value: format_files(event.Payload.GetFiles().GetModified()),
		Short: false,
	}
}

func create_attachment_rename_files(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Rename Files",
		Value: format_renamed_files(event.Payload.GetFiles().GetRenamed()),
		Short: false,
	}
}

func create_attachment_current_head(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Current HEAD",
		Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, event.Payload.HeadBranch, event.Payload.HeadBranch),
		Short: true,
	}
}

func create_attachment_conflict_head() slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Conflict HEAD",
		Value: fmt.Sprintf("<%s|%s>", "", ""),
		Short: true,
	}
}

func create_attachment_affected_files() slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Affected Files",
		Value: fmt.Sprintf("%s", ""), // nolint: gosimple
		Short: false,
	}
}
