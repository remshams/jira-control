package tui_last_updated

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type LastUpdatedKeymap struct {
	global common.GlobalKeyMap
}

func (m LastUpdatedKeymap) ShortHelp() []key.Binding {
	return m.global.KeyBindings()
}

func (m LastUpdatedKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

var LastUpdatedKeys = LastUpdatedKeymap{
	global: common.GlobalKeys,
}

type Model struct {
	adapter tui_jira.JiraAdapter
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(title.CreateSetPageTitleMsg("Last Updated Issues"), help.CreateSetKeyMapMsg(LastUpdatedKeys))
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return "Last Updated Issues View"
}
