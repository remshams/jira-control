package worklog_list

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
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
	goBack key.Binding
}

func (m WorklogListKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.goBack,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m WorklogListKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var WorklogListKeys = WorklogListKeyMap{
	global: common.GlobalKeys,
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
	state    utils.ViewState
}

func New(adapter tui_jira.JiraAdapter, issue jira.Issue) Model {
	spinner := spinner.New(spinner.WithSpinner(spinner.Dot))
	spinner.Style = lipgloss.NewStyle().Foreground(styles.SelectedColor)
	return Model{
		issue:    issue,
		worklogs: []jira.Worklog{},
		spinner:  spinner,
		state:    worklogListStateLoading,
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, WorklogListKeys.goBack):
			cmd = CreateGoBackAction
		}
	}
	return cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case LoadWorklogsSuccessAction:
		m.worklogs = msg.Worklogs
		m.state = worklogListStateLoaded
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
		return fmt.Sprintf("Loaded %d worklogs", len(m.worklogs))
	case worklogListStateError:
		return "Error loading worklogs"
	default:
		return ""
	}
}
