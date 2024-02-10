package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
)

func main() {
	logger.PrepareLogger()
	jiraAdapter, err := issue_worklog.WorklogJiraAdapterFromEnv()
	if err != nil {
		log.Errorf("Could not create JiraAdapter: %v", err)
		os.Exit(1)
	}
	worklog := issue_worklog.NewWorklog(jiraAdapter, "NC-40", 4.5)
	worklog.Log()
}
