package tempo_workloglistmodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SwitchToWorklogListView struct {
	Toast *tea.Cmd
}

func CreateSwitchToWorklogListView(toast *tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return SwitchToWorklogListView{
			Toast: toast,
		}
	}
}

type LoadWorklogListAction struct{}

func CreateLoadWorklogListAction() tea.Msg {
	return LoadWorklogListAction{}
}
