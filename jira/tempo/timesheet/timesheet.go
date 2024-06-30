package timesheet

import "github.com/remshams/jira-control/jira/user"

type Approver = user.User

type TimesheetAdapter interface {
	Approvers(accountId string) ([]user.User, error)
}
