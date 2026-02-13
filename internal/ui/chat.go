package ui

import (
	"fmt"
	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type ChatModel struct {
	user      *model.User
	store     *model.Store
	input     textinput.Model
	messages  []model.Message
	errMsg    string
	logout    bool
	scrollIdx int
	Width     int
	Height    int
}

func NewChatModel(user *model.User, store *model.Store) *ChatModel {
	in := textinput.New()
	in.Placeholder = "Type message..."
	in.Focus()
	return &ChatModel{
		user:     user,
		store:    store,
		input:    in,
		messages: store.ListMessages(),
	}
}

func (m *ChatModel) Init() tea.Cmd { return nil }

func (m *ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			m.logout = true
			return m, nil
		} else if msg.String() == "enter" {
			val := m.input.Value()
			if val != "" {
				m.store.AddMessage(model.Message{
					Sender:    m.user.Username,
					Timestamp: time.Now(),
					Text:      val,
				})
				m.input.SetValue("")
				m.messages = m.store.ListMessages()
				m.scrollIdx = 0 // auto-scroll to bottom when a message is sent
			}
		} else if msg.String() == "up" {
			maxScroll := len(m.messages) - m.visibleMsgCount()
			if maxScroll < 0 {
				maxScroll = 0
			}
			if m.scrollIdx < maxScroll {
				m.scrollIdx++
			}
		} else if msg.String() == "down" {
			if m.scrollIdx > 0 {
				m.scrollIdx--
			}
		}
	}
	m.input, _ = m.input.Update(msg)
	return m, nil
}

// visibleMsgCount returns how many messages can be shown
func (m *ChatModel) visibleMsgCount() int {
	const inputRows, footerRows, errorRows, headRows = 1, 2, 1, 2
	room := m.Height - (inputRows + footerRows + headRows)
	if m.errMsg != "" {
		room -= errorRows
	}
	if room < 1 {
		room = 1
	}
	return room
}

func (m *ChatModel) View() string {
	// --- Layout parameters ---
	const headerRows = 2 // header + separator
	const inputRows = 1
	const footerRows = 2
	const errRows = 1

	minChatWidth := 60
	maxChatWidth := m.Width - 8
	if maxChatWidth < minChatWidth {
		maxChatWidth = minChatWidth
	}
	leftPad := (m.Width - maxChatWidth) / 2
	leftPad = max(leftPad, 0)

	head := Header("Chat - User: "+RenderUsername(m.user)) + "\n" + Separator() + "\n"

	bottomRows := inputRows + footerRows
	if m.errMsg != "" {
		bottomRows += errRows
	}
	availableLines := m.Height - headerRows - bottomRows
	availableLines = max(availableLines, 1)

	// Prepare visible messages
	visibleMessages := m.messages
	if len(visibleMessages) > availableLines {
		visibleMessages = visibleMessages[len(visibleMessages)-availableLines:]
	}

	msgArea := ""
	for _, msg := range visibleMessages {
		sender, found := m.store.Users[strings.ToLower(msg.Sender)]
		senderAvatar := "[?] "
		senderName := msg.Sender
		if found {
			senderAvatar = Avatar(sender)
			senderName = sender.Color.Render(sender.Username)
		}
		mine := strings.EqualFold(msg.Sender, m.user.Username)
		prefix := senderAvatar + senderName
		if mine {
			prefix = senderAvatar + senderName + " (you)"
		}
		msgArea += fmt.Sprintf("[%s] %s: %s\n", msg.Timestamp.Format("15:04:05"), prefix, msg.Text)
	}
	msgArea = strings.TrimRight(msgArea, "\n")

	msgAreaBoxed := ChatBoxStyle.Width(maxChatWidth).Height(availableLines).Render(msgArea)
	leftSpace := strings.Repeat(" ", leftPad)
	inputF := leftSpace + m.input.View() + "\n"
	err := ""
	if m.errMsg != "" {
		err = leftSpace + RenderError(m.errMsg) + "\n"
	}
	footer := leftSpace + Footer("Enter: Send   Esc/Ctrl+C: Logout   Up/Down: (reserved)")

	return head + leftSpace + msgAreaBoxed + "\n" + inputF + err + footer
}

// LoggedOut reports whether the user has requested to log out (esc/ctrl+c).
func (m *ChatModel) LoggedOut() bool {
	return m.logout
}
