package utils

import (
	"time"

	"github.com/charmbracelet/log"
)

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
