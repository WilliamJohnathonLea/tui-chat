package model_test

import (
	"testing"
	"time"

	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestFormatForDisplay(t *testing.T) {
	msg := model.Message{
		Sender:    "carol",
		Timestamp: time.Date(2024, 2, 13, 17, 22, 12, 0, time.UTC),
		Text:      "Greetings!",
	}
	colored := "Car0l" // This would be colored normally; just use a placeholder
	formatted := model.FormatForDisplay(msg, colored)
	assert.Contains(t, formatted, "17:22:12")
	assert.Contains(t, formatted, "Car0l")
	assert.Contains(t, formatted, "Greetings!")
}
