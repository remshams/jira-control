package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	jira "github.com/remshams/jira-control/jira/public"
)

func main() {
	logger.PrepareLogger()
	app, err := jira.PrepareApplication()
	if err != nil {
		log.Errorf("Could not create JiraAdapter: %v", err)
		os.Exit(1)
	}
	worklog := issue_worklog.NewWorklog(app.IssueWorklogAdapter, "NC-40", 4.5)
	worklog.Log()
}
