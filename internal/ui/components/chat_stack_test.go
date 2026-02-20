package components

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func Test_EmptyStackView(t *testing.T) {
	stack := ChatStack{height: 5, width: 5, messages: []string{}}

	rendered := stack.View()
	expected := strings.Repeat("\n", 5)

	if !strings.EqualFold(rendered, expected) {
		t.Fail()
	}
}

func Test_OneMsg(t *testing.T) {
	stack := ChatStack{height: 5, width: 5, messages: []string{"hello"}}

	rendered := stack.View()
	expected := strings.Repeat("\n", 4) + "hello"

	if !strings.EqualFold(rendered, expected) {
		t.Fail()
	}
}

func Test_MsgsEqualHeight(t *testing.T) {
	stack := ChatStack{height: 5, width: 5, messages: []string{"hello", "world", "foo", "bar", "baz"}}

	rendered := stack.View()
	expected := strings.Join(stack.messages, "\n")

	if !strings.EqualFold(rendered, expected) {
		fmt.Printf("Rendered: %s\n===\n", rendered)
		fmt.Printf("Expected: %s\n\n", expected)
		t.Fail()
	}
}

func Test_MsgsGreaterThanHeight(t *testing.T) {
	stack := ChatStack{height: 5, width: 5, messages: []string{"hello", "world", "foo", "bar", "baz", "notseen"}}

	rendered := stack.View()
	expected := strings.Join(stack.messages[:5], "\n")

	if !strings.EqualFold(rendered, expected) {
		fmt.Printf("Rendered: %s\n===\n", rendered)
		fmt.Printf("Expected: %s\n\n", expected)
		t.Fail()
	}
}

func Test_UpdateWindowSize(t *testing.T) {
	stack := ChatStack{height: 5, width: 5, messages: []string{}}
	updateMsg := tea.WindowSizeMsg {Width: 5, Height: 10}

	updated, _ := stack.Update(updateMsg)

	rendered := updated.View()
	expected := strings.Repeat("\n", 10)

	if !strings.EqualFold(rendered, expected) {
		t.Fail()
	}
}
