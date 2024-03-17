package tui_last_updated_issue_list

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
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

type loadIssuesSuccessAction struct {
	issues []jira.Issue
}

func createLoadIssuesSuccessAction(issues []jira.Issue) tea.Msg {
	return loadIssuesSuccessAction{
		issues: issues,
	}
}

type loadIssuesErrorAction struct{}

func createLoadIssuesAction(adapter tui_jira.JiraAdapter) tea.Cmd {
	return func() tea.Msg {
		issuesChan := make(chan []jira.Issue)
		errorChan := make(chan error)
		go loadIssues(adapter, issuesChan, errorChan)
		select {
		case issues := <-issuesChan:
			return createLoadIssuesSuccessAction(issues)
		case <-errorChan:
			return loadIssuesErrorAction{}
		}
	}
}

func loadIssues(adapter tui_jira.JiraAdapter, issuesChan chan []jira.Issue, errorChan chan error) {
	defer close(issuesChan)
	defer close(errorChan)
	issueSearchRequest := jira.NewIssueSearchRequest(adapter.App.IssueAdapter).
		WithOrderBy(jira.NewOrderBy([]string{"updated"}, jira.SortingDesc)).
		WithUpdatedBy(adapter.App.Username)
	issues, err := issueSearchRequest.Search()
	if err != nil {
		errorChan <- err
	} else {
		issuesChan <- issues
	}
}

type LastUpdatedKeymap struct {
	global common.GlobalKeyMap
	table  table.KeyMap
}

func (m LastUpdatedKeymap) ShortHelp() []key.Binding {
	return m.global.KeyBindings()
}

func (m LastUpdatedKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.ShortHelp(),
		table.DefaultKeyBindings,
	}
}

var LastUpdatedKeys = LastUpdatedKeymap{
	global: common.GlobalKeys,
	table:  table.DefaultKeyMap,
}

const (
	lastUpdatedListStateLoading utils.ViewState = "lastUpdatedListStateLoading"
	lastUpdatedListStateLoaded  utils.ViewState = "lastUpdatedListStateLoaded"
)

type Model struct {
	adapter tui_jira.JiraAdapter
	table   table.Model[[]jira.Issue]
	spinner spinner.Model
	state   utils.ViewState
	issues  []jira.Issue
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter: adapter,
		table: table.
			New(createTableColumns, createTableRows, 5, 10).
			WithNotDataMessage("No issues"),
		spinner: spinner.New().WithLabel("Loading issues"),
		state:   lastUpdatedListStateLoading,
		issues:  []jira.Issue{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Last Updated Issues"),
		help.CreateSetKeyMapMsg(LastUpdatedKeys),
		m.spinner.Tick(),
		createLoadIssuesAction(m.adapter),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case lastUpdatedListStateLoading:
		cmd = m.processLoadingUpdate(msg)
	case lastUpdatedListStateLoaded:
		cmd = m.proccessLoadedUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case loadIssuesSuccessAction:
		m.state = lastUpdatedListStateLoaded
		cmd = table.CreateTableDataUpdatedAction(msg.issues)
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return cmd
}

func (m *Model) proccessLoadedUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return cmd
}

func (m Model) View() string {
	switch m.state {
	case lastUpdatedListStateLoading:
		return m.spinner.View()
	case lastUpdatedListStateLoaded:
		return m.table.View()
	default:
		return ""
	}
}

func createTableColumns(tableWidth int) []table.Column {
	return []table.Column{
		{Title: "Key", Width: styles.CalculateDimensionsFromPercentage(20, tableWidth, 20)},
		{Title: "Updated", Width: styles.CalculateDimensionsFromPercentage(20, tableWidth, 20)},
		{Title: "Summary", Width: styles.CalculateDimensionsFromPercentage(60, tableWidth, styles.UnlimitedDimension)},
	}
}

func createTableRows(issues []jira.Issue) []table.Row {
	rows := []table.Row{}

	for _, issue := range issues {
		rows = append(rows, table.Row{
			issue.Key,
			issue.Project.Updated.Format("2006-01-02 15:04:05"),
			issue.Summary,
		})
	}
	return rows
}
