package issue

import (
	"time"

	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	"github.com/remshams/jira-control/jira/utils"
)

type IssueSearchRequest struct {
	adapter     IssueAdapter
	projectName string
	key         string
	summary     string
	fields      []string
	updatedBy   string
	orderBy     utils.OrderBy
}

func NewIssueSearchRequest(adapter IssueAdapter) IssueSearchRequest {
	return IssueSearchRequest{
		adapter: adapter,
	}
}

func (request IssueSearchRequest) WithProjectName(projectName string) IssueSearchRequest {
	request.projectName = projectName
	request.fields = []string{"id", "key", "summary", "updated"}
	return request
}

func (request IssueSearchRequest) WithKey(key string) IssueSearchRequest {
	request.key = key
	return request
}

func (request IssueSearchRequest) WithSummary(summary string) IssueSearchRequest {
	request.summary = summary
	return request
}

func (request IssueSearchRequest) WithUpdatedBy(updatedBy string) IssueSearchRequest {
	request.updatedBy = updatedBy
	return request
}

func (request IssueSearchRequest) WithOrderBy(orderBy utils.OrderBy) IssueSearchRequest {
	request.orderBy = orderBy
	return request
}

func (issueSearchRequest IssueSearchRequest) Search() ([]Issue, error) {
	return issueSearchRequest.adapter.searchIssues(issueSearchRequest)
}

type IssueAdapter interface {
	searchIssues(request IssueSearchRequest) ([]Issue, error)
	worklogs(query issue_worklog.WorklogListQuery) (issue_worklog.WorklogList, error)
}

type IssueProject struct {
	id      string
	Key     string
	Name    string
	Updated time.Time
}

func NewIssueProject(id string, key string, name string, updated time.Time) IssueProject {
	return IssueProject{
		id:      id,
		Key:     key,
		Name:    name,
		Updated: updated,
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
