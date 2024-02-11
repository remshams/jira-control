package issue_search_result

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	table_utils "github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/jira-control/jira/issue"
)

type SearchResultKeyMap struct {
	table table.KeyMap
}

func (m SearchResultKeyMap) ShortHelp() []key.Binding {
	return table_utils.TableKeyBindings()
}

func (m SearchResultKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var SearchResultKeys = SearchResultKeyMap{
	table: table.DefaultKeyMap(),
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
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
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
