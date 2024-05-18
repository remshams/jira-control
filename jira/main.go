package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	jira "github.com/remshams/jira-control/jira/public"
	tempo_worklog "github.com/remshams/jira-control/jira/tempo/worklog"
)

func main() {
	logger.PrepareLogger()
	app, err := jira.PrepareApplication()
	if err != nil {
		log.Errorf("Could not create JiraAdapter: %v", err)
		os.Exit(1)
	}
	query := tempo_worklog.NewWorklistQuery()
	queries, err := app.TempoWorklogAdapter.List(query)
	if err != nil {
		log.Error("Could not load tempo worklogs")
		os.Exit(1)
	}
	log.Debugf("Number of worklogs: %d", len(queries))
}
