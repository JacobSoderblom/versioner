package changelog

import (
	"versioner/internal/detect"
	"versioner/internal/tui/confirm"
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
	confirmChangeset
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
		confirm: confirm.New("Is this your desired changeset? (Y/n)"),
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
		return m.result, m.err
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

			if len(m.result.Summary) == 0 {
				m.err = errors.New("summary cannot be empty")
				return m, tea.Quit
			}

			m.state = confirmChangeset
		}

		m.summary = s
		cmd = c

	case confirmChangeset:
		s, c := m.confirm.Update(msg)

		bm, ok := s.(confirm.Model)
		if !ok {
			m.err = errors.New("could not assert Markdown Model")
			return m, tea.Quit
		}

		if ok && bm.Done() {
			m.result.Confirmed = bm.Confirm()

			return m, tea.Quit
		}

		m.confirm = s
		cmd = c
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
	case confirmChangeset:
		return m.confirm.View()
	}

	return ""
}
