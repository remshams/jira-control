package jira

import (
	"github.com/remshams/jira-control/jira/app"
	"github.com/remshams/jira-control/jira/issue"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
)

type Worklog = issue_worklog.Worklog
type WorklogAdapter = issue_worklog.WorklogAdapter
type WorklogMockAdapter = issue_worklog.WorklogMockAdatpter
type WorklogJiraAdapter = issue_worklog.WorklogJiraAdapter
type IssueAdapter = issue.IssueAdapter
type Issue = issue.Issue
type IssueSearchRequest = issue.IssueSearchRequest
type App = app.App

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

func PrepareApplication() (*app.App, error) {
	return app.AppFromEnv()
}
