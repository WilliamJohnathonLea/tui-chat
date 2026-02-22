package ui

import (
	"time"

	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/WilliamJohnathonLea/tui-chat/internal/ui/components"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const chatWindowHeightOffset = 9 // space taken by header, input, and footer

type ChatModel struct {
	user       *model.User
	store      *model.Store
	input      textinput.Model
	chatWindow components.ChatStack
	messages   []model.Message
	logout     bool
	room       string // room name ("main" default)
	Width      int
	Height     int
}

func NewChatModel(user *model.User, store *model.Store, width, height int) *ChatModel {
	in := textinput.New()
	in.Placeholder = "Type message..."
	in.Focus()
	room := "main"

	msgs := store.ListMessages(room)

	chatWindow := components.New()
	chatWindow.SetWidth(width)
	chatWindow.SetHeight(height - chatWindowHeightOffset) // leave space for header, input, and footer
	for _, msg := range msgs {
		colouredUsername := user.Color.Render(msg.Sender)
		chatWindow.AddMessage(model.FormatForDisplay(msg, colouredUsername))
	}

	return &ChatModel{
		user:       user,
		store:      store,
		input:      in,
		chatWindow: chatWindow,
		messages:   msgs,
		room:       room,
		Width:      width,
		Height:     height,
	}
}

func (m *ChatModel) Init() tea.Cmd { return nil }

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.chatWindow.SetWidth(m.Width)
		m.chatWindow.SetHeight(m.Height - chatWindowHeightOffset) // leave space for header, input, and footer
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			m.logout = true
			return m, nil
		}
		if msg.String() == "enter" {
			val := m.input.Value()
			if val != "" {
				msg := model.Message{
					Room:      m.room,
					Sender:    m.user.Username,
					Timestamp: time.Now(),
					Text:      val,
				}
				m.store.AddMessage(msg)
				colouredUsername := m.user.Color.Render(m.user.Username)
				m.chatWindow.AddMessage(model.FormatForDisplay(msg, colouredUsername))
				m.input.SetValue("")
				m.messages = m.store.ListMessages(m.room)
			}
		}
	}
	m.chatWindow, _ = m.chatWindow.Update(msg)
	m.input, _ = m.input.Update(msg)
	return m, nil
}

func (m *ChatModel) View() string {
	minChatWidth := 60
	maxChatWidth := m.Width - 8
	maxChatWidth = max(maxChatWidth, minChatWidth)

	headW := maxChatWidth
	if m.Width > 0 {
		headW = maxChatWidth
	}
	headerText := "  Room: " + m.room + " | User: " + RenderUsername(m.user)
	head := HeaderStyle.Width(headW).Render(headerText) + "\n"

	msgArea := m.chatWindow.View()
	msgArea = ChatBoxStyle.Width(m.Width - 4).Render(msgArea)

	inputF := ChatInputStyle.Width(m.Width-4).Render(m.input.View()) + "\n"

	footer := Footer("Enter: Send   Esc/Ctrl+C: Logout   Up/Down/Mouse Wheel: Scroll")

	return head + msgArea + "\n" + inputF + footer
}

// LoggedOut reports whether the user has requested to log out (esc/ctrl+c).
func (m *ChatModel) LoggedOut() bool {
	return m.logout
}
