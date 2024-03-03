package worklog_list

import (
	"fmt"
	"math"
	"time"

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
		table.DefaultKeyBindings,
	}
}

var WorklogListKeys = WorklogListKeyMap{
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap,
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
	table    table.Model[[]jira.Worklog]
	state    utils.ViewState
}

func New(adapter tui_jira.JiraAdapter, issue jira.Issue) Model {
	spinner := spinner.New().WithLabel("Loading worklogs...")
	model := Model{
		issue:    issue,
		worklogs: []jira.Worklog{},
		spinner:  spinner,
		state:    worklogListStateLoading,
	}
	model.table = table.New[[]jira.Worklog](createTableColumns, createTableRows, 5, 10)
	return model
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg(fmt.Sprintf("Worklogs for %s", m.issue.Key)),
		help.CreateSetKeyMapMsg(WorklogListKeys),
		createLoadWorklogsAction(m.adapter, m.issue),
		m.spinner.Tick(),
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
		cmd = table.CreateTableDataUpdatedAction(m.worklogs)
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
		return m.spinner.View()
	case worklogListStateLoaded:
		return m.table.View()
	case worklogListStateError:
		return "Error loading worklogs"
	default:
		return ""
	}
}

func createTableColumns(tableWidth int) []table.Column {
	return []table.Column{
		{Title: "Start", Width: styles.CalculateDimensionsFromPercentage(40, tableWidth, 20)},
		{Title: "Time Spent", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Description", Width: styles.CalculateDimensionsFromPercentage(50, tableWidth, styles.UnlimitedDimension)},
	}
}

func createTableRows(worklogs []jira.Worklog) []table.Row {
	rows := []table.Row{}

	log.Debugf("Worklogs: %v", len(worklogs))
	for _, worklog := range worklogs {
		hoursSpent := math.Ceil(float64(worklog.TimeSpentInSeconds / 3600))
		rows = append(rows, table.Row{
			worklog.Start.Format("2006-01-02 15:04"),
			fmt.Sprintf("%d h", int(hoursSpent)),
			worklog.Description,
		})
	}
	return rows
}
