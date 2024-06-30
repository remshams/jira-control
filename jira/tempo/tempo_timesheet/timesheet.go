package tempo_timesheet

import "github.com/remshams/jira-control/jira/user"

type TimesheetStatus = string

type TimesheetAdapter interface {
	Reviewers(accountId string) ([]user.User, error)
	Status(accountId string) (string, error)
}
