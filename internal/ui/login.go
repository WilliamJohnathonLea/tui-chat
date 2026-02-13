package ui

import (
	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

// LoginModel manages login screen state.
type LoginModel struct {
	Username textinput.Model
	Password textinput.Model
	ErrMsg   string
	Success  bool
	users    map[string]*model.User
	Width    int
	Height   int
}

// LoginSuccessMsg allows screen transitions upon success.
type LoginSuccessMsg struct {
	User *model.User
}

func NewLoginModel(users map[string]*model.User) *LoginModel {
	u := textinput.New()
	u.Placeholder = "Username"
	u.Focus()
	p := textinput.New()
	p.Placeholder = "Password"
	p.EchoMode = textinput.EchoPassword
	return &LoginModel{
		Username: u,
		Password: p,
		users:    users,
	}
}

func (m *LoginModel) Init() tea.Cmd { return nil }

func (m *LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "tab" {
			if m.Username.Focused() {
				m.Username.Blur()
				m.Password.Focus()
			} else {
				m.Password.Blur()
				m.Username.Focus()
			}
		} else if msg.String() == "enter" {
			user, err := model.Authenticate(m.users, m.Username.Value(), m.Password.Value())
			if err != nil {
				m.ErrMsg = err.Error()
			} else {
				m.ErrMsg = ""
				m.Success = true
				return m, func() tea.Msg { return LoginSuccessMsg{User: user} }
			}
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	m.Username, _ = m.Username.Update(msg)
	m.Password, _ = m.Password.Update(msg)
	return m, nil
}

func (m *LoginModel) View() string {
	title := AppTitle("TUI Chat", m.Width)
	// Horizontal centering for input fields, error, and footer
	minFieldWidth := 40
	fieldPad := (m.Width - minFieldWidth) / 2
	if fieldPad < 0 {
		fieldPad = 0
	}
	fields := strings.Repeat(" ", fieldPad) + m.Username.View() + "\n" + strings.Repeat(" ", fieldPad) + m.Password.View()
	err := ""
	if m.ErrMsg != "" {
		err = strings.Repeat(" ", fieldPad) + RenderError(m.ErrMsg) + "\n"
	}
	footer := strings.Repeat(" ", fieldPad) + Footer("Tab: Switch Field   Enter: Log In   Ctrl+C: Quit")

	content := title + Separator() + "\n" + fields + "\n" + err + footer

	// Vertical centering
	lines := 6
	if m.ErrMsg != "" {
		lines++
	}
	pad := 0
	if m.Height > lines {
		pad = (m.Height - lines) / 2
	}
	return strings.Repeat("\n", pad) + content
}
