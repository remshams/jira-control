package issue_home

import (
	tea "github.com/charmbracelet/bubbletea"
	issue_search_home "github.com/remshams/jira-control/tui/issue/search/home"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type Model struct {
	adapter tui_jira.JiraAdapter
	search  issue_search_home.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		search: issue_search_home.New(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.search.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.search, cmd = m.search.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.search.View()
}
