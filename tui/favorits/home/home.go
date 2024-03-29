package favorite_home

import (
	tea "github.com/charmbracelet/bubbletea"
	favorite_list "github.com/remshams/jira-control/tui/favorits/list"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type Model struct {
	adapter tui_jira.JiraAdapter
	list    favorite_list.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
		list:    favorite_list.New(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.list.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
