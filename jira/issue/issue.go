package issue

import (
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
)

type IssueSearchRequest struct {
	adapter     IssueAdapter
	ProjectName string
	Key         string
	Summary     string
	Fields      []string
}

func NewIssueSearchRequest(adapter IssueAdapter) IssueSearchRequest {
	return IssueSearchRequest{
		adapter: adapter,
	}
}

func (issueSearchRequest IssueSearchRequest) Search() ([]Issue, error) {
	return issueSearchRequest.adapter.searchIssues(issueSearchRequest)
}

type IssueAdapter interface {
	searchIssues(request IssueSearchRequest) ([]Issue, error)
	worklogs(query issue_worklog.WorklogListQuery) (issue_worklog.WorklogList, error)
}

type IssueProject struct {
	id   string
	Key  string
	Name string
}

func NewIssueProject(id string, key string, name string) IssueProject {
	return IssueProject{
		id:   id,
		Key:  key,
		Name: name,
	}
}

type Issue struct {
	adapter IssueAdapter
	id      string
	Project IssueProject
	Key     string
	Summary string
}

func NewIssue(adapter IssueAdapter, id string, project IssueProject, key string, summary string) Issue {
	return Issue{
		adapter: adapter,
		id:      id,
		Project: project,
		Key:     key,
		Summary: summary,
	}
}

func (issue Issue) WorklogsQuery() issue_worklog.WorklogListQuery {
	return issue_worklog.NewWorklogListQuery(issue.Key)
}

func (issue Issue) Worklogs(query issue_worklog.WorklogListQuery) ([]issue_worklog.Worklog, error) {
	worklogList, err := issue.adapter.worklogs(query)
	if err != nil {
		return nil, err
	}
	return worklogList.SortByStart(true), nil
}
