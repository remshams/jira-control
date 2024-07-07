package app_store

import jira "github.com/remshams/jira-control/jira/public"

type Layout struct {
	Height int
	Width  int
}

type AppData struct {
	Account jira.User
}

var LayoutStore = Layout{
	Height: 0,
	Width:  0,
}

var AppDataStore = AppData{
	Account: jira.User{},
}
