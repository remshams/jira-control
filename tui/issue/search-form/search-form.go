package issue_search_form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/cursor"
	"github.com/remshams/common/tui/bubbles/textinput"
)

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
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return m.searchTerm.View()
}
