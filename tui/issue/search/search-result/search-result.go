package issue_search_result

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	table_utils "github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/jira-control/jira/issue"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
	app_store "github.com/remshams/jira-control/tui/store"
)

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
		table_utils.TableKeyBindings(),
	}
}

var SearchResultKeys = SearchResultKeyMap{
	logWork: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Log work"),
	),
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap(),
	switchView: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Switch to search"),
	),
}

type Model struct {
	issues []issue.Issue
	table  table.Model
}

func New() Model {
	m := Model{
		issues: []issue.Issue{},
		table: table.New(
			table.WithKeyMap(table.DefaultKeyMap()),
			table.WithFocused(true),
		),
	}
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
	case tea.WindowSizeMsg:
		m.recalculateTableLayout()
	case SetSearchResultAction:
		m.issues = msg.issues
		m.recalculateTableLayout()
		m.table.GotoTop()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SearchResultKeys.switchView):
			cmd = tea.Batch(CreateSwitchViewAction(), help.CreateSetKeyMapMsg(SearchResultKeys))
		case key.Matches(msg, SearchResultKeys.help.Help):
			cmd = help.CreateToggleFullHelpMsg()
		case key.Matches(msg, SearchResultKeys.logWork):
			issue := m.findIssue(m.table.SelectedRow()[0])
			if issue == nil {
				cmd = toast.CreateErrorToastAction("Selected issue could not be found")
			}
			cmd = CreateLogWorkAction(*issue)
		default:
			m.table, cmd = m.table.Update(msg)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	if len(m.table.Rows()) > 0 {
		return m.table.View()
	} else {
		style := lipgloss.NewStyle().
			Foreground(styles.SelectedColor).
			Width(app_store.LayoutStore.Width).
			Align(lipgloss.Center)
		return style.Render("No issues")
	}
}

func (m Model) createTable(columns []table.Column, rows []table.Row) table.Model {
	return table_utils.CreateTable(columns, rows)
}

func (m Model) createTableColumns() []table.Column {
	tableWidth, _ := m.calculateTableDimensions()
	return []table.Column{
		{Title: "Key", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth)},
		{Title: "Summary", Width: styles.CalculateDimensionsFromPercentage(50, tableWidth)},
		{Title: "ProjectName", Width: styles.CalculateDimensionsFromPercentage(30, tableWidth)},
		{Title: "ProjectKey", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth)},
	}
}

func (m Model) createTableRows() []table.Row {
	rows := []table.Row{}

	for _, issue := range m.issues {
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

func (m Model) calculateTableDimensions() (int, int) {
	width := app_store.LayoutStore.Width - 5
	height := app_store.LayoutStore.Height - 11
	if height < 0 {
		height = styles.CalculateDimensionsFromPercentage(80, app_store.LayoutStore.Height)
	}
	return width, height
}

func (m *Model) recalculateTableLayout() {
	width, height := m.calculateTableDimensions()
	log.Debugf("width: %d, height: %d", width, height)
	m.table.SetWidth(width)
	m.table.SetHeight(height)
	m.table.SetColumns(m.createTableColumns())
	m.table.SetRows(m.createTableRows())
}
