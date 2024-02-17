package issue_home

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	issue_search_form "github.com/remshams/jira-control/tui/issue/search-form"
	issue_search_result "github.com/remshams/jira-control/tui/issue/search-result"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

const (
	stateSearchForm   utils.ViewState = "search-form"
	stateSearchResult utils.ViewState = "search-result"
)

type Model struct {
	adapter      tui_jira.JiraAdapter
	searchForm   issue_search_form.Model
	searchResult issue_search_result.Model
	state        utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter:      adapter,
		searchForm:   issue_search_form.New(),
		searchResult: issue_search_result.New(),
		state:        stateSearchForm,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Issue"),
		m.searchForm.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateSearchForm:
		cmd = m.processSearchFormUpdate(msg)
	case stateSearchResult:
		cmd = m.processSearchResultUpdate(msg)
	}
	return m, cmd
}

func (m *Model) processSearchFormUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case issue_search_form.ApplySearchAction:
		searchRequest := jira.NewIssueSearchRequest(m.adapter.IssueAdapter)
		issues, err := searchRequest.Search()
		if err != nil {
			cmd = toast.CreateErrorToastAction("Could not search for issues")
		}
		cmd = issue_search_result.CreateSearchResultAction(issues)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, issue_search_form.SearchFormKeys.SwitchView):
			m.state = stateSearchResult
			cmd = m.searchResult.Init()
		default:
			m.searchForm, cmd = m.searchForm.Update(msg)
		}
	default:
		m.searchForm, cmd = m.searchForm.Update(msg)
	}
	return cmd
}

func (m *Model) processSearchResultUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, issue_search_result.SearchResultKeys.SwitchView):
			m.state = stateSearchForm
			cmd = m.searchForm.Init()
		default:
			m.searchResult, cmd = m.searchResult.Update(msg)
		}
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
		(m.searchResult.View()),
	)
}
