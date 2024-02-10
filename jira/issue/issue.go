package issue

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
	id      string
	Project IssueProject
	Key     string
	Summary string
}

func NewIssue(id string, project IssueProject, key string, summary string) Issue {
	return Issue{
		id:      id,
		Project: project,
		Key:     key,
		Summary: summary,
	}
}
