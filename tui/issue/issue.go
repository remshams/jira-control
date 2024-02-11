package issue

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	common "github.com/remshams/jira-control/tui/_common"
	issue_search_form "github.com/remshams/jira-control/tui/issue/search-form"
	issue_search_result "github.com/remshams/jira-control/tui/issue/search-result"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

const (
	stateSearchForm   utils.ViewState = "search-form"
	stateSearchResult utils.ViewState = "search-result"
)

type issueKeyMap struct {
	global       common.GlobalKeyMap
	help         help.KeyMap
	searchForm   *issue_search_form.SearchFormKeyMap
	searchResult *issue_search_result.SearchResultKeyMap
	switchView   key.Binding
}

func (m issueKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.switchView,
	}
	if m.searchForm != nil {
		shortHelp = append(shortHelp, issue_search_form.SearchFormKeys.ShortHelp()...)
	}
	if m.searchResult != nil {
		shortHelp = append(shortHelp, issue_search_result.SearchResultKeys.ShortHelp()...)
	}
	shortHelp = append(shortHelp, help.HelpKeys.Help)
	shortHelp = append(
		shortHelp,
		m.global.Tab.Tab,
		m.global.Quit,
	)
	return shortHelp
}

func (m issueKeyMap) FullHelp() [][]key.Binding {
	firstRow := m.ShortHelp()
	secondRow := []key.Binding{}
	if m.searchResult != nil {
		secondRow = append(secondRow, issue_search_result.SearchResultKeys.FullHelp()...)
	}
	if m.searchForm != nil {
		secondRow = append(secondRow, issue_search_form.SearchFormKeys.FullHelp()...)
	}
	return [][]key.Binding{
		firstRow,
		secondRow,
	}
}

var issueSearchKeys = issueKeyMap{
	global:     common.GlobalKeys,
	searchForm: &issue_search_form.SearchFormKeys,
	switchView: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "switch view"),
	),
}

var issueResultKeys = issueKeyMap{
	global:       common.GlobalKeys,
	help:         help.HelpKeys,
	searchResult: &issue_search_result.SearchResultKeys,
	switchView: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "switch view"),
	),
}

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
		help.CreateSetKeyMapMsg(issueSearchKeys),
		m.searchForm.Init(),
		m.searchResult.Init(),
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
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, issueSearchKeys.switchView):
			m.state = stateSearchResult
			cmd = help.CreateSetKeyMapMsg(issueResultKeys)
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
		case key.Matches(msg, issueResultKeys.help.Help):
			cmd = help.CreateToggleFullHelpMsg()
		case key.Matches(msg, issueResultKeys.switchView):
			m.state = stateSearchForm
			cmd = help.CreateSetKeyMapMsg(issueSearchKeys)
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
