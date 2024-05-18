package tempo_home

import (
	tea "github.com/charmbracelet/bubbletea"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	tempo_workloglist "github.com/remshams/jira-control/tui/tempo/list"
)

type Model struct {
	adapter     tui_jira.JiraAdapter
	worklogList tempo_workloglist.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter:     adapter,
		worklogList: tempo_workloglist.New(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.worklogList.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.worklogList, cmd = m.worklogList.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.worklogList.View()
}
