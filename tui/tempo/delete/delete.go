package tempo_worklogdelete

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	jira "github.com/remshams/jira-control/jira/public"
	"github.com/remshams/jira-control/tui/common"
	tempo_workloglistmodel "github.com/remshams/jira-control/tui/tempo/model"
)

type WorklogDeleteKeymap struct {
	global common.GlobalKeyMap
	help   help.KeyMap
	yes    key.Binding
	no     key.Binding
}

func (m WorklogDeleteKeymap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.help.Help,
		m.yes,
		m.no,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m WorklogDeleteKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
	}
}

var WorklogDeleteKeys = WorklogDeleteKeymap{
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	yes: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "yes"),
	),
	no: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "No"),
	),
}

type Model struct {
	worklog jira.TempoWorklog
}

func New() Model {
	return Model{}
}

func (m *Model) Init(worklog jira.TempoWorklog) tea.Cmd {
	m.worklog = worklog
	return help.CreateSetKeyMapMsg(WorklogDeleteKeys)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, WorklogDeleteKeys.yes):
			error := m.worklog.Delete()
			if error != nil {
				cmd = toast.CreateErrorToastAction(fmt.Sprintf("Worklog with id %d could not be deleted", m.worklog.Id))
			}
			cmd = tea.Batch(toast.CreateSuccessToastAction("Worklog deleted"), tempo_workloglistmodel.CreateSwitchWorklogListView)
		case key.Matches(msg, WorklogDeleteKeys.no):
			cmd = tea.Batch(tempo_workloglistmodel.CreateSwitchWorklogListView, help.CreateResetKeyMapMsg())
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s%s",
		styles.TextWarningColor.Render("Delete worklog with id: "),
		styles.TextInfoColor.Render(strconv.Itoa(m.worklog.Id)),
	)
}
