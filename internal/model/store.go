package model

import (
	"sync"
)

const MaxMessages = 10000

type Store struct {
	Users    map[string]*User
	Messages []Message
	mu       sync.Mutex
}

func NewStore(users map[string]*User) *Store {
	return &Store{
		Users:    users,
		Messages: []Message{},
	}
}

func (s *Store) AuthUser(username, password string) (*User, error) {
	return Authenticate(s.Users, username, password)
}

func (s *Store) AddMessage(msg Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.Messages) >= MaxMessages {
		s.Messages = s.Messages[1:]
	}
	s.Messages = append(s.Messages, msg)
}

func (s *Store) ListMessages() []Message {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]Message, len(s.Messages))
	copy(result, s.Messages)
	return result
}
