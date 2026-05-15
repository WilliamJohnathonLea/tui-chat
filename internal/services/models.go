package services

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

// ChannelChatMessageSub represents a channel.chat.message subsription request
func ChannelChatMessageSub(userID, sessionID string) map[string]any {
	return map[string]any {
		"type": "channel.chat.message",
		"version": "1",
		"condition": map[string]any {
			"broadcaster_user_id": userID,
			"user_id": userID,
		},
		"transport": map[string]any {
			"method": "websocket",
			"session_id": sessionID,
		},
	}
}
