package user

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
	jira_common_http "github.com/remshams/jira-control/jira/common"
)

type UserQueryParamType = string

const (
	AccountId UserQueryParamType = "accountId"
)

const userPath = "rest/api/3/user"
const usersPath = "rest/api/3/user/bulk"
const myselfPath = "rest/api/3/myself"

type userDto struct {
	AccountId string `json:"accountId"`
	Email     string `json:"emailAddress"`
	Name      string `json:"displayName"`
}

func (userDto userDto) toUser() User {
	return NewUser(userDto.AccountId, userDto.Name, userDto.Email)
}

func userFromJson(body []byte) (User, error) {
	var userDto userDto
	err := json.Unmarshal(body, &userDto)
	if err != nil {
		return User{}, err
	}
	return userDto.toUser(), nil

}

type usersDto struct {
	Values []userDto `json:"values"`
}

func (usersDto usersDto) toUsers() []User {
	users := []User{}
	for _, userDto := range usersDto.Values {
		users = append(users, userDto.toUser())
	}
	return users
}

func usersFromJson(body []byte) ([]User, error) {
	var usersDto usersDto
	err := json.Unmarshal(body, &usersDto)
	if err != nil {
		return nil, err
	}
	return usersDto.toUsers(), nil
}

type JiraUserAdapter struct {
	url      url.URL
	username string
	apiToken string
}

func NewJiraUserAdapter(url url.URL, username string, apiToken string) JiraUserAdapter {
	return JiraUserAdapter{
		url:      url,
		username: username,
		apiToken: apiToken,
	}
}

func (jiraUserAdapter JiraUserAdapter) Myself() (User, error) {
	log.Debug("JiraUserAdapter: Requesting myself")
	path := jiraUserAdapter.url.JoinPath(myselfPath)
	_, body, err := utils_http.PerformRequest(
		"JiraUserAdapter",
		path.String(),
		http.MethodGet,
		jira_common_http.CreateDefaultHttpHeaders(jiraUserAdapter.username, jiraUserAdapter.apiToken),
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not perform request: %v", err)
		return User{}, err
	}
	user, err := userFromJson(body)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not parse response body %v", err)
		return User{}, err
	}
	return user, nil

}

func (jiraUserAdapter JiraUserAdapter) User(accountId string) (User, error) {
	log.Debugf("JiraUserAdapter: Requesting user with accountId: %s", accountId)
	path := jiraUserAdapter.url.JoinPath(userPath)
	headers := jira_common_http.CreateDefaultHttpHeaders(jiraUserAdapter.username, jiraUserAdapter.apiToken)
	queryParams := []utils_http.QueryParam{
		{
			Key:   AccountId,
			Value: accountId,
		},
	}
	_, body, err := utils_http.PerformRequest(
		"JiraUserAdapter",
		path.String(),
		http.MethodGet,
		headers,
		queryParams,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not perform request: %v", err)
		return User{}, nil
	}
	user, err := userFromJson(body)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not parse response body: %v", err)
		return User{}, nil
	}
	return user, nil
}

func (jiraUserAdapter JiraUserAdapter) Users(accountIds []string) ([]User, error) {
	log.Debugf("JiraUserAdapter: Requesting users with #accountIds: %d", len(accountIds))
	path := jiraUserAdapter.url.JoinPath(usersPath)
	headers := jira_common_http.CreateDefaultHttpHeaders(jiraUserAdapter.username, jiraUserAdapter.apiToken)
	queryParams := []utils_http.QueryParam{}
	for _, accountId := range accountIds {
		queryParams = append(queryParams, utils_http.QueryParam{
			Key:   AccountId,
			Value: accountId,
		})
	}
	_, body, err := utils_http.PerformRequest(
		"JiraUserAdapter",
		path.String(),
		http.MethodGet,
		headers,
		queryParams,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not perform request: %v", err)
		return nil, err
	}
	users, err := usersFromJson(body)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not parse response body: %v", err)
	}
	return users, nil
}
