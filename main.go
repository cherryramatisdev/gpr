package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/cherryramatisdev/gpr/pkg/gh"
)

type pr struct {
	title       string
	state       string
	description string
	org         string
	repo        string
	url         string
	number      int
	prrStatus   string
}

func (p pr) Title() string {
	pr_state_color := "#64FC00"

	var prStateStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(pr_state_color))

	var prrStatusStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FEA400"))

	return fmt.Sprintf("%s %s %s", prStateStyle.Render("[", p.state, "]"), p.title, prrStatusStyle.Render("[", p.prrStatus, "]"))
}

func (p pr) Description() string { return p.description }
func (p pr) FilterValue() string { return p.title }

type model struct {
	prs list.Model
}

func initialModel() model {
	ghPrs, _ := gh.ListAllPrs()

	items := make([]list.Item, len(ghPrs))

	for i, ghPr := range ghPrs {
		u, _ := url.Parse(ghPr.URL)
		parts := strings.Split(u.Path, "/")
		org, repo := parts[2], parts[3]

		items[i] = pr{
			title:       ghPr.Name,
			state:       ghPr.State,
			description: ghPr.Desc,
			org:         org,
			repo:        repo,
			number:      ghPr.Number,
			url:         ghPr.URL,
			prrStatus:   ghPr.Status,
		}
	}

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
		case "o":
			currentItem := m.prs.SelectedItem().(pr)

			openCmd := exec.Command("open", fmt.Sprintf("https://github.com/%s/%s/pull/%d", currentItem.org, currentItem.repo, currentItem.number))

			cmd := tea.ExecProcess(openCmd, func(err error) tea.Msg {
				return nil
			})

			return m, cmd
		case "s":
			currentItem := m.prs.SelectedItem().(pr)

			submitCmd := exec.Command("prr", "submit", fmt.Sprintf("%s/%s/%d", currentItem.org, currentItem.repo, currentItem.number))

			cmd := tea.ExecProcess(submitCmd, func(err error) tea.Msg {
				return nil
			})

			return m, cmd
		case "enter":
			currentItem := m.prs.SelectedItem().(pr)

			getCmd := exec.Command("prr", "get", fmt.Sprintf("%s/%s/%d", currentItem.org, currentItem.repo, currentItem.number))
			getCmd.Run()

			editCmd := exec.Command("prr", "edit", fmt.Sprintf("%s/%s/%d", currentItem.org, currentItem.repo, currentItem.number))

			cmd := tea.ExecProcess(editCmd, func(err error) tea.Msg {
				return nil
			})

			return m, cmd
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
		fmt.Printf("Deu merda menor %v\n", err)
		os.Exit(1)
	}
}
