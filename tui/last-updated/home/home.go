package tui_last_updated

import (
	tea "github.com/charmbracelet/bubbletea"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	tui_last_updated_issue_list "github.com/remshams/jira-control/tui/last-updated/issue-list"
)

type Model struct {
	adapter   tui_jira.JiraAdapter
	issueList tui_last_updated_issue_list.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter:   adapter,
		issueList: tui_last_updated_issue_list.New(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.issueList.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.issueList, cmd = m.issueList.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.issueList.View()
}
