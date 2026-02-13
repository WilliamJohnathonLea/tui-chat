package main

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
	screen appScreen
	user   *model.User
	store  *model.Store
	login  tea.Model
	chat   tea.Model
}

func NewApp(userStore map[string]*model.User) *AppModel {
	store := model.NewStore(userStore)
	login := ui.NewLoginModel(userStore)
	return &AppModel{screen: loginScreen, store: store, login: login}
}

func (m *AppModel) Init() tea.Cmd {
	return m.login.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.screen {
	case loginScreen:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if success, ok := msg.(ui.LoginSuccessMsg); ok {
			m.user = success.User
			m.chat = ui.NewChatModel(m.user, m.store)
			m.screen = chatScreen
			return m, m.chat.Init()
		}
		model_, cmd := m.login.Update(msg)
		m.login = model_
		return m, cmd

	case chatScreen:
		chatModel, cmd := m.chat.Update(msg)
		// Handle logout by checking for quit flag
		if chat, ok := chatModel.(*ui.ChatModel); ok && chat.LoggedOut() {
			m.user = nil
			m.login = ui.NewLoginModel(m.store.Users)
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
	app := NewApp(users)
	if _, err := tea.NewProgram(app, tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}
}
