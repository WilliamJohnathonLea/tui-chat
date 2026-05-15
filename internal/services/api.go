package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const clientID = "8pbsu0inj1huddl1inp1800p4vtmwy"
const twitchChatMessagesURL = "https://api.twitch.tv/helix/chat/messages"
const twitchUsersURL = "https://api.twitch.tv/helix/users"
const twitchEventSubSubscriptionsURL = "https://api.twitch.tv/helix/eventsub/subscriptions"

// SendMessage sends a chat message to Twitch using the API.
func SendMessage(client *http.Client, accessToken, senderId, message string) error {
	// Prepare HTTP request
	payload := map[string]any{
		"broadcaster_id": senderId, // broadcaster and sender are the same user
		"sender_id":      senderId,
		"message":        message,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Construct request
	req, err := http.NewRequest("POST", twitchChatMessagesURL, bytes.NewBuffer(body))
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

// GetUsers retrieves info about one or more users from the Twitch API.
func GetUsers(httpClient *http.Client, accessToken string, userIDs ...string) ([]UserInfo, error) {
	base := twitchUsersURL
	q := url.Values{}
	for _, id := range userIDs {
		q.Add("id", id)
	}

	req, err := http.NewRequest("GET", base+"?"+q.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	// Auth headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Client-Id", clientID)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case 200:
		var out struct {
			Data []UserInfo `json:"data"`
		}
		if err := json.Unmarshal(body, &out); err != nil {
			return nil, fmt.Errorf("failed to decode user response: %w", err)
		}
		return out.Data, nil
	case 400:
		return nil, fmt.Errorf("twitch API: bad request; possibly invalid ID parameter(s)")
	case 401:
		return nil, fmt.Errorf("twitch API: unauthorized; check accessToken")
	}
	return nil, fmt.Errorf("twitch API: unexpected status %d: %s", resp.StatusCode, string(body))
}

// CreateEventSub creates a new EventSub subscription via Twitch API.
func CreateEventSub(client *http.Client, accessToken, sessionID string, subscriptionMsg map[string]any) error {
	body, err := json.Marshal(subscriptionMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", twitchEventSubSubscriptionsURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusAccepted: // 202
		return nil // success
	case http.StatusBadRequest: // 400
		return fmt.Errorf("twitch API: bad request (400): %s", string(body))
	case http.StatusUnauthorized: // 401
		return fmt.Errorf("twitch API: unauthorized (401): %s", string(body))
	case http.StatusForbidden: // 403
		return fmt.Errorf("twitch API: forbidden (403): %s", string(body))
	case http.StatusConflict: // 409
		return fmt.Errorf("twitch API: conflict (409): %s", string(body))
	case http.StatusTooManyRequests: // 429
		return fmt.Errorf("twitch API: too many requests (429): %s", string(body))
	}
	return fmt.Errorf("twitch API: unexpected status %d: %s", resp.StatusCode, string(body))
}
