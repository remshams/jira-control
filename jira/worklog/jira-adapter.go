package issue_worklog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	dc_http "github.com/remshams/common/utils/http"
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
	url url.URL
}

func New(url url.URL) WorklogJiraAdapter {
	return WorklogJiraAdapter{
		url: url,
	}
}

func (w WorklogJiraAdapter) logWork(worklog Worklog) error {
	log.Debugf("WorklogJiraAdapter: Logging work %v", worklog)
	path := w.url.JoinPath(fmt.Sprintf(path, worklog.issueKey))
	worklogJson, err := worklogDtoFromWorklog(worklog).toJson()
	if err != nil {
		return err
	}
	headers := []dc_http.HttpHeader{
		{
			Type:  dc_http.ContentType,
			Value: "application/json",
		},
	}
	_, err = dc_http.PerformRequest(
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
