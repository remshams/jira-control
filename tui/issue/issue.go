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
	searchForm   *issue_search_form.SearchFormKeyMap
	searchResult *issue_search_result.SearchResultKeyMap
	switchView   key.Binding
}

func (m issueKeyMap) ShortHelp() []key.Binding {
	help := []key.Binding{}
	if m.searchForm != nil {
		help = append(help, issue_search_form.SearchFormKeys.ShortHelp()...)
	}
	if m.searchResult != nil {
		help = append(help, issue_search_result.SearchResultKeys.ShortHelp()...)
	}
	help = append(
		help,
		m.global.Tab.Tab,
		m.global.Quit,
		m.switchView,
	)
	return help
}

func (m issueKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var issueSearchKeys = issueKeyMap{
	global:     common.GlobalKeys,
	searchForm: &issue_search_form.SearchFormKeys,
	switchView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch view"),
	),
}

var issueResultKeys = issueKeyMap{
	global:       common.GlobalKeys,
	searchResult: &issue_search_result.SearchResultKeys,
	switchView: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "switch view"),
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
	return m, cmd
}

func (m Model) View() string {
	styles := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf(
		"%s\n%s",
		styles.Render(m.searchForm.View()),
		(m.searchResult.View()),
	)
}
