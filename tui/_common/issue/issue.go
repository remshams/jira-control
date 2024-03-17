package common_issue

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/jira-control/jira/issue"
	jira "github.com/remshams/jira-control/jira/public"
)

type LogWorkAction struct {
	Issue jira.Issue
}

func CreateLogWorkAction(issue jira.Issue) tea.Cmd {
	return func() tea.Msg {
		return LogWorkAction{
			Issue: issue,
		}
	}
}

func FindIssue(issues []jira.Issue, key string) *issue.Issue {
	for _, issue := range issues {
		if issue.Key == key {
			return &issue
		}
	}
	return nil
}
