package issue

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	common "github.com/remshams/jira-control/tui/_common"
)

type keyMap struct {
	global common.GlobalKeyMap
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.global.Tab.Tab,
		k.global.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var issueKeys = keyMap{
	global: common.GlobalKeys,
}

type Model struct {
}

func New() Model {
	return Model{}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Issue"),
		help.CreateSetKeyMapMsg(issueKeys),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return "Issue"
}
