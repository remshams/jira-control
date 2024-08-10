package tempo_worklogdelete

import (
	tea "github.com/charmbracelet/bubbletea"
	jira "github.com/remshams/jira-control/jira/public"
)

type Model struct {
	worklog jira.TempoWorklog
}

func New() Model {
	return Model{}
}

func (m *Model) Init(worklog jira.TempoWorklog) tea.Cmd {
	m.worklog = worklog
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return "Delete worklog"
}
