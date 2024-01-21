package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remshams/common/utils/logger"
	"github.com/remshams/jira-control/tui/home"
)

func main() {
	logger.PrepareLogger()
	f, err := tea.LogToFileWith("debug.log", "jira-control", log.Default())
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(home.New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
