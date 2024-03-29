package common_worklog

import (
	tea "github.com/charmbracelet/bubbletea"
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
