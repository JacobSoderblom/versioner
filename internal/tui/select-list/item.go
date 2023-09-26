package selectlist

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Item string

func (i Item) FilterValue() string { return "" }

type SingleChoiceDelegate struct{}

func (d SingleChoiceDelegate) Height() int                             { return 1 }
func (d SingleChoiceDelegate) Spacing() int                            { return 0 }
func (d SingleChoiceDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d SingleChoiceDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}
