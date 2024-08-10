package tempo_workloglist

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	tempo_workloglistmodel "github.com/remshams/jira-control/tui/tempo/model"
	tempo_status "github.com/remshams/jira-control/tui/tempo/status"
)

type SwitchToSubmitViewAction struct{}

func createSwitchToSubmitViewAction() tea.Msg {
	return SwitchToSubmitViewAction{}
}

type SwitchToDeleteWorklogViewAction struct {
	Worklog jira.TempoWorklog
}

func createSwitchToDeleteWorklogViewAction(worklog jira.TempoWorklog) tea.Cmd {
	return func() tea.Msg {
		return SwitchToDeleteWorklogViewAction{
			Worklog: worklog,
		}
	}
}

type initAction struct{}

func createInitAction() tea.Msg {
	return initAction{}
}

type LoadWorklogsSuccessAction struct {
	Worklogs []jira.TempoWorklog
}

type LoadWorklogsErrorAction struct{}

func (m Model) createLoadWorklogsAction() tea.Cmd {
	return func() tea.Msg {
		worklogsChan := make(chan []jira.TempoWorklog)
		errorChan := make(chan error)
		go m.loadWorklogs(worklogsChan, errorChan)
		select {
		case worklogs := <-worklogsChan:
			return LoadWorklogsSuccessAction{Worklogs: worklogs}
		case <-errorChan:
			return LoadWorklogsErrorAction{}
		}
	}
}

func (m Model) loadWorklogs(worklogsChan chan []jira.TempoWorklog, errorChan chan error) {
	worklogs, err := m.timesheet.Worklogs(jira.NewTempoWorklogListQuery().WithSortDescending())
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
	submit key.Binding
	delete key.Binding
	reload key.Binding
}

func (m WorklogListKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.submit,
		m.delete,
		m.reload,
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
	submit: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Submit timesheet"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Delete worklog"),
	),
	reload: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "Reload worklogs"),
	),
}

const (
	tempoWorklogStateLoaded  utils.ViewState = "tempoWorklogStateLoaded"
	tempoWorklogStateLoading utils.ViewState = "tempoWorklogStateLoading"
	tempoWorklogStateError   utils.ViewState = "tempoWorklogStateError"
)

type Model struct {
	adapter         tui_jira.JiraAdapter
	state           utils.ViewState
	timesheet       jira.Timesheet
	timesheetStatus tempo_status.Model
	worklogs        []jira.TempoWorklog
	table           table.Model[[]jira.TempoWorklog]
	spinner         spinner.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	model := Model{
		adapter:         adapter,
		state:           tempoWorklogStateLoading,
		spinner:         spinner.New().WithLabel("Loading worklogs..."),
		timesheetStatus: tempo_status.New(),
	}
	model.table = table.
		New(createTableColumns, createTableRows, 5, 15).
		WithNotDataMessage("No worklogs")
	return model
}

func (m Model) Init(timesheet jira.Timesheet, timesheetStatus jira.TimesheetStatus) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var timesheetStatusCmd tea.Cmd
	m.timesheet = timesheet
	m.timesheetStatus, timesheetStatusCmd = m.timesheetStatus.Init(timesheetStatus)
	cmd = tea.Batch(
		title.CreateSetPageTitleMsg("Tempo worklog list"),
		help.CreateSetKeyMapMsg(WorklogListKeys),
		m.spinner.Tick(),
		createInitAction,
		timesheetStatusCmd,
	)
	return m, cmd
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
		cmd = m.createLoadWorklogsAction()
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
	switch msg := msg.(type) {
	case initAction:
		m.state = tempoWorklogStateLoading
		cmd = m.createLoadWorklogsAction()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, WorklogListKeys.submit):
			cmd = createSwitchToSubmitViewAction
		case key.Matches(msg, WorklogListKeys.delete):
			worklog := m.findWorklog(m.table.SelectedRowCell(0))
			if worklog != nil {
				cmd = createSwitchToDeleteWorklogViewAction(*worklog)
			} else {
				cmd = toast.CreateErrorToastAction(fmt.Sprintf("Could not find worklog with id %d", worklog.Id))
			}
		case key.Matches(msg, WorklogListKeys.reload):
			cmd = tempo_workloglistmodel.CreateLoadWorklogListAction
		default:
			m.table, cmd = m.table.Update(msg)
		}
	default:
		m.table, cmd = m.table.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	styles := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	switch m.state {
	case tempoWorklogStateLoaded:
		return fmt.Sprintf("%s\n%s", styles.Render(m.timesheetStatus.View()), styles.Render(m.table.View()))
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
		{Title: "Id", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Issue Id", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Start", Width: styles.CalculateDimensionsFromPercentage(20, tableWidth, 40)},
		{Title: "Time Spent", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Description", Width: styles.CalculateDimensionsFromPercentage(60, tableWidth, styles.UnlimitedDimension)},
	}
}

func createTableRows(worklogs []jira.TempoWorklog) []table.Row {
	rows := []table.Row{}

	log.Debugf("Worklogs: %d", len(worklogs))
	for _, worklog := range worklogs {
		description := worklog.Description
		if description == "" {
			description = "No worklog comment"
		}
		rows = append(rows, table.Row{
			strconv.Itoa(worklog.Id),
			strconv.Itoa(worklog.IssueKey),
			worklog.Start.Format("2006-01-02 15:04 (Mon)"),
			fmt.Sprintf("%.1f h", float64(worklog.TimeSpentSeconds)/3600),
			description,
		})
	}
	return rows
}

func (m Model) findWorklog(id string) *jira.TempoWorklog {
	for _, worklog := range m.worklogs {
		if strconv.Itoa(worklog.Id) == id {
			return &worklog
		}
	}
	return nil
}
