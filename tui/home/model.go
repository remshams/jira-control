package home

import "github.com/charmbracelet/bubbles/key"

type GlobalKeyMap struct {
	Quit key.Binding
}

var GlobalKeys = GlobalKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q/esc/C-c", "Quit"),
	),
}
