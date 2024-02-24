package issue_worklog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
)

const path = "rest/api/3/issue/%s/worklog"

type worklogResponseDto struct {
	Total    int          `json:"total"`
	Worklogs []worklogDto `json:"worklogs"`
}

func (w worklogResponseDto) getWorklogs() []Worklog {
	var worklogs []Worklog
	for _, worklog := range w.Worklogs {
		worklogs = append(worklogs)
	}
	return worklogs
}

func worklogResponseDtoFromJson(worklogResponseJson []byte) (*worklogResponseDto, error) {
	var worklogResponseDto worklogResponseDto
	err := json.Unmarshal(worklogResponseJson, &worklogResponseDto)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not unmarshal worklogReponse: %v", err)
		return nil, err
	}
	return &worklogResponseDto, nil
}

type worklogDto struct {
	TimeSpent        string `json:"timeSpent"`
	TimeSpentSeconds int    `json:"timeSpentSeconds"`
	Start            string `json:"started,omitempty"`
	Description      string `json:"description,omitempty"`
}

func (w worklogDto) toWorklog() Worklog {
	hoursSpent := w.TimeSpentSeconds / (60 * 60)
	return Worklog{
		adapter:     nil,
		issueKey:    "",
		HoursSpent:  fmt.Sprintf("%d", hoursSpent),
		Start:       utils_http.ParseStart(w.Start),
		Description: w.Description,
	}
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

func (w WorklogJiraAdapter) list(query WorklogListQuery) ([]Worklog, error) {
	log.Debugf("WorklogJiraAdapter: Requesting worklog with query %v", query)
	path := w.url.JoinPath(fmt.Sprintf(path, query.issueKey))
	headers := []utils_http.HttpHeader{
		utils_http.CreateBasicAuthHeader(w.username, w.apiToken),
	}
	queryParams := []utils_http.QueryParam{
		{
			Key:   "startedAfter",
			Value: fmt.Sprintf("%d", query.startedAfter),
		},
		{
			Key:   "startedBefore",
			Value: fmt.Sprintf("%d", query.startedBefore),
		},
	}
	res, worklogResponseJson, err := utils_http.PerformRequest(
		"Worklog List",
		path.String(),
		http.MethodGet,
		headers,
		queryParams,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not perform request: %v", err)
		return nil, err
	}
	worklogResponseDto, err := worklogResponseDtoFromJson(worklogResponseJson)
	if err != nil {
		return nil, err
	}
	return worklogDtoToWorklog(worklogResponseDto.Worklogs), nil
}
