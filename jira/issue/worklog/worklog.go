package issue_worklog

import "time"

type WorklogAdapter interface {
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
