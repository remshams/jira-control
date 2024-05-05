package issue_worklog

import (
	"sort"
	"time"
)

type WorklogList []Worklog

func (w WorklogList) SortByStart(descending bool) WorklogList {
	sort.Slice(w, func(i, j int) bool {
		if descending {
			return w[i].Start.After(w[j].Start)
		}
		return w[i].Start.Before(w[j].Start)
	})
	return w
}

type WorklogListQuery struct {
	issueKey       string
	startedAfter   time.Time
	startedBefore  time.Time
	sortDescending bool
}

func NewWorklogListQuery(issueKey string) WorklogListQuery {
	return WorklogListQuery{
		issueKey:      issueKey,
		startedAfter:  time.Now(),
		startedBefore: time.Now(),
	}
}

func (w WorklogListQuery) WithIssueKey(issueKey string) WorklogListQuery {
	w.issueKey = issueKey
	return w
}

func (w WorklogListQuery) WithStartedAfter(time time.Time) WorklogListQuery {
	w.startedAfter = time
	return w
}

func (w WorklogListQuery) WithStartedBefore(time time.Time) WorklogListQuery {
	w.startedBefore = time
	return w
}

func (w WorklogListQuery) WithSortDescending(descending bool) WorklogListQuery {
	w.sortDescending = descending
	return w
}

type WorklogAdapter interface {
	List(query WorklogListQuery) (WorklogList, error)
	logWork(worklog Worklog) error
}

type Worklog struct {
	adapter            WorklogAdapter
	issueKey           string
	Id                 string
	TimeSpentInSeconds int
	HoursSpent         float64
	Start              time.Time
	Description        string
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
