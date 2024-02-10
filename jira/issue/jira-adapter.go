package issue

import (
	"fmt"
	"net/url"
	"strings"
)

type issueSearchRequestDto struct {
	Jql    string `json:"jql"`
	Fields string `json:"fields,omitempty"`
}

func fromIssueSearchRequest(request IssueSearchRequest) issueSearchRequestDto {
	return issueSearchRequestDto{
		Jql:    jqlFromSearchRequest(request),
		Fields: strings.Join(request.Fields, ","),
	}
}

func jqlFromSearchRequest(request IssueSearchRequest) string {
	jql := []string{}
	if request.Summary != "" {
		jql = append(jql, fmt.Sprintf("summary ~ \"%s\"", request.Summary))
	}
	if request.Key != "" {
		jql = append(jql, fmt.Sprintf("key = \"%s\"", request.Key))
	}
	if request.ProjectName != "" {
		jql = append(jql, fmt.Sprintf("project = \"%s\"", request.ProjectName))
	}
	return strings.Join(jql, " AND ")
}

type JiraIssueAdapter struct {
	url      url.URL
	username string
	apiToken string
}

func NewJiraIssueAdapter(url url.URL, username string, apiToken string) JiraIssueAdapter {
	return JiraIssueAdapter{
		url:      url,
		username: username,
		apiToken: apiToken,
	}
}

func searchIssues(request IssueSearchRequest) ([]Issue, error) {
	return []Issue{}, nil
}
