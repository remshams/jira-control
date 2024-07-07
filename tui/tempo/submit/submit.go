package tempo_submit

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/table"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type SwitchToWorklogListView struct{}

func createSwitchToWorklogListView() tea.Msg {
	return SwitchToWorklogListView{}
}

type SubmitKeymap struct {
	global      common.GlobalKeyMap
	help        help.KeyMap
	table       table.KeyMap
	worklogList key.Binding
}

func (m SubmitKeymap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.worklogList,
		m.help.Help,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m SubmitKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table.DefaultKeyBindings,
	}
}

var SubmitKeys = SubmitKeymap{
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap,
	worklogList: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Show worklog list"),
	),
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
	return tea.Batch(
		title.CreateSetPageTitleMsg("Submit timesheet"),
		help.CreateSetKeyMapMsg(SubmitKeys),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	cmd = m.processLoadedUpdate(msg)
	return m, cmd
}

func (m *Model) processLoadedUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SubmitKeys.worklogList):
			cmd = createSwitchToWorklogListView
		}
	}
	return cmd
}

func (m Model) View() string {
	return "Submit timesheet"
}
