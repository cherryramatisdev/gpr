package main

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/cherryramatisdev/gpr/pkg/gh"
)

const listHeight = 14

type pr struct {
	title       string
	state       string
	description string
	org         string
	repo        string
	number      int
}

func (p pr) Title() string       { return p.title }
func (p pr) Description() string { return p.description }
func (p pr) FilterValue() string { return p.title }

type model struct {
	prs list.Model
}

func initialModel() model {
	ghPrs, _ := gh.ListAllPrs()

	items := make([]list.Item, len(ghPrs))

	for i, ghPr := range ghPrs {
		items[i] = pr{
			title:       ghPr.Name,
			state:       ghPr.State,
			description: ghPr.Desc,
			org:         "org",
			repo:        "repo",
			number:      ghPr.Number,
		}
	}

	const defaultWidth = 20

	width, height, _ := term.GetSize(0)

	prs := list.New(items, list.NewDefaultDelegate(), width, height)
	prs.Title = "Pull requests"

	return model{prs: prs}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.prs.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			fmt.Println(m.prs.SelectedItem())
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.prs, cmd = m.prs.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.prs.View()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Deu merda menor %v", err)
		os.Exit(1)
	}
}
