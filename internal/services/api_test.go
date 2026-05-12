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

// --- GetUsers tests ---

func TestGetUsers_LoggedInUser(t *testing.T) {
	mockRT := new(MockRoundTripper)
	respBody := `{"data": [{"id": "1234", "login": "demo", "display_name": "Demo", "type": "", "broadcaster_type": "", "description": "", "profile_image_url": "url", "offline_image_url": "url_offline", "created_at": "2020-10-10T00:00:00Z"}]}`
	mockRT.On("RoundTrip", mock.Anything).Return(makeResp(http.StatusOK, respBody), nil)
	client := buildMockClient(mockRT)
	users, err := GetUsers(client, "tok", "")
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "1234", users[0].ID)
	mockRT.AssertExpectations(t)
}

func TestGetUsers_ByID(t *testing.T) {
	mockRT := new(MockRoundTripper)
	respBody := `{"data": [
		{"id": "111", "login": "foo", "display_name": "FooBar", "type": "", "broadcaster_type": "", "description": "", "profile_image_url": "", "offline_image_url": "", "created_at": "2016-02-20T00:00:00Z"},
		{"id": "222", "login": "bar", "display_name": "BarBaz", "type": "", "broadcaster_type": "", "description": "", "profile_image_url": "", "offline_image_url": "", "created_at": "2017-03-25T00:00:00Z"}
	]}`
	mockRT.On("RoundTrip", mock.Anything).Return(makeResp(http.StatusOK, respBody), nil)
	client := buildMockClient(mockRT)
	users, err := GetUsers(client, "tok", "111", "222")
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "foo", users[0].Login)
	assert.Equal(t, "BarBaz", users[1].DisplayName)
	mockRT.AssertExpectations(t)
}

func TestGetUsers_BadRequest(t *testing.T) {
	mockRT := new(MockRoundTripper)
	respBody := `{"error": "bad_request"}`
	mockRT.On("RoundTrip", mock.Anything).Return(makeResp(http.StatusBadRequest, respBody), nil)
	client := buildMockClient(mockRT)
	users, err := GetUsers(client, "tok", "badid")
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "bad request")
	mockRT.AssertExpectations(t)
}

func TestGetUsers_Unauthorized(t *testing.T) {
	mockRT := new(MockRoundTripper)
	respBody := `{"error": "unauthorized"}`
	mockRT.On("RoundTrip", mock.Anything).Return(makeResp(http.StatusUnauthorized, respBody), nil)
	client := buildMockClient(mockRT)
	users, err := GetUsers(client, "bad_token")
	assert.Error(t, err)
	assert.Nil(t, users)
	assert.Contains(t, err.Error(), "unauthorized")
	mockRT.AssertExpectations(t)
}
