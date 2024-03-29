package common_worklog

import (
	tea "github.com/charmbracelet/bubbletea"
	jira "github.com/remshams/jira-control/jira/public"
)

type LogWorkAction struct {
	IssueKey string
}

func CreateLogWorkAction(issueKey string) tea.Cmd {
	return func() tea.Msg {
		return LogWorkAction{
			IssueKey: issueKey,
		}
	}
}

type ShowWorklogsAction struct {
	Issue jira.Issue
}

func CreateShowWorklogsAction(issue jira.Issue) tea.Cmd {
	return func() tea.Msg {
		return ShowWorklogsAction{
			Issue: issue,
		}
	}
}
