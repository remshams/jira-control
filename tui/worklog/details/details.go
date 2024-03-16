package worklog_details

import (
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/cursor"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/textinput"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	jira "github.com/remshams/jira-control/jira/public"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

type SetIssueKeyAction struct {
	IssueKey string
}

func CreateSetIssueKeyAction(issueKey string) tea.Cmd {
	return func() tea.Msg {
		return SetIssueKeyAction{
			IssueKey: issueKey,
		}
	}
}

type keyMap struct {
	global    common.GlobalKeyMap
	cursor    cursor.KeyMap
	textinput textinput.KeyMap
	save      key.Binding
}

var worklogKeys = keyMap{
	global:    common.GlobalKeys,
	cursor:    cursor.CursorKeyMap,
	textinput: textinput.TextInputKeyMap,
	save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "save"),
	),
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.cursor.Up,
		k.cursor.Down,
		textinput.TextInputKeyMap.Edit,
		k.save,
		k.global.Tab.Tab,
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

type Model struct {
	adapter  tui_jira.JiraAdapter
	issueKey textinput.Model
	work     textinput.Model
	start    textinput.Model
	comment  textinput.Model
	cursor   cursor.CursorState
	state    utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		adapter:  adapter,
		issueKey: textinput.New("Issue key", ""),
		work:     textinput.New("Work", ""),
		start:    textinput.New("Start", ""),
		comment:  textinput.New("Comment", ""),
		cursor:   cursor.New(4),
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
	case SetIssueKeyAction:
		m.issueKey.Input.SetValue(msg.IssueKey)
	case tea.KeyMsg:
		if m.state == navigate {
			switch {
			case key.Matches(msg, worklogKeys.textinput.Edit):
				cmd = m.updateSelection(msg)
				m.state = edit
			case key.Matches(msg, worklogKeys.save):
				cmd = m.logWorkInJira()
			default:
				m.cursor = m.cursor.Update(msg)
			}
		} else {
			switch {
			case key.Matches(msg, worklogKeys.textinput.Apply):
				cmd = m.updateSelection(msg)
				m.state = navigate
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
		m.start, cmd = m.start.Update(msg)
	case 3:
		m.comment, cmd = m.comment.Update(msg)
	}
	return cmd
}

func (m *Model) logWorkInJira() tea.Cmd {
	hoursSpent, err := strconv.ParseFloat(m.work.Input.Value(), 64)
	if err != nil {
		return toast.CreateErrorToastAction("Invalid work value")
	}
	worklog := jira.NewWorklog(m.adapter.App.IssueWorklogAdapter, m.issueKey.Input.Value(), hoursSpent)
	worklog.Description = m.comment.Input.Value()
	if m.start.Input.Value() != "" {
		// Currently only date is supported as start
		worklog.Start, err = time.Parse(
			time.RFC3339,
			fmt.Sprintf(
				"%sT%sZ",
				m.start.Input.Value(),
				time.Now().Format("15:04:05")),
		)
	}
	if err != nil {
		return toast.CreateErrorToastAction("Invalid start time")
	}
	err = worklog.Log()
	if err != nil {
		return toast.CreateErrorToastAction("Could not save worklog in jira")
	}
	return toast.CreateSuccessToastAction("Worklog updated")
}

func (m Model) View() string {
	styles := lipgloss.NewStyle().PaddingBottom(styles.Padding)
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		styles.Render(cursor.RenderLine(m.issueKey.View(), m.cursor.Index() == 0, m.issueKey.Input.Focused())),
		styles.Render(cursor.RenderLine(m.work.View(), m.cursor.Index() == 1, m.work.Input.Focused())),
		styles.Render(cursor.RenderLine(m.start.View(), m.cursor.Index() == 2, m.start.Input.Focused())),
		cursor.RenderLine(m.comment.View(), m.cursor.Index() == 3, m.comment.Input.Focused()),
	)
}
