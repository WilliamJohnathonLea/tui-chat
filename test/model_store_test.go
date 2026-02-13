package model_test

import (
	"testing"
	"time"

	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestStore_AddAndListMessages(t *testing.T) {
	users := map[string]*model.User{
		"alice": {Username: "alice"},
	}
	store := model.NewStore(users)

	msg1 := model.Message{Sender: "alice", Timestamp: time.Now(), Text: "Hello"}
	msg2 := model.Message{Sender: "alice", Timestamp: time.Now(), Text: "World"}
	store.AddMessage(msg1)
	store.AddMessage(msg2)
	msgs := store.ListMessages()
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, "Hello", msgs[0].Text)
	assert.Equal(t, "World", msgs[1].Text)
}

func TestStore_MaxMessagesFIFO(t *testing.T) {
	users := map[string]*model.User{}
	store := model.NewStore(users)
	for i := 0; i < model.MaxMessages+5; i++ {
		store.AddMessage(model.Message{Sender: "x", Timestamp: time.Now(), Text: string('A' + rune(i%26))})
	}
	msgs := store.ListMessages()
	assert.Equal(t, model.MaxMessages, len(msgs))
	// The first message should now be the 6th inserted message
}

func TestStore_AuthUser(t *testing.T) {
	users := map[string]*model.User{
		"bob": {Username: "bob", Password: "secret"},
	}
	store := model.NewStore(users)
	u, err := store.AuthUser("bob", "secret")
	assert.NoError(t, err)
	assert.Equal(t, "bob", u.Username)
	_, err = store.AuthUser("bob", "wrong")
	assert.Error(t, err)
}
