package model

import "github.com/charmbracelet/lipgloss"

// AssignColor is a public wrapper for assignColor for use in tests and TUI.
func AssignColor(username string) lipgloss.Style {
	return assignColor(username)
}
