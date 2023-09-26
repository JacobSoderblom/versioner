package selectlist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	Single Mode = iota
	Multiple
)

type SingleModel struct {
	list   list.Model
	choice string
	mode   Mode
	done   bool
}

func New(title string, mode Mode, items []list.Item) tea.Model {
	l := list.New(items, SingleChoiceDelegate{}, defaultWidth, listHeight)
	l.Title = title

	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return SingleModel{
		list: l,
		mode: mode,
	}
}

func (m SingleModel) Init() tea.Cmd {
	return nil
}

func (m SingleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			i, ok := m.list.SelectedItem().(Item)
			if ok {
				m.choice = string(i)
			}

			m.done = true

			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SingleModel) View() string {
	return docStyles.Render(m.list.View())
}

func (m SingleModel) Done() bool {
	return m.done
}

func (m SingleModel) Choice() string {
	return m.choice
}
