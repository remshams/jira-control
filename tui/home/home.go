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
	common_worklog "github.com/remshams/jira-control/tui/_common/worklog"
	favorite_home "github.com/remshams/jira-control/tui/favorits/home"
	issue_home "github.com/remshams/jira-control/tui/issue/home"
	tui_jira "github.com/remshams/jira-control/tui/jira"
	tui_last_updated "github.com/remshams/jira-control/tui/last-updated/home"
	app_store "github.com/remshams/jira-control/tui/store"
	tempo_home "github.com/remshams/jira-control/tui/tempo/home"
	worklog_details "github.com/remshams/jira-control/tui/worklog/details"
)

const (
	stateIssue       utils.ViewState = "issue"
	stateWorklog     utils.ViewState = "worklog"
	stateLastUpdated utils.ViewState = "last_updated"
	stateFavorites   utils.ViewState = "favorites"
	stateTempo       utils.ViewState = "tempo"
)

type Model struct {
	adapter     tui_jira.JiraAdapter
	tab         tabs.Model
	title       title.Model
	toast       toast.Model
	help        help.Model
	issue       issue_home.Model
	worklog     worklog_details.Model
	tempo       tempo_home.Model
	lastUpdated tui_last_updated.Model
	favorites   favorite_home.Model
	state       utils.ViewState
}

func New(adapter tui_jira.JiraAdapter) Model {
	return Model{
		tab: tabs.New(
			[]string{"Worklog", "Issues", "Last Updated", "Favorites", "Tempo"},
		),
		title:       title.New(),
		toast:       toast.New(),
		help:        help.New(),
		worklog:     worklog_details.New(adapter),
		tempo:       tempo_home.New(adapter),
		issue:       issue_home.New(adapter),
		lastUpdated: tui_last_updated.New(adapter),
		favorites:   favorite_home.New(adapter),
		state:       stateWorklog,
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
	case tea.WindowSizeMsg:
		app_store.LayoutStore.Height = msg.Height
		app_store.LayoutStore.Width = msg.Width
		cmd = m.processUpdate(msg)
	case tabs.TabSelectedMsg:
		cmd = m.processTab(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, common.GlobalKeys.Quit):
			cmd = tea.Quit
		default:
			cmd = m.processUpdate(msg)
		}
	default:
		cmd = m.processUpdate(msg)
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
	case 2:
		cmd = m.lastUpdated.Init()
		m.state = stateLastUpdated
	case 3:
		cmd = m.favorites.Init()
		m.state = stateFavorites
	case 4:
		cmd = m.tempo.Init()
		m.state = stateTempo
	}
	return cmd
}

func (m *Model) processUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.state {
	case stateIssue:
		cmd = m.processIssueUpdate(msg)
	case stateWorklog:
		m.worklog, cmd = m.worklog.Update(msg)
	case stateLastUpdated:
		cmd = m.processLastUpdatedUpdate(msg)
	case stateFavorites:
		cmd = m.processFavoritesUpdate(msg)
	case stateTempo:
		cmd = m.processTempoUpdate(msg)
	}
	return cmd
}

func (m *Model) processIssueUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case common_worklog.LogWorkAction:
		cmd = m.logWork(msg.IssueKey, nil)
	default:
		m.issue, cmd = m.issue.Update(msg)
	}
	return cmd
}

func (m *Model) processLastUpdatedUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case common_worklog.LogWorkAction:
		cmd = m.logWork(msg.IssueKey, nil)
	default:
		m.lastUpdated, cmd = m.lastUpdated.Update(msg)
	}
	return cmd
}

func (m *Model) processFavoritesUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case common_worklog.LogWorkAction:
		cmd = m.logWork(msg.IssueKey, msg.HoursSpent)
	default:
		m.favorites, cmd = m.favorites.Update(msg)
	}
	return cmd
}

func (m *Model) processTempoUpdate(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	m.tempo, cmd = m.tempo.Update(msg)
	return cmd
}

func (m *Model) logWork(issueKey string, hoursSpent *float64) tea.Cmd {
	var cmd tea.Cmd
	m.state = stateWorklog
	cmd = tea.Batch(
		tabs.CreateSelectTabAction(0),
		m.worklog.Init(),
		worklog_details.CreateSetIssueKeyAction(issueKey),
	)
	if hoursSpent != nil {
		cmd = tea.Batch(cmd, worklog_details.CreateHoursSpentAction(*hoursSpent))
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
	case stateLastUpdated:
		return style.Render(m.lastUpdated.View())
	case stateFavorites:
		return style.Render(m.favorites.View())
	case stateTempo:
		return style.Render(m.tempo.View())
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
