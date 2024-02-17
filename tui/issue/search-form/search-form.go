package issue_search_form

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/cursor"
	"github.com/remshams/common/tui/bubbles/help"
	"github.com/remshams/common/tui/bubbles/textinput"
	common "github.com/remshams/jira-control/tui/_common"
)

type SearchFormKeyMap struct {
	global     common.GlobalKeyMap
	cursor     cursor.KeyMap
	textinput  textinput.KeyMap
	SwitchView key.Binding
}

func (m SearchFormKeyMap) ShortHelp() []key.Binding {
	shortHelp := []key.Binding{
		m.cursor.Up,
		m.cursor.Down,
		textinput.TextInputKeyMap.Edit,
		textinput.TextInputKeyMap.Discard,
		m.SwitchView,
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
	SwitchView: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Switch view"),
	),
}

type Model struct {
	searchTerm textinput.Model
	cursor     cursor.CursorState
}

func New() Model {
	return Model{
		searchTerm: textinput.New("Search", ""),
		cursor:     cursor.New(0),
	}
}

func (m Model) Init() tea.Cmd {
	return help.CreateSetKeyMapMsg(SearchFormKeys)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SearchFormKeys.SwitchView):
			cmd = help.CreateSetKeyMapMsg(SearchFormKeys)
		default:
			m.searchTerm, cmd = m.searchTerm.Update(msg)
		}
	default:
		m.searchTerm, cmd = m.searchTerm.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	return m.searchTerm.View()
}
