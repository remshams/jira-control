package worklog

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/tui/bubbles/cursor"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/textinput"
)

type Model struct {
	issueKey textinput.Model
	work     textinput.Model
	comment  textinput.Model
	cursor   cursor.CursorState
}

func New() Model {
	return Model{
		issueKey: textinput.New("Issue key", ""),
		work:     textinput.New("Work", ""),
		comment:  textinput.New("Comment", ""),
		cursor:   cursor.New(1),
	}
}

func (m Model) Init() tea.Cmd {
	return title.CreateSetPageTitleMsg("Log work")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	return m, cmd
}

func (m Model) View() string {
	return "Log work"
}
