package tempo_timesheet

import "github.com/remshams/jira-control/jira/user"

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
	Status(accountId string) (TimesheetStatus, error)
}
