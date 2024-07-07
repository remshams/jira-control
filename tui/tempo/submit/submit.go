package tempo_submit

import (
	tea "github.com/charmbracelet/bubbletea"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type Model struct {
	adapter tui_jira.JiraAdapter
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
	}
}

func (m Model) Init() tea.Cmd {
	var cmd tea.Cmd
	return cmd
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return "Submit timesheet"
}
