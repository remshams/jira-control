package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	"github.com/remshams/jira-control/jira/issue"
	jira "github.com/remshams/jira-control/jira/public"
	"github.com/remshams/jira-control/jira/utils"
)

func main() {
	logger.PrepareLogger()
	app, err := jira.PrepareApplication()
	if err != nil {
		log.Errorf("Could not create JiraAdapter: %v", err)
		os.Exit(1)
	}
	// worklog := issue_worklog.NewWorklog(app.IssueWorklogAdapter, "NC-40", 4.5)
	// worklog.Log()
	issueSearchRequest := issue.NewIssueSearchRequest(app.IssueAdapter)
	issueSearchRequest = issueSearchRequest.WithUpdatedBy(app.Username())
	issueSearchRequest = issueSearchRequest.WithOrderBy(utils.NewOrderBy([]string{"updated"}, utils.SortingDesc))
	issues, err := issueSearchRequest.Search()
	fmt.Println(issues[0].Project.Updated)

}
