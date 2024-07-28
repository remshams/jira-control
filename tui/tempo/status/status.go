package tempo_status

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/styles"
	jira "github.com/remshams/jira-control/jira/public"
	tui_common_utils "github.com/remshams/jira-control/tui/common/utils"
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
		tui_common_utils.RenderKeyValue(
			"Required hours",
			fmt.Sprintf("%d hours", m.timesheetStatus.RequiredHours),
		),
		tui_common_utils.RenderKeyValue(
			"Spent hours",
			fmt.Sprintf("%s hours", spentHoursColor.Render(strconv.Itoa(m.timesheetStatus.SpentHours))),
		),
		tui_common_utils.RenderKeyValue(
			"Status",
			m.timesheetStatus.Status,
		),
	)
}
