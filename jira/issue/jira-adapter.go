package issue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
)

const path = "rest/api/3/search"
const worklogPath = "rest/api/3/issue/%s/worklog"

type issueSearchRequestDto struct {
	Jql    string `json:"jql"`
	Fields string `json:"fields,omitempty"`
}

func (issueSearchRequestDto issueSearchRequestDto) toJson() ([]byte, error) {
	json, err := json.Marshal(issueSearchRequestDto)
	if err != nil {
		log.Errorf("JiraIssueAdapter: Could not marshal issueSearchRequestDto: %v", err)
		return nil, err
	}
	return json, nil
}

func fromIssueSearchRequest(request IssueSearchRequest) issueSearchRequestDto {
	return issueSearchRequestDto{
		Jql:    jqlFromSearchRequest(request),
		Fields: strings.Join(request.Fields, ","),
	}
}

type issueProjectDto struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type issueFieldDto struct {
	Summary string          `json:"summary"`
	Project issueProjectDto `json:"project"`
}

type issueDto struct {
	Id     string        `json:"id"`
	Key    string        `json:"key"`
	Fields issueFieldDto `json:"fields"`
}

type issueSearchResponseDto struct {
	Issues []issueDto `json:"issues"`
}

func fromJson(body []byte) (issueSearchResponseDto, error) {
	var issueSearchResponse issueSearchResponseDto
	err := json.Unmarshal(body, &issueSearchResponse)
	if err != nil {
		log.Errorf("JiraIssueAdapter: Could not unmarshal issueSearchResponseDto: %v", err)
		return issueSearchResponseDto{}, err
	}
	return issueSearchResponse, nil
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
	return strings.Join(jql, " OR ")
}

type JiraIssueAdapter struct {
	worklogAdapter issue_worklog.WorklogAdapter
	url            url.URL
	username       string
	apiToken       string
}

func NewJiraIssueAdapter(worklogAdapter issue_worklog.WorklogAdapter, url url.URL, username string, apiToken string) JiraIssueAdapter {
	return JiraIssueAdapter{
		worklogAdapter: worklogAdapter,
		url:            url,
		username:       username,
		apiToken:       apiToken,
	}
}

func (jiraIssueAdapter JiraIssueAdapter) searchIssues(request IssueSearchRequest) ([]Issue, error) {
	searchRequestDto, err := fromIssueSearchRequest(request).toJson()
	if err != nil {
		return nil, err
	}
	path := jiraIssueAdapter.url.JoinPath(path)
	headers := []utils_http.HttpHeader{
		{
			Type:  utils_http.ContentType,
			Value: "application/json",
		},
		utils_http.CreateBasicAuthHeader(jiraIssueAdapter.username, jiraIssueAdapter.apiToken),
	}
	timeout := 5 * time.Second
	_, body, err := utils_http.PerformRequest(
		"JiraIssueAdapter",
		path.String(),
		http.MethodPost,
		headers,
		[]utils_http.QueryParam{},
		searchRequestDto,
		&timeout,
	)
	if err != nil {
		log.Errorf("JiraIssueAdapter: Could not perform request: %v", err)
		return nil, err
	}
	issueSearchResponseDto, err := fromJson(body)
	if err != nil {
		return nil, err
	}
	return issuesFromDto(jiraIssueAdapter, issueSearchResponseDto.Issues), nil
}

func issuesFromDto(adapter JiraIssueAdapter, issuesDto []issueDto) []Issue {
	issues := []Issue{}
	for _, issue := range issuesDto {
		issues = append(
			issues,
			NewIssue(
				adapter,
				issue.Id,
				NewIssueProject(
					issue.Fields.Project.Key,
					issue.Fields.Project.Key,
					issue.Fields.Project.Name,
				),
				issue.Key,
				issue.Fields.Summary,
			))
	}
	return issues
}

func (jiraIssueAdapter JiraIssueAdapter) worklogs(query issue_worklog.WorklogListQuery) (issue_worklog.WorklogList, error) {
	return jiraIssueAdapter.worklogAdapter.List(query)
}
