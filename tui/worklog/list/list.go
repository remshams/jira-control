package worklog_list

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	table_utils "github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	app_store "github.com/remshams/jira-control/tui/store"
)

type GoBackAction struct{}

func CreateGoBackAction() tea.Msg {
	return GoBackAction{}
}

type LoadWorklogsSuccessAction struct {
	Worklogs []jira.Worklog
}

type LoadWorklogsErrorAction struct{}

func createLoadWorklogsAction(adapter tui_jira.JiraAdapter, issue jira.Issue) tea.Cmd {
	return func() tea.Msg {
		worklogsChan := make(chan []jira.Worklog)
		errorChan := make(chan error)
		go loadWorklogs(adapter, issue, worklogsChan, errorChan)
		select {
		case worklogs := <-worklogsChan:
			return LoadWorklogsSuccessAction{Worklogs: worklogs}
		case <-errorChan:
			return LoadWorklogsErrorAction{}
		}
	}
}

func loadWorklogs(adapter tui_jira.JiraAdapter, issue jira.Issue, worklogsChan chan []jira.Worklog, errorChan chan error) {
	startedAfter := time.Now()
	// Load worklogs from the last 2 months
	startedAfter = startedAfter.Add(-7 * 24 * 4 * 2 * time.Hour)
	query := issue.WorklogsQuery().WithStartedAfter(startedAfter)
	worklogs, err := issue.Worklogs(query)
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
	goBack key.Binding
}

func (m WorklogListKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.help.Help,
		m.goBack,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m WorklogListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table_utils.TableKeyBindings(),
	}
}

var WorklogListKeys = WorklogListKeyMap{
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap(),
	goBack: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Go back"),
	),
}

const (
	worklogListStateLoading utils.ViewState = "workLogListStateLoading"
	worklogListStateLoaded  utils.ViewState = "workLogListStateLoaded"
	worklogListStateError   utils.ViewState = "workLogListStateError"
)

type Model struct {
	adapter  tui_jira.JiraAdapter
	issue    jira.Issue
	worklogs []jira.Worklog
	spinner  spinner.Model
	table    table.Model
	state    utils.ViewState
}

func New(adapter tui_jira.JiraAdapter, issue jira.Issue) Model {
	spinner := spinner.New(spinner.WithSpinner(spinner.Dot))
	spinner.Style = lipgloss.NewStyle().Foreground(styles.SelectedColor)
	return Model{
		issue:    issue,
		worklogs: []jira.Worklog{},
		spinner:  spinner,
		table: table.New(
			table.WithFocused(true),
		),
		state: worklogListStateLoading,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg(fmt.Sprintf("Worklog List for %s", m.issue.Key)),
		help.CreateSetKeyMapMsg(WorklogListKeys),
		createLoadWorklogsAction(m.adapter, m.issue),
		m.spinner.Tick,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case worklogListStateLoading:
		cmd = m.processLoadingUpdate(msg)
	case worklogListStateLoaded:
		cmd = m.processWorkListUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processWorkListUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.recalculateTableLayout()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, WorklogListKeys.help.Help):
			cmd = help.CreateToggleFullHelpMsg()
		case key.Matches(msg, WorklogListKeys.goBack):
			cmd = CreateGoBackAction
		default:
			m.table, cmd = m.table.Update(msg)
		}
	default:
		m.table, cmd = m.table.Update(msg)
	}
	return cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case LoadWorklogsSuccessAction:
		m.worklogs = msg.Worklogs
		m.state = worklogListStateLoaded
		m.recalculateTableLayout()
	case LoadWorklogsErrorAction:
		m.state = worklogListStateError
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case worklogListStateLoading:
		styles := lipgloss.NewStyle().Foreground(styles.SelectedColor)
		return fmt.Sprintf("%s %s", m.spinner.View(), styles.Render("Loading worklogs..."))
	case worklogListStateLoaded:
		return m.table.View()
	case worklogListStateError:
		return "Error loading worklogs"
	default:
		return ""
	}
}

func (m Model) createTable(columns []table.Column, rows []table.Row) table.Model {
	return table_utils.CreateTable(columns, rows)
}

func (m Model) createTableColumns() []table.Column {
	tableWidth, _ := m.calculateTableDimensions()
	return []table.Column{
		{Title: "Start", Width: styles.CalculateDimensionsFromPercentage(40, tableWidth, 20)},
		{Title: "Time Spent", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Description", Width: styles.CalculateDimensionsFromPercentage(50, tableWidth, styles.UnlimitedDimension)},
	}
}

func (m Model) calculateTableDimensions() (int, int) {
	width := app_store.LayoutStore.Width - 5
	height := app_store.LayoutStore.Height - 8
	if height < 0 {
		height = styles.CalculateDimensionsFromPercentage(80, app_store.LayoutStore.Height, styles.UnlimitedDimension)
	}
	return width, height
}

func (m Model) createTableRows() []table.Row {
	rows := []table.Row{}

	for _, worklog := range m.worklogs {
		hoursSpent := math.Ceil(float64(worklog.TimeSpentInSeconds / 3600))
		rows = append(rows, table.Row{
			worklog.Start.Format("2006-01-02 15:04"),
			fmt.Sprintf("%d h", int(hoursSpent)),
			worklog.Description,
		})
	}
	return rows
}

func (m *Model) recalculateTableLayout() {
	width, height := m.calculateTableDimensions()
	m.table.SetWidth(width)
	m.table.SetHeight(height)
	m.table.SetColumns(m.createTableColumns())
	m.table.SetRows(m.createTableRows())
}
