package utils

import (
	"github.com/charmbracelet/bubbles/textinput"
	"go.breu.io/ctrlplane/cmd/cli/styles"
)

func SetFocusedState(ti *textinput.Model) {
	ti.Focus()
	ti.PromptStyle = styles.FocusedStyle
	ti.TextStyle = styles.FocusedStyle
}

func RemoveFocusedState(ti *textinput.Model) {
	ti.Blur()
	ti.PromptStyle = styles.NoStyle
	ti.TextStyle = styles.NoStyle
}
