package main

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ----- THEME -----

var (
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
)

// ----- MODEL -----

type model struct {
	input        string
	output       []string
	history      []string
	historyIndex int
	theme        string
}

func initialModel() model {
	return model{
		output:       []string{},
		history:      []string{},
		historyIndex: 0,
		theme:        "dark",
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// ----- UPDATE -----

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			command := m.input
			m.output = append(m.output, "aero > "+command)
			if command != "" {
				m.history = append(m.history, command)
				m.historyIndex = len(m.history)
			}
			m.input = ""
			return runCommand(m, command)
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		case "up":
			if m.historyIndex > 0 {
				m.historyIndex--
				m.input = m.history[m.historyIndex]
			}
		case "down":
			if m.historyIndex < len(m.history)-1 {
				m.historyIndex++
				m.input = m.history[m.historyIndex]
			} else {
				m.historyIndex = len(m.history)
				m.input = ""
			}
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

// ----- VIEW -----

func (m model) View() string {
	var b strings.Builder

	b.WriteString("\n")

	for _, line := range m.output {
		b.WriteString(outputStyle.Render("  " + line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(promptStyle.Render("  aero > "))
	b.WriteString(m.input)

	b.WriteString("\n")
	b.WriteString(statusStyle.Render("  theme: " + m.theme + "   ctrl+c to quit"))

	return b.String()
}

// ----- MAIN -----

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		os.Exit(1)
	}
}

// ----- COMMANDS -----

func runCommand(m model, input string) (tea.Model, tea.Cmd) {
	switch input {
	case "exit":
		return m, tea.Quit
	case "help":
		m.output = append(m.output,
			"Available commands:",
			"  help         - show this message",
			"  exit         - leave Aero",
			"  theme dark   - dark theme",
			"  theme light  - light theme",
		)
		return m, nil
	case "theme dark":
		m.theme = "dark"
		promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
		outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
		return m, nil
	case "theme light":
		m.theme = "light"
		promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("25")).Bold(true)
		outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
		return m, nil
	default:
		return runSystemCommand(m, input)
	}
}

func runSystemCommand(m model, input string) (tea.Model, tea.Cmd) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return m, nil
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.output = append(m.output, err.Error())
	}
	if len(output) > 0 {
		m.output = append(m.output, string(output))
	}
	return m, nil
}
