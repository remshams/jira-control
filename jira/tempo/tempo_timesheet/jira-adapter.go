package tempo_timesheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
	jira_common_http "github.com/remshams/jira-control/jira/common"
	"github.com/remshams/jira-control/jira/user"
)

const reviewersPath = "/4/timesheet-approvals/user/%s/reviewers"

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
