package tempo_worklog

type MockWorklogAdapter struct {
}

func List(query WorklogListQuery) ([]Worklog, error) {
	return []Worklog{
		NewWorklog("issue-0", "0", 3600, 1),
		NewWorklog("issue-1", "1", 7200, 2),
	}, nil
}
