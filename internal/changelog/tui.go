package changelog

import (
	"versioner/internal/detect"
	"versioner/internal/tui/markdown"
	selectlist "versioner/internal/tui/select-list"
	"versioner/internal/version"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
)

type view int

const (
	bump view = iota
	summary
	confirm
	done
)

type mainModel struct {
	state   view
	bump    tea.Model
	summary tea.Model
	confirm tea.Model
	result  Changelog
	err     error
}

func startTea(project detect.Project) (Changelog, error) {
	items := []list.Item{
		selectlist.Item(version.Major),
		selectlist.Item(version.Minor),
		selectlist.Item(version.Patch),
	}

	releases := map[string]string{}
	releases[project.Name] = ""

	model := mainModel{
		state:   bump,
		bump:    selectlist.New("What kind of release is this?", selectlist.Single, items),
		summary: markdown.New("Please enter a summary for this change (this will be in the changelogs).", "Submit empty line to open external editor"),
		result: Changelog{
			Releases: releases,
		},
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return model.result, err
	}

	m, ok := result.(mainModel)
	if !ok {
		return model.result, errors.New("could not assert to main model")
	}

	if m.err != nil {
		return m.result, err
	}

	return m.result, nil
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	switch m.state {
	case bump:
		b, c := m.bump.Update(msg)

		bm, ok := b.(selectlist.SingleModel)
		if !ok {
			m.err = errors.New("could not assert Single List Model")
			return m, tea.Quit
		}

		if ok && bm.Done() {
			for k := range m.result.Releases {
				m.result.Releases[k] = bm.Choice()
				break
			}

			m.state = summary
		}

		m.bump = b
		cmd = c

	case summary:
		s, c := m.summary.Update(msg)

		bm, ok := s.(markdown.Model)
		if !ok {
			m.err = errors.New("could not assert Markdown Model")
			return m, tea.Quit
		}

		if ok && bm.Done() {

			m.result.Summary = bm.Value()

			if m.result.Summary == "" {
				summary, err := openEditor()
				m.err = err
				m.result.Summary = summary
			}

			if m.result.Summary == "" {
				m.err = errors.New("summary cannot be empty")
			}

			for _, v := range m.result.Releases {
				if v == version.Major {
					m.state = confirm
				}
			}

			if m.state != confirm {
				return m, tea.Quit
			}

		}

		m.summary = s
		cmd = c

	case confirm:
		_, cmd = m.confirm.Update(msg)
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.state {
	case bump:
		return m.bump.View()
	case summary:
		return m.summary.View()
	case confirm:
		return m.confirm.View()
	}

	return ""
}
