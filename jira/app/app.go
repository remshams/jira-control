package app

import (
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/charmbracelet/log"
	"github.com/remshams/jira-control/jira/favorite"
	"github.com/remshams/jira-control/jira/issue"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	tempo_worklog "github.com/remshams/jira-control/jira/tempo/worklog"
	"github.com/remshams/jira-control/jira/user"
)

type App struct {
	Production          bool
	Url                 url.URL
	TempoUrl            url.URL
	Username            string
	ApiToken            string
	TempoApiToken       string
	IssueAdapter        issue.IssueAdapter
	IssueWorklogAdapter issue_worklog.WorklogAdapter
	FavoriteAdapter     favorite.FavoriteAdapter
	TempoWorklogAdapter tempo_worklog.WorklogListAdapter
	UserAdapter         user.UserAdapter
}

func AppFromEnv() (*App, error) {
	isProduction, err := strconv.ParseBool(os.Getenv("IS_PRODUCTION"))
	if err != nil {
		log.Debug("IS_PRODUCTION is not set, defaulting to false")
		isProduction = false
	}
	url, err := url.Parse(os.Getenv("JIRA_URL"))
	if err != nil || url.String() == "" {
		log.Errorf("App: JIRA_URL not set or invalid: %v", err)
		return nil, errors.New("JIRA_URL not set or invalid")
	}
	tempoUrl, err := url.Parse(os.Getenv("TEMPO_URL"))
	if err != nil || url.String() == "" {
		log.Errorf("App: TEMPO_URL not set or invalid: %v", err)
		return nil, errors.New("TEMPO_URL not set or invalid")
	}
	username := os.Getenv("JIRA_USERNAME")
	if username == "" {
		log.Errorf("App: JIRA_USERNAME is not set")
		return nil, errors.New("JIRA_USERNAME is not set")
	}
	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		log.Errorf("App: JIRA_API_TOKEN is not set")
		return nil, errors.New("JIRA_API_TOKEN is not set")
	}
	tempoApiToken := os.Getenv("TEMPO_API_TOKEN")
	if tempoApiToken == "" {
		log.Errorf("App: TEMPO_API_TOKEN is not set")
		return nil, errors.New("TEMPO_API_TOKEN not set or invalid")
	}
	app := App{
		Production:    isProduction,
		Url:           *url,
		TempoUrl:      *tempoUrl,
		Username:      username,
		ApiToken:      apiToken,
		TempoApiToken: tempoApiToken,
	}
	app.addAdapers()
	return &app, nil
}

func (app *App) addAdapers() {
	if app.Production == true {
		app.IssueWorklogAdapter = issue_worklog.NewWorklogJiraAdapter(app.Url, app.Username, app.ApiToken)
		app.IssueAdapter = issue.NewJiraIssueAdapter(app.IssueWorklogAdapter, app.Url, app.Username, app.ApiToken)
		app.FavoriteAdapter = favorite.NewFavoriteJsonAdapter("favorites.json")
		app.TempoWorklogAdapter = tempo_worklog.NewJiraWorklogAdapter(app.TempoUrl, app.TempoApiToken)
		app.UserAdapter = user.NewJiraUserAdapter(app.Url, app.Username, app.ApiToken)
	} else {
		app.IssueWorklogAdapter = issue_worklog.NewWorklogMockAdapter()
		app.IssueAdapter = issue.NewMockIssueAdapter(app.IssueWorklogAdapter)
		app.FavoriteAdapter = favorite.NewFavoriteJsonAdapter("favorites.json")
		app.TempoWorklogAdapter = tempo_worklog.NewMockWorklogAdapter()
		app.UserAdapter = user.NewMockUserAdapter()
	}
}
