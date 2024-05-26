package tempo_worklog

import (
	"slices"
	"time"
)

type WorklogListQuery struct {
	adapter        WorklogListAdapter
	from           time.Time
	to             time.Time
	sortDescending bool
}

func NewWorkloglistQuery(adapter WorklogListAdapter) WorklogListQuery {
	now := time.Now()
	return WorklogListQuery{
		adapter:        adapter,
		from:           time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()),
		to:             time.Date(now.Year(), now.Month(), now.Day(), 24, 59, 59, 0, now.Location()),
		sortDescending: false,
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

func (w WorklogListQuery) WithSortDescending() WorklogListQuery {
	w.sortDescending = true
	return w
}

func (w WorklogListQuery) Search() ([]Worklog, error) {
	worklogs, err := w.adapter.List(w)
	if err != nil {
		return []Worklog{}, err
	}
	if w.sortDescending {
		slices.SortFunc(worklogs, func(base, compare Worklog) int {
			if base.Start.After(compare.Start) {
				return -1
			} else {
				return 1
			}
		})
	}
	return worklogs, nil
}

type WorklogListAdapter interface {
	List(query WorklogListQuery) ([]Worklog, error)
	Delete(id int) error
}

type Worklog struct {
	IssueKey         int
	Id               int
	TimeSpentSeconds int
	BillableSeconds  int
	Start            time.Time
	Description      string
}

func NewWorklog(issueKey int, id int, timeSpentSeconds int, billableSeconds int, start time.Time, description string) Worklog {
	return Worklog{
		IssueKey:         issueKey,
		Id:               id,
		TimeSpentSeconds: timeSpentSeconds,
		BillableSeconds:  billableSeconds,
		Start:            start,
		Description:      description,
	}
}
