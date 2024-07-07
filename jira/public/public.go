package jira

import (
	"github.com/remshams/jira-control/jira/app"
	"github.com/remshams/jira-control/jira/favorite"
	"github.com/remshams/jira-control/jira/issue"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	"github.com/remshams/jira-control/jira/tempo/tempo_timesheet"
	tempo_worklog "github.com/remshams/jira-control/jira/tempo/worklog"
	"github.com/remshams/jira-control/jira/user"
	"github.com/remshams/jira-control/jira/utils"
)

type Worklog = issue_worklog.Worklog
type WorklogAdapter = issue_worklog.WorklogAdapter
type WorklogMockAdapter = issue_worklog.WorklogMockAdapter
type WorklogJiraAdapter = issue_worklog.WorklogJiraAdapter
type IssueAdapter = issue.IssueAdapter
type Issue = issue.Issue
type IssueSearchRequest = issue.IssueSearchRequest
type OrderBy = utils.OrderBy
type Sorting = utils.Sorting
type Favorite = favorite.Favorite
type TempoWorklog = tempo_worklog.Worklog
type TempoWorklogListQuery = tempo_worklog.WorklogListQuery
type User = user.User
type UserAdapter = user.UserAdapter
type Timesheet = tempo_timesheet.Timesheet
type TimesheetStatus = tempo_timesheet.TimesheetStatus
type TimesheetAdapter = tempo_timesheet.TimesheetAdapter
type App = app.App

const (
	SortingAsc  = utils.SortingAsc
	SortingDesc = utils.SortingDesc
)

func NewWorklogJiraAdapter() WorklogJiraAdapter {
	return issue_worklog.WorklogJiraAdapter{}
}

func NewWorklogMockAdapter() WorklogMockAdapter {
	return issue_worklog.WorklogMockAdapter{}
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

func NewFavorite(adapter favorite.FavoriteAdapter, issueKey string, hoursSpent float64) Favorite {
	return favorite.NewFavorite(adapter, issueKey, hoursSpent)
}

func PrepareApplication() (*app.App, error) {
	return app.AppFromEnv()
}

func NewTempoWorklogListQuery(adapter tempo_worklog.WorklogListAdapter) tempo_worklog.WorklogListQuery {
	return tempo_worklog.NewWorkloglistQuery(adapter)
}

func NewTimesheet(adapter tempo_timesheet.TimesheetAdapter, accountId string) tempo_timesheet.Timesheet {
	return tempo_timesheet.NewTimesheet(adapter, accountId)
}
