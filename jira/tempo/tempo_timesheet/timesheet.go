package tempo_timesheet

import (
	"time"

	utils_time "github.com/remshams/common/utils/time"
	"github.com/remshams/jira-control/jira/user"
)

type TimesheetStatus struct {
	status        string
	requiredHours int
	spentHours    int
}

func NewTimesheetStatus(status string, requiredHours int, spentHours int) TimesheetStatus {
	return TimesheetStatus{
		status:        status,
		requiredHours: requiredHours,
		spentHours:    spentHours,
	}
}

type TimesheetAdapter interface {
	Reviewers(accountId string) ([]user.User, error)
	Status(accountId string, from time.Time, to time.Time) (TimesheetStatus, error)
}

type Timesheet struct {
	adapter   TimesheetAdapter
	AccountId string
}

func NewTimesheet(adapter TimesheetAdapter, accountId string) Timesheet {
	return Timesheet{
		adapter:   adapter,
		AccountId: accountId,
	}
}

func (timesheet Timesheet) Reviewers() ([]user.User, error) {
	return timesheet.adapter.Reviewers(timesheet.AccountId)
}

func (timesheet Timesheet) Status() (TimesheetStatus, error) {
	startOfMonth, endOfMonth := utils_time.GetStartAndEndOfMonth(time.Now())
	return timesheet.adapter.Status(timesheet.AccountId, startOfMonth, endOfMonth)
}
