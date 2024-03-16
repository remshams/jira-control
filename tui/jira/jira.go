package tui_jira

import jira "github.com/remshams/jira-control/jira/public"

type JiraAdapter struct {
	App *jira.App
}

func NewJiraAdapter(app *jira.App) JiraAdapter {
	return JiraAdapter{
		App: app,
	}
}
