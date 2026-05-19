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
	"github.com/gorilla/websocket"
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
	wsConn              *websocket.Conn
	accessToken         string
	loggedInUser        string // The authenticated user's ID
	sessionID           string // The EventSub Session ID
}

type ChatMsgSent struct {
	err error
}

type ChatInit struct {
	userID string
	conn   *websocket.Conn
	err    error
}

type SessionIDReceived struct {
	sessionID string
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
		inputFocused:        false,
		httpClient:          httpClient,
		accessToken:         accessToken,
	}
}

func (m *ChatModel) Init() tea.Cmd {
	return func() tea.Msg {
		conn, _, err := websocket.DefaultDialer.Dial("wss://eventsub.wss.twitch.tv/ws", nil)
		if err != nil {
			return ChatInit{err: err}
		}

		users, err := services.GetUsers(m.httpClient, m.accessToken)
		if err != nil || len(users) == 0 {
			return ChatInit{err: err}
		}
		return ChatInit{
			userID: users[0].ID,
			conn:   conn,
			err:    nil,
		}
	}
}

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ChatInit:
		if msg.err != nil {
			return m, func() tea.Msg {
				log.Println(msg.err)
				return nil
			}
		}
		m.loggedInUser = msg.userID
		m.wsConn = msg.conn
		return m, func() tea.Msg {
			sessionID, err := services.HandleWelcomeEvent(m.wsConn)
			if err != nil {
				log.Println("could not get session ID for Twitch")
				return nil
			}
			return SessionIDReceived{sessionID: sessionID}
		}
	case SessionIDReceived:
		m.sessionID = msg.sessionID
		chatSubReq := services.ChannelChatMessageSub(m.loggedInUser, m.sessionID)

		return m, func() tea.Msg {
			err := services.CreateEventSub(m.httpClient, m.accessToken, m.sessionID, chatSubReq)
			if err != nil {
				log.Printf("failed to create event subscription %s", err.Error())
			}
			return nil
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		chatInputHeight := ChatInputStyle.GetVerticalFrameSize() +
			FooterStyle.GetVerticalFrameSize() +
			ChatBoxStyle.GetVerticalFrameSize() + 2

		m.toggleChatWidth()

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
			m.toggleChatWidth()

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

func (m *ChatModel) toggleChatWidth() {
	if m.participantsVisible {
		m.chat.SetWidth(m.Width - m.participants.Width())
	} else {
		offset := ChatBoxStyle.GetHorizontalMargins()
		m.chat.SetWidth(m.Width - offset)
	}
}

// LoggedOut reports whether the user has requested to log out (esc).
func (m *ChatModel) LoggedOut() bool {
	return m.logout
}
