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
	promptStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Bold(true)
	outputStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	statusStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Italic(true)
	searchPromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	searchResultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
)

const nl = "\n"

// ----- MODEL -----

type model struct {
	input        string
	output       []string
	history      []string
	historyIndex int
	theme        string
	mode         string // shell | search
	searchQuery  string
	instruction  string
	results      []string
}

func initialModel() model {
	return model{
		output:       []string{},
		history:      []string{},
		historyIndex: 0,
		theme:        "dark",
		mode:         "shell",
		results:      []string{},
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
		case "esc":
			if m.mode == "search" {
				m.mode = "shell"
				m.searchQuery = ""
				m.instruction = ""
				m.results = []string{}
				return m, nil
			}
			return m, tea.Quit
		case "/":
			m.mode = "search"
			m.searchQuery = ""
			m.instruction = ""
			m.results = []string{}
			return m, nil
		case "enter":
			if m.mode == "search" {
				query, instr := parseSearch(m.searchQuery)
				m.instruction = instr
				m.output = append(m.output, searchPromptStyle.Render("search > ")+query)
				for _, line := range mockAISearch(query, instr) {
					m.output = append(m.output, line)
				}
				m.searchQuery = ""
				return m, nil
			}
			command := m.input
			m.output = append(m.output, "aero > "+command)
			if command != "" {
				m.history = append(m.history, command)
				m.historyIndex = len(m.history)
			}
			m.input = ""
			return runCommand(m, command)
		case "backspace":
			if m.mode == "search" {
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
				}
				return m, nil
			}
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if len(msg.String()) == 1 {
				if m.mode == "search" {
					m.searchQuery += msg.String()
				} else {
					m.input += msg.String()
				}
			}
		}
	}
	return m, nil
}

// ----- VIEW -----

func (m model) View() string {
	var b strings.Builder

	b.WriteString(nl)

	if m.mode == "search" {
		b.WriteString(searchPromptStyle.Render("  search > "))
		b.WriteString(m.searchQuery)
		b.WriteString(nl + nl)
		for _, r := range m.results {
			b.WriteString(searchResultStyle.Render("  " + r))
			b.WriteString(nl)
		}
		b.WriteString(nl)
		b.WriteString(statusStyle.Render("  search mode   esc to exit"))
		return b.String()
	}

	for _, line := range m.output {
		b.WriteString(outputStyle.Render("  " + line))
		b.WriteString(nl)
	}

	b.WriteString(nl)
	b.WriteString(promptStyle.Render("  aero > "))
	b.WriteString(m.input)

	b.WriteString(nl)
	b.WriteString(statusStyle.Render("  theme: " + m.theme + "   / search"))

	return b.String()
}

// ----- MAIN -----

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		os.Exit(1)
	}
}

// ----- SEARCH + AI (MOCK) -----

func parseSearch(input string) (query string, instruction string) {
	parts := strings.SplitN(input, "|", 2)
	query = strings.TrimSpace(parts[0])
	if len(parts) == 2 {
		instruction = strings.TrimSpace(parts[1])
	} else {
		instruction = "Summarize the following search results clearly and concisely."
	}
	return
}

func mockAISearch(query, instruction string) []string {
	if query == "" {
		return []string{}
	}
	return []string{
		"AI Summary for: " + query,
		"Instruction: " + instruction,
		"• This is where Pollinations output will go.",
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
			"  help   - show this message",
			"  exit   - leave Aero",
		)
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
