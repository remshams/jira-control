package main

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
	jira "github.com/remshams/jira-control/jira/public"
	"github.com/remshams/jira-control/tui/home"
	tui_jira "github.com/remshams/jira-control/tui/jira"

	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/remshams/common/utils/logger"
)

func main() {
	logger.PrepareLogger()
	f, err := tea.LogToFileWith("debug.log", "jira-control", log.Default())
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	host := os.Getenv("HOST")
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 2289
	}
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithMiddleware(
			bm.Middleware(teaHandler),
			lm.Middleware(),
		),
	)
	if err != nil {
		log.Error("could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", "localhost", "port", 2289)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("could not stop server", "error", err)
	}
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	app, err := jira.PrepareApplication()
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	jiraAdapter := tui_jira.NewJiraAdapter(
		app.IssueAdapter,
		app.IssueWorklogAdapter,
	)
	m := home.New(jiraAdapter)
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}
