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
	"github.com/remshams/common/tui/utils"
	common "github.com/remshams/jira-control/tui/_common"
)

type keyMap struct {
	global    common.GlobalKeyMap
	cursor    cursor.KeyMap
	textinput textinput.KeyMap
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.cursor.Up,
		k.cursor.Down,
		textinput.TextInputKeyMap.Edit,
		k.global.Quit,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

const (
	navigate utils.ViewState = "navigate"
	edit     utils.ViewState = "edit"
)

var worklogKeys = keyMap{
	cursor:    cursor.CursorKeyMap,
	global:    common.GlobalKeys,
	textinput: textinput.TextInputKeyMap,
}

type Model struct {
	issueKey textinput.Model
	work     textinput.Model
	comment  textinput.Model
	cursor   cursor.CursorState
	state    utils.ViewState
}

func New() Model {
	return Model{
		issueKey: textinput.New("Issue key", ""),
		work:     textinput.New("Work", ""),
		comment:  textinput.New("Comment", ""),
		cursor:   cursor.New(3),
		state:    navigate,
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
		if m.state == navigate {
			switch {
			case key.Matches(msg, worklogKeys.textinput.Edit):
				cmd = m.updateSelection(msg)
				m.state = edit
			default:
				m.cursor = m.cursor.Update(msg)
			}
		} else {
			switch {
			case key.Matches(msg, worklogKeys.textinput.Discard):
				cmd = m.updateSelection(msg)
				m.state = navigate
			default:
				cmd = m.updateSelection(msg)
			}
		}
	}
	return m, cmd
}

func (m *Model) updateSelection(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.cursor.Index() {
	case 0:
		m.issueKey, cmd = m.issueKey.Update(msg)
	case 1:
		m.work, cmd = m.work.Update(msg)
	case 2:
		m.comment, cmd = m.comment.Update(msg)
	}
	return cmd

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
