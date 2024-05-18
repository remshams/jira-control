package tempo_workloglist

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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

type initAction struct{}

func createInitAction() tea.Msg {
	return initAction{}
}

type LoadWorklogsSuccessAction struct {
	Worklogs []jira.TempoWorklog
}

type LoadWorklogsErrorAction struct{}

func createLoadWorklogsAction(adapter tui_jira.JiraAdapter) tea.Cmd {
	return func() tea.Msg {
		worklogsChan := make(chan []jira.TempoWorklog)
		errorChan := make(chan error)
		go loadWorklogs(adapter, worklogsChan, errorChan)
		select {
		case worklogs := <-worklogsChan:
			return LoadWorklogsSuccessAction{Worklogs: worklogs}
		case <-errorChan:
			return LoadWorklogsErrorAction{}
		}
	}
}

func loadWorklogs(adapter tui_jira.JiraAdapter, worklogsChan chan []jira.TempoWorklog, errorChan chan error) {
	worklogs, err := jira.NewTempoWorklogListQuery(adapter.App.TempoWorklogAdapter).Search()
	if err != nil {
		errorChan <- err
	} else {
		worklogsChan <- worklogs
	}
}

type WorklogListKeyMap struct {
	global common.GlobalKeyMap
	help   help.KeyMap
	table  table.KeyMap
}

func (m WorklogListKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.help.Help,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m WorklogListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table.DefaultKeyBindings,
	}
}

var WorklogListKeys = WorklogListKeyMap{
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap,
}

const (
	tempoWorklogStateLoaded  utils.ViewState = "tempoWorklogStateLoaded"
	tempoWorklogStateLoading utils.ViewState = "tempoWorklogStateLoading"
	tempoWorklogStateError   utils.ViewState = "tempoWorklogStateError"
)

type Model struct {
	adapter  tui_jira.JiraAdapter
	state    utils.ViewState
	worklogs []jira.TempoWorklog
	table    table.Model[[]jira.TempoWorklog]
	spinner  spinner.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	model := Model{
		adapter: adapter,
		state:   tempoWorklogStateLoading,
		spinner: spinner.New().WithLabel("Loading worklogs..."),
	}
	model.table = table.
		New(createTableColumns, createTableRows, 5, 10).
		WithNotDataMessage("No worklogs")
	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Tempo worklog list"),
		help.CreateSetKeyMapMsg(WorklogListKeys),
		m.spinner.Tick(),
		createInitAction,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case tempoWorklogStateLoading:
		cmd = m.processLoadingUpdate(msg)
	case tempoWorklogStateLoaded:
		cmd = m.processLoadedUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initAction:
		m.state = tempoWorklogStateLoading
		cmd = createLoadWorklogsAction(m.adapter)
	case LoadWorklogsSuccessAction:
		m.state = tempoWorklogStateLoaded
		m.worklogs = msg.Worklogs
		cmd = table.CreateTableDataUpdatedAction(m.worklogs)
	case LoadWorklogsErrorAction:
		m.state = tempoWorklogStateError
		m.worklogs = []jira.TempoWorklog{}
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return cmd
}

func (m *Model) processLoadedUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case tempoWorklogStateLoaded:
		return m.table.View()
	case tempoWorklogStateLoading:
		return m.spinner.View()
	case tempoWorklogStateError:
		return "Error loading worklogs"
	default:
		return ""
	}
}
func createTableColumns(tableWidth int) []table.Column {
	return []table.Column{
		{Title: "Id", Width: styles.CalculateDimensionsFromPercentage(20, tableWidth, 20)},
		{Title: "Start", Width: styles.CalculateDimensionsFromPercentage(40, tableWidth, 20)},
		{Title: "Time Spent", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Description", Width: styles.CalculateDimensionsFromPercentage(30, tableWidth, styles.UnlimitedDimension)},
	}
}

func createTableRows(worklogs []jira.TempoWorklog) []table.Row {
	rows := []table.Row{}

	log.Debugf("Worklogs: %v", len(worklogs))
	for _, worklog := range worklogs {
		rows = append(rows, table.Row{
			strconv.Itoa(worklog.Id),
			worklog.Start.Format("2006-01-02 15:04"),
			fmt.Sprintf("%.1f h", float64(worklog.TimeSpentInSeconds)/3600),
			"",
		})
	}
	return rows
}
