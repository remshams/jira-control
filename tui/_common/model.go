package common

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/remshams/common/tui/bubbles/tabs"
)

type GlobalKeyMap struct {
	Tab  tabs.TabKeyMap
	Quit key.Binding
}

var GlobalKeys = GlobalKeyMap{
	Tab: tabs.TabKeys,
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/C-c", "Quit"),
	),
}

func (k GlobalKeyMap) KeyBindings() []key.Binding {
	return []key.Binding{
		k.Tab.Tab,
		k.Quit,
	}
}
