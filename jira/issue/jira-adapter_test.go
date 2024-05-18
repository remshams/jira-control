package issue

import (
	"fmt"
	"strings"
	"testing"

	issue_worklog "github.com/remshams/jira-control/jira/issue/worklog"
	"github.com/remshams/jira-control/jira/utils"
	"github.com/stretchr/testify/assert"
)

const summary = "summary"

var summaryJql = fmt.Sprintf("summary ~ \"%s\"", summary)

func TestJqlFromSearchRequest_Summary(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	request.summary = summary
	expected := summaryJql

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Key(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	request.key = "key"
	expected := fmt.Sprintf("key = \"%s\"", request.key)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Project(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	request.projectName = "project"
	expected := fmt.Sprintf("project = \"%s\"", request.projectName)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Combined(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	request.summary = summary
	request.key = "key"
	request.projectName = "project"
	expected := fmt.Sprintf(
		"%s OR key = \"%s\" OR project = \"%s\"",
		summaryJql, request.key, request.projectName,
	)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_UpdatedBy(t *testing.T) {
	updatedBy := "updatedBy"
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	request = request.WithUpdatedBy(updatedBy)
	expected := fmt.Sprintf("issueKey IN updatedBy(\"%s\")", updatedBy)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_OrderBy(t *testing.T) {
	orderBy := utils.NewOrderBy([]string{"updated"}, utils.SortingDesc)
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	request = request.WithSummary(summary)
	request = request.WithOrderBy(orderBy)
	expected := fmt.Sprintf("%s ORDER BY %s %s", summaryJql, strings.Join(orderBy.Fields, ","), orderBy.Sorting)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_NoFieldsSelected(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter(issue_worklog.NewWorklogMockAdapter()))
	expected := ""

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}
