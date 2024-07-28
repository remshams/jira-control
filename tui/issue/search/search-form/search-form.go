package issue_search_form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/cursor"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/textinput"
	"github.com/remshams/common/tui/utils"
	common "github.com/remshams/jira-control/tui/common"
)

type ApplySearchAction struct {
	SearchTerm string
}

func CreateApplySearchAction(searchTerm string) tea.Cmd {
	return func() tea.Msg {
		return ApplySearchAction{SearchTerm: searchTerm}
	}
}

type SwitchViewAction struct{}

func CreateSwitchViewAction() tea.Cmd {
	return func() tea.Msg {
		return SwitchViewAction{}
	}
}

type SearchFormKeyMap struct {
	global     common.GlobalKeyMap
	cursor     cursor.KeyMap
	textinput  textinput.KeyMap
	switchView key.Binding
}

func (m SearchFormKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.cursor.Up,
		m.cursor.Down,
		textinput.TextInputKeyMap.Edit,
		textinput.TextInputKeyMap.Discard,
		m.switchView,
	}
	return append(shortHelp, m.global.KeyBindings()...)
}

func (m SearchFormKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var SearchFormKeys = SearchFormKeyMap{
	global:    common.GlobalKeys,
	cursor:    cursor.CursorKeyMap,
	textinput: textinput.TextInputKeyMap,
	switchView: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Switch to result"),
	),
}

const (
	searchFormNavigate utils.ViewState = "navigate"
	searchFormEdit     utils.ViewState = "edit"
)

type Model struct {
	searchTerm textinput.Model
	cursor     cursor.CursorState
	state      utils.ViewState
}

func New() Model {
	return Model{
		searchTerm: textinput.New("Search", ""),
		cursor:     cursor.New(0),
		state:      searchFormNavigate,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Issue Search"),
		help.CreateSetKeyMapMsg(SearchFormKeys),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case searchFormEdit:
		cmd = m.processEdit(msg)
	case searchFormNavigate:
		cmd = m.processNavigate(msg)
	}
	return m, cmd
}

func (m *Model) processNavigate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SearchFormKeys.textinput.Edit):
			m.state = searchFormEdit
			m.searchTerm, cmd = m.searchTerm.Update(msg)
		case key.Matches(msg, SearchFormKeys.switchView):
			cmd = tea.Batch(help.CreateSetKeyMapMsg(SearchFormKeys), CreateSwitchViewAction())
		}
	}
	return cmd
}

func (m *Model) processEdit(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SearchFormKeys.textinput.Discard):
			m.state = searchFormNavigate
			m.searchTerm, cmd = m.searchTerm.Update(msg)
		case key.Matches(msg, SearchFormKeys.textinput.Apply):
			var searchTermCmd tea.Cmd
			m.state = searchFormNavigate
			m.searchTerm, searchTermCmd = m.searchTerm.Update(msg)
			cmd = CreateApplySearchAction(m.searchTerm.Input.Value())
			cmd = tea.Batch(cmd, searchTermCmd)
		default:
			m.searchTerm, cmd = m.searchTerm.Update(msg)
		}
	}
	return cmd
}

func (m Model) View() string {
	return m.searchTerm.View()
}
