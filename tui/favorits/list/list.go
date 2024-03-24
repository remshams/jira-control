package favorite_list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type FavoritesKeymap struct {
	global common.GlobalKeyMap
	table  table.KeyMap
	help   help.KeyMap
}

func (m FavoritesKeymap) ShortHelp() []key.Binding {
	keyBindings := []key.Binding{
		m.help.Help,
	}
	return append(keyBindings, m.global.KeyBindings()...)
}

func (m FavoritesKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table.DefaultKeyBindings,
	}
}

var FavoritesKeys = FavoritesKeymap{
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap,
	global: common.GlobalKeys,
}

const (
	favoritesStateLoaded  utils.ViewState = "favoritesStateLoaded"
	favoritesStateLoading utils.ViewState = "favoritesStateLoading"
	favoritesStateError   utils.ViewState = "favoritesStateError"
)

type Model struct {
	adapter tui_jira.JiraAdapter
	table   table.Model[[]jira.Favorite]
	spinner spinner.Model
	state   utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
		table: table.
			New(createTableColumns, createTableRows, 5, 10).
			WithNotDataMessage("No favorites"),
		spinner: spinner.New().WithLabel("Loading favorites..."),
		state:   favoritesStateLoading,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Favorites"),
		help.CreateSetKeyMapMsg(FavoritesKeys),
		m.spinner.Tick(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case favoritesStateLoading:
		cmd = m.processLoadingUpdate(msg)
	case favoritesStateLoaded:
		cmd = m.processLoadedUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	return cmd
}

func (m *Model) processLoadedUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case favoritesStateLoading:
		return m.spinner.View()
	case favoritesStateLoaded:
		return m.table.View()
	default:
		return ""
	}
}

func createTableColumns(tableWidth int) []table.Column {
	return []table.Column{
		{Title: "Key", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Time Spent", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Last Updated At", Width: styles.CalculateDimensionsFromPercentage(40, tableWidth, 20)},
		{Title: "Created At", Width: styles.CalculateDimensionsFromPercentage(40, tableWidth, 20)},
	}
}

func createTableRows(favorits []jira.Favorite) []table.Row {
	timeFormat := "2006-01-02 15:04"
	rows := []table.Row{}
	for _, favorite := range favorits {
		rows = append(rows, table.Row{
			favorite.IssueKey,
			fmt.Sprintf("%d h", int(favorite.HoursSpent)),
			favorite.LastUsedAt.Format(timeFormat),
			favorite.CreatedAt.Format(timeFormat),
		})
	}
	return rows
}
