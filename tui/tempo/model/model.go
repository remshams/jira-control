package tempo_workloglistmodel

import tea "github.com/charmbracelet/bubbletea"

type SwitchToWorklogListView struct{}

func CreateSwitchToWorklogListView() tea.Msg {
	return SwitchToWorklogListView{}
}
