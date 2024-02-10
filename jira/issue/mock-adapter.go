package issue

type MockIssueAdapter struct{}

func (m MockIssueAdapter) searchIssues(request IssueSearchRequest) ([]Issue, error) {
	return []Issue{
		NewIssue("1", NewIssueProject("1", "P1", "Project 1"), "KEY-1", "Summary 1"),
		NewIssue("2", NewIssueProject("2", "P2", "Project 2"), "KEY-2", "Summary 2"),
		NewIssue("3", NewIssueProject("3", "P3", "Project 3"), "KEY-3", "Summary 3"),
	}, nil
}
