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
	favorites, err := app.FavoriteAdapter.Load()
	favorite := favorites[0]
	favorite.HoursSpent = 10.5
	// favorite := favorite.NewFavorite(app.FavoriteAdapter, "NX-Testing", 9)
	err = favorite.Store()
	if err != nil {
		log.Errorf("Could not store favorite: %v", err)
		os.Exit(1)
	}
}
