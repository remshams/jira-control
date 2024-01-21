package issue_worklog

import "github.com/charmbracelet/log"

type WorklogJiraAdapter struct {
}

func (w WorklogJiraAdapter) logWork(worklog Worklog) error {
	log.Debugf("WorklogJiraAdapter: Logging work %v", worklog)
	return nil
}
