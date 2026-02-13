package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Color    lipgloss.Style
}

func LoadUsers(path string) (map[string]*User, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var users []User
	if err := json.Unmarshal(bytes, &users); err != nil {
		return nil, err
	}
	userMap := make(map[string]*User)
	for i := range users {
		users[i].Color = assignColor(users[i].Username)
		userMap[strings.ToLower(users[i].Username)] = &users[i]
	}
	return userMap, nil
}

func Authenticate(userMap map[string]*User, username, password string) (*User, error) {
	u, ok := userMap[strings.ToLower(username)]
	if !ok {
		return nil, fmt.Errorf("invalid username or password")
	}
	if u.Password != password {
		return nil, fmt.Errorf("invalid username or password")
	}
	return u, nil
}

// assignColor deterministically returns a lipgloss style for a username.
func assignColor(username string) lipgloss.Style {
	palette := []lipgloss.Color{
		lipgloss.Color("#1A73E8"), // Google Blue
		lipgloss.Color("#34A853"), // Google Green
		lipgloss.Color("#FBBC05"), // Google Yellow
		lipgloss.Color("#EA4335"), // Google Red
		lipgloss.Color("#A142F4"), // Purple Accent
		lipgloss.Color("#F95F62"), // Coral Pink
		lipgloss.Color("#01B9C1"), // Cyan
		lipgloss.Color("#FF9900"), // Accessible Orange
		lipgloss.Color("#2DD4BF"), // Aqua Green
		lipgloss.Color("#F72585"), // Magenta
	}
	idx := 0
	for _, c := range username {
		idx += int(c)
	}
	idx = idx % len(palette)
	return lipgloss.NewStyle().Foreground(palette[idx])
}
