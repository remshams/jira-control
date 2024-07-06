package tempo_timesheet

import (
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/remshams/jira-control/jira/user"
)

type MockTimesheetAdapter struct{}

func NewMockTimesheetAdapter() MockTimesheetAdapter {
	return MockTimesheetAdapter{}
}

func (_ MockTimesheetAdapter) Reviewers(accountId string) ([]user.User, error) {
	log.Debugf("MockTimesheetAdapter: Request approvers for: %s", accountId)
	return []user.User{user.NewUser("0", "mock user", fmt.Sprintf("mock.%s@mock.com", accountId))}, nil
}

func (_ MockTimesheetAdapter) Status(accountId string) (TimesheetStatus, error) {
	log.Debugf("MockTimesheetAdapter) Request status for account: %s", accountId)
	return NewTimesheetStatus("OPEN", 120, 120), nil
}
