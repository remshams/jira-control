package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	jira "github.com/remshams/jira-control/jira/public"
)

func main() {
	logger.PrepareLogger()
	app, err := jira.PrepareApplication()
	if err != nil {
		log.Errorf("Could not create JiraAdapter: %v", err)
		os.Exit(1)
	}
	user, err := app.UserAdapter.Myself()
	if err != nil {
		log.Error("Could not load myself")
		os.Exit(1)
	}
	log.Debugf("User: %s %s", user.AccountId, user.Name)
}
