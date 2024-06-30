package timesheet

import (
	"net/url"

	"github.com/remshams/jira-control/jira/user"
)

type JiraTimesheetAdapter struct {
	url      url.URL
	username string
	apiToken string
}

func NewJiraTimesheetAdapter(url url.URL, username string, apiToken string) JiraTimesheetAdapter {
	return JiraTimesheetAdapter{
		url:      url,
		username: username,
		apiToken: apiToken,
	}
}

func (jiraTimesheetAdapter JiraTimesheetAdapter) Approvers(accountId string) ([]user.User, error) {
	return nil, nil
}
