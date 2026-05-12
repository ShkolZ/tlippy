package app

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// ─────────────────────────────────────────────
//  Public result types (use these after Run())
// ─────────────────────────────────────────────

type DownloadMode int

const (
	ModeBulk   DownloadMode = iota
	ModeSingle DownloadMode = iota
)

type BulkType int

const (
	BulkByCategory BulkType = iota
	BulkByStreamer  BulkType = iota
)

type TimeRange string

const (
	TimeRange24h  TimeRange = "24h"
	TimeRange7d   TimeRange = "7d"
	TimeRangeAll  TimeRange = "all"
)

// UserInput holds everything collected from the GUI.
type UserInput struct {
	Mode DownloadMode

	// Bulk fields
	BulkType  BulkType
	QueryName string // category or streamer name
	TimeRange TimeRange
	ClipCount string // number of clips as string

	// Single fields
	ClipID string
}

// ─────────────────────────────────────────────
//  Screen / State machine
// ─────────────────────────────────────────────

type screenID int

const (
	screenMain       screenID = iota // "Bulk Download" / "Single Download"
	screenBulkType   screenID = iota // "By Category" / "By Streamer"
	screenBulkName   screenID = iota // text input: category or streamer name
	screenTimeRange  screenID = iota // "24 Hours" / "7 Days" / "All Time"
	screenClipCount  screenID = iota // text input: how many clips
	screenSingleID   screenID = iota // text input: clip ID
	screenDone       screenID = iota // finished – quit program
)

// menuChoice is a label+value pair for list screens.
type menuChoice struct {
	label string
	value string
}

// ─────────────────────────────────────────────
//  Model
// ─────────────────────────────────────────────

type model struct {
	screen    screenID
	cursor    int
	input     textinput.Model
	collected UserInput

	// pre-built choice lists
	mainChoices      []menuChoice
	bulkTypeChoices  []menuChoice
	timeRangeChoices []menuChoice

	err string
}

func newModel() model {
	ti := textinput.New()
	ti.CharLimit = 256

	return model{
		screen: screenMain,
		cursor: 0,
		input:  ti,

		mainChoices: []menuChoice{
			{"Bulk Download", "bulk"},
			{"Single Download", "single"},
		},
		bulkTypeChoices: []menuChoice{
			{"By Twitch Category", "category"},
			{"By Streamer", "streamer"},
		},
		timeRangeChoices: []menuChoice{
			{"Last 24 Hours", "24h"},
			{"Last 7 Days", "7d"},
			{"All Time", "all"},
		},
	}
}

// ─────────────────────────────────────────────
//  Init
// ─────────────────────────────────────────────

func (m model) Init() tea.Cmd {
	return nil
}

// ─────────────────────────────────────────────
//  Update
// ─────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_ = cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()

		// Global quit
		if key == "ctrl+c" {
			return m, tea.Quit
		}

		// Input screens handle keys through the textinput component
		if m.isInputScreen() {
			switch key {
			case "esc":
				m = m.goBack()
				return m, nil
			case "enter":
				val := strings.TrimSpace(m.input.Value())
				if val == "" {
					m.err = "Input cannot be empty!"
					return m, nil
				}
				m.err = ""
				m = m.applyInput(val)
				return m, textinput.Blink
			}
			// Forward all other keys to textinput
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}

		// List / menu screens
		switch key {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			choices := m.currentChoices()
			if m.cursor < len(choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m = m.applyChoice()
			if m.isInputScreen() {
				m.input.SetValue("")
				m.input.Focus()
				return m, textinput.Blink
			}
			return m, nil
		case "esc", "backspace":
			m = m.goBack()
			return m, nil
		}
	default:
		// Pass non-key messages to textinput when on an input screen
		if m.isInputScreen() {
			m.input, cmd = m.input.Update(msg)
			return m, cmd
		}
	}

	return m, cmd
}

// ─────────────────────────────────────────────
//  View
// ─────────────────────────────────────────────

