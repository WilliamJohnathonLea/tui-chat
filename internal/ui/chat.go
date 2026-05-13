package ui

import (
	"log"
	"net/http"

	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/WilliamJohnathonLea/tui-chat/internal/services"
	"github.com/WilliamJohnathonLea/tui-chat/internal/ui/components"
)

type ChatModel struct {
	chat                components.ChatStack
	input               textinput.Model
	participants        list.Model
	participantsVisible bool
	logout              bool
	Width               int
	Height              int
	inputFocused        bool
	httpClient          *http.Client
	accessToken         string
	loggedInUser        string // The authenticated user's ID
}

type ChatMsgSent struct {
	err error
}

type UserIdReceived struct {
	userID string
}

func NewChatModel(httpClient *http.Client, accessToken string) *ChatModel {
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
		chat:                components.New(),
		input:               in,
		participants:        participants,
		participantsVisible: true,
		inputFocused:        true,
		httpClient:          httpClient,
		accessToken:         accessToken,
	}
}

func (m *ChatModel) Init() tea.Cmd {
	// On first initialize, look up the authenticated user's info
	return func() tea.Msg {
		users, err := services.GetUsers(m.httpClient, m.accessToken)
		if err != nil || len(users) == 0 {
			return UserIdReceived{userID: ""}
		}
		return UserIdReceived{userID: users[0].ID}
	}
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle UserIdReceived for sender ID
	if uid, ok := msg.(UserIdReceived); ok {
		m.loggedInUser = uid.userID
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		chatInputHeight := ChatInputStyle.GetVerticalFrameSize() +
			FooterStyle.GetVerticalFrameSize() +
			ChatBoxStyle.GetVerticalFrameSize() + 2

		if m.participantsVisible {
			m.chat.SetWidth(m.Width - m.participants.Width())
		} else {
			offset := ChatBoxStyle.GetHorizontalMargins()
			m.chat.SetWidth(m.Width - offset)
		}

		m.chat.SetHeight(m.Height - chatInputHeight)
		m.participants.SetHeight(m.Height - chatInputHeight)

		return m, nil

	case tea.KeyMsg:
		if msg.String() == "esc" {
			m.logout = true
			return m, nil
		}
		if msg.String() == "tab" {
			m.inputFocused = !m.inputFocused
			if m.inputFocused {
				m.input.Focus()
			} else {
				m.input.Blur()
			}
			return m, nil
		}
		// Only toggle participants if input is unfocused
		if !m.inputFocused && (msg.String() == "c" || msg.String() == "C") {
			m.participantsVisible = !m.participantsVisible

			if m.participantsVisible {
				m.chat.SetWidth(m.Width - m.participants.Width())
			} else {
				offset := ChatBoxStyle.GetHorizontalMargins()
				m.chat.SetWidth(m.Width - offset)
			}

			return m, nil
		}
		if m.inputFocused {
			if msg.String() == "enter" {
				val := m.input.Value()
				if val != "" {
					m.input.SetValue("")
					return m, func() tea.Msg {
						err := services.SendMessage(m.httpClient, m.accessToken, m.loggedInUser, val)
						if err != nil {
							log.Println(err)
						}
						return ChatMsgSent{err: err}
					}
				}
			}
		}
	}

	m.input, _ = m.input.Update(msg)
	m.participants, _ = m.participants.Update(msg)
	return m, nil
}

func (m *ChatModel) View() tea.View {
	chatView := ChatBoxStyle.Width(m.chat.Width()).Render(m.chat.View())
	participantsView := ""
	if m.participantsVisible {
		participantsView = m.participants.View()
	} // else leave blank
	chatAndParticipantsView := lipgloss.JoinHorizontal(lipgloss.Top, chatView, participantsView)
	roomView := lipgloss.PlaceHorizontal(m.Width, lipgloss.Left, chatAndParticipantsView) + "\n"

	widthOffset := ChatInputStyle.GetHorizontalMargins()
	inputField := ChatInputStyle.Width(m.Width - widthOffset).Render(m.input.View())
	if !m.inputFocused {
		inputField = ChatInputDisabledStyle.Width(m.Width - widthOffset).Render("Tab to chat")
	}

	footer := ""
	if m.inputFocused {
		footer = FooterStyle.Render("tab: toggle input   enter: send   esc: logout")
	} else {
		footer = FooterStyle.Render("tab: toggle input   c: toggle chatters   esc: logout")
	}
	

	inputAndFooter := lipgloss.PlaceVertical(m.Height, lipgloss.Bottom, roomView+inputField+footer)

	view := tea.NewView(inputAndFooter)
	view.AltScreen = true

	return view
}

// LoggedOut reports whether the user has requested to log out (esc).
func (m *ChatModel) LoggedOut() bool {
	return m.logout
}
