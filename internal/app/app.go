package app

import (
	"fmt"
	"log"
	"os"

	"github.com/ShkolZ/tlippy/internal/config"
	"github.com/ShkolZ/tlippy/internal/oauth"

	tea "charm.land/bubbletea/v2"
)

type stage int

const (
	MainStage stage = iota
	BulkDownload
	SingleDownload
)

type model struct {
	state      stage
	choices    []string
	output     string
	clipAmount int
	cursor     int
}

func initialModel() model {

	return model{
		state:   MainStage,
		choices: []string{"Bulk Download", "Single Download"},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.state {
		case MainStage:
			switch msg.String() {
			case "up":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down":
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			case "enter":
				switch m.choices[m.cursor] {
				case "Bulk Download":
					m.state = BulkDownload
					m.cursor = 0
				case "Single Download":
					m.state = SingleDownload
					m.cursor = 0
			}

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
