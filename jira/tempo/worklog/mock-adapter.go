package tempo_worklog

import (
	"time"

	"github.com/charmbracelet/log"
)

type MockWorklogAdapter struct {
}

func NewMockWorklogAdapter() MockWorklogAdapter {
	return MockWorklogAdapter{}
}

func (w MockWorklogAdapter) List(query WorklogListQuery) ([]Worklog, error) {
	log.Debugf("Query worklog list %v", query)
	return []Worklog{
		NewWorklog(0, 0, 3600, 1, time.Now(), "Description"),
		NewWorklog(1, 1, 7200, 2, time.Now(), "Description"),
	}, nil
}

func (w MockWorklogAdapter) Delete(id int) error {
	log.Debugf("Delete worklog with id %d", id)
	return nil
}
