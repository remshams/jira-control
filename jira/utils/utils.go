package utils

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

type Sorting string

const (
	SortingAsc  Sorting = "ASC"
	SortingDesc Sorting = "DESC"
)

type OrderBy struct {
	Fields  []string
	Sorting Sorting
}

func NewOrderBy(fields []string, sorting Sorting) OrderBy {
	return OrderBy{
		Fields:  fields,
		Sorting: sorting,
	}
}

func JiraDateToTime(timeString string) (time.Time, error) {
	t, err := time.Parse(
		"2006-01-02T15:04:05.999-0700",
		timeString,
	)
	if err != nil {
		log.Errorf("JiraDateToTime: Could not parse time: %v", err)
		return time.Time{}, err
	}
	return t, nil
}

func TimeToJiraDate(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.999-0700")
}

func TempoDateToTime(tempoDate string, tempoTime string) (time.Time, error) {
	return time.Parse(
		"2006-01-02T15:04:05Z",
		fmt.Sprintf("%sT%sZ", tempoDate, tempoTime),
	)
}

func TimeToTempoDate(t time.Time) string {
	return t.Format("2006-01-02")
}
