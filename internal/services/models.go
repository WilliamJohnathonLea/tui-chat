package services

import "time"

// UserInfo represents a user from the Twitch API.
type UserInfo struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	DisplayName     string `json:"display_name"`
	Type            string `json:"type"`
	BroadcasterType string `json:"broadcaster_type"`
	Description     string `json:"description"`
	ProfileImageURL string `json:"profile_image_url"`
	OfflineImageURL string `json:"offline_image_url"`
	Email           string `json:"email,omitempty"`
	CreatedAt       string `json:"created_at"`
}

// Metadata represents the metadata in EventSub messages
// - Welcome
// - KeepAlive
// - Reconnect
type Metadata struct {
	MessageID        string    `json:"message_id"`
	MessageType      string    `json:"message_type"`
	MessageTimestamp time.Time `json:"message_timestamp"`
}

// SubscriptionMetadata represents the metadata in EventSub messages for a subcription
// - Notification
// - Revocation
type SubscriptionMetadata struct {
	MessageID           string    `json:"message_id"`
	MessageType         string    `json:"message_type"`
	MessageTimestamp    time.Time `json:"message_timestamp"`
	SubscriptionType    string    `json:"subscription_type"`
	SubscriptionVersion string    `json:"subscription_version"`
}

// ChannelChatMessageSub represents a channel.chat.message subsription request
func ChannelChatMessageSub(userID, sessionID string) map[string]any {
	return map[string]any{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]any{
			"broadcaster_user_id": userID,
			"user_id":             userID,
		},
		"transport": map[string]any{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}
}
