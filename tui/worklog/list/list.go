package worklog_list

import (
	tea "github.com/charmbracelet/bubbletea"
	jira "github.com/remshams/jira-control/jira/public"
)

type Model struct {
	worklogs []jira.Worklog
}

func New(worklogs []jira.Worklog) Model {
	return Model{
		worklogs: worklogs,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return ""
}
