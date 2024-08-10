package tempo_workloglistmodel

import tea "github.com/charmbracelet/bubbletea"

type SwitchWorklogListView struct{}

func CreateSwitchWorklogListView() tea.Msg {
	return SwitchWorklogListView{}
}
