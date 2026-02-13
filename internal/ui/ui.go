package ui

import (
	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/charmbracelet/lipgloss"
)

var (
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#F44336")).Bold(true) // Accessible red
	SeparatorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#B0B3B8"))            // Higher contrast gray
	HeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#5B3CC4")) // Indigo/Purple for professional accent
	FooterStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E4E6EB")).Background(lipgloss.Color("#262626")).Italic(true).MarginTop(1).Padding(0, 1)
	TitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFF")).Background(lipgloss.Color("#5B3CC4")).Padding(1, 8).MarginBottom(1).Align(lipgloss.Center)
	ChatBoxStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#5B3CC4")).Padding(0, 1).Margin(0, 1) // Accent border
)

func RenderUsername(u *model.User) string {
	return u.Color.Render(u.Username)
}

func RenderError(msg string) string {
	return ErrorStyle.Render(msg)
}

func Separator() string {
	return SeparatorStyle.Render("──────────────────────────────")
}

func Header(text string) string {
	return HeaderStyle.Render(text)
}

func Footer(text string) string {
	return FooterStyle.Render(text)
}

func AppTitle(text string, width int) string {
	return TitleStyle.Width(width).Render(text)
}

// Avatar renders a colored, padded initial as an avatar for the user.
func Avatar(u *model.User) string {
	initial := "?"
	if len(u.Username) > 0 {
		initial = string([]rune(u.Username)[0])
	}
	return u.Color.Bold(true).Background(lipgloss.Color("#F5F6FA")).Foreground(lipgloss.Color("#222222")).Padding(0, 1).MarginRight(1).Render(initial)
}
