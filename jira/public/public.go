package jira

import (
	"github.com/remshams/jira-control/jira/app"
	"github.com/remshams/jira-control/jira/favorite"
	"github.com/remshams/jira-control/jira/issue"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	"github.com/remshams/jira-control/jira/utils"
)

type Worklog = issue_worklog.Worklog
type WorklogAdapter = issue_worklog.WorklogAdapter
type WorklogMockAdapter = issue_worklog.WorklogMockAdatpter
type WorklogJiraAdapter = issue_worklog.WorklogJiraAdapter
type IssueAdapter = issue.IssueAdapter
type Issue = issue.Issue
type IssueSearchRequest = issue.IssueSearchRequest
type OrderBy = utils.OrderBy
type Sorting = utils.Sorting
type Favorite = favorite.Favorite
type App = app.App

const (
	SortingAsc  = utils.SortingAsc
	SortingDesc = utils.SortingDesc
)

func NewWorklogJiraAdapter() WorklogJiraAdapter {
	return issue_worklog.WorklogJiraAdapter{}
}

func NewWorklogMockAdapter() WorklogMockAdapter {
	return issue_worklog.WorklogMockAdatpter{}
}

func NewWorklog(adapter WorklogAdapter, issueKey string, hoursSpent float64) Worklog {
	return issue_worklog.NewWorklog(adapter, issueKey, hoursSpent)
}

func NewIssueSearchRequest(adapter IssueAdapter) IssueSearchRequest {
	return issue.NewIssueSearchRequest(adapter)
}

func NewOrderBy(fields []string, sorting utils.Sorting) OrderBy {
	return utils.NewOrderBy(fields, sorting)
}

func PrepareApplication() (*app.App, error) {
	return app.AppFromEnv()
}
