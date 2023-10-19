package confirm

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	textInput textinput.Model
	title     string
	done      bool
	confirm   bool
}

func New(title string) tea.Model {
	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 1
	ti.Width = 20

	return Model{
		textInput: ti,
		title:     title,
		confirm:   true,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.done = true
			return m, nil

		case tea.KeyRunes:
			switch strings.ToLower(string(msg.Runes)) {
			case "y":
				m.done = true
				return m, nil

			case "n":
				m.done = true
				m.confirm = false
				return m, nil
			}
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.title,
		m.textInput.View(),
	)
}

func (m Model) Done() bool {
	return m.done
}

func (m Model) Confirm() bool {
	return m.confirm
}
