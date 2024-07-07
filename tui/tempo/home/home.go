package tempo_home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/utils"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	tempo_workloglist "github.com/remshams/jira-control/tui/tempo/list"
	tempo_submit "github.com/remshams/jira-control/tui/tempo/submit"
)

const (
	stateWorklog utils.ViewState = "worklog"
	stateSubmit  utils.ViewState = "submit"
)

type Model struct {
	worklogList tempo_workloglist.Model
	submit      tempo_submit.Model
	state       utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		worklogList: tempo_workloglist.New(adapter),
		submit:      tempo_submit.New(adapter),
		state:       stateWorklog,
	}
}

func (m Model) Init() tea.Cmd {
	return m.worklogList.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateWorklog:
		cmd = m.processWorklogListUpdate(msg)
	case stateSubmit:
		cmd = m.processSubmitUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processWorklogListUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.worklogList, cmd = m.worklogList.Update(msg)
	return cmd
}

func (m *Model) processSubmitUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.submit, cmd = m.submit.Update(msg)
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case stateWorklog:
		return m.worklogList.View()
	case stateSubmit:
		return m.submit.View()
	default:
		return "Invalid view state"
	}
}
