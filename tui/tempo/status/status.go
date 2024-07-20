package tempo_status

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/styles"
	jira "github.com/remshams/jira-control/jira/public"
)

type Model struct {
	timesheetStatus jira.TimesheetStatus
}

func New() Model {
	return Model{}
}

func (m Model) Init(timeSheetStatus jira.TimesheetStatus) (Model, tea.Cmd) {
	m.timesheetStatus = timeSheetStatus
	return m, nil
}

func (m Model) View() string {
	spentHoursColor := styles.TextSuccessColor
	if m.timesheetStatus.RequiredHours-m.timesheetStatus.SpentHours > 0 {
		spentHoursColor = styles.TextErrorColor
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		renderKeyValue(
			"Required hours",
			fmt.Sprintf("%d hours", m.timesheetStatus.RequiredHours),
		),
		renderKeyValue(
			"Spent hours",
			fmt.Sprintf("%s hours", spentHoursColor.Render(strconv.Itoa(m.timesheetStatus.SpentHours))),
		),
		renderKeyValue(
			"Status",
			m.timesheetStatus.Status,
		),
	)
}

func renderKeyValue(key string, value string) string {
	return fmt.Sprintf("%s%s %s", styles.TextAccentColor.Render(key), styles.TextAccentColor.Render(":"), value)
}
