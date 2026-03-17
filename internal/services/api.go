package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const clientID = "8pbsu0inj1huddl1inp1800p4vtmwy"

const twitchAPIURL = "https://api.twitch.tv/helix/chat/messages"

// SendMessage sends a chat message to Twitch using the API.
// client: HTTP client (mockable)
// accessToken: OAuth token for authenticating Twitch user
// senderId: The user ID sending the message
// message: The text to send
func SendMessage(client *http.Client, accessToken, senderId, message string) error {
	// Prepare HTTP request
	payload := map[string]any{
		"broadcaster_id": senderId, // broadcaster and sender are the same user
		"sender_id":      senderId,
		"text":           message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Construct request
	req, err := http.NewRequest("POST", twitchAPIURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("twitch API returned error status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse response for data/is_sent and drop_reason
	var apiResp struct {
		Data []struct {
			MessageID  string `json:"message_id"`
			IsSent     bool   `json:"is_sent"`
			DropReason *struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			} `json:"drop_reason"`
		} `json:"data"`
	}
	err = json.Unmarshal(respBody, &apiResp)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Logic for 200 OK: check is_sent (at least one message)
	if len(apiResp.Data) == 0 {
		return fmt.Errorf("twitch API 200 OK but no message data returned")
	}

	msg := apiResp.Data[0]
	if !msg.IsSent {
		if msg.DropReason != nil {
			return fmt.Errorf("twitch API drop_reason (code=%d): %s", msg.DropReason.Code, msg.DropReason.Message)
		}
		return fmt.Errorf("twitch API: message was not sent and no drop_reason provided")
	}

	return nil
}
