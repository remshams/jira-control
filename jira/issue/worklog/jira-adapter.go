package issue_worklog

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
	"github.com/remshams/jira-control/jira/utils"
)

const issuePath = "rest/api/3/issue/%s/worklog"
const deleteWorklogPath = "rest/api/3/worklog/%s"

type worklogResponseDto struct {
	Total    int          `json:"total"`
	Worklogs []worklogDto `json:"worklogs"`
}

func (w worklogResponseDto) toWorklogs(issueKey string) []Worklog {
	var worklogs []Worklog
	for _, worklog := range w.Worklogs {
		worklogs = append(worklogs, worklog.toWorklog(issueKey))
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
	ID               string `json:"id,omitempty"`
	TimeSpent        string `json:"timeSpent"`
	TimeSpentSeconds int    `json:"timeSpentSeconds,omitempty"`
	Start            string `json:"started,omitempty"`
	Description      string `json:"description,omitempty"`
}

func (w worklogDto) toWorklog(issueKey string) Worklog {
	hoursSpent := float64(w.TimeSpentSeconds) / 3600
	start, err := utils.JiraDateToTime(w.Start)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not parse start time, falling back to unix start time: %v", err)
		start = time.Unix(0, 0)
	}
	return Worklog{
		adapter:            nil,
		issueKey:           issueKey,
		Id:                 w.ID,
		TimeSpentInSeconds: w.TimeSpentSeconds,
		HoursSpent:         hoursSpent,
		Start:              start,
		Description:        w.Description,
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
	path := w.url.JoinPath(fmt.Sprintf(issuePath, worklog.issueKey))
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
		[]utils_http.QueryParam{},
		worklogJson,
		nil,
	)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not log work: %v", err)
		return err
	}
	return nil
}

func (w WorklogJiraAdapter) List(query WorklogListQuery) (WorklogList, error) {
	log.Debugf("WorklogJiraAdapter: Requesting worklog with query %v", query)
	path := w.url.JoinPath(fmt.Sprintf(issuePath, query.issueKey))
	headers := []utils_http.HttpHeader{
		utils_http.CreateBasicAuthHeader(w.username, w.apiToken),
	}
	queryParams := []utils_http.QueryParam{
		{
			Key:   "startedAfter",
			Value: strconv.FormatInt(query.startedAfter.UnixMilli(), 10),
		},
		{
			Key:   "startedBefore",
			Value: strconv.FormatInt(query.startedBefore.UnixMilli(), 10),
		},
	}
	log.Debugf("WorklogJiraAdapter: Query params %v", queryParams)
	_, worklogResponseJson, err := utils_http.PerformRequest(
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
	return worklogResponseDto.toWorklogs(query.issueKey), nil
}

func (w WorklogJiraAdapter) DeleteWorklog(worklog Worklog) error {
	path := w.url.JoinPath(fmt.Sprintf(deleteWorklogPath, worklog.Id))
	headers := []utils_http.HttpHeader{
		utils_http.CreateBasicAuthHeader(w.username, w.apiToken),
	}
	log.Debugf("WorklogJiraAdapter: Perform worklog delete for %s", worklog.Id)
	_, _, err := utils_http.PerformRequest(
		"Worklog Delete",
		path.String(),
		http.MethodDelete,
		headers,
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Could not delete worklog %s", worklog.Id)
		return err
	}
	return nil
}
