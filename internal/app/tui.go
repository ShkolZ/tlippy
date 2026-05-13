package app

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/ShkolZ/tlippy/internal/config"
	"github.com/ShkolZ/tlippy/internal/download"
)

// ─────────────────────────────────────────────
//  Screen / State machine
// ─────────────────────────────────────────────

type screenID int

const (
	screenMain         screenID = iota // "Bulk Download" / "Single Download"
	screenBulkType     screenID = iota // "By Category" / "By Streamer"
	screenBulkName     screenID = iota // text input: category or streamer name
	screenTimeRange    screenID = iota // "24 Hours" / "7 Days" / "All Time"
	screenClipCount    screenID = iota // text input: how many clips
	screenSingleID     screenID = iota // text input: clip ID
	screenDownloadPath screenID = iota // text input: where to save clips
	screenDownloading  screenID = iota // live progress view
	screenDone         screenID = iota // finished – quit program
)

// menuChoice is a label+value pair for list screens.
type menuChoice struct {
	label string
	value string
}

// ─────────────────────────────────────────────
//  Messages
// ─────────────────────────────────────────────

// progressMsg carries a single download update from the goroutine.
type progressMsg download.Progress

// waitForProgress reads one value from the channel and wraps it as a Cmd.
func waitForProgress(ch <-chan download.Progress) tea.Cmd {
	return func() tea.Msg {
		p, ok := <-ch
		if !ok {
			// channel closed without a Done flag — treat as finished
			return progressMsg{Done: true}
		}
		return progressMsg(p)
	}
}

// ─────────────────────────────────────────────
//  Model
// ─────────────────────────────────────────────

type model struct {
	screen    screenID
	cursor    int
	input     textinput.Model
	collected config.UserInput

	// progress state
	progressBar progress.Model
	progressCh  <-chan download.Progress
	dlCurrent   int
	dlTotal     int
	dlLastClip  string
	dlErr       string

	// pre-built choice lists
	mainChoices      []menuChoice
	bulkTypeChoices  []menuChoice
	timeRangeChoices []menuChoice

	err string
}

func newModel() model {
	ti := textinput.New()
	ti.CharLimit = 256

	pb := progress.New(progress.WithDefaultBlend())

	return model{
		screen:      screenMain,
		cursor:      0,
		input:       ti,
		progressBar: pb,

		mainChoices: []menuChoice{
			{"Bulk Download", "bulk"},
			{"Single Download", "single"},
		},
		bulkTypeChoices: []menuChoice{
			{"By Twitch Category", "category"},
			{"By Streamer", "streamer"},
		},
		timeRangeChoices: []menuChoice{
			{"Last 24 Hours", string(config.TimeRange24h)},
			{"Last 7 Days", string(config.TimeRange7d)},
			{"All Time", string(config.TimeRangeAll)},
		},
	}
}

// Initialization
func (m model) Init() tea.Cmd {
	return nil
}

// Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	_ = cmd

	switch msg := msg.(type) {

	case progressMsg:
		if msg.Err != nil {
			m.dlErr = fmt.Sprintf("Error: %v", msg.Err)
		}
		m.dlCurrent = msg.Current
		m.dlTotal = msg.Total
		if msg.Name != "" {
			m.dlLastClip = msg.Name
		}

		if msg.Done {
			m.screen = screenDone
			return m, nil // wait for user to press Enter/Q
		}

		var pct float64
		if m.dlTotal > 0 {
			pct = float64(m.dlCurrent) / float64(m.dlTotal)
		}
		pbCmd := m.progressBar.SetPercent(pct)

		// Wait for the next update
		return m, tea.Batch(pbCmd, waitForProgress(m.progressCh))

	// ── Progress bar internal animation ────────
	case progress.FrameMsg:
		progressModel, pbCmd := m.progressBar.Update(msg)
		m.progressBar = progressModel
		return m, pbCmd

	// ── Keyboard ────────────────────────────────
	case tea.KeyPressMsg:
		key := msg.String()

		// Global quit
		if key == "ctrl+c" {
			return m, tea.Quit
		}

		// Block input while downloading
		if m.screen == screenDownloading {
			return m, nil
		}

		// On the done screen, any key quits
		if m.screen == screenDone {
			if key == "enter" || key == "q" {
				return m, tea.Quit
			}
			return m, nil
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
				if m.screen == screenDownloading {
					// Kick off downloads and start listening for progress
					ch := download.StartDownloadChan(&m.collected)
					m.progressCh = ch
					return m, waitForProgress(ch)
				}
				m.input.SetValue("")
				m.input.Focus()
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
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(banner())
	b.WriteString("\n")

	switch m.screen {
	case screenDownloading:

		if m.dlTotal > 0 {
			b.WriteString("  " + m.progressBar.View() + "\n\n")
			b.WriteString(fmt.Sprintf("  %d / %d clips\n", m.dlCurrent, m.dlTotal))
		} else {

		}

		if m.dlLastClip != "" {
			b.WriteString(fmt.Sprintf("\n  ✓ %s\n", m.dlLastClip))
		}
		if m.dlErr != "" {
			b.WriteString(fmt.Sprintf("\n  ✗ %s\n", m.dlErr))
		}
		b.WriteString("\n  [Ctrl+C] cancel\n")

	case screenDone:
		b.WriteString("  ✓ All done!\n\n")
		b.WriteString(fmt.Sprintf("  Downloaded %d clips.\n", m.dlCurrent))
		b.WriteString("\n  [Enter] exit\n")

	default:
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
	}

	return tea.NewView(b.String())
}

// ─────────────────────────────────────────────
//  Helpers
// ─────────────────────────────────────────────

func banner() string {
	return "╔══════════════════════════════╗\n" +
		"║      tlippy  –  clip tool    ║\n" +
		"╚══════════════════════════════╝\n"
}

func (m model) screenTitle() string {
	switch m.screen {
	case screenMain:
		return "Choose a download mode:"
	case screenBulkType:
		return "Download by:"
	case screenBulkName:
		if m.collected.BulkType == config.BulkByCategory {
			return "Enter Twitch category name:"
		}
		return "Enter streamer name:"
	case screenTimeRange:
		return "Select time range:"
	case screenClipCount:
		return "How many clips do you want to download?"
	case screenSingleID:
		return "Enter clip ID:"
	case screenDownloadPath:
		return "Enter download path (folder where clips will be saved):"
	default:
		return ""
	}
}

func (m model) isInputScreen() bool {
	return m.screen == screenBulkName ||
		m.screen == screenClipCount ||
		m.screen == screenSingleID ||
		m.screen == screenDownloadPath
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
			m.collected.Mode = config.ModeBulk
			m.screen = screenBulkType
		} else {
			m.collected.Mode = config.ModeSingle
			m.screen = screenSingleID
		}

	case screenBulkType:
		m.cursor = 0
		if chosen.value == "category" {
			m.collected.BulkType = config.BulkByCategory
		} else {
			m.collected.BulkType = config.BulkByStreamer
		}
		m.screen = screenBulkName

	case screenTimeRange:
		m.cursor = 0
		m.collected.TimeRange = config.TimeRange(chosen.value)
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
		m.screen = screenDownloadPath

	case screenSingleID:
		m.collected.ClipID = val
		m.screen = screenDownloadPath

	case screenDownloadPath:
		m.collected.DownloadPath = val
		m.screen = screenDownloading // → triggers download start
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
	case screenDownloadPath:
		if m.collected.Mode == config.ModeBulk {
			m.screen = screenClipCount
		} else {
			m.screen = screenSingleID
		}
	}
	return m
}

// ─────────────────────────────────────────────
//  App entry point
// ─────────────────────────────────────────────

type TUI struct {
	Input *config.UserInput
}

func NewTUI() TUI {
	return TUI{}
}

// Run launches the BubbleTea TUI, shows download progress inline, and
// returns the collected UserInput when everything is done.
// Returns nil Input if the user cancelled.
func (a *TUI) Run() (*config.UserInput, error) {
	p := tea.NewProgram(newModel())
	finalModel, err := p.Run()
	if err != nil {
		return nil, fmt.Errorf("tui error: %w", err)
	}

	m, ok := finalModel.(model)
	if !ok {
		return nil, fmt.Errorf("unexpected model type")
	}

	if m.screen != screenDone {
		return nil, nil // user cancelled
	}

	result := m.collected
	a.Input = &result
	return &result, nil
}
