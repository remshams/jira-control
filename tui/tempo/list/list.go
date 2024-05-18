package tempo_workloglist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
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
	worklogs, err := adapter.App.TempoWorklogAdapter.List(jira.NewTempoWorklistQuery())
	if err != nil {
		errorChan <- err
	} else {
		worklogsChan <- worklogs
	}
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
	spinner  spinner.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
		state:   tempoWorklogStateLoading,
		spinner: spinner.New().WithLabel("Loading worklogs..."),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Tempo worklog list"),
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
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case tempoWorklogStateLoaded:
		return fmt.Sprintf("Worklogs: %d", len(m.worklogs))
	case tempoWorklogStateLoading:
		return m.spinner.View()
	case tempoWorklogStateError:
		return "Error loading worklogs"
	default:
		return ""
	}

}
