package common

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/tabs"
	jira "github.com/remshams/jira-control/jira/public"
)

type GlobalKeyMap struct {
	Tab  tabs.TabKeyMap
	Quit key.Binding
}

var GlobalKeys = GlobalKeyMap{
	Tab: tabs.TabKeys,
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/C-c", "Quit"),
	),
}

func (k GlobalKeyMap) KeyBindings() []key.Binding {
	return []key.Binding{
		k.Tab.Tab,
		k.Quit,
	}
}

type LogWorkAction struct {
	Issue jira.Issue
}

func CreateLogWorkAction(issue jira.Issue) tea.Cmd {
	return func() tea.Msg {
		return LogWorkAction{
			Issue: issue,
		}
	}
}
