package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	"github.com/remshams/jira-control/jira/favorite"
	jira "github.com/remshams/jira-control/jira/public"
)

func main() {
	logger.PrepareLogger()
	app, err := jira.PrepareApplication()
	if err != nil {
		log.Errorf("Could not create JiraAdapter: %v", err)
		os.Exit(1)
	}
	// favorites, err := app.FavoriteAdapter.Load()
	// favorite1 := favorites[0]
	// favorite1.HoursSpent = 10.5
	favorite1 := favorite.NewFavorite(app.FavoriteAdapter, "NX-Testing3", 9)
	favorite2 := favorite.NewFavorite(app.FavoriteAdapter, "NX-Testing4", 10)
	err = favorite1.Store()
	err = favorite2.Store()
	if err != nil {
		log.Errorf("Could not store favorite: %v", err)
		os.Exit(1)
	}
}
