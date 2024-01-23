package worklog_log

import (
	tea "github.com/charmbracelet/bubbletea"
	title "github.com/remshams/common/tui/bubbles/page_title"
	"github.com/remshams/common/tui/utils"
)

const (
	issue utils.ViewState = "issue"
)

type Model struct {
	issueKey string
	state    utils.ViewState
}

func New() Model {
	return Model{
		issueKey: "",
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
