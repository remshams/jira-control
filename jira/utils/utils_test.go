package utils

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJiraDateToTime_Valid(t *testing.T) {
	year := 2024
	month := time.February
	day := 19
	hour := 19
	minute := 3
	second := 50
	nanosecond := 0
	location, _ := time.LoadLocation(time.Local.String())
	expected := time.Date(year, month, day, hour, minute, second, nanosecond, location)
	convertedTime, err := JiraDateToTime(
		fmt.Sprintf(
			"%04d-%02d-%02dT%02d:%02d:%02d.%03d+0100",
			year,
			month,
			day,
			hour,
			minute,
			second,
			nanosecond,
		),
	)

	assert.NoError(t, err)
	assert.Equal(t, expected, convertedTime)
}

func TestJiraDateToTime_Invalid(t *testing.T) {
	_, err := JiraDateToTime("2024-02-1920:03:50.00+0100")

	assert.Error(t, err)
}

func TestTimeToJiraDate_Valid(t *testing.T) {
	year := 2024
	month := time.February
	day := 19
	hour := 19
	minute := 3
	second := 50
	nanosecond := 0
	location, _ := time.LoadLocation(time.Local.String())
	input := time.Date(year, month, day, hour, minute, second, nanosecond, location)
	expected := fmt.Sprintf(
		"%04d-%02d-%02dT%02d:%02d:%02d+0100",
		year,
		month,
		day,
		hour,
		minute,
		second,
	)
	convertedString := TimeToJiraDate(input)

	assert.Equal(t, expected, convertedString)
}

func TestTempoDateToTime(t *testing.T) {
	year := 2024
	month := time.May
	day := 12
	hour := 18
	minutes := 12
	seconds := 12
	tempoDate := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	tempoTime := fmt.Sprintf("%02d:%02d:%02d", hour, minutes, seconds)
	expected := time.Date(year, month, day, hour, minutes, seconds, 0, time.UTC)

	dateReturned, _ := TempoDateToTime(tempoDate, tempoTime)

	assert.Equal(t, expected, dateReturned)
}

func TestTempoDateToTimeInvalidDate(t *testing.T) {
	_, err := TempoDateToTime("2024-5-12", "18:12:12")

	assert.NotNil(t, err)
}
