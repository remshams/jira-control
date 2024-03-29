package favorite_list

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
	common_worklog "github.com/remshams/jira-control/tui/_common/worklog"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type LoadFavoritesSuccessAction struct {
	Favorites []jira.Favorite
}

type LoadFavoritesErrorAction struct{}

func createLoadFavoritesAction(adapter tui_jira.JiraAdapter) tea.Cmd {
	return func() tea.Msg {
		favoritesChan := make(chan []jira.Favorite)
		errorChan := make(chan error)
		go loadFavorites(adapter, favoritesChan, errorChan)
		select {
		case favorites := <-favoritesChan:
			return LoadFavoritesSuccessAction{Favorites: favorites}
		case <-errorChan:
			return LoadFavoritesErrorAction{}
		}
	}
}

func loadFavorites(adapter tui_jira.JiraAdapter, favoritesChan chan []jira.Favorite, errorChan chan error) {
	favorites, err := adapter.App.FavoriteAdapter.Load()
	if err != nil {
		errorChan <- err
	} else {
		favoritesChan <- favorites
	}
}

type FavoritesKeymap struct {
	global  common.GlobalKeyMap
	logWork key.Binding
	table   table.KeyMap
	help    help.KeyMap
}

func (m FavoritesKeymap) ShortHelp() []key.Binding {
	keyBindings := []key.Binding{
		m.logWork,
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
	logWork: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Log work"),
	),
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
	adapter   tui_jira.JiraAdapter
	table     table.Model[[]jira.Favorite]
	spinner   spinner.Model
	state     utils.ViewState
	favorites []jira.Favorite
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
		table: table.
			New(createTableColumns, createTableRows, 5, 10).
			WithNotDataMessage("No favorites"),
		spinner:   spinner.New().WithLabel("Loading favorites..."),
		state:     favoritesStateLoading,
		favorites: []jira.Favorite{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Favorites"),
		help.CreateSetKeyMapMsg(FavoritesKeys),
		m.spinner.Tick(),
		createLoadFavoritesAction(m.adapter),
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
	switch msg := msg.(type) {
	case LoadFavoritesSuccessAction:
		m.state = favoritesStateLoaded
		m.favorites = msg.Favorites
		cmd = table.CreateTableDataUpdatedAction(m.favorites)
	case LoadFavoritesErrorAction:
		m.state = favoritesStateError
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return cmd
}

func (m *Model) processLoadedUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, FavoritesKeys.logWork):
			favorite := m.findFavorite(m.table.SelectedRowCell(0))
			if favorite == nil {
				cmd = toast.CreateErrorToastAction("Selected favorite could not be found")
			} else {
				cmd = common_worklog.CreateLogWorkAction(favorite.IssueKey)
			}
		default:
			m.table, cmd = m.table.Update(msg)
		}
	default:
		m.table, cmd = m.table.Update(msg)
	}
	return cmd
}

func (m Model) findFavorite(id string) *jira.Favorite {
	for _, favorite := range m.favorites {
		if favorite.Id.String() == id {
			return &favorite
		}
	}
	return nil
}

func (m Model) View() string {
	switch m.state {
	case favoritesStateLoading:
		return m.spinner.View()
	case favoritesStateLoaded:
		return m.table.View()
	case favoritesStateError:
		return "Error loading favorites"
	default:
		return ""
	}
}

func createTableColumns(tableWidth int) []table.Column {
	return []table.Column{
		{Title: "Id", Width: styles.CalculateDimensionsFromPercentage(5, tableWidth, 5)},
		{Title: "Key", Width: styles.CalculateDimensionsFromPercentage(35, tableWidth, 20)},
		{Title: "Time Spent", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Last Updated At", Width: styles.CalculateDimensionsFromPercentage(25, tableWidth, 20)},
		{Title: "Created At", Width: styles.CalculateDimensionsFromPercentage(25, tableWidth, 20)},
	}
}

func createTableRows(favorits []jira.Favorite) []table.Row {
	timeFormat := "2006-01-02 15:04"
	rows := []table.Row{}
	for _, favorite := range favorits {
		rows = append(rows, table.Row{
			favorite.Id.String(),
			favorite.IssueKey,
			fmt.Sprintf("%.1f h", favorite.HoursSpent),
			favorite.LastUsedAt.Format(timeFormat),
			favorite.CreatedAt.Format(timeFormat),
		})
	}
	return rows
}
