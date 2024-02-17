package issue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJqlFromSearchRequest_Summary(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.Summary = "summary"
	expected := fmt.Sprintf("summary ~ \"%s\"", request.Summary)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Key(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.Key = "key"
	expected := fmt.Sprintf("key = \"%s\"", request.Key)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Project(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.ProjectName = "project"
	expected := fmt.Sprintf("project = \"%s\"", request.ProjectName)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Combined(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.Summary = "summary"
	request.Key = "key"
	request.ProjectName = "project"
	expected := fmt.Sprintf(
		"summary ~ \"%s\" OR key = \"%s\" OR project = \"%s\"",
		request.Summary, request.Key, request.ProjectName,
	)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}
