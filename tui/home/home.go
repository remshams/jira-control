package home

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/tabs"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	"github.com/remshams/common/tui/utils"
	common "github.com/remshams/jira-control/tui/_common"
	"github.com/remshams/jira-control/tui/issue"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	"github.com/remshams/jira-control/tui/worklog"
)

const (
	stateIssue   utils.ViewState = "issue"
	stateWorklog utils.ViewState = "worklog"
)

type Model struct {
	adapter tui_jira.JiraAdapter
	tab     tabs.Model
	title   title.Model
	toast   toast.Model
	help    help.Model
	issue   issue.Model
	worklog worklog.Model
	state   utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		tab: tabs.New(
			[]string{"Worklog", "Issues"},
		),
		title:   title.New(),
		toast:   toast.New(),
		help:    help.New(),
		worklog: worklog.New(adapter),
		issue:   issue.New(),
		state:   stateWorklog,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.tab.Init(),
		m.worklog.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var tabsCmd tea.Cmd
	m.tab, tabsCmd = m.tab.Update(msg)
	m.toast, _ = m.toast.Update(msg)
	m.help, _ = m.help.Update(msg)
	m.title, _ = m.title.Update(msg)
	switch msg := msg.(type) {
	case tabs.TabSelectedMsg:
		cmd = m.processTab(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.GlobalKeys.Quit):
			cmd = tea.Quit
		default:
			cmd = m.processKeyPress(msg)
		}
	}
	return m, tea.Batch(cmd, tabsCmd)
}

func (m *Model) processTab(msg tabs.TabSelectedMsg) tea.Cmd {
	var cmd tea.Cmd
	switch msg {
	case 0:
		cmd = m.worklog.Init()
		m.state = stateWorklog
	case 1:
		cmd = m.issue.Init()
		m.state = stateIssue
	}
	return cmd
}

func (m *Model) processKeyPress(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case stateIssue:
		m.issue, cmd = m.issue.Update(msg)
	case stateWorklog:
		m.worklog, cmd = m.worklog.Update(msg)
	}
	return cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s",
		m.title.View(),
		m.renderTab(),
		m.renderContent(),
		m.renderHelp(),
		m.renderToast(),
	)
}

func (m Model) renderTab() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return style.Render(m.tab.View())
}

func (m Model) renderContent() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	switch m.state {
	case stateIssue:
		return style.Render(m.issue.View())
	case stateWorklog:
		return style.Render(m.worklog.View())
	default:
		return "View does not exist"
	}
}

func (m Model) renderToast() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return style.Render(m.toast.View())
}

func (m Model) renderHelp() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return style.Render(m.help.View())
}
