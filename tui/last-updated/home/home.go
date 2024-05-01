package tui_last_updated

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/utils"
	common_worklog "github.com/remshams/jira-control/tui/_common/worklog"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	tui_last_updated_issue_list "github.com/remshams/jira-control/tui/last-updated/issue-list"
	worklog_list "github.com/remshams/jira-control/tui/worklog/list"
)

const (
	stateLastUpdatedList utils.ViewState = "list"
	stateWorklogs        utils.ViewState = "worklogs"
)

type Model struct {
	adapter     tui_jira.JiraAdapter
	issueList   tui_last_updated_issue_list.Model
	worklogList worklog_list.Model
	state       utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter:   adapter,
		issueList: tui_last_updated_issue_list.New(adapter),
		state:     stateLastUpdatedList,
	}
}

func (m Model) Init() tea.Cmd {
	return m.issueList.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateLastUpdatedList:
		cmd = m.processLastUpdateList(msg)
	case stateWorklogs:
		cmd = m.processWorklogList(msg)
	}
	return m, cmd
}

func (m *Model) processLastUpdateList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case common_worklog.ShowWorklogsAction:
		worklogList := worklog_list.New(m.adapter, msg.Issue)
		cmd = worklogList.Init()
		m.worklogList = worklogList
		m.state = stateWorklogs
	default:
		m.issueList, cmd = m.issueList.Update(msg)
	}
	return cmd
}

func (m *Model) processWorklogList(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case worklog_list.GoBackAction:
		m.state = stateLastUpdatedList
	default:
		m.worklogList, cmd = m.worklogList.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case stateLastUpdatedList:
		return m.issueList.View()
	case stateWorklogs:
		return m.worklogList.View()
	default:
		return ""
	}
}
