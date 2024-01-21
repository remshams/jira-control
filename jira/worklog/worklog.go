package issue_worklog

type WorklogAdapter interface {
	logWork(worklog Worklog) error
}

type Worklog struct {
	adapter    WorklogAdapter
	issueKey   string
	hoursSpent float64
}

func NewWorklog(adapter WorklogAdapter, issueKey string, hoursSpent float64) Worklog {
	return Worklog{
		adapter:    adapter,
		issueKey:   issueKey,
		hoursSpent: hoursSpent,
	}
}

func (worklog Worklog) hours() float64 {
	return worklog.hoursSpent
}

func (w Worklog) Log() error {
	return (w.adapter).logWork(w)
}
