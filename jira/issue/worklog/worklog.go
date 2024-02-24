package issue_worklog

import "time"

type WorklogListQuery struct {
	issueKey      string
	startedAfter  time.Time
	startedBefore time.Time
}

func NewWorklogListQuery(issueKey string) WorklogListQuery {
	return WorklogListQuery{
		issueKey:      issueKey,
		startedAfter:  time.Now(),
		startedBefore: time.Now(),
	}
}

func (w WorklogListQuery) withstartedAfter(time time.Time) WorklogListQuery {
	w.startedAfter = time
	return w
}

func (w WorklogListQuery) withstartedBefore(time time.Time) WorklogListQuery {
	w.startedBefore = time
	return w
}

type WorklogAdapter interface {
	list(issueKey string, query WorklogListQuery) ([]Worklog, error)
	logWork(worklog Worklog) error
}

type Worklog struct {
	adapter     WorklogAdapter
	issueKey    string
	HoursSpent  float64
	Start       time.Time
	Description string
}

func NewWorklog(adapter WorklogAdapter, issueKey string, hoursSpent float64) Worklog {
	return Worklog{
		adapter:     adapter,
		issueKey:    issueKey,
		HoursSpent:  hoursSpent,
		Start:       time.Now(),
		Description: "",
	}
}

func (w Worklog) Log() error {
	return (w.adapter).logWork(w)
}
