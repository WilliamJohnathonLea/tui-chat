package ui

import (
	"fmt"

	"github.com/WilliamJohnathonLea/tui-chat/internal/model"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type RoomSelectedMsg struct{ Room string }

type RoomListModel struct {
	store        *model.Store
	rooms        []string
	selected     int
	creatingRoom bool
	newRoomInput textinput.Model
	roomError    string
	Width        int
	Height       int
}

func NewRoomListModel(store *model.Store, width int, height int) *RoomListModel {
	rooms := store.ListRooms()
	if len(rooms) == 0 {
		rooms = []string{"main"}
	}
	ti := textinput.New()
	ti.Width = 24
	ti.CharLimit = 32
	ti.Placeholder = "Room name"
	return &RoomListModel{
		store:        store,
		rooms:        rooms,
		selected:     0,
		creatingRoom: false,
		newRoomInput: ti,
		roomError:    "",
		Width:        width,
		Height:       height,
	}
}

func (m *RoomListModel) Init() tea.Cmd { return nil }

func (m *RoomListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width, m.Height = msg.Width, msg.Height
	case tea.KeyMsg:
		if m.creatingRoom {
			// In create mode: handle text input
			var cmd tea.Cmd
			m.newRoomInput, cmd = m.newRoomInput.Update(msg)
			switch msg.String() {
			case "esc":
				m.creatingRoom = false
				m.roomError = ""
				return m, nil
			case "enter":
				input := m.newRoomInput.Value()
				err := m.store.AddRoom(input)
				if err != nil {
					m.roomError = err.Error()
					return m, nil
				}
				// Success: reload room list, select new room, exit create mode
				m.rooms = m.store.ListRooms()
				for i, r := range m.rooms {
					if r == input {
						m.selected = i
						break
					}
				}
				m.creatingRoom = false
				m.roomError = ""
				m.newRoomInput.SetValue("")
				return m, func() tea.Msg { return RoomSelectedMsg{Room: m.rooms[m.selected]} }
			}
			return m, cmd
		} else {
			switch msg.String() {
			case "up":
				if m.selected > 0 {
					m.selected--
				}
			case "down":
				if m.selected < len(m.rooms)-1 {
					m.selected++
				}
			case "enter":
				return m, func() tea.Msg { return RoomSelectedMsg{Room: m.rooms[m.selected]} }
			case "esc":
				return m, func() tea.Msg { return RoomSelectedMsg{Room: ""} }
			case "n":
				m.creatingRoom = true
				m.roomError = ""
				ti := textinput.New()
				ti.Placeholder = "Room name"
				ti.Width = 24
				ti.CharLimit = 32
				ti.Focus()
				m.newRoomInput = ti
				return m, nil
			}
		}
	}
	return m, nil
}

func (m *RoomListModel) View() string {
	lines := []string{"Choose chat room:"}
	if m.creatingRoom {
		lines = append(lines, "")
		lines = append(lines, "Create new room:")
		lines = append(lines, m.newRoomInput.View())
		if m.roomError != "" {
			lines = append(lines, "[!] "+m.roomError)
		}
		lines = append(lines, "")
		lines = append(lines, "Enter: Create   Esc: Cancel")
	} else {
		for idx, room := range m.rooms {
			cursor := "  "
			if idx == m.selected {
				cursor = "> "
			}
			lines = append(lines, fmt.Sprintf("%s%s", cursor, room))
		}
		lines = append(lines, "") // blank line for visual separation
		lines = append(lines, "↑/↓: Move   Enter: Join   N: New Room   Esc: Cancel")
	}
	if len(lines) == 0 {
		return "[room view: NO CONTENT]"
	}

	return joinLines(lines)
}

func joinLines(lines []string) string {
	result := ""
	for _, l := range lines {
		result += l + "\n"
	}
	return result
}
