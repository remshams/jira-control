package issue_worklog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
)

const path = "rest/api/3/issue/%s/worklog"

type worklogDto struct {
	TimeSpent   string `json:"timeSpent"`
	Start       string `json:"started,omitempty"`
	Description string `json:"description,omitempty"`
}

func worklogDtoFromWorklog(worklog Worklog) worklogDto {
	layout := "2006-01-02T15:04:05.000+0000"
	log.Debugf("Worklog start %s", worklog.Start.Format(layout))
	return worklogDto{
		TimeSpent:   fmt.Sprintf("%fh", worklog.HoursSpent),
		Start:       worklog.Start.Format(layout),
		Description: worklog.Description,
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
	_, _, err = utils_http.PerformRequest(
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