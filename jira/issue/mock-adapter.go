package issue

import (
	"time"

	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
)

type MockIssueAdapter struct {
	worklogAdapter issue_worklog.WorklogAdapter
}

func NewMockIssueAdapter(worklogAdapter issue_worklog.WorklogAdapter) MockIssueAdapter {
	return MockIssueAdapter{
		worklogAdapter,
	}
}

func (m MockIssueAdapter) searchIssues(request IssueSearchRequest) ([]Issue, error) {
	return []Issue{
		NewIssue(m, "1", NewIssueProject("1", "P1", "Project 1", time.Now()), "KEY-1", "Summary 1"),
		NewIssue(m, "2", NewIssueProject("2", "P2", "Project 2", time.Now()), "KEY-2", "Summary 2"),
		NewIssue(m, "3", NewIssueProject("3", "P3", "Project 3", time.Now()), "KEY-3", "Summary 3"),
	}, nil
}

func (m MockIssueAdapter) worklogs(query issue_worklog.WorklogListQuery) (issue_worklog.WorklogList, error) {
	return m.worklogAdapter.List(query)
}
