// inspired by choose in gum https://github.com/charmbracelet/gum/blob/main/choose/choose.go
package choose

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	subduedStyle     = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#847A85", Dark: "#979797"})
	verySubduedStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#DDDADA", Dark: "#3C3C3C"})
)

type Model struct {
	height           int
	cursor           string
	items            []item
	done             bool
	index            int
	paginator        paginator.Model
	selectedPrefix   string
	unselectedPrefix string
	cursorPrefix     string
	header           string

	// styles
	cursorStyle       lipgloss.Style
	headerStyle       lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
}

type item struct {
	text     string
	selected bool
	order    int
}

func New(options []string) Model {
	m := Model{
		height:            10,
		cursor:            ">",
		cursorPrefix:      "○ ",
		selectedPrefix:    "◉ ",
		unselectedPrefix:  "○ ",
		header:            "",
		cursorStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
		headerStyle:       lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		itemStyle:         lipgloss.NewStyle(),
		selectedItemStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
	}

	items := make([]item, len(options))

	for i, opt := range options {
		items[i] = item{text: opt, selected: false}
	}

	m.items = items

	pager := paginator.New()
	pager.SetTotalPages((len(items) + m.height - 1) / m.height)
	pager.PerPage = m.height
	pager.Type = paginator.Dots
	pager.ActiveDot = subduedStyle.Render("•")
	pager.InactiveDot = verySubduedStyle.Render("•")
	pager.KeyMap = paginator.KeyMap{}
	pager.Page = 1 / m.height

	m.paginator = pager

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
		start, end := m.paginator.GetSliceBounds(len(m.items))
		switch keypress := msg.String(); keypress {
		case "down", "j", "ctrl+j", "ctrl+n":
			m.index++
			if m.index >= len(m.items) {
				m.index = 0
				m.paginator.Page = 0
			}
			if m.index >= end {
				m.paginator.NextPage()
			}
		case "up", "k", "ctrl+k", "ctrl+p":
			m.index--
			if m.index < 0 {
				m.index = len(m.items) - 1
				m.paginator.Page = m.paginator.TotalPages - 1
			}
			if m.index < start {
				m.paginator.PrevPage()
			}
		case "right", "l", "ctrl+f":
			m.index = clamp(m.index+m.height, 0, len(m.items)-1)
			m.paginator.NextPage()
		case "left", "h", "ctrl+b":
			m.index = clamp(m.index-m.height, 0, len(m.items)-1)
			m.paginator.PrevPage()
		case "G", "end":
			m.index = len(m.items) - 1
			m.paginator.Page = m.paginator.TotalPages - 1
		case "g", "home":
			m.index = 0
			m.paginator.Page = 0
		case "enter":
			m.done = true
			m.items[m.index].selected = true
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.paginator, cmd = m.paginator.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.done {
		return ""
	}

	var s strings.Builder

	start, end := m.paginator.GetSliceBounds(len(m.items))
	for i, item := range m.items[start:end] {
		if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursor))
		} else {
			s.WriteString(strings.Repeat(" ", lipgloss.Width(m.cursor)))
		}

		if item.selected {
			s.WriteString(m.selectedItemStyle.Render(m.selectedPrefix + item.text))
		} else if i == m.index%m.height {
			s.WriteString(m.cursorStyle.Render(m.cursorPrefix + item.text))
		} else {
			s.WriteString(m.itemStyle.Render(m.unselectedPrefix + item.text))
		}
		if i != m.height {
			s.WriteRune('\n')
		}
	}

	if m.paginator.TotalPages > 1 {
		s.WriteString(strings.Repeat("\n", m.height-m.paginator.ItemsOnPage(len(m.items))+1))
		s.WriteString("  " + m.paginator.View())
	}

	if m.header != "" {
		header := m.headerStyle.Render(m.header)
		return lipgloss.JoinVertical(lipgloss.Left, header, s.String())
	}

	return s.String()
}

func (m Model) Selected() string {
	selected := ""
	for _, opt := range m.items {
		if opt.selected {
			selected = opt.text
			break
		}
	}

	return selected
}

func (m Model) Done() bool {
	return m.done
}

//nolint:unparam
func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}
