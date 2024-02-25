package issue_home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/utils"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	issue_search_home "github.com/remshams/jira-control/tui/issue/search/home"
	issue_search_result "github.com/remshams/jira-control/tui/issue/search/search-result"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	worklog_list "github.com/remshams/jira-control/tui/worklog/list"
)

const (
	stateIssueSearch utils.ViewState = "issue-search"
	stateWorklogs    utils.ViewState = "worklogs"
)

type Model struct {
	adapter  tui_jira.JiraAdapter
	search   issue_search_home.Model
	worklogs worklog_list.Model
	state    utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		search:   issue_search_home.New(adapter),
		worklogs: worklog_list.New("", []issue_worklog.Worklog{}),
		state:    stateIssueSearch,
	}
}

func (m Model) Init() tea.Cmd {
	return m.search.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateIssueSearch:
		cmd = m.processIssueSearchUpdate(msg)
	case stateWorklogs:
		cmd = m.processWorklogsUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processIssueSearchUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case issue_search_result.ShowWorklogsAction:
		m.state = stateWorklogs
		m.worklogs = worklog_list.New(msg.Issue.Key, []issue_worklog.Worklog{})
		cmd = m.worklogs.Init()
	default:
		m.search, cmd = m.search.Update(msg)
	}
	return cmd
}

func (m *Model) processWorklogsUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case worklog_list.GoBackAction:
		m.state = stateIssueSearch
		cmd = m.search.Init()
	default:
		m.worklogs, cmd = m.worklogs.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case stateIssueSearch:
		return m.search.View()
	case stateWorklogs:
		return m.worklogs.View()
	default:
		return ""
	}
}
