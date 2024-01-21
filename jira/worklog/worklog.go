package issue_worklog

type Worklog struct {
	issueKey   string
	hoursSpent float64
}

func NewWorklog(issueKey string, hoursSpent float64) Worklog {
	return Worklog{
		issueKey:   issueKey,
		hoursSpent: hoursSpent,
	}
}
