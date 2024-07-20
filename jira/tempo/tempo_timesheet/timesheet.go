package tempo_timesheet

import (
	"time"

	utils_time "github.com/remshams/common/utils/time"
	tempo_worklog "github.com/remshams/jira-control/jira/tempo/worklog"
	"github.com/remshams/jira-control/jira/user"
)

type TimesheetStatus struct {
	Status        string
	RequiredHours int
	SpentHours    int
}

func NewTimesheetStatus(status string, requiredHours int, spentHours int) TimesheetStatus {
	return TimesheetStatus{
		Status:        status,
		RequiredHours: requiredHours,
		SpentHours:    spentHours,
	}
}

type TimesheetAdapter interface {
	Reviewers(accountId string) ([]user.User, error)
	Status(accountId string, from time.Time, to time.Time) (TimesheetStatus, error)
	Submit(accountId string, reviewerAccountId string, from time.Time, to time.Time) error
}

type Timesheet struct {
	adapter            TimesheetAdapter
	worklogListAdapter tempo_worklog.WorklogListAdapter
	AccountId          string
}

func NewTimesheet(adapter TimesheetAdapter, worklogListAdapter tempo_worklog.WorklogListAdapter, accountId string) Timesheet {
	return Timesheet{
		adapter:            adapter,
		worklogListAdapter: worklogListAdapter,
		AccountId:          accountId,
	}
}

func (timesheet Timesheet) Reviewers() ([]user.User, error) {
	return timesheet.adapter.Reviewers(timesheet.AccountId)
}

func (timesheet Timesheet) Status() (TimesheetStatus, error) {
	startOfMonth, endOfMonth := utils_time.GetStartAndEndOfMonth(time.Now())
	return timesheet.adapter.Status(timesheet.AccountId, startOfMonth, endOfMonth)
}

func (timesheet Timesheet) Submit(reviewerAccountId string) error {
	startOfMonth, endOfMonth := utils_time.GetStartAndEndOfMonth(time.Now())
	return timesheet.adapter.Submit(timesheet.AccountId, reviewerAccountId, startOfMonth, endOfMonth)
}

func (timesheet Timesheet) Worklogs(query tempo_worklog.WorklogListQuery) ([]tempo_worklog.Worklog, error) {
	return query.WithSortDescending().Search(timesheet.worklogListAdapter)
}
