package user

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	client "go.breu.io/ctrlplane/cmd/cli/apiClient"
	"go.breu.io/ctrlplane/cmd/cli/styles"
	"go.breu.io/ctrlplane/cmd/cli/utils"
	"go.breu.io/ctrlplane/internal/auth"
)

type InputModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
	fields     InputFields
}

type InputFields interface {
	// GetFields returns field names of the struct from json tag
	GetFields() []string
	BindFields(inputs []textinput.Model)
	RunCmd()
}

type (
	registerOptions struct {
		auth.RegisterationRequest
	}
)

func NewCmdUserRegister() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Registers a user",
		Long:  `registers a user with quantum`,

		RunE: func(cmd *cobra.Command, args []string) error {
			f := &registerOptions{}
			if err := tea.NewProgram(initialModel(f)).Start(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func (o *registerOptions) RunCmd() {
	o.RunRegister()
}

func (o *registerOptions) GetFields() []string {
	fields := []string{"First Name", "Last Name", "Email", "Team", "Password", "Confirm Password"}
	return fields
}

func (o *registerOptions) BindFields(inputs []textinput.Model) {
	for i := 0; i < len(inputs); i++ {
		switch inputs[i].Placeholder {
		case "Email":
			o.RegisterationRequest.Email = inputs[i].Value()
		case "First Name":
			o.RegisterationRequest.FirstName = inputs[i].Value()
		case "Last Name":
			o.RegisterationRequest.LastName = inputs[i].Value()
		case "Team":
			o.RegisterationRequest.TeamName = inputs[i].Value()
		case "Password":
			o.RegisterationRequest.Password = inputs[i].Value()
		case "Confirm Password":
			o.RegisterationRequest.ConfirmPassword = inputs[i].Value()
		}
	}
}

func initialModel(f InputFields) InputModel {
	prompts := f.GetFields()
	m := InputModel{inputs: make([]textinput.Model, len(prompts)), fields: f}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.CursorStyle = styles.CursorStyle
		t.Placeholder = prompts[i]
		m.inputs[i] = t
	}

	// focus on first index
	utils.SetFocusedState(&m.inputs[0])
	return m
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				m.fields.BindFields(m.inputs)
				m.fields.RunCmd()
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					utils.SetFocusedState(&m.inputs[i])
					continue
				}
				utils.RemoveFocusedState(&m.inputs[i])
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m InputModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m InputModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &styles.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &styles.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return b.String()
}

func (o registerOptions) RunRegister() {

	o.RegisterationRequest.ConfirmPassword = o.RegisterationRequest.Password

	c := client.Client.AuthClient
	r, err := c.Register(context.Background(), o.RegisterationRequest)

	if err != nil {
		if strings.Contains(err.Error(), "No connection") {
			fmt.Print("Quantum server is not running\n")
		} else {
			fmt.Printf("failed to register user: %v", err.Error())
		}
		os.Exit(1)
	}

	pr, err := auth.ParseRegisterResponse(r)
	if err != nil {
		panic("Error: Unable to parse register response")
	}

	switch r.StatusCode {
	case 200:
		fmt.Println("User added")
		return
	case 400:
		err, ok := pr.JSON400.Errors.Get("email")
		if ok && err == "already exists" {
			fmt.Println("User already exists")
		} else {
			fmt.Println("Unable to register user")
		}
	case 201:
		fmt.Println("User Registered")
	default:
		fmt.Println("Unable to register user")
	}
}

// GetFields gets the fields for prompt from structure tags
// (problem: the names will appear in alphabetical order)
// func (o *registerOptions) GetFields() []string {
// 	t := reflect.TypeOf(o.RegisterationRequest)
// 	fields := make([]string, t.NumField())

// 	for i := 0; i < t.NumField(); i++ {
// 		if val, ok := t.Field(i).Tag.Lookup("json"); ok {
// 			fields[i] = val
// 		}
// 	}

// 	return fields
// }

// BindFields binds the user input with structure
// func (o *registerOptions) BindFields(inputs []textinput.Model) {
// 	v := reflect.ValueOf(o.RegisterationRequest)
// 	t := reflect.TypeOf(o.RegisterationRequest)
// 	for i := 0; i < len(inputs); i++ {
// 		for j := 0; j < t.NumField(); j++ {
// 			tag := t.Field(i).Tag.Get("json")
// 			if tag == inputs[i].Placeholder {
// 				val := v.Field(i)
// 				if val.CanSet() {
// 					val.SetString(inputs[i].Value())
// 				}
// 			}
// 		}
// 	}
// }
