package tempo_home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	app_store "github.com/remshams/jira-control/tui/store"
	tempo_workloglist "github.com/remshams/jira-control/tui/tempo/list"
	tempo_submit "github.com/remshams/jira-control/tui/tempo/submit"
)

type loadTimesheetStatusSuccessAction struct {
	Status jira.TimesheetStatus
}

type loadTimesheetStatusErrorAction struct {
	Error error
}

func (m Model) createLoadTimesheetStatusAction() tea.Cmd {
	return func() tea.Msg {
		statusChan := make(chan jira.TimesheetStatus)
		errorChan := make(chan error)
		go m.loadTimesheetStatus(statusChan, errorChan)
		select {
		case status := <-statusChan:
			return loadTimesheetStatusSuccessAction{Status: status}
		case error := <-errorChan:
			return loadTimesheetStatusErrorAction{Error: error}
		}
	}
}

func (m Model) loadTimesheetStatus(statusChan chan jira.TimesheetStatus, errorChan chan error) {
	status, err := m.timesheet.Status()
	if err != nil {
		errorChan <- err
	} else {
		statusChan <- status
	}
}

const (
	stateWorklog      utils.ViewState = "worklog"
	stateSubmit       utils.ViewState = "submit"
	stateLoading      utils.ViewState = "loading"
	stateLoadingError utils.ViewState = "loadingError"
)

type Model struct {
	worklogList     tempo_workloglist.Model
	submit          tempo_submit.Model
	adapter         tui_jira.JiraAdapter
	timesheet       jira.Timesheet
	timesheetStatus jira.TimesheetStatus
	state           utils.ViewState
	spinner         spinner.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter:     adapter,
		worklogList: tempo_workloglist.New(adapter),
		submit:      tempo_submit.New(adapter),
		state:       stateLoading,
		spinner:     spinner.New().WithLabel("Loading timesheet"),
	}
}

func (m *Model) Init() tea.Cmd {
	m.timesheet = jira.NewTimesheet(m.adapter.App.TempoTimesheetAdapter, m.adapter.App.TempoWorklogAdapter, app_store.AppDataStore.Account.AccountId)
	return tea.Batch(m.createLoadTimesheetStatusAction(), m.spinner.Tick())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateLoading:
		cmd = m.processLoadingUpdate(msg)
	case stateWorklog:
		cmd = m.processWorklogListUpdate(msg)
	case stateSubmit:
		cmd = m.processSubmitUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case loadTimesheetStatusSuccessAction:
		m.timesheetStatus = msg.Status
		m.state = stateWorklog
		m.worklogList, cmd = m.worklogList.Init(m.timesheet)
	case loadTimesheetStatusErrorAction:
		m.timesheetStatus = jira.TimesheetStatus{}
		m.state = stateLoadingError
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return cmd
}

func (m *Model) processWorklogListUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.(type) {
	case tempo_workloglist.SwitchToSubmitViewAction:
		m.state = stateSubmit
		m.submit, cmd = m.submit.Init(m.timesheet, m.timesheetStatus)
	default:
		m.worklogList, cmd = m.worklogList.Update(msg)
	}
	return cmd
}

func (m *Model) processSubmitUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg.(type) {
	case tempo_submit.SwitchToWorklogListView:
		m.state = stateWorklog
		m.worklogList, cmd = m.worklogList.Init(m.timesheet)
	default:
		m.submit, cmd = m.submit.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.spinner.View()
	case stateWorklog:
		return m.worklogList.View()
	case stateSubmit:
		return m.submit.View()
	case stateLoadingError:
		return "Error loading timesheet"
	default:
		return "Invalid view state"
	}
}
