package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/WilliamJohnathonLea/tui-chat/internal/ui"
	"log"
)

type appScreen int

const (
	loginScreen appScreen = iota
	chatScreen
)

type AppModel struct {
	screen appScreen
	login  tea.Model
	chat   tea.Model
	Width  int
	Height int
}

func NewApp() *AppModel {
	login := ui.NewLoginModel(0, 0) // Dimensions will be set once available
	chat := ui.NewChatModel()
	return &AppModel{
		screen: chatScreen,
		login:  login,
		chat:   chat,
		Width:  0,
		Height: 0,
	}
}

func (m *AppModel) Init() tea.Cmd {
	return m.login.Init()
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.Width, m.Height = sizeMsg.Width, sizeMsg.Height
		// Propagate to current submodel so their fields are always current
		m.chat, _ = m.chat.Update(msg)
		m.login, _ = m.login.Update(msg)
	}

	switch m.screen {
	case loginScreen:
		if key, ok := msg.(tea.KeyMsg); ok && key.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if _, ok := msg.(ui.LoginSuccessMsg); ok {
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
			m.login = ui.NewLoginModel(m.Width, m.Height)
			m.chat = nil
			m.screen = loginScreen
			return m, m.login.Init()
		}
		m.chat = chatModel
		return m, cmd
	}
	return m, nil
}

func (m *AppModel) View() tea.View {
	switch m.screen {
	case loginScreen:
		return m.login.View()
	case chatScreen:
		return m.chat.View()
	}
	return tea.NewView("")
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal("Error setting up log file:", err)
	}
	defer f.Close()

	app := NewApp()
	if _, err := tea.NewProgram(app).Run(); err != nil {
		log.Fatal(err)
	}
}
