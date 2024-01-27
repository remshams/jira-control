package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/utils/logger"
	jira "github.com/remshams/jira-control/jira/public"
	"github.com/remshams/jira-control/tui/home"
	tui_jira "github.com/remshams/jira-control/tui/jira"
)

func main() {
	logger.PrepareLogger()
	f, err := tea.LogToFileWith("debug.log", "jira-control", log.Default())
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	jiraAdapter := tui_jira.NewJiraAdapter(jira.NewWorklogMockAdapter())
	p := tea.NewProgram(home.New(jiraAdapter))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
