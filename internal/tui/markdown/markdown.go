package markdown

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type Model struct {
	textarea textarea.Model
	err      error
	title    string
	done     bool
}

func New(title, placeholder string) tea.Model {
	ti := textarea.New()
	ti.Placeholder = placeholder
	ti.Focus()

	return Model{
		textarea: ti,
		title:    title,
		err:      nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if m.textarea.Focused() {
				m.textarea.Blur()
			}
		case tea.KeyEnter:
			m.done = true
			return m, nil
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s.\n\n%s\n\n",
		m.title,
		m.textarea.View(),
	)
}

func (m Model) Value() string {
	return m.textarea.Value()
}

func (m Model) Error() error {
	return m.err
}

func (m Model) Done() bool {
	return m.done
}
