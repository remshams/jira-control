package issue_worklog

import "time"

type WorklogAdapter interface {
	logWork(worklog Worklog) error
}

type Worklog struct {
	adapter    WorklogAdapter
	issueKey   string
	hoursSpent float64
	start      time.Time
}

func NewWorklog(adapter WorklogAdapter, issueKey string, hoursSpent float64) Worklog {
	return Worklog{
		adapter:    adapter,
		issueKey:   issueKey,
		hoursSpent: hoursSpent,
		start:      time.Now(),
	}
}

func (worklog *Worklog) withStart(start time.Time) *Worklog {
	worklog.start = start
	return worklog
}

func (w Worklog) Log() error {
	return (w.adapter).logWork(w)
}
