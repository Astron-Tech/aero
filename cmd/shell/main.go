package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	input string
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			m.input = ""
			return m, nil

		default:
			m.input += msg.String()
		}
	}

	return m, nil
}

func (m model) View() string {
	return "\n  aero › " + m.input
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		os.Exit(1)
	}
}
