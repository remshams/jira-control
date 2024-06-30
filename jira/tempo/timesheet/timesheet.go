package timesheet

import "github.com/remshams/jira-control/jira/user"

type TimesheetAdapter interface {
	Reviewers(accountId string) ([]user.User, error)
}
