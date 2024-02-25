package worklog_list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
)

type GoBackAction struct{}

func CreateGoBackAction() tea.Msg {
	return GoBackAction{}
}

type WorklogListKeyMap struct {
	global common.GlobalKeyMap
	goBack key.Binding
}

func (m WorklogListKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.goBack,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m WorklogListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var WorklogListKeys = WorklogListKeyMap{
	global: common.GlobalKeys,
	goBack: key.NewBinding(
		key.WithKeys("b"),
		key.WithHelp("b", "Go back"),
	),
}

type Model struct {
	issueKey string
	worklogs []jira.Worklog
}

func New(issueKey string, worklogs []jira.Worklog) Model {
	return Model{
		issueKey: issueKey,
		worklogs: worklogs,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg(fmt.Sprintf("Worklog List for %s", m.issueKey)),
		help.CreateSetKeyMapMsg(WorklogListKeys),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, WorklogListKeys.goBack):
			cmd = CreateGoBackAction
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return ""
}
