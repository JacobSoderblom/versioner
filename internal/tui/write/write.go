package write

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	autoWidth   bool
	header      string
	headerStyle lipgloss.Style
	done        bool
	textarea    textarea.Model
}

func New(placeholder string) Model {
	m := Model{}

	a := textarea.New()
	a.Focus()

	a.Prompt = "┃ "
	a.Placeholder = placeholder
	a.ShowLineNumbers = false
	a.CharLimit = 0

	style := textarea.Style{
		Placeholder:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		CursorLineNumber: lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		EndOfBuffer:      lipgloss.NewStyle().Foreground(lipgloss.Color("0")),
		LineNumber:       lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
		Prompt:           lipgloss.NewStyle().Foreground(lipgloss.Color("7")),
	}

	a.BlurredStyle = style
	a.FocusedStyle = style
	a.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	a.SetWidth(50)
	a.SetHeight(5)

	m.textarea = a

	return m
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) View() string {
	if m.done {
		return ""
	}

	if m.header != "" {
		header := m.headerStyle.Render(m.header)
		return lipgloss.JoinVertical(lipgloss.Left, header, m.textarea.View())
	}

	return m.textarea.View()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if m.autoWidth {
			m.textarea.SetWidth(msg.Width)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			m.done = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m Model) Value() string {
	return m.textarea.Value()
}

func (m Model) Done() bool {
	return m.done
}
