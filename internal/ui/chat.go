package ui

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/WilliamJohnathonLea/tui-chat/internal/ui/components"
)

type ChatModel struct {
	chat         components.ChatStack
	input        textinput.Model
	participants list.Model
	logout       bool
	Width        int
	Height       int
}

func NewChatModel() *ChatModel {
	in := textinput.New()
	in.Focus()

	items := []list.Item{
		components.Participant{Name: "GrazhProtiv"},
	}
	participants := list.New(items, list.NewDefaultDelegate(), 20, 0)
	participants.SetShowTitle(false)
	participants.SetShowHelp(false)
	participants.SetShowStatusBar(false)

	return &ChatModel{
		chat:         components.New(),
		input:        in,
		participants: participants,
	}
}

func (m *ChatModel) Init() tea.Cmd { return nil }

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		chatInputHeight := ChatInputStyle.GetVerticalFrameSize() +
			FooterStyle.GetVerticalFrameSize() +
			ChatBoxStyle.GetVerticalFrameSize() + 2

		m.chat.SetWidth(m.Width - m.participants.Width())
		m.chat.SetHeight(m.Height - chatInputHeight)
		m.participants.SetHeight(m.Height - chatInputHeight)

		return m, nil
	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.logout = true
			return m, nil
		}
		if msg.String() == "enter" {
			val := m.input.Value()
			if val != "" {
				m.input.SetValue("")
			}
		}
	}
	m.input, _ = m.input.Update(msg)
	m.participants, _ = m.participants.Update(msg)
	return m, nil
}

func (m *ChatModel) View() tea.View {
	chatView := ChatBoxStyle.Width(m.chat.Width()).Render(m.chat.View())
	participantsView := m.participants.View()
	chatAndParticipantsView := lipgloss.JoinHorizontal(lipgloss.Top, chatView, participantsView)
	roomView := lipgloss.PlaceHorizontal(m.Width, lipgloss.Left, chatAndParticipantsView) + "\n"

	widthOffset := ChatInputStyle.GetHorizontalMargins()
	input := ChatInputStyle.Width(m.Width - widthOffset).Render(m.input.View())

	footer := FooterStyle.Render("Enter: Send   Esc: Logout   Up/Down/Mouse Wheel: Scroll")

	inputAndFooter := lipgloss.PlaceVertical(m.Height, lipgloss.Bottom, roomView+input+footer)

	view := tea.NewView(inputAndFooter)
	view.AltScreen = true

	return view
}

// LoggedOut reports whether the user has requested to log out (esc).
func (m *ChatModel) LoggedOut() bool {
	return m.logout
}
