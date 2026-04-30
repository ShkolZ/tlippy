package app

import (
	"fmt"
	"log"
	"os"

	"github.com/ShkolZ/tlippy/internal/config"
	"github.com/ShkolZ/tlippy/internal/oauth"

	tea "charm.land/bubbletea/v2"
)

type screen struct {
	title   string
	cursor  int
	choices []string
	isInput bool
}

type model struct {
	stack []screen
}

func (m *model) current() *screen {
	return &m.stack[len(m.stack)-1]
}

func (m *model) pop() {
	if len(m.stack) > 1 {
		m.stack = m.stack[len(m.stack):]
	}

}

func (m *model) push(s screen) {
	m.stack = append(m.stack, s)
}

func initialModel() model {
	curScreen := screen{
		title:   "Pick Download Method",
		choices: []string{"Clips in Bulk", "Single Clip"},
		isInput: false,
	}
	tStack := make([]screen, 0)
	tStack = append(tStack, curScreen)
	return model{
		stack: tStack,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cur := m.current()

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "up":
			if cur.cursor > 0 {
				cur.cursor--
			}

		case "down":
			if cur.cursor < len(cur.choices)-1 {
				cur.cursor++
			}
		case "enter":
			choice := cur.choices[cur.cursor]
			switch choice {
			case "Clips in Bulk":
				m.push(screen{})
			}

		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
}

func (m model) View() tea.View {

}

type App struct {
}

func NewApp() App {
	return App{}
}

func (a *App) Run() error {

	p := tea.NewProgram()

	if len(os.Args) < 3 {
		log.Fatalln("not enough arguments")
	}

	cfg, err := config.SetConfig(os.Args[1], os.Args[2])
	if err != nil {
		return fmt.Errorf("problem setting config: %v\n", err)
	}

	token, err := oauth.GetToken()
	if err != nil {
		return fmt.Errorf("error with getting token: %v\n", err)
	}

	clips, err := GetClips(token, cfg)
	if err != nil {
		return fmt.Errorf("error with getting clips: %v\n", err)
	}

	if err := DownloadClips(token, clips, cfg); err != nil {
		return fmt.Errorf("error with downloading: %v\n", err)
	}

	return nil
}
