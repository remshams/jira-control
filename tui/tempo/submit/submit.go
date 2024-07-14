package tempo_submit

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
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	app_store "github.com/remshams/jira-control/tui/store"
)

type initAction struct{}

func createInitAction() tea.Msg {
	return initAction{}
}

type loadTimesheetReviewersSuccessAction struct {
	Reviewers []jira.User
}

type loadReviewersErrorAction struct{}

func (m Model) createLoadReviewersAction() tea.Cmd {
	return func() tea.Msg {
		var reviewers []jira.User
		var err error
		reviewersChan := make(chan []jira.User)
		errorChan := make(chan error)
		go m.loadReviewers(reviewersChan, errorChan)
		select {
		case reviewers = <-reviewersChan:
		case err = <-errorChan:
		}
		if err != nil {
			return loadReviewersErrorAction{}
		} else {
			return loadTimesheetReviewersSuccessAction{
				Reviewers: reviewers,
			}
		}
	}
}

func (m Model) loadReviewers(reviewersChan chan []jira.User, errorChan chan error) {
	reviewers, err := m.timesheet.Reviewers()
	if err != nil {
		errorChan <- err
	} else {
		reviewersChan <- reviewers
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
	submit      key.Binding
}

func (m SubmitKeymap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.worklogList,
		m.help.Help,
		m.submit,
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
	submit: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "Submit and close timesheet"),
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
	table           table.Model[[]jira.User]
}

func New(adapter tui_jira.JiraAdapter) Model {
	model := Model{
		adapter: adapter,
		state:   stateLoading,
		spinner: spinner.New().WithLabel("Loading timesheet details..."),
	}
	model.table = table.
		New(createTableColumns, createTableRows, 5, 20).
		WithNotDataMessage("No reviewers")
	return model
}

func (m *Model) Init(timesheet jira.Timesheet, timesheetStatus jira.TimesheetStatus) tea.Cmd {
	m.timesheet = timesheet
	m.timesheetStatus = timesheetStatus
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
		cmd = m.createLoadReviewersAction()
	case loadTimesheetReviewersSuccessAction:
		m.reviewers = msg.Reviewers
		m.state = stateLoaded
		cmd = table.CreateTableDataUpdatedAction(m.reviewers)
	case loadReviewersErrorAction:
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
		case key.Matches(msg, SubmitKeys.submit):
			reviewerAccountId := m.table.SelectedRowCell(0)
			err := m.timesheet.Submit(reviewerAccountId)
			if err != nil {
				cmd = toast.CreateErrorToastAction("Could not submit timesheet")
			} else {
				cmd = tea.Batch(toast.CreateSuccessToastAction("Timesheet submitted"), m.createLoadReviewersAction())
			}
		default:
			m.table, cmd = m.table.Update(msg)
		}
	default:
		m.table, cmd = m.table.Update(msg)
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
			"%s\n%s\n%s",
			styles.Render(m.renderAccountInfo()),
			styles.Render(m.renderTimesheetInfo()),
			m.table.View(),
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

func createTableColumns(tableWidth int) []table.Column {
	return []table.Column{
		{Title: "AccountId", Width: styles.CalculateDimensionsFromPercentage(10, tableWidth, 10)},
		{Title: "Name", Width: styles.CalculateDimensionsFromPercentage(45, tableWidth, 40)},
		{Title: "Email", Width: styles.CalculateDimensionsFromPercentage(45, tableWidth, 40)},
	}
}

func createTableRows(reviewers []jira.User) []table.Row {
	rows := []table.Row{}

	log.Debugf("Reviewers: %d", len(reviewers))
	for _, reviewer := range reviewers {
		rows = append(rows, table.Row{
			reviewer.AccountId,
			reviewer.Name,
			reviewer.Email,
		})
	}
	return rows
}
