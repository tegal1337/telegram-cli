package auth

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/tegal1337/telegram-cli/internal/ui/widgets"
)

// Step represents the current authentication step.
type Step int

const (
	StepPhone Step = iota
	StepCode
	StepPassword
	StepQR
	StepLoading
	StepDone
)

// Model is the authentication flow component.
type Model struct {
	theme      *theme.Theme
	authorizer *telegram.TUIAuthorizer
	step       Step
	input      widgets.TextArea
	error      string
	qrLink     string
	width      int
	height     int
	hint       string
}

// New creates a new auth model.
func New(th *theme.Theme, authorizer *telegram.TUIAuthorizer) Model {
	ta := widgets.NewTextArea()
	ta.Focused = true
	ta.Placeholder = "+1234567890"
	ta.Style = th.AuthInput

	return Model{
		theme:      th,
		authorizer: authorizer,
		step:       StepPhone,
		input:      ta,
		hint:       "Enter your phone number with country code",
	}
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = 40
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case telegram.AuthStateMsg:
		return m.handleAuthState(msg)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			return m.submit()
		case "escape":
			m.error = ""
		default:
			m.input.Update(msg)
		}
	}

	return m, nil
}

func (m Model) handleAuthState(msg telegram.AuthStateMsg) (Model, tea.Cmd) {
	switch msg.State.(type) {
	case *telegram.AuthStateMsg:
		// Handled by the authorizer directly.
	}

	// The TUIAuthorizer handles state transitions internally.
	// We update the UI step based on what the authorizer needs.
	return m, nil
}

func (m Model) submit() (Model, tea.Cmd) {
	value := m.input.Value
	if value == "" {
		m.error = "This field is required"
		return m, nil
	}

	m.error = ""

	switch m.step {
	case StepPhone:
		m.authorizer.SubmitPhone(value)
		m.step = StepCode
		m.input.Reset()
		m.input.Placeholder = "12345"
		m.hint = "Enter the verification code sent to your phone"

	case StepCode:
		m.authorizer.SubmitCode(value)
		m.step = StepLoading
		m.input.Reset()
		m.hint = "Verifying..."

	case StepPassword:
		m.authorizer.SubmitPassword(value)
		m.step = StepLoading
		m.input.Reset()
		m.hint = "Verifying..."
	}

	return m, nil
}

// SetStep changes the auth step (called when TDLib auth state changes).
func (m *Model) SetStep(step Step) {
	m.step = step
	m.input.Reset()

	switch step {
	case StepPhone:
		m.input.Placeholder = "+1234567890"
		m.hint = "Enter your phone number with country code"
	case StepCode:
		m.input.Placeholder = "12345"
		m.hint = "Enter the verification code"
	case StepPassword:
		m.input.Placeholder = "password"
		m.hint = "Enter your two-step verification password"
	case StepQR:
		m.hint = "Scan the QR code with your phone"
	}
}

// View renders the auth screen.
func (m Model) View() string {
	logo := `
 _____ ___ _    ___   _____ _   _ ___
|_   _| __| |  | __| |_   _| | | |_ _|
  | | | _|| |__| _|    | | | |_| || |
  |_| |___|____|___|   |_|  \___/|___|
`

	var content string

	switch m.step {
	case StepQR:
		qr := widgets.RenderQRCode(m.qrLink, 256)
		content = fmt.Sprintf("%s\n\n%s\n\n%s",
			m.theme.AuthTitle.Render("QR Code Login"),
			qr,
			m.theme.AuthLabel.Render(m.hint),
		)
	case StepLoading:
		sp := widgets.NewSpinner("Authenticating...")
		sp.Style = m.theme.Spinner
		content = sp.View()
	case StepDone:
		content = m.theme.AuthTitle.Render("Authenticated! Loading chats...")
	default:
		stepLabels := []string{"Phone", "Code", "Password"}
		stepIdx := int(m.step)
		if stepIdx >= len(stepLabels) {
			stepIdx = 0
		}

		title := m.theme.AuthTitle.Render(fmt.Sprintf("Step %d: %s", stepIdx+1, stepLabels[stepIdx]))
		label := m.theme.AuthLabel.Render(m.hint)
		input := m.input.View()

		errMsg := ""
		if m.error != "" {
			errMsg = "\n" + lipgloss.NewStyle().Foreground(m.theme.Error).Render(m.error)
		}

		content = fmt.Sprintf("%s\n\n%s\n%s%s", title, label, input, errMsg)
	}

	fullContent := fmt.Sprintf("%s\n\n%s",
		m.theme.AuthTitle.Render(logo),
		content,
	)

	return m.theme.AuthPane.
		Width(m.width).
		Height(m.height).
		Render(fullContent)
}
