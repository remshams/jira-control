package jira_common_http

import utils_http "github.com/remshams/common/utils/http"

func CreateDefaultHttpHeaders(username string, apiToken string) []utils_http.HttpHeader {
	return []utils_http.HttpHeader{utils_http.CreateBasicAuthHeader(username, apiToken)}
}

func CreateDefaultTempoHttpHeaders(apiToken string) []utils_http.HttpHeader {
	return []utils_http.HttpHeader{utils_http.CreateBearerTokenHeader(apiToken)}
}
