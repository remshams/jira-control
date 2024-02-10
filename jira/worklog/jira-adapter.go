package issue_worklog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
)

const path = "rest/api/3/issue/%s/worklog"

type worklogDto struct {
	TimeSpent   string `json:"timeSpent"`
	Start       string `json:"started,omitempty"`
	Description string `json:"comment,omitempty"`
}

func worklogDtoFromWorklog(worklog Worklog) worklogDto {
	return worklogDto{
		TimeSpent:   fmt.Sprintf("%fh", worklog.hoursSpent),
		Start:       worklog.start.Format(time.RFC3339),
		Description: worklog.description,
	}
}

func (worklogDto worklogDto) toJson() ([]byte, error) {
	json, err := json.Marshal(worklogDto)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not marshal worklogDto: %v", err)
		return nil, err
	}
	return json, nil
}

type WorklogJiraAdapter struct {
	url      url.URL
	username string
	apiToken string
}

func NewWorklogJiraAdapter(url url.URL, username string, apiToken string) WorklogJiraAdapter {
	return WorklogJiraAdapter{
		url:      url,
		username: username,
		apiToken: apiToken,
	}
}

func WorklogJiraAdapterFromEnv() (*WorklogJiraAdapter, error) {
	url, err := url.Parse(os.Getenv("JIRA_URL"))
	if err != nil || url.String() == "" {
		log.Errorf("WorklogJiraAdapter: JIRA_URL not set or invalid: %v", err)
		return nil, errors.New("JIRA_URL not set or invalid")
	}
	username := os.Getenv("JIRA_USERNAME")
	if username == "" {
		log.Errorf("WorklogJiraAdapter: JIRA_USERNAME is not set")
		return nil, errors.New("JIRA_USERNAME is not set")
	}
	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		log.Errorf("WorklogJiraAdapter: JIRA_API_TOKEN is not set")
		return nil, errors.New("JIRA_API_TOKEN is not set")
	}
	adapter := NewWorklogJiraAdapter(*url, username, apiToken)
	return &adapter, nil
}

func (w WorklogJiraAdapter) logWork(worklog Worklog) error {
	log.Debugf("WorklogJiraAdapter: Logging work %v", worklog)
	path := w.url.JoinPath(fmt.Sprintf(path, worklog.issueKey))
	worklogJson, err := worklogDtoFromWorklog(worklog).toJson()
	if err != nil {
		return err
	}
	headers := []utils_http.HttpHeader{
		{
			Type:  utils_http.ContentType,
			Value: "application/json",
		},
		utils_http.CreateBasicAuthHeader(w.username, w.apiToken),
	}
	_, err = utils_http.PerformRequest(
		"Worklog",
		path.String(),
		http.MethodPost,
		headers,
		worklogJson,
		nil,
	)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not log work: %v", err)
		return err
	}
	return nil
}
