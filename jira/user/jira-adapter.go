package user

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/charmbracelet/log"
	utils_http "github.com/remshams/common/utils/http"
)

const path = "rest/api/3/user"
const myselfPath = "rest/api/3/myself"

type userDto struct {
	AccountId string `json:"accountId"`
	Email     string `json:"emailAddress"`
	Name      string `json:"displayName"`
}

func fromJson(body []byte) (User, error) {
	var userDto userDto
	err := json.Unmarshal(body, &userDto)
	if err != nil {
		return User{}, err
	}
	return NewUser(userDto.AccountId, userDto.Name, userDto.Email), err

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
	headers := []utils_http.HttpHeader{
		utils_http.CreateBasicAuthHeader(jiraUserAdapter.username, jiraUserAdapter.apiToken),
	}
	_, body, err := utils_http.PerformRequest(
		"JiraUserAdapter",
		path.String(),
		http.MethodGet,
		headers,
		nil,
		nil,
		nil,
	)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not perform request: %v", err)
		return User{}, err
	}
	user, err := fromJson(body)
	if err != nil {
		log.Errorf("JiraUserAdapter: Could not parse response body %v", err)
		return User{}, err
	}
	return user, nil

}
