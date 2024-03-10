package issue

import (
	"fmt"
	"strings"
	"testing"

	"github.com/remshams/jira-control/jira/utils"
	"github.com/stretchr/testify/assert"
)

const summary = "summary"

var summaryJql = fmt.Sprintf("summary ~ \"%s\"", summary)

func TestJqlFromSearchRequest_Summary(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.summary = summary
	expected := summaryJql

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Key(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.key = "key"
	expected := fmt.Sprintf("key = \"%s\"", request.key)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Project(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request.projectName = "project"
	expected := fmt.Sprintf("project = \"%s\"", request.projectName)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_Combined(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
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
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request = request.WithUpdatedBy(updatedBy)
	expected := fmt.Sprintf("issueKey IN updatedBy(\"%s\")", updatedBy)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_OrderBy(t *testing.T) {
	orderBy := utils.OrderBy{Fields: []string{"summary"}, Sorting: utils.SortingDesc}
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request = request.WithSummary(summary)
	request = request.WithOrderBy(orderBy)
	expected := fmt.Sprintf("%s ORDER BY %s %s", summaryJql, strings.Join(orderBy.Fields, ","), orderBy.Sorting)

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_OrderByNoFields(t *testing.T) {
	orderBy := utils.OrderBy{Sorting: utils.SortingDesc}
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	request = request.WithSummary(summary)
	request = request.WithOrderBy(orderBy)
	expected := summaryJql

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}

func TestJqlFromSearchRequest_NoFieldsSelected(t *testing.T) {
	request := NewIssueSearchRequest(NewMockIssueAdapter())
	expected := ""

	assert.Equal(t, expected, jqlFromSearchRequest(request))
}
