package services

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func HandleWelcomeEvent(conn *websocket.Conn) (string, error) {
	var welcome struct {
		Metadata Metadata `json:"metadata"`
		Payload  struct {
			Session struct {
				ID string `json:"id"`
			} `json:"session"`
		} `json:"payload"`
	}

	err := conn.ReadJSON(&welcome)
	if err != nil {
		return "", fmt.Errorf("failed to decode JSON from Twitch EventSub")
	}

	if welcome.Metadata.MessageType != "session_welcome" {
		return "", fmt.Errorf("expected 'session_welcome' got '%s'", welcome.Metadata.MessageType)
	}

	return welcome.Payload.Session.ID, nil
}

func HandleEvent(conn *websocket.Conn) ([]byte, error) {
	_, msg, err := conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}

	return msg, nil
}