func (m model) View() tea.View {
	if m.screen == screenDone {
		return tea.NewView("\n  ✓ All done! Starting download...\n\n")
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(banner())
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s\n\n", m.screenTitle()))

	if m.isInputScreen() {
		b.WriteString("  ")
		b.WriteString(m.input.View())
		b.WriteString("\n")
		if m.err != "" {
			b.WriteString(fmt.Sprintf("\n  ✗ %s\n", m.err))
		}
		b.WriteString("\n  [Enter] confirm   [Esc] back   [Ctrl+C] quit\n")
	} else {
		choices := m.currentChoices()
		for i, c := range choices {
			cursor := "  "
			if m.cursor == i {
				cursor = "► "
			}
			b.WriteString(fmt.Sprintf("  %s%s\n", cursor, c.label))
		}
		b.WriteString("\n  [↑/↓] navigate   [Enter] select   [Esc] back   [Ctrl+C] quit\n")
	}

	return tea.NewView(b.String())
}

// ─────────────────────────────────────────────
//  Helpers
// ─────────────────────────────────────────────

func banner() string {
	return "  ╔══════════════════════════════╗\n" +
		"  ║      tlippy  –  clip tool    ║\n" +
		"  ╚══════════════════════════════╝\n"
}

func (m model) screenTitle() string {
	switch m.screen {
	case screenMain:
		return "Choose a download mode:"
	case screenBulkType:
		return "Download by:"
	case screenBulkName:
		if m.collected.BulkType == BulkByCategory {
			return "Enter Twitch category name:"
		}
		return "Enter streamer name:"
	case screenTimeRange:
		return "Select time range:"
	case screenClipCount:
		return "How many clips do you want to download?"
	case screenSingleID:
		return "Enter clip ID:"
	default:
		return ""
	}
}

func (m model) isInputScreen() bool {
	return m.screen == screenBulkName ||
		m.screen == screenClipCount ||
		m.screen == screenSingleID
}

func (m model) currentChoices() []menuChoice {
	switch m.screen {
	case screenMain:
		return m.mainChoices
	case screenBulkType:
		return m.bulkTypeChoices
	case screenTimeRange:
		return m.timeRangeChoices
	default:
		return nil
	}
}

// applyChoice advances the state machine based on the selected menu item.
func (m model) applyChoice() model {
	choices := m.currentChoices()
	if len(choices) == 0 {
		return m
	}
	chosen := choices[m.cursor]

	switch m.screen {
	case screenMain:
		m.cursor = 0
		if chosen.value == "bulk" {
			m.collected.Mode = ModeBulk
			m.screen = screenBulkType
		} else {
			m.collected.Mode = ModeSingle
			m.screen = screenSingleID
		}

	case screenBulkType:
		m.cursor = 0
		if chosen.value == "category" {
			m.collected.BulkType = BulkByCategory
		} else {
			m.collected.BulkType = BulkByStreamer
		}
		m.screen = screenBulkName

	case screenTimeRange:
		m.cursor = 0
		m.collected.TimeRange = TimeRange(chosen.value)
		m.screen = screenClipCount
	}

	return m
}

// applyInput stores the typed value and advances the state machine.
func (m model) applyInput(val string) model {
	switch m.screen {
	case screenBulkName:
		m.collected.QueryName = val
		m.cursor = 0
		m.screen = screenTimeRange

	case screenClipCount:
		m.collected.ClipCount = val
		m.screen = screenDone

	case screenSingleID:
		m.collected.ClipID = val
		m.screen = screenDone
	}
	return m
}

// goBack steps back one screen in the flow.
func (m model) goBack() model {
	m.err = ""
	m.cursor = 0
	switch m.screen {
	case screenBulkType:
		m.screen = screenMain
	case screenBulkName:
		m.screen = screenBulkType
	case screenTimeRange:
		m.screen = screenBulkName
	case screenClipCount:
		m.screen = screenTimeRange
	case screenSingleID:
		m.screen = screenMain
	}
	return m
}

// ─────────────────────────────────────────────
//  App entry point
// ─────────────────────────────────────────────

type App struct{}

func NewApp() App {
	return App{}
}

// Run launches the BubbleTea TUI and returns the collected UserInput.
// Call your download logic with the returned UserInput.
func (a *App) Run() (*UserInput, error) {
	p := tea.NewProgram(newModel())
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("tui error: %w", err)
	}

	m, ok := finalModel.(model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}

	// If user quit before finishing, screen won't be screenDone
	if m.screen != screenDone {
		return nil, nil // user cancelled
	}

	result := m.collected
	return &result, nil
}
