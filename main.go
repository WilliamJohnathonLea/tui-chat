package main

// WARNING: Now closes chat log on shutdown via Store.Close()

import (
	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/WilliamJohnathonLea/tui-chat/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"log"
)

type appScreen int

const (
	loginScreen appScreen = iota
	chatScreen
)

type AppModel struct {
	screen   appScreen
	user     *model.User
	store    *model.Store
	login    tea.Model
	chat     tea.Model
	Width    int
	Height   int
}

func NewApp(userStore map[string]*model.User) *AppModel {
	store := model.NewStore(userStore)
	login := ui.NewLoginModel(userStore, 0, 0) // Dimensions will be set once available
	return &AppModel{screen: loginScreen, store: store, login: login, Width: 0, Height: 0}
}

func (m *AppModel) Init() tea.Cmd {
	return m.login.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width, m.Height = sizeMsg.Width, sizeMsg.Height
		// Propagate to current submodel so their fields are always current
		switch m.screen {
		case chatScreen:
			m.chat, _ = m.chat.Update(msg)
		case loginScreen:
			m.login, _ = m.login.Update(msg)
		}
	}

	switch m.screen {
	case loginScreen:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if success, ok := msg.(ui.LoginSuccessMsg); ok {
			m.user = success.User
			m.chat = ui.NewChatModel(m.user, m.store, m.Width, m.Height)
			m.screen = chatScreen
			return m, m.chat.Init()
		}
		model_, cmd := m.login.Update(msg)
		m.login = model_
		return m, cmd

	case chatScreen:
		if m.chat == nil {
			// Should never happen, but avoid panic
			return m, nil
		}
		chatModel, cmd := m.chat.Update(msg)
		// Handle logout by checking for quit flag
		if chat, ok := chatModel.(*ui.ChatModel); ok && chat.LoggedOut() {
			m.user = nil
			m.login = ui.NewLoginModel(m.store.Users, m.Width, m.Height)
			m.chat = nil
			m.screen = loginScreen
			return m, m.login.Init()
		}
		m.chat = chatModel
		return m, cmd
	}
	return m, nil
}

func (m *AppModel) View() string {
	switch m.screen {
	case loginScreen:
		return m.login.View()
	case chatScreen:
		return m.chat.View()
	}
	return ""
}

func main() {
	users, err := model.LoadUsers("testdata/users.json")
	if err != nil {
		log.Fatal("Error loading users:", err)
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("Error setting up log file:", err)
	}
	defer f.Close()

	app := NewApp(users)
	defer app.store.Close() // Cleanly closes the log file
	if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
