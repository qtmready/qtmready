package attach

import (
	"fmt"
	"strings"

	"github.com/slack-go/slack"

	"go.breu.io/quantm/internal/events"
	eventsv1 "go.breu.io/quantm/internal/proto/ctrlplane/events/v1"
)

// Repo creates an attachment field for the repository.
func Repo[E events.Payload](event *events.Event[eventsv1.ChatHook, E]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Repository*",
		Value: fmt.Sprintf("<%s|%s>", event.Context.Source, extract_repo(event.Context.Source)),
		Short: true,
	}
}

// Branch creates an attachment field for the branch.
func Branch(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, "", ""), // TODO - verify
		Short: true,
	}
}

// BranchMerge creates an attachment field for the branch in merge context.
func BranchMerge(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Branch*",
		Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, event.Payload.BaseBranch, event.Payload.BaseBranch),
		Short: true,
	}
}

// Threshold creates an attachment field for the threshold.
func Threshold() slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Threshold*",
		Value: fmt.Sprintf("%d", 0),
		Short: true,
	}
}

// TotalLinesCount creates an attachment field for the total lines count.
func TotalLinesCount(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Total Lines Count*",
		Value: fmt.Sprintf("%d", event.Payload.GetLines().GetAdded()+event.Payload.GetLines().GetRemoved()),
		Short: true,
	}
}

// LinesAdded creates an attachment field for lines added.
func LinesAdded(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Lines Added*",
		Value: fmt.Sprintf("%d", event.Payload.GetLines().GetAdded()),
		Short: true,
	}
}

// LinesDeleted creates an attachment field for lines deleted.
func LinesDeleted(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "*Lines Deleted*",
		Value: fmt.Sprintf("%d", event.Payload.GetLines().GetRemoved()),
		Short: true,
	}
}

// AddedFiles creates an attachment field for added files.
func AddedFiles(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Added Files",
		Value: format_files(event.Payload.GetFiles().GetAdded()),
		Short: false,
	}
}

// DeletedFiles creates an attachment field for deleted files.
func DeletedFiles(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Deleted Files",
		Value: format_files(event.Payload.GetFiles().GetDeleted()),
		Short: false,
	}
}

// ModifiedFiles creates an attachment field for modified files.
func ModifiedFiles(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Modified Files",
		Value: format_files(event.Payload.GetFiles().GetModified()),
		Short: false,
	}
}

// RenameFiles creates an attachment field for renamed files.
func RenameFiles(event *events.Event[eventsv1.ChatHook, eventsv1.Diff]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Rename Files",
		Value: fomrat_rename(event.Payload.GetFiles().GetRenamed()),
		Short: false,
	}
}

// CurrentHead creates an attachment field for current head in merge context.
func CurrentHead(event *events.Event[eventsv1.ChatHook, eventsv1.Merge]) slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Current HEAD",
		Value: fmt.Sprintf("<%s/tree/%s|%s>", event.Context.Source, event.Payload.HeadBranch, event.Payload.HeadBranch),
		Short: true,
	}
}

// ConflictHead creates an attachment field for conflict head in merge context.
func ConflictHead() slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Conflict HEAD",
		Value: fmt.Sprintf("<%s|%s>", "", ""), // TODO - verify
		Short: true,
	}
}

// AffectedFiles creates an attachment field for affected files in merge context.
func AffectedFiles() slack.AttachmentField {
	return slack.AttachmentField{
		Title: "Affected Files",
		Value: fmt.Sprintf("%s", ""), // nolint: gosimple
		Short: false,
	}
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

func fomrat_rename(files []*eventsv1.RenamedFile) string {
	result := ""
	for _, file := range files {
		result += fmt.Sprintf("- %s -> %s\n", file.GetOld(), file.GetNew())
	}

	return result
}
