package tempo_timesheet

import (
	"fmt"
	"time"

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

func (_ MockTimesheetAdapter) Status(accountId string, from time.Time, to time.Time) (TimesheetStatus, error) {
	log.Debugf("MockTimesheetAdapter: Request status for account: %s, from: %v, to %v", accountId, from, to)
	return NewTimesheetStatus("OPEN", 120, 120), nil
}

func (_ MockTimesheetAdapter) Submit(accountId string, reviewerAccountId string, from time.Time, to time.Time) error {
	log.Debugf(
		"MockTimesheetAdapter: Submit timesheet for accountId: %s with reviewerId: %s from %v to %v",
		accountId,
		reviewerAccountId,
		from,
		to,
	)
	return nil
}
