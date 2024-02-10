package issue

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
)

const path = "rest/api/3/search"

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
	res, err := utils_http.PerformRequest(
		"JiraIssueAdapter",
		path.String(),
		http.MethodPost,
		headers,
		searchRequestDto,
		nil,
	)
	if err != nil {
		log.Errorf("JiraIssueAdapter: Could not perform request: %v", err)
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Errorf("JiraIssueAdapter: Could not read response body: %v", err)
		return nil, err
	}
	issueSearchResponseDto, err := fromJson(body)
	if err != nil {
		return nil, err
	}
	return issuesFromDto(issueSearchResponseDto.Issues), nil
}

func issuesFromDto(issuesDto []issueDto) []Issue {
	issues := []Issue{}
	for _, issue := range issuesDto {
		issues = append(
			issues,
			NewIssue(
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
