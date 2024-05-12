package tempo_worklog

import "time"

type WorklogListQuery struct {
	from           time.Time
	to             time.Time
	sortDescending bool
}

func NewWorklistQuery() WorklogListQuery {
	return WorklogListQuery{
		from: time.Now(),
		to:   time.Now().Add(time.Duration(time.Hour * 24)),
	}
}

func (w WorklogListQuery) WithFrom(from time.Time) WorklogListQuery {
	w.from = from
	return w
}

func (w WorklogListQuery) WithTo(to time.Time) WorklogListQuery {
	w.to = to
	return w
}

func (w WorklogListQuery) WithSortDescending(sortDescending bool) WorklogListQuery {
	w.sortDescending = sortDescending
	return w
}

type WorklistAdapter interface {
	List(query WorklogListQuery) ([]Worklog, error)
}

type Worklog struct {
	IssueKey           string
	Id                 string
	TimeSpentInSeconds int
	HoursSpent         float64
}

func NewWorklog(issueKey string, id string, timeSpentInSeconds int, hoursSpent float64) Worklog {
	return Worklog{
		IssueKey:           issueKey,
		Id:                 id,
		TimeSpentInSeconds: timeSpentInSeconds,
		HoursSpent:         hoursSpent,
	}
}
