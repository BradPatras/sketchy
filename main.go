package main

import (
	"log"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const chars = "━ ┃ ┓ ┣ ┳ ┗ ┛ ┫ ╋ ┻ ┏"
const CLOCKWISE_SPIN = -1
const ANTICLOCKWISE_SPIN = 1
const NO_SPIN = 0

func main() {
	p := tea.NewProgram(&model{
		leftKnob: &knob{
			leftSpinKey:  "d",
			rightSpinKey: "f",
		},
		rightKnob: &knob{
			leftSpinKey:  "j",
			rightSpinKey: "k",
		},
	})
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type loc struct {
	x float64
	y float64
}

type cell struct {
	n bool
	e bool
	s bool
	w bool
}

// a model of the knob that controls one axis of stylus movement
type knob struct {
	leftSpinKey      string
	rightSpinKey     string
	leftSpinKeyDown  bool
	rightSpinKeyDown bool
	loc              int
	size             int
	knobDirection    int // anticlockwise == -1, clockwise == 1, none == 0
}

// update the knob state and return it's new location
func (k *knob) gameStep() int {
	if k.knobDirection == ANTICLOCKWISE_SPIN {
		k.loc = max(0, k.loc-1)
	} else if k.knobDirection == CLOCKWISE_SPIN {
		k.loc = min(k.loc+1, k.size-1)
	}

	return k.loc
}

func (k *knob) handleKeyPress(key string) {
	if key == k.leftSpinKey {
		k.leftSpinKeyDown = true
	} else if key == k.rightSpinKey {
		k.rightSpinKeyDown = true
	}

	// pressing both keys deadlocks the knob
	if k.leftSpinKeyDown && k.rightSpinKeyDown {
		k.knobDirection = NO_SPIN
	} else if k.leftSpinKeyDown {
		k.knobDirection = ANTICLOCKWISE_SPIN
	} else if k.rightSpinKeyDown {
		k.knobDirection = CLOCKWISE_SPIN
	} else {
		k.knobDirection = NO_SPIN
	}
}

func (k *knob) handleKeyRelease(key string) {
	if key != k.leftSpinKey && key != k.rightSpinKey {
		return
	}

	if key == k.leftSpinKey {
		k.leftSpinKeyDown = false
	} else if key == k.rightSpinKey {
		k.rightSpinKeyDown = false
	}

	if k.leftSpinKeyDown {
		k.knobDirection = ANTICLOCKWISE_SPIN
	} else if k.rightSpinKeyDown {
		k.knobDirection = CLOCKWISE_SPIN
	} else {
		k.knobDirection = NO_SPIN
	}
}

type model struct {
	hasStarted    bool
	cells         [][]cell
	leftKnob      *knob
	rightKnob     *knob
	stepTimeStamp time.Time
	wWidth        int
	wHeight       int
}

type TickMsg time.Time

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyReleaseMsg:
		switch msg.String() {
		case "d", "f", "j", "k":
			return m.handleDialKeyReleased(msg.String())
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "d", "f", "j", "k":
			return m.handleDialKeyPressed(msg.String())
		}
	case TickMsg:
		m.gameStep()
		return m, m.gameTick()
	case tea.WindowSizeMsg:
		m.wHeight = msg.Height
		m.wWidth = msg.Width
		m.leftKnob.size = msg.Width
		m.rightKnob.size = msg.Height
		if !m.hasStarted {
			m.hasStarted = true
			return m, m.gameTick()
		}
	}

	return m, cmd
}

func (m *model) gameTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m *model) gameStep() {
	m.leftKnob.gameStep()
	m.rightKnob.gameStep()
}

func (m *model) handleDialKeyPressed(input string) (tea.Model, tea.Cmd) {
	m.leftKnob.handleKeyPress(input)
	m.rightKnob.handleKeyPress(input)
	return m, nil
}

func (m *model) handleDialKeyReleased(input string) (tea.Model, tea.Cmd) {
	m.leftKnob.handleKeyRelease(input)
	m.rightKnob.handleKeyRelease(input)
	return m, nil
}

func (m *model) View() tea.View {
	if m.wHeight == 0 || m.wWidth == 0 {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#6cacd4"))
		return tea.NewView(style.Render("Loading..."))
	}

	return m.screenView()
}

func (m *model) screenView() tea.View {
	var sb strings.Builder
	for i := range m.wHeight {
		for j := range m.wWidth {
			y := m.leftKnob.loc
			x := m.rightKnob.loc
			if j == y && i == x {
				sb.WriteRune('•')
			} else {
				sb.WriteRune(' ')
			}
		}
		sb.WriteRune('\n')
	}
	v := tea.NewView(lipgloss.Place(m.wWidth, m.wHeight, lipgloss.Center, lipgloss.Center, sb.String()))
	v.KeyboardEnhancements.ReportEventTypes = true
	v.AltScreen = true
	return v
}
