package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const listHeight = 14

type pr struct {
	title  string
	desc   string
	org    string
	repo   string
	number int
}

func (p pr) Title() string       { return p.title }
func (p pr) Description() string { return p.desc }
func (p pr) FilterValue() string { return p.title }

type model struct {
	prs list.Model
}

func initialModel() model {
	items := []list.Item{
		pr{
			title:  "PR TITLE 1",
			org:    "cherryramatisdev",
			repo:   "gpr",
			number: 700,
		},
		pr{
			title:  "PR TITLE 2",
			org:    "cherryramatisdev",
			repo:   "gpr",
			number: 700,
		},
		pr{
			title:  "PR TITLE 3",
			org:    "cherryramatisdev",
			repo:   "gpr",
			number: 700,
		},
	}

	const defaultWidth = 20

	prs := list.New(items, list.NewDefaultDelegate(), defaultWidth, listHeight)
	prs.Title = "Title"

	return model{prs: prs}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
