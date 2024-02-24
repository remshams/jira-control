package issue_worklog

import (
	"github.com/charmbracelet/log"
)

type WorklogMockAdatpter struct{}

func NewWorklogMockAdapter() WorklogMockAdatpter {
	return WorklogMockAdatpter{}
}

func (w WorklogMockAdatpter) logWork(worklog Worklog) error {
	log.Debugf("WorklogMockAdapter: Saving worklog %v", worklog)
	return nil
}

func (w WorklogMockAdatpter) list(issueKey string, query WorklogListQuery) ([]Worklog, error) {
	log.Debugf("WorklogMockAdapter: Listing worklogs for issue %s", issueKey)
	return []Worklog{}, nil
}
