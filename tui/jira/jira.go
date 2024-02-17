package tui_jira

import jira "github.com/remshams/jira-control/jira/public"

type JiraAdapter struct {
	IssueAdapter   jira.IssueAdapter
	WorklogAdapter jira.WorklogAdapter
}

func NewJiraAdapter(issueAdapter jira.IssueAdapter, worklogAdapter jira.WorklogAdapter) JiraAdapter {
	return JiraAdapter{
		IssueAdapter:   issueAdapter,
		WorklogAdapter: worklogAdapter,
	}
}
