package model

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

const MaxMessages = 10000

type Store struct {
	Users    map[string]*User
	Messages map[string][]Message // room -> messages
	mu       sync.Mutex
	logFile  *os.File
}

func NewStoreWithLog(users map[string]*User, logPath string) *Store {
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open chat log file: " + err.Error())
	}
	store := &Store{
		Users:    users,
		Messages: make(map[string][]Message),
		logFile:  logFile,
	}
	store.loadMessagesFromLog(logPath)
	return store
}

// Backward compatible: uses chatlog.jsonl
func NewStore(users map[string]*User) *Store {
	return NewStoreWithLog(users, "chatlog.jsonl")
}

// loadMessagesFromLog loads messages from the log at startup for restoration
func (s *Store) loadMessagesFromLog(logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		return // not fatal if log doesn't exist
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	for {
		var msg Message
		err := dec.Decode(&msg)
		if err != nil {
			break // EOF or any parse error
		}
		room := msg.Room
		if room == "" {
			room = "main"
		}
		s.Messages[room] = append(s.Messages[room], msg)
	}
}

func (s *Store) AuthUser(username, password string) (*User, error) {
	return Authenticate(s.Users, username, password)
}

func (s *Store) AddMessage(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room := msg.Room
	if room == "" {
		room = "main"
	}
	if len(s.Messages[room]) >= MaxMessages {
		s.Messages[room] = s.Messages[room][1:]
	}
	s.Messages[room] = append(s.Messages[room], msg)
	if s.logFile != nil {
		enc := json.NewEncoder(s.logFile)
		err := enc.Encode(msg)
		if err != nil {
			println("[WARN] Could not write to chat log:", err.Error())
		}
	}
}

// ListRooms returns a sorted list of all room names
func (s *Store) ListRooms() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.Messages))
	for k := range s.Messages {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ListMessages returns messages for a given room (default "main")
func (s *Store) ListMessages(room ...string) []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := "main"
	if len(room) > 0 && room[0] != "" {
		key = room[0]
	}
	msgs := s.Messages[key]
	result := make([]Message, len(msgs))
	copy(result, msgs)
	return result
}

// AddRoom creates a new room with the given name if it does not already exist and name is not blank.
func (s *Store) AddRoom(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if name == "" {
		return fmt.Errorf("room name cannot be blank")
	}
	if _, exists := s.Messages[name]; exists {
		return fmt.Errorf("room already exists")
	}
	s.Messages[name] = []Message{}
	return nil
}

// Close closes the logFile if open. Recommended to call on app shutdown.
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.logFile != nil {
		err := s.logFile.Close()
		s.logFile = nil
		return err
	}
	return nil
}
