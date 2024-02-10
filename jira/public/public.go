package jira

import (
	"github.com/remshams/jira-control/jira/app"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
)

type Worklog = issue_worklog.Worklog
type WorklogAdapter = issue_worklog.WorklogAdapter
type WorklogMockAdapter = issue_worklog.WorklogMockAdatpter
type WorklogJiraAdapter = issue_worklog.WorklogJiraAdapter
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

func PrepareApplication() (*app.App, error) {
	return app.AppFromEnv()
}
