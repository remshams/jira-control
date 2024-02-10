package issue

type IssueSearchRequest struct {
	ProjectName string
	Key         string
	Summary     string
}

func NewIssueSearchRequest() IssueSearchRequest {
	return IssueSearchRequest{}
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
