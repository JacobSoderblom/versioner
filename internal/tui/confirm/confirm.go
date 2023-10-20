package confirm

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	prompt      string
	affirmative string
	negative    string
	done        bool

	confirmation bool

	defaultSelection bool

	// styles
	promptStyle     lipgloss.Style
	selectedStyle   lipgloss.Style
	unselectedStyle lipgloss.Style
}

func New(prompt string) Model {
	m := Model{
		affirmative:      "Yes",
		negative:         "No",
		confirmation:     false,
		defaultSelection: false,
		prompt:           prompt,
		selectedStyle:    lipgloss.NewStyle().Background(lipgloss.Color("212")),
		unselectedStyle:  lipgloss.NewStyle().Background(lipgloss.Color("235")),
		promptStyle:      lipgloss.NewStyle().Margin(1, 0, 0, 0),
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "N":
			m.confirmation = false
			m.done = true
			return m, nil
		case "left", "h", "ctrl+p", "tab",
			"right", "l", "ctrl+n", "shift+tab":
			if m.negative == "" {
				break
			}
			m.confirmation = !m.confirmation
		case "enter":
			m.done = true
			return m, nil
		case "y", "Y":
			m.done = true
			m.confirmation = true
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.done {
		return ""
	}

	var aff, neg string

	if m.confirmation {
		aff = m.selectedStyle.Render(m.affirmative)
		neg = m.unselectedStyle.Render(m.negative)
	} else {
		aff = m.unselectedStyle.Render(m.affirmative)
		neg = m.selectedStyle.Render(m.negative)
	}

	// If the option is intentionally empty, do not show it.
	if m.negative == "" {
		neg = ""
	}

	return lipgloss.JoinVertical(lipgloss.Center, m.promptStyle.Render(m.prompt), lipgloss.JoinHorizontal(lipgloss.Left, aff, neg))
}

func (m Model) Done() bool {
	return m.done
}

func (m Model) Confirmation() bool {
	return m.confirmation
}
