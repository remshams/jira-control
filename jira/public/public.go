package jira

import (
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	issue_worklog "github.com/remshams/jira-control/jira/worklog"
)

type Worklog = issue_worklog.Worklog
type WorklogAdapter = issue_worklog.WorklogAdapter
type WorklogMockAdapter = issue_worklog.WorklogMockAdatpter
type WorklogJiraAdapter = issue_worklog.WorklogJiraAdapter

func NewWorklogJiraAdapter() WorklogJiraAdapter {
	return issue_worklog.WorklogJiraAdapter{}
}

func NewWorklogMockAdapter() WorklogMockAdapter {
	return issue_worklog.WorklogMockAdatpter{}
}

func NewWorklog(adapter WorklogAdapter, issueKey string, hoursSpent float64) Worklog {
	return issue_worklog.NewWorklog(adapter, issueKey, hoursSpent)
}

func PrepareApplication() (issue_worklog.WorklogAdapter, error) {
	isProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	if err != nil {
		log.Debug("IS_PRODUCTION is not set, defaulting to false")
		isProduction = false
	}
	var worklogAdapter issue_worklog.WorklogAdapter
	if isProduction == true {
		worklogAdapter, err = issue_worklog.WorklogJiraAdapterFromEnv()
		if err != nil {
			log.Errorf("Could not create JiraAdapter: %v", err)
			return nil, err
		}
	} else {
		worklogAdapter = issue_worklog.WorklogMockAdatpter{}
	}
	return worklogAdapter, nil
}
