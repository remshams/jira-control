package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
	"github.com/remshams/jira-control/jira/issue"
	jira "github.com/remshams/jira-control/jira/public"
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
	issueSearchRequest.Summary = "Project Management"
	issues, err := issueSearchRequest.Search()
	var issue issue.Issue
	for _, currentIssue := range issues {
		if currentIssue.Key == "NXNXTUI-45" {
			issue = currentIssue
			break
		}
	}
	startedAfter := time.Now()
	startedAfter = startedAfter.Add(-4 * 24 * time.Hour)
	query := issue.WorklogsQuery().WithStartedAfter(startedAfter)
	worklogs, _ := issue.Worklogs(query)
	fmt.Println(len(worklogs))

}
