package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(&model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type model struct {
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "up", "down", "left", "right":
			return m, nil
		
		}
	}

	return m, cmd
}

func (m *model) handleArrowInput(input string) (tea.Model, tea.Cmd) {

	return m, nil
}

func (m *model) View() tea.View {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	return tea.NewView(style.Render("Hello, world!"))
}
