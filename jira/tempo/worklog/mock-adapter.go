package tempo_worklog

import "time"

type MockWorklogAdapter struct {
}

func (w MockWorklogAdapter) List(query WorklogListQuery) ([]Worklog, error) {
	return []Worklog{
		NewWorklog(0, 0, 3600, 1, time.Now()),
		NewWorklog(1, 1, 7200, 2, time.Now()),
	}, nil
}
