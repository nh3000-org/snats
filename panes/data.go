package panes

import (
	"fyne.io/fyne/v2"
)

// Pane defines the data structure
type MyPane struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Panes defines the metadata
	MyPanes = map[string]MyPane{
		"logon":    {"Logon", "", logonScreen, true},
		"messages": {"Messages", "", messagesScreen, true},
	}

	// PanesIndex  defines how our panes should be laid out in the index tree
	MyPanesIndex = map[string][]string{
		"": {"logon", "messages"},
	}
)
