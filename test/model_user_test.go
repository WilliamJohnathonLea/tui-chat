package model_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
)

func TestLoadUsers_AndAuthentication(t *testing.T) {
	// Ensure relative path works regardless of workspace location
	path := filepath.Join("..", "testdata", "users.json")
	userMap, err := model.LoadUsers(path)
	assert.NoError(t, err)

	// Your starter json might be [ { "username": "alice", ... }, ...]
	alice, ok := userMap["alice"]
	assert.True(t, ok)
	assert.Equal(t, "alice", alice.Username)
	assert.NotEmpty(t, alice.Password)

	// Try correct password
	user, err := model.Authenticate(userMap, "Alice", alice.Password)
	assert.NoError(t, err)
	assert.Equal(t, alice, user)
	// Try wrong password
	_, err = model.Authenticate(userMap, "alice", "wrong")
	assert.Error(t, err)
	// Try nonexistent username
	_, err = model.Authenticate(userMap, "bob", "whatever")
	assert.Error(t, err)
}

func TestAssignColor_Deterministic(t *testing.T) {
	c1 := model.AssignColor("userA")
	c2 := model.AssignColor("userA")
	c3 := model.AssignColor("userB")
	assert.Equal(t, c1.GetForeground(), c2.GetForeground())
	assert.NotEqual(t, c1.GetForeground(), c3.GetForeground())
}

// Wrapping assignColor for test (since it's not exported in model):
// Add the following to your user.go or update test accordingly
