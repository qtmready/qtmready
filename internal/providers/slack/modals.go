package slack

import (
	"github.com/slack-go/slack"
)

// createModal creates a modal view.
func createModal(blockset []slack.Block, callbackID string, externalID string, title string) slack.ModalViewRequest {
	return slack.ModalViewRequest{
		Type:       slack.VTModal,
		CallbackID: callbackID,
		ExternalID: externalID,
		Title:      slack.NewTextBlockObject(slack.PlainTextType, title, false, false), // Update the title text
		Blocks:     slack.Blocks{BlockSet: blockset},                                   // Wrap blocks in a slack.Blocks struct
		Close:      slack.NewTextBlockObject(slack.PlainTextType, "Close", false, false),
		Submit:     slack.NewTextBlockObject(slack.PlainTextType, "Submit", false, false),
	}
}
