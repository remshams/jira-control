package issue_worklog

import (
	"github.com/charmbracelet/log"
)

type WorklogMockAdapter struct{}

func NewWorklogMockAdapter() WorklogMockAdapter {
	return WorklogMockAdapter{}
}

func (w WorklogMockAdapter) logWork(worklog Worklog) error {
	log.Debugf("WorklogMockAdapter: Saving worklog %v", worklog)
	return nil
}

func (w WorklogMockAdapter) List(query WorklogListQuery) (WorklogList, error) {
	log.Debugf("WorklogMockAdapter: Listing worklogs for query %v", query)
	return []Worklog{}, nil
}

func (w WorklogMockAdapter) DeleteWorklog(worklog Worklog) error {
	log.Debugf("WorklogMockAdapter: Delete worklog: %s", worklog.Id)
	return nil
}
