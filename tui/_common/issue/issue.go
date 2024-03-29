package common_issue

import (
	"github.com/remshams/jira-control/jira/issue"
	jira "github.com/remshams/jira-control/jira/public"
)

func FindIssue(issues []jira.Issue, key string) *issue.Issue {
	for _, issue := range issues {
		if issue.Key == key {
			return &issue
		}
	}
	return nil
}
