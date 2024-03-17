package issue_search_result

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/jira-control/jira/issue"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
)

type ShowWorklogsAction struct {
	Issue jira.Issue
}

func CreateShowWorklogsAction(issue jira.Issue) tea.Cmd {
	return func() tea.Msg {
		return ShowWorklogsAction{
			Issue: issue,
		}
	}
}

type SwitchViewAction struct {
}

func CreateSwitchViewAction() tea.Cmd {
	return func() tea.Msg {
		return SwitchViewAction{}
	}
}

type SetSearchResultAction struct {
	issues []issue.Issue
}

func CreateSearchResultAction(issues []issue.Issue) tea.Cmd {
	return func() tea.Msg {
		return SetSearchResultAction{
			issues: issues,
		}
	}
}

type SearchResultKeyMap struct {
	global       common.GlobalKeyMap
	help         help.KeyMap
	table        table.KeyMap
	showWorklogs key.Binding
	logWork      key.Binding
	switchView   key.Binding
}

func (m SearchResultKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.logWork,
		m.showWorklogs,
		m.switchView,
		m.help.Help,
		m.global.Tab.Tab,
		m.global.Quit,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m SearchResultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table.DefaultKeyBindings,
	}
}

var SearchResultKeys = SearchResultKeyMap{
	logWork: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Log work"),
	),
	showWorklogs: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "Show worklogs"),
	),
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap,
	switchView: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Switch to search"),
	),
}

type Model struct {
	issues []issue.Issue
	table  table.Model[[]issue.Issue]
}

func New() Model {
	m := Model{
		issues: []issue.Issue{},
	}
	m.table = table.
		New(createTableColumns, createTableRows, 5, 11).
		WithNotDataMessage("No issues")
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Issue Search Result"),
		help.CreateSetKeyMapMsg(SearchResultKeys),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case SetSearchResultAction:
		m.issues = msg.issues
		cmd = table.CreateTableDataUpdatedAction(m.issues)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SearchResultKeys.switchView):
			cmd = tea.Batch(CreateSwitchViewAction(), help.CreateSetKeyMapMsg(SearchResultKeys))
		case key.Matches(msg, SearchResultKeys.showWorklogs):
			issue := m.findIssue(m.table.SelectedRowCell(0))
			if issue == nil {
				cmd = toast.CreateErrorToastAction("Selected issue could not be found")
			}
			cmd = CreateShowWorklogsAction(*issue)
		case key.Matches(msg, SearchResultKeys.help.Help):
			cmd = help.CreateToggleFullHelpMsg()
		case key.Matches(msg, SearchResultKeys.logWork):
			issue := m.findIssue(m.table.SelectedRowCell(0))
			if issue == nil {
				cmd = toast.CreateErrorToastAction("Selected issue could not be found")
			}
			cmd = common.CreateLogWorkAction(*issue)
		default:
			m.table, cmd = m.table.Update(msg)
		}
	default:
		m.table, cmd = m.table.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func createTableColumns(width int) []table.Column {
	return []table.Column{
		{Title: "Key", Width: styles.CalculateDimensionsFromPercentage(10, width, styles.UnlimitedDimension)},
		{Title: "Summary", Width: styles.CalculateDimensionsFromPercentage(50, width, styles.UnlimitedDimension)},
		{Title: "ProjectName", Width: styles.CalculateDimensionsFromPercentage(30, width, styles.UnlimitedDimension)},
		{Title: "ProjectKey", Width: styles.CalculateDimensionsFromPercentage(10, width, styles.UnlimitedDimension)},
	}
}

func createTableRows(issues []issue.Issue) []table.Row {
	rows := []table.Row{}

	for _, issue := range issues {
		rows = append(rows, table.Row{
			issue.Key,
			issue.Summary,
			issue.Project.Name,
			issue.Project.Key,
		})
	}
	return rows
}

func (m Model) findIssue(key string) *issue.Issue {
	for _, issue := range m.issues {
		if issue.Key == key {
			return &issue
		}
	}
	return nil
}
