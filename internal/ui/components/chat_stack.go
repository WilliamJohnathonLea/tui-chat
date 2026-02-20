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
	case tea.WindowSizeMsg:
		cs.width = msg.Width
		cs.height = msg.Height
	}
	return cs, nil
}

func (cs ChatStack) View() string {
	b := strings.Builder{}
	limit := min(len(cs.messages), cs.height)

	topFill := strings.Repeat("\n", max(0, cs.height-len(cs.messages)))

	for _, msg := range cs.messages[:limit] {
		b.WriteString(msg)
		b.WriteString("\n")
	}

	return topFill + strings.TrimSpace(b.String())
}
