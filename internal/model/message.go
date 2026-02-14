package model

import (
	"fmt"
	"time"
)

type Message struct {
	Room      string    `json:"room,omitempty"`
	Sender    string    `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Text      string    `json:"text"`
}

// FormatForDisplay returns a formatted chat message string.
func FormatForDisplay(msg Message, coloredUsername string) string {
	return fmt.Sprintf("[%s] %s: %s", msg.Timestamp.Format("15:04:05"), coloredUsername, msg.Text)
}
