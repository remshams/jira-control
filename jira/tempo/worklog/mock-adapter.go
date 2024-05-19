package tempo_worklog

import "time"

type MockWorklogAdapter struct {
}

func NewMockWorklogAdapter() MockWorklogAdapter {
	return MockWorklogAdapter{}
}

func (w MockWorklogAdapter) List(query WorklogListQuery) ([]Worklog, error) {
	return []Worklog{
		NewWorklog(0, 0, 3600, 1, time.Now(), "Description"),
		NewWorklog(1, 1, 7200, 2, time.Now(), "Description"),
	}, nil
}
