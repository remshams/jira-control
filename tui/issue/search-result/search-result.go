package issue_search_result

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	table_utils "github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/jira-control/jira/issue"
	common "github.com/remshams/jira-control/tui/_common"
)

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
	global     common.GlobalKeyMap
	help       help.KeyMap
	table      table.KeyMap
	SwitchView key.Binding
}

func (m SearchResultKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.SwitchView,
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
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap(),
	SwitchView: key.NewBinding(
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
			table.WithRows([]table.Row{}),
			table.WithColumns(createTableColumns()),
			table.WithKeyMap(table.DefaultKeyMap()),
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SearchResultKeys.SwitchView):
			cmd = help.CreateSetKeyMapMsg(SearchResultKeys)
		case key.Matches(msg, SearchResultKeys.help.Help):
			cmd = help.CreateToggleFullHelpMsg()
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) createTable(columns []table.Column, rows []table.Row) table.Model {
	return table_utils.CreateTable(columns, rows)
}

func createInitialTable() table.Model {
	return table_utils.CreateTable(createTableColumns(), []table.Row{})
}

func createTableColumns() []table.Column {
	return []table.Column{
		{Title: "Key", Width: 10},
		{Title: "Summary", Width: 40},
		{Title: "ProjectName", Width: 40},
		{Title: "ProjectKey", Width: 20},
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
