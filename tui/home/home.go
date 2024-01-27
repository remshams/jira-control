package home

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/remshams/common/tui/bubbles/help"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/bubbles/toast"
	"github.com/remshams/common/tui/styles"
	common "github.com/remshams/jira-control/tui/_common"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	"github.com/remshams/jira-control/tui/worklog"
)

type Model struct {
	adapter tui_jira.JiraAdapter
	title   title.Model
	toast   toast.Model
	help    help.Model
	worklog worklog.Model
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		title:   title.New(),
		toast:   toast.New(),
		help:    help.New(),
		worklog: worklog.New(adapter),
	}
}

func (m Model) Init() tea.Cmd {
	return m.worklog.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.toast, _ = m.toast.Update(msg)
	m.help, _ = m.help.Update(msg)
	m.title, _ = m.title.Update(msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.GlobalKeys.Quit):
			cmd = tea.Quit
		default:
			m.worklog, cmd = m.worklog.Update(msg)
		}
	}
	return m, cmd
}

func (m Model) View() string {
	return fmt.Sprintf(
		"%s\n%s\n%s\n%s",
		m.title.View(),
		m.renderContent(),
		m.renderHelp(),
		m.renderToast(),
	)
}

func (m Model) renderContent() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return style.Render(m.worklog.View())
}

func (m Model) renderToast() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return style.Render(m.toast.View())
}

func (m Model) renderHelp() string {
	style := lipgloss.NewStyle().PaddingTop(styles.Padding)
	return style.Render(m.help.View())
}
