package tempo_worklog

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
	"github.com/remshams/jira-control/jira/utils"
)

const worklogPath = "/4/worklogs"

type worklogResponseMetadataDto struct {
	Count  int `json:"count"`
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type worklogIssueDto struct {
	Id int `json:"id"`
}

type worklogDto struct {
	TempoWorklogId   int             `json:"tempoWorklogId"`
	Issue            worklogIssueDto `json:"worklogIssueDto"`
	TimeSpentSeconds int             `json:"timeSpentSeconds"`
	BillableSeconds  int             `json:"billableSeconds"`
	StartDate        string          `json:"startDate"`
	StartTime        string          `json:"startTime"`
	CreatedAt        string          `json:"createdAt"`
	UpdatedAt        string          `json:"updatedAt"`
	Description      string          `json:"description"`
}

func (w worklogDto) toWorklog() (*Worklog, error) {
	start, err := utils.TempoDateToTime(w.StartDate, w.StartTime)
	if err != nil {
		return nil, err
	}
	worklog := NewWorklog(
		w.Issue.Id,
		w.TempoWorklogId,
		w.TimeSpentSeconds,
		w.BillableSeconds,
		start,
		w.Description,
	)
	return &worklog, nil
}

type worklogResponseDto struct {
	Metadata worklogResponseMetadataDto `json:"metadata"`
	Results  []worklogDto               `json:"results"`
}

func (w worklogResponseDto) toWorklogs() ([]Worklog, error) {
	worklogs := []Worklog{}
	for _, worklogDto := range w.Results {
		worklog, err := worklogDto.toWorklog()
		if err != nil {
			return []Worklog{}, err
		}
		worklogs = append(worklogs, *worklog)
	}
	return worklogs, nil
}

func fromJson(jsonBytes []byte) (worklogResponseDto, error) {
	var jsonResponseDto worklogResponseDto
	err := json.Unmarshal(jsonBytes, &jsonResponseDto)
	return jsonResponseDto, err
}

type JiraWorklogAdapter struct {
	url      url.URL
	apiToken string
}

func NewJiraWorklogAdapter(url url.URL, apiToken string) JiraWorklogAdapter {
	return JiraWorklogAdapter{
		url:      url,
		apiToken: apiToken,
	}
}

func (w JiraWorklogAdapter) List(query WorklogListQuery) ([]Worklog, error) {
	log.Debug("WorklogJiraAdapter: Query worklogs with: %v", query)
	headers := []utils_http.HttpHeader{utils_http.CreateBearerTokenHeader(w.apiToken)}
	queryParams := []utils_http.QueryParam{
		{
			Key:   "from",
			Value: utils.TimeToTempoDate(query.from),
		},
		{
			Key:   "to",
			Value: utils.TimeToTempoDate(query.to),
		},
	}
	log.Debugf("WorklogJiraAdapter: Query params: %v", queryParams)
	_, worklogResponseBytes, err := utils_http.PerformRequest(
		"Tempo Worklog List",
		w.url.JoinPath(worklogPath).String(),
		http.MethodGet,
		headers,
		queryParams,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("WorklogJiraAdapter: Worklog request failed: %v", err)
		return []Worklog{}, err
	}
	worklogResponseDto, err := fromJson(worklogResponseBytes)
	if err != nil {
		log.Error("Could not parse json result")
		return []Worklog{}, err
	}
	worklogs, err := worklogResponseDto.toWorklogs()
	if err != nil {
		log.Error("Could not parse worklogs: %v", worklogs)
	}
	return worklogs, nil
}
