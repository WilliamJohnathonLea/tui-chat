package model_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/stretchr/testify/assert"
)

func newTempLogFile(t *testing.T) (string, func()) {
	temp, err := ioutil.TempFile("", "chatlog-test-*.jsonl")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	name := temp.Name()
	temp.Close()
	cleanup := func() { os.Remove(name) }
	return name, cleanup
}

func TestStore_AddAndLogMessages_MultiRoom(t *testing.T) {
	logPath, cleanup := newTempLogFile(t)
	defer cleanup()
	users := map[string]*model.User{"alice": {Username: "alice"}, "bob": {Username: "bob"}}
	store := model.NewStoreWithLog(users, logPath)
	msgMain := model.Message{Room: "main", Sender: "alice", Timestamp: time.Now(), Text: "Hi in main"}
	msgOther := model.Message{Room: "other", Sender: "bob", Timestamp: time.Now(), Text: "Yo in other"}
	store.AddMessage(msgMain)
	store.AddMessage(msgOther)
	store.Close()
	store2 := model.NewStoreWithLog(users, logPath)
	msgsMain := store2.ListMessages("main")
	msgsOther := store2.ListMessages("other")
	assert.Equal(t, 1, len(msgsMain))
	assert.Equal(t, 1, len(msgsOther))
	assert.Equal(t, "Hi in main", msgsMain[0].Text)
	assert.Equal(t, "Yo in other", msgsOther[0].Text)
}
