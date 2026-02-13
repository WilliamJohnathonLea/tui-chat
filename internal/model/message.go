package model

import (
	"fmt"
	"time"
)

type Message struct {
	Sender    string
	Timestamp time.Time
	Text      string
}

// FormatForDisplay returns a formatted chat message string.
func FormatForDisplay(msg Message, coloredUsername string) string {
	return fmt.Sprintf("[%s] %s: %s", msg.Timestamp.Format("15:04:05"), coloredUsername, msg.Text)
}
