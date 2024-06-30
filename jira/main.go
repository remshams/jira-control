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
	user, err = app.UserAdapter.User(user.AccountId)
	log.Debugf("User: %s %s", user.AccountId, user.Name)
	users, err := app.UserAdapter.Users([]string{user.AccountId})
	log.Debugf("Users: %v", users)
	reviewers, err := app.TempoTimesheetAdapter.Reviewers(user.AccountId)
	log.Debugf("Reviewers: %v", reviewers)
}
