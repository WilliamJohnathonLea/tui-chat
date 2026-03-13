package ui

import (
	"charm.land/lipgloss/v2"
)

var (
	ErrorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Red).Bold(true)
	HeaderStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Magenta)
	FooterStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E4E6EB")).Background(lipgloss.Color("#262626")).Italic(true).MarginTop(1).Padding(0, 1)
	TitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.White).Background(lipgloss.Magenta).Padding(1, 8).MarginBottom(1).Align(lipgloss.Center)
	ChatBoxStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Magenta).Padding(0, 1).Margin(0, 1) // Accent border
	ChatInputStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Magenta).Padding(0, 1).Margin(0, 1) // Accent border
)

func RenderError(msg string) string {
	return ErrorStyle.Render(msg)
}

func Header(text string) string {
	return HeaderStyle.Render(text)
}

func AppTitle(text string, width int) string {
	return TitleStyle.Width(width).Render(text)
}
