package issue_search_home

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/spinner"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	issue_search_form "github.com/remshams/jira-control/tui/issue/search/search-form"
	issue_search_result "github.com/remshams/jira-control/tui/issue/search/search-result"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type searchSuccessAction struct {
	issues []jira.Issue
}

type searchErrorAction struct{}

func createSearchIssueAction(searchRequest jira.IssueSearchRequest) tea.Cmd {
	return func() tea.Msg {
		issuesChan := make(chan []jira.Issue)
		searchErrorChan := make(chan error)
		go search(searchRequest, issuesChan, searchErrorChan)
		select {
		case issues := <-issuesChan:
			return searchSuccessAction{
				issues: issues,
			}
		case <-searchErrorChan:
			return searchErrorAction{}
		}
	}
}

func search(searchRequest jira.IssueSearchRequest, issuesChan chan []jira.Issue, searchErrorChan chan error) {
	defer close(issuesChan)
	defer close(searchErrorChan)
	issues, err := searchRequest.Search()
	if err != nil {
		searchErrorChan <- err
	} else {
		issuesChan <- issues
	}
}

const (
	stateSearchForm    utils.ViewState = "search-form"
	stateSearchResult  utils.ViewState = "search-result"
	stateSearchLoading utils.ViewState = "search-loading"
)

type Model struct {
	adapter      tui_jira.JiraAdapter
	searchForm   issue_search_form.Model
	searchResult issue_search_result.Model
	spinner      spinner.Model
	state        utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	spinner := spinner.New().WithLabel("Loading issues...")
	return Model{
		adapter:      adapter,
		searchForm:   issue_search_form.New(),
		searchResult: issue_search_result.New(),
		spinner:      spinner,
		state:        stateSearchForm,
	}
}

func (m Model) Init() tea.Cmd {
	switch m.state {
	case stateSearchForm:
		return m.searchForm.Init()
	case stateSearchResult:
		return m.searchResult.Init()
	default:
		return nil
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateSearchForm:
		cmd = m.processSearchFormUpdate(msg)
	case stateSearchResult:
		cmd = m.processSearchResultUpdate(msg)
	case stateSearchLoading:
		cmd = m.processLoadingUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processSearchFormUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case issue_search_form.ApplySearchAction:
		searchRequest := jira.NewIssueSearchRequest(m.adapter.IssueAdapter)
		searchRequest = searchRequest.WithSummary(msg.SearchTerm)
		m.state = stateSearchLoading
		cmd = tea.Batch(m.spinner.Tick(), createSearchIssueAction(searchRequest))
	case issue_search_form.SwitchViewAction:
		m.state = stateSearchResult
		cmd = m.searchResult.Init()
	default:
		m.searchForm, cmd = m.searchForm.Update(msg)
	}
	return cmd
}

func (m *Model) processLoadingUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case searchSuccessAction:
		m.state = stateSearchResult
		return tea.Batch(
			m.searchResult.Init(),
			issue_search_result.CreateSearchResultAction(msg.issues),
		)
	case searchErrorAction:
		m.state = stateSearchForm
		cmd = toast.CreateErrorToastAction("Error while searching issues")
	default:
		m.spinner, cmd = m.spinner.Update(msg)
	}
	return cmd
}

func (m *Model) processSearchResultUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case issue_search_result.SwitchViewAction:
		m.state = stateSearchForm
		cmd = m.searchForm.Init()
	default:
		m.searchResult, cmd = m.searchResult.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	styles := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf(
		"%s\n%s",
		styles.Render(m.searchForm.View()),
		m.renderSearchResult(),
	)
}

func (m Model) renderSearchResult() string {
	if m.state == stateSearchLoading {
		return m.spinner.View()
	} else {
		return m.searchResult.View()
	}
}
