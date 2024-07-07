package tempo_submit

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/bubbles/table"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	app_store "github.com/remshams/jira-control/tui/store"
)

var wg sync.WaitGroup

type initAction struct{}

func createInitAction() tea.Msg {
	return initAction{}
}

type loadTimesheetInfoSuccessAction struct {
	Status    jira.TimesheetStatus
	Reviewers []jira.User
}

type loadTimesheetInfoErrorAction struct{}

func (m Model) loadTimesheetInfo() tea.Cmd {
	return func() tea.Msg {
		var status jira.TimesheetStatus
		var reviewers []jira.User
		var err error
		statusChan := make(chan jira.TimesheetStatus)
		reviewersChan := make(chan []jira.User)
		errorChan := make(chan error)
		wg.Add(2)
		go m.loadTimesheetStatus(statusChan, errorChan)
		go m.loadReviewers(reviewersChan, errorChan)
		select {
		case status = <-statusChan:
		case err = <-errorChan:
		}
		select {
		case reviewers = <-reviewersChan:
		case err = <-errorChan:
		}
		wg.Wait()
		if err != nil {
			return loadTimesheetInfoErrorAction{}
		} else {
			return loadTimesheetInfoSuccessAction{
				Status:    status,
				Reviewers: reviewers,
			}
		}
	}
}

func (m Model) loadTimesheetStatus(statusChan chan jira.TimesheetStatus, errorChan chan error) {
	status, err := m.timesheet.Status()
	if err != nil {
		errorChan <- err
		wg.Done()
	} else {
		statusChan <- status
		wg.Done()
	}
}

func (m Model) loadReviewers(reviewersChan chan []jira.User, errorChan chan error) {
	reviewers, err := m.timesheet.Reviewers()
	if err != nil {
		errorChan <- err
		wg.Done()
	} else {
		reviewersChan <- reviewers
		wg.Done()
	}
}

type SwitchToWorklogListView struct{}

func createSwitchToWorklogListView() tea.Msg {
	return SwitchToWorklogListView{}
}

type SubmitKeymap struct {
	global      common.GlobalKeyMap
	help        help.KeyMap
	table       table.KeyMap
	worklogList key.Binding
}

func (m SubmitKeymap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.worklogList,
		m.help.Help,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m SubmitKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table.DefaultKeyBindings,
	}
}

var SubmitKeys = SubmitKeymap{
	global: common.GlobalKeys,
	help:   help.HelpKeys,
	table:  table.DefaultKeyMap,
	worklogList: key.NewBinding(
		key.WithKeys("l"),
		key.WithHelp("l", "Show worklog list"),
	),
}

const (
	stateLoading      utils.ViewState = "loading"
	stateLoaded       utils.ViewState = "loaded"
	stateLoadingError utils.ViewState = "loadingError"
)

type Model struct {
	adapter         tui_jira.JiraAdapter
	timesheet       jira.Timesheet
	timesheetStatus jira.TimesheetStatus
	reviewers       []jira.User
	state           utils.ViewState
	spinner         spinner.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
		state:   stateLoading,
		spinner: spinner.New().WithLabel("Loading timesheet details..."),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Submit timesheet"),
		help.CreateSetKeyMapMsg(SubmitKeys),
		m.spinner.Tick(),
		createInitAction,
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateLoading:
		cmd = m.processLoadingUpdate(msg)
	case stateLoaded:
		cmd = m.processLoadedUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case initAction:
		m.timesheet = jira.NewTimesheet(m.adapter.App.TempoTimesheetAdapter, app_store.AppDataStore.Account.AccountId)
		cmd = m.loadTimesheetInfo()
	case loadTimesheetInfoSuccessAction:
		m.timesheetStatus = msg.Status
		m.reviewers = msg.Reviewers
		m.state = stateLoaded
	case loadTimesheetInfoErrorAction:
		m.state = stateLoadingError
		m.timesheetStatus = jira.TimesheetStatus{}
		m.reviewers = []jira.User{}
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
		case key.Matches(msg, SubmitKeys.worklogList):
			cmd = createSwitchToWorklogListView
		}
	}
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case stateLoading:
		return m.spinner.View()
	case stateLoaded:
		styles := lipgloss.NewStyle().PaddingBottom(styles.Padding)
		return fmt.Sprintf(
			"%s\n%s",
			styles.Render(m.renderAccountInfo()),
			m.renderTimesheetInfo(),
		)
	case stateLoadingError:
		return "Could not load timesheet details"
	default:
		return ""
	}
}

func (m Model) renderAccountInfo() string {
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.renderKeyValue(
			"AccountId",
			app_store.AppDataStore.Account.AccountId,
		),
		m.renderKeyValue(
			"Name",
			app_store.AppDataStore.Account.Name,
		),
		m.renderKeyValue(
			"Email",
			app_store.AppDataStore.Account.Email,
		),
	)
}

func (m Model) renderTimesheetInfo() string {
	spentHoursColor := styles.TextSuccessColor
	if m.timesheetStatus.RequiredHours-m.timesheetStatus.SpentHours > 0 {
		spentHoursColor = styles.TextErrorColor
	}
	return fmt.Sprintf(
		"%s\n%s\n%s",
		m.renderKeyValue(
			"Required hours",
			fmt.Sprintf("%d hours", m.timesheetStatus.RequiredHours),
		),
		m.renderKeyValue(
			"Spent hours",
			fmt.Sprintf("%s hours", spentHoursColor.Render(strconv.Itoa(m.timesheetStatus.SpentHours))),
		),
		m.renderKeyValue(
			"Status",
			m.timesheetStatus.Status,
		),
	)
}

func (m Model) renderKeyValue(key string, value string) string {
	return fmt.Sprintf("%s%s %s", styles.TextAccentColor.Render(key), styles.TextAccentColor.Render(":"), value)
}
