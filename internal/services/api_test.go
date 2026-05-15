package services

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRoundTripper for http.Transport mocking with testify
// Allows custom responses for SendMessage/GetUsers tests

type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// Helper to build a mocked HTTP client
func buildMockClient(rt http.RoundTripper) *http.Client {
	return &http.Client{Transport: rt}
}

func makeResp(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

// --- CreateEventSub tests ---

func TestCreateEventSub(t *testing.T) {
	tests := []struct {
		name         string
		respCode     int
		respBody     string
		token        string
		subscription map[string]any
		wantErr      bool
		errContains  string
	}{
		{
			name:     "202 Accepted",
			respCode: http.StatusAccepted,
			respBody: `{"data": {}}`,
			token:    "goodtoken",
			subscription: map[string]any{
				"type":    "stream.online",
				"version": "1",
				"condition": map[string]any{
					"broadcaster_user_id": "123",
				},
				"transport": map[string]any{
					"method":     "websocket",
					"session_id": "mysessionid",
				},
			},
			wantErr: false,
		},
		{ // 400 Bad Request (missing condition)
			name:     "400 Bad Request - missing condition",
			respCode: http.StatusBadRequest,
			respBody: `{"error":"Missing required field: condition"}`,
			token:    "token",
			subscription: map[string]any{
				"type":    "x",
				"version": "1",
			},
			wantErr:     true,
			errContains: "bad request",
		},
		{ // 401 Unauthorized
			name:     "401 Unauthorized - invalid token",
			respCode: http.StatusUnauthorized,
			respBody: `{"error":"Invalid access token"}`,
			token:    "badtoken",
			subscription: map[string]any{
				"type":      "x",
				"version":   "1",
				"condition": map[string]any{},
			},
			wantErr:     true,
			errContains: "unauthorized",
		},
		{ // 403 Forbidden (scope)
			name:     "403 Forbidden - missing scopes",
			respCode: http.StatusForbidden,
			respBody: `{"error":"Missing required scope"}`,
			token:    "token",
			subscription: map[string]any{
				"type":      "x",
				"version":   "1",
				"condition": map[string]any{},
			},
			wantErr:     true,
			errContains: "forbidden",
		},
		{ // 409 Conflict
			name:     "409 Conflict - already exists",
			respCode: http.StatusConflict,
			respBody: `{"error":"Subscription already exists"}`,
			token:    "token",
			subscription: map[string]any{
				"type":      "x",
				"version":   "1",
				"condition": map[string]any{},
			},
			wantErr:     true,
			errContains: "conflict",
		},
		{ // 429 Too Many Requests
			name:     "429 Too Many Requests",
			respCode: http.StatusTooManyRequests,
			respBody: `{"error":"Rate limit exceeded"}`,
			token:    "token",
			subscription: map[string]any{
				"type":      "x",
				"version":   "1",
				"condition": map[string]any{},
			},
			wantErr:     true,
			errContains: "too many requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRT := new(MockRoundTripper)
			mockRT.On("RoundTrip", mock.Anything).Return(makeResp(tt.respCode, tt.respBody), nil)
			client := buildMockClient(mockRT)

			err := CreateEventSub(client, tt.token, "unusedSession", tt.subscription)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRT.AssertExpectations(t)
		})
	}
}
