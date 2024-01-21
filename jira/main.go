package main

import (
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	issue_worklog "github.com/remshams/jira-control/jira/worklog"
)

func main() {
	logger.PrepareLogger()
	isProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	if err != nil {
		log.Debug("IS_PRODUCTION is not set, defaulting to false")
		isProduction = false
	}
	var worklogAdapter issue_worklog.WorklogAdapter
	if isProduction == true {
		worklogAdapter = issue_worklog.WorklogJiraAdapter{}
	} else {
		worklogAdapter = issue_worklog.WorklogMockAdatpter{}
	}
	worklog := issue_worklog.NewWorklog(worklogAdapter, "NC-40", 4.5)
	worklog.Log()
}
