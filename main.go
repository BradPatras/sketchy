package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const chars = "━ ┃ ┓ ┣ ┳ ┗ ┛ ┫ ╋ ┻ ┏"
const CLOCKWISE_SPIN = -1
const ANTICLOCKWISE_SPIN = 1
const NO_SPIN = 0

func main() {
	p := tea.NewProgram(&model{})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type loc struct {
	x int
	y int
}

type cell struct {
	n bool
	e bool
	s bool
	w bool
}

type model struct {
	cursor loc
	cells  [][]cell
	// user input is modeled using an int that represents which way the dial is spinning
	// -1: spinning left, 1: spinning right, 0: not spinning
	// just like an etch a sketch, the left dial controls horizontal movement
	// and the right dial controls vertical movement
	lInput int
	rInput int
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
		case "d", "f", "j", "k":
			return m.handleDialKeyPressed(msg.String())
		}
	case tea.KeyReleaseMsg:
		switch msg.String() {
		case "d", "f", "j", "k":
			return m.handleDialKeyReleased(msg.String())

		}
	}

	return m, cmd
}

func (m *model) handleDialKeyPressed(input string) (tea.Model, tea.Cmd) {
	if input == "j" {
		m.rInput = ANTICLOCKWISE_SPIN
	} else if input == "k" {
		m.rInput = CLOCKWISE_SPIN
	}
	if input == "d" {
		m.lInput = ANTICLOCKWISE_SPIN
	} else if input == "f" {
		m.lInput = CLOCKWISE_SPIN
	}
	return m, nil
}

func (m *model) handleDialKeyReleased(input string) (tea.Model, tea.Cmd) {
	if input == "j" || input == "k" {
		m.rInput = NO_SPIN
	}
	if input == "d" || input == "f" {
		m.lInput = NO_SPIN
	}
	return m, nil
}

func (m *model) View() tea.View {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
	return tea.NewView(style.Render("Hello, world!"))
}
