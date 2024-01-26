package worklog

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/cursor"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/textinput"
	"github.com/remshams/common/tui/styles"
	common "github.com/remshams/jira-control/tui/_common"
)

type keyMap struct {
	global common.GlobalKeyMap
	cursor cursor.KeyMap
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.cursor.Up,
		k.cursor.Down,
		k.global.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

var worklogKeys = keyMap{
	cursor: cursor.CursorKeyMap,
	global: common.GlobalKeys,
}

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
		cursor:   cursor.New(3),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		title.CreateSetPageTitleMsg("Log work"),
		help.CreateSetKeyMapMsg(worklogKeys),
	)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.cursor = m.cursor.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	styles := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf(
		"%s\n%s\n%s",
		styles.Render(cursor.RenderLine(m.issueKey.View(), m.cursor.Index() == 0, m.issueKey.Input.Focused())),
		styles.Render(cursor.RenderLine(m.work.View(), m.cursor.Index() == 1, m.work.Input.Focused())),
		cursor.RenderLine(m.comment.View(), m.cursor.Index() == 2, m.comment.Input.Focused()),
	)
}
