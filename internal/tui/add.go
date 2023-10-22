package tui

import (
	"versioner/internal/changeset"
	"versioner/internal/detect"
	"versioner/internal/tui/choose"
	"versioner/internal/tui/confirm"
	"versioner/internal/tui/write"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/errors"
)

type view int

const (
	convType view = iota
	breaking
	summary
	done
)

type mainModel struct {
	state            view
	conventionalType tea.Model
	breaking         tea.Model
	summary          tea.Model
	result           changeset.Changeset
	err              error
	aborting         bool
}

func NewAddProgram(project detect.Project) (changeset.Changeset, bool, error) {
	items := make([]string, len(changeset.Types))

	for i, t := range changeset.Types {
		items[i] = t.Type
	}

	model := mainModel{
		state:            convType,
		result:           changeset.Changeset{},
		conventionalType: choose.New(items),
		breaking:         confirm.New("Are your change/changes breaking?"),
		summary:          write.New("Summary of this change"),
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return model.result, model.aborting, err
	}

	m, ok := result.(mainModel)
	if !ok {
		return model.result, model.aborting, errors.New("could not assert to main model")
	}

	if m.err != nil {
		return m.result, m.aborting, m.err
	}

	return m.result, m.aborting, nil
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
			m.aborting = true
			return m, tea.Quit
		}
	}

	switch m.state {
	case convType:
		ct, c := m.conventionalType.Update(msg)

		cm, ok := ct.(choose.Model)
		if !ok {
			m.err = errors.New("could not assert Choose Model")
			return m, tea.Quit
		}

		if ok && cm.Done() {
			m.result.Type = cm.Selected()

			if changeset.Types.CanBeBreaking(cm.Selected()) {
				m.state = breaking
			} else {
				m.state = summary
			}
		}

		m.conventionalType = ct
		cmd = c

	case breaking:
		b, c := m.breaking.Update(msg)

		bm, ok := b.(confirm.Model)
		if !ok {
			m.err = errors.New("could not assert Confirm Model")
			return m, tea.Quit
		}

		if ok && bm.Done() {
			m.result.Breaking = bm.Confirmation()
			m.state = summary
		}

		m.breaking = b
		cmd = c

	case summary:
		s, c := m.summary.Update(msg)

		bm, ok := s.(write.Model)
		if !ok {
			m.err = errors.New("could not assert Write Model")
			return m, tea.Quit
		}

		if ok && bm.Done() {

			m.result.Summary = bm.Value()

			if len(m.result.Summary) == 0 {
				m.err = errors.New("summary cannot be empty")
				return m, tea.Quit
			}

			return m, tea.Quit
		}

		m.summary = s
		cmd = c
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.state {
	case convType:
		return m.conventionalType.View()
	case breaking:
		return m.breaking.View()
	case summary:
		return m.summary.View()
	}

	return ""
}
