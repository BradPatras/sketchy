package main

import (
	"log"
	"math/rand/v2"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/imacks/bitflags-go"
)

const CLOCKWISE_SPIN = 1
const ANTICLOCKWISE_SPIN = -1
const NO_SPIN = 0

const ERASE_THRESHOLD = 10

const EDGE_PADDING_H = 8
const EDGE_PADDING_v = 5

type TickMsg time.Time

type CellFlags int

const (
	w CellFlags = 1 << iota
	s
	e
	n
)

type loc struct {
	x int
	y int
}

type model struct {
	hasStarted    bool
	pixelGrid     [][]CellFlags
	leftKnob      *knob
	rightKnob     *knob
	stepTimeStamp time.Time
	wWidth        int
	wHeight       int
	shakeVal      int
}

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
		case "space":
			m.shakeVal += 8
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
			return m, m.startGame()
		}
	}

	return m, cmd
}

func createEmptyGrid(width int, height int) [][]CellFlags {
	grid := make([][]CellFlags, width)
	for x := range width {
		grid[x] = make([]CellFlags, height)
		for y := range height {
			grid[x][y] = 0
		}
	}
	return grid
}

func (m *model) startGame() tea.Cmd {
	m.pixelGrid = createEmptyGrid(m.wWidth, m.wHeight)
	m.leftKnob.loc = m.wWidth / 2
	m.rightKnob.loc = m.wHeight / 2

	return m.gameTick()
}

func (m *model) gameTick() tea.Cmd {
	return tea.Tick(75*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// since flags are just an int, I can fetch the corresponding rune by making an
// array that lists them in numerical order.
// So this array is basically: 0000, 0001, 0010, 0011, etc
var runesInFlagOrder = [16]rune{' ', '╸', '╻', '┓', '╺', '━', '┏', '┳', '╹', '┛', '┃', '┫', '┗', '┻', '┣', '╋'}

func getRuneForFlags(flags CellFlags) rune {
	if int(flags) <= 15 && int(flags) >= 0 {
		return runesInFlagOrder[flags]
	} else {
		return '!'
	}
}

func (m *model) gameStep() {
	// handle screen shake/erase
	m.shakeVal = max(0, m.shakeVal-3)
	if m.shakeVal > ERASE_THRESHOLD {
		m.pixelGrid = createEmptyGrid(m.wWidth, m.wHeight)
	} else if m.shakeVal > 0 {
		// randomly erase some cells
		for range (m.wWidth * m.wHeight) / 2 {
			m.pixelGrid[rand.IntN(len(m.pixelGrid))][rand.IntN(len(m.pixelGrid[0]))] = 0
		}
	}

	// take current stylus location
	oldLoc := loc{m.leftKnob.loc, m.rightKnob.loc}
	oldVal := m.pixelGrid[oldLoc.x][oldLoc.y]
	// step stylus state and collect the horizontal and vertical deltas
	hDelta := m.leftKnob.gameStep()
	vDelta := m.rightKnob.gameStep()
	// take new stylus position
	newLoc := loc{m.leftKnob.loc, m.rightKnob.loc}
	newVal := m.pixelGrid[newLoc.x][newLoc.y]
	// use stylus movement deltas to update flags in previous cell and in new location cell
	// first handle all the diagonal cases, those require the most work.
	if hDelta == 1 && vDelta == 1 {
		// SE
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, s)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, w)
		diagVal := m.pixelGrid[oldLoc.x][oldLoc.y+1]
		m.pixelGrid[oldLoc.x][oldLoc.y+1] = bitflags.Set(diagVal, n, e)
	} else if hDelta == -1 && vDelta == -1 {
		// NW
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, n)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, e)
		diagVal := m.pixelGrid[oldLoc.x][oldLoc.y-1]
		m.pixelGrid[oldLoc.x][oldLoc.y-1] = bitflags.Set(diagVal, s, w)
	} else if hDelta == -1 && vDelta == 1 {
		// SW
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, s)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, e)
		diagVal := m.pixelGrid[oldLoc.x][oldLoc.y+1]
		m.pixelGrid[oldLoc.x][oldLoc.y+1] = bitflags.Set(diagVal, n, w)
	} else if hDelta == 1 && vDelta == -1 {
		// NE
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, n)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, w)
		diagVal := m.pixelGrid[oldLoc.x][oldLoc.y-1]
		m.pixelGrid[oldLoc.x][oldLoc.y-1] = bitflags.Set(diagVal, s, e)
	} else if hDelta == 1 {
		// E
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, e)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, w)
	} else if hDelta == -1 {
		// W
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, w)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, e)
	} else if vDelta == 1 {
		// S
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, s)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, n)
	} else if vDelta == -1 {
		// N
		m.pixelGrid[oldLoc.x][oldLoc.y] = bitflags.Set(oldVal, n)
		m.pixelGrid[newLoc.x][newLoc.y] = bitflags.Set(newVal, s)
	}
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
	for y := range len(m.pixelGrid[0]) {
		for x := range len(m.pixelGrid) {
			sb.WriteRune(getRuneForFlags(m.pixelGrid[x][y]))
		}
		sb.WriteRune('\n')
	}

	lv := m.leftKnob.getView()
	rv := m.rightKnob.getView()
	leftKnob := lipgloss.NewLayer(lv).X(6).Y(m.wHeight - (lipgloss.Height(lv) + 3))
	rightKnob := lipgloss.NewLayer(rv).X(m.wWidth - (lipgloss.Width(rv) + 6)).Y(m.wHeight - (lipgloss.Height(lv) + 3))
	drawView := lipgloss.NewLayer(lipgloss.Place(m.wWidth, m.wHeight, lipgloss.Center, lipgloss.Center, sb.String()))
	composed := lipgloss.NewCompositor(drawView, leftKnob, rightKnob).Render()

	// screen shake
	if m.shakeVal > 0 {
		edge := rand.IntN(4)
		switch edge {
		case 0:
			composed = lipgloss.NewStyle().Padding(0, 0, 0, 1).Render(composed)
		case 1:
			composed = lipgloss.NewStyle().Padding(0, 0, 1, 0).Render(composed)
		case 2:
			composed = lipgloss.NewStyle().Padding(0, 1, 0, 0).Render(composed)
		default:
			composed = lipgloss.NewStyle().Padding(1, 0, 0, 0).Render(composed)
		}
	}

	v := tea.NewView(composed)
	v.KeyboardEnhancements.ReportEventTypes = true
	v.AltScreen = true
	return v
}
