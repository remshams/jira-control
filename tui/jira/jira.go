package tui_jira

import jira "github.com/remshams/jira-control/jira/public"

type JiraAdapter struct {
	WorklogAdapter jira.WorklogAdapter
}

func NewJiraAdapter(worklogAdapter jira.WorklogAdapter) JiraAdapter {
	return JiraAdapter{
		WorklogAdapter: worklogAdapter,
	}
}
