package tempo_timesheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
	jira_common_http "github.com/remshams/jira-control/jira/common"
	"github.com/remshams/jira-control/jira/user"
)

const reviewersPath = "/4/timesheet-approvals/user/%s/reviewers"
const statusPath = "/4/timesheet-approvals/user/%s"
const submitPath = "/4/timesheet-approvals/user/%s/submit"

type reviewerDto struct {
	AccountId string `json:"accountId"`
}

type reviewersDto struct {
	Results []reviewerDto `json:"results"`
}

func reviewersFromJson(body []byte) ([]reviewerDto, error) {
	var reviewersDto reviewersDto
	err := json.Unmarshal(body, &reviewersDto)
	if err != nil {
		return nil, err
	}
	return reviewersDto.Results, nil
}

type statusDto struct {
	Key string `json:"key"`
}

type timesheetStatusDto struct {
	RequiredSeconds  int       `json:"requiredSeconds"`
	TimeSpentSeconds int       `json:"timeSpentSeconds"`
	Status           statusDto `json:"status"`
}

func (timesheetStatusDto timesheetStatusDto) toTimeSheetStatus() TimesheetStatus {
	return NewTimesheetStatus(
		timesheetStatusDto.Status.Key,
		timesheetStatusDto.RequiredSeconds/3600,
		timesheetStatusDto.TimeSpentSeconds/3600,
	)
}

func timesheetStatusFromJson(body []byte) (timesheetStatusDto, error) {
	var timesheetStatusDto timesheetStatusDto
	err := json.Unmarshal(body, &timesheetStatusDto)
	if err != nil {
		return timesheetStatusDto, err
	}
	return timesheetStatusDto, nil
}

type submitTimesheetDto struct {
	ReviewerAccountId string `json:"reviewerAccountId"`
}

func NewSubmitTimesheetDto(reviewerAccountId string) submitTimesheetDto {
	return submitTimesheetDto{
		ReviewerAccountId: reviewerAccountId,
	}
}

func (submitTimesheetDto submitTimesheetDto) toJson() ([]byte, error) {
	return json.Marshal(submitTimesheetDto)
}

type JiraTimesheetAdapter struct {
	userAdapter user.UserAdapter
	url         url.URL
	apiToken    string
}

func NewJiraTimesheetAdapter(url url.URL, apiToken string, userAdapter user.UserAdapter) JiraTimesheetAdapter {
	return JiraTimesheetAdapter{
		userAdapter: userAdapter,
		url:         url,
		apiToken:    apiToken,
	}
}

func (jiraTimesheetAdapter JiraTimesheetAdapter) Reviewers(accountId string) ([]user.User, error) {
	log.Debugf("JiraTimesheetAdapter: Request reviewers for account: %s", accountId)
	path := jiraTimesheetAdapter.url.JoinPath(fmt.Sprintf(reviewersPath, accountId))
	_, body, err := utils_http.PerformRequest(
		"JiraTimesheetAdapter",
		path.String(),
		http.MethodGet,
		jira_common_http.CreateDefaultTempoHttpHeaders(jiraTimesheetAdapter.apiToken),
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("JiraTimesheetAdapter: Could not perform request %v", err)
		return nil, err
	}
	reviewerDtos, err := reviewersFromJson(body)
	if err != nil {
		log.Errorf("JiraTimesheetAdapter: Could not parse reviewers: %v", err)
		return nil, err
	}
	reviewerAccountIds := []string{}
	for _, reviewerDto := range reviewerDtos {
		reviewerAccountIds = append(reviewerAccountIds, reviewerDto.AccountId)
	}
	reviewers, err := jiraTimesheetAdapter.userAdapter.Users(reviewerAccountIds)
	if err != nil {
		log.Errorf("JiraTimesheetAdapter: Could not request accounts of reviewers: %v", err)
		return nil, err
	}
	return reviewers, nil
}

func (jiraTimesheetAdapter JiraTimesheetAdapter) Status(accountId string, from time.Time, to time.Time) (TimesheetStatus, error) {
	log.Debugf("JiraTimesheetAdapter: Request timesheet status for accountId: %s from %v to %v", accountId, from, to)
	path := jiraTimesheetAdapter.url.JoinPath(fmt.Sprintf(statusPath, accountId))
	params := []utils_http.QueryParam{
		{
			Key:   "from",
			Value: from.Format("2006-01-02"),
		},
		{
			Key:   "to",
			Value: to.Format("2006-01-02"),
		},
	}
	_, body, err := utils_http.PerformRequest(
		"JiraTimesheetAdapter",
		path.String(),
		http.MethodGet,
		jira_common_http.CreateDefaultTempoHttpHeaders(jiraTimesheetAdapter.apiToken),
		params,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("JiraTimesheetAdapter: Could not perform request %v", err)
		return TimesheetStatus{}, err
	}
	timesheetDto, err := timesheetStatusFromJson(body)
	if err != nil {
		log.Errorf("JiraTimesheetAdapter: Could not parse timesheet status: %v", err)
		return TimesheetStatus{}, err
	}
	return timesheetDto.toTimeSheetStatus(), nil
}

func (jiraTimesheetAdapter JiraTimesheetAdapter) Submit(
	accountId string,
	reviewerAccountId string,
	from time.Time,
	to time.Time,
) error {
	log.Debugf("JiraTimesheetAdapter: Submit timesheet for accountId %s and reviewerId %s from %v to %v",
		accountId,
		reviewerAccountId,
		from,
		to,
	)
	path := jiraTimesheetAdapter.url.JoinPath(fmt.Sprintf(submitPath, accountId))
	body, err := NewSubmitTimesheetDto(reviewerAccountId).toJson()
	params := []utils_http.QueryParam{
		{
			Key:   "from",
			Value: from.Format("2006-01-02"),
		},
		{
			Key:   "to",
			Value: to.Format("2006-01-02"),
		},
	}
	if err != nil {
		log.Debugf("Could not marshal request body %v", err)
		return err
	}
	_, _, err = utils_http.PerformRequest(
		"JiraTimesheetAdapter",
		path.String(),
		http.MethodPost,
		jira_common_http.CreateDefaultTempoHttpHeaders(jiraTimesheetAdapter.apiToken),
		params,
		body,
		nil,
	)
	if err != nil {
		log.Errorf("JiraTimesheetAdapter: Could not perform request %v", err)
		return err
	}
	return nil
}
