package issue_worklog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/http"
)

const path = "rest/api/3/issue/%s/worklog"

type worklogDto struct {
	TimeSpent string `json:"timeSpent"`
}

func worklogDtoFromWorklog(worklog Worklog) worklogDto {
	return worklogDto{
		TimeSpent: fmt.Sprintf("%fh", worklog.hours()),
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
	apiToken string
}

func NewWorklogJiraAdapter(url url.URL, apiToken string) WorklogJiraAdapter {
	return WorklogJiraAdapter{
		url:      url,
		apiToken: apiToken,
	}
}

func WorklogJiraAdapterFromEnv() (*WorklogJiraAdapter, error) {
	url, err := url.Parse(os.Getenv("JIRA_URL"))
	if err != nil || url.String() == "" {
		log.Errorf("WorklogJiraAdapter: JIRA_URL not set or invalid: %v", err)
		return nil, errors.New("JIRA_URL not set or invalid")
	}
	apiToken := os.Getenv("JIRA_API_TOKEN")
	if apiToken == "" {
		log.Errorf("WorklogJiraAdapter: JIRA_API_TOKEN is not set")
		return nil, errors.New("JIRA_API_TOKEN is not set")
	}
	adapter := NewWorklogJiraAdapter(*url, apiToken)
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
		{
			Type:  utils_http.Authorization,
			Value: fmt.Sprintf("Bearer %s", w.apiToken),
		},
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
