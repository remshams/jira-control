package tempo_worklog

import "time"

type WorklogListQuery struct {
	adapter        WorklogListAdapter
	from           time.Time
	to             time.Time
	sortDescending bool
}

func NewWorkloglistQuery(adapter WorklogListAdapter) WorklogListQuery {
	now := time.Now()
	return WorklogListQuery{
		adapter: adapter,
		from:    time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
		to:      time.Date(now.Year(), now.Month(), now.Day(), 24, 59, 59, 0, now.Location()),
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

func (w WorklogListQuery) Search() ([]Worklog, error) {
	return w.adapter.List(w)
}

type WorklogListAdapter interface {
	List(query WorklogListQuery) ([]Worklog, error)
}

type Worklog struct {
	IssueKey           int
	Id                 int
	TimeSpentInSeconds int
	BillableSeconds    int
	Start              time.Time
}

func NewWorklog(issueKey int, id int, timeSpentInSeconds int, billableSeconds int, start time.Time) Worklog {
	return Worklog{
		IssueKey:           issueKey,
		Id:                 id,
		TimeSpentInSeconds: timeSpentInSeconds,
		BillableSeconds:    billableSeconds,
		Start:              start,
	}
}
