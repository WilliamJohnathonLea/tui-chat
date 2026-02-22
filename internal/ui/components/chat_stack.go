package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// ChatStack is a Bubble Tea model representing a stack window
// of chat messages
type ChatStack struct {
	width  int
	height int

	// The start offset for the most recent message
	msgOffset int

	// The messages ordered most-to-least recent
	messages []string
}

func New() ChatStack {
	return ChatStack{}
}

func (cs ChatStack) Init() tea.Cmd {
	return nil
}

func (cs ChatStack) Update(msg tea.Msg) (ChatStack, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down": // Scroll forwards through chat
			if len(cs.messages) > cs.height {
				cs.msgOffset++
				cs.msgOffset = min(len(cs.messages)-cs.height, cs.msgOffset)
			}
		case "up": // Scroll backwards through chat
			cs.msgOffset--
			cs.msgOffset = max(0, cs.msgOffset)
		}
	}
	return cs, nil
}

func (cs ChatStack) View() string {
	b := strings.Builder{}
	limit := min(len(cs.messages), cs.height)
	offset := cs.msgOffset

	topFill := strings.Repeat("\n", max(0, cs.height-len(cs.messages)))

	for _, msg := range cs.messages[offset:offset+limit] {
		b.WriteString(msg)
		b.WriteString("\n")
	}

	return topFill + strings.TrimSpace(b.String())
}

func (cs *ChatStack) AddMessage(msg string) {
	if len(cs.messages) >= cs.height {
		cs.msgOffset++
	}
	cs.messages = append(cs.messages, msg)
}

func (cs *ChatStack) SetWidth(width int) {
	cs.width = width
}

func (cs *ChatStack) SetHeight(height int) {
	cs.height = height
}
