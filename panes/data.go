package panes

import (
	"fyne.io/fyne/v2"
)

// Tutorial defines the data structure for a tutorial
type Panes struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var (
	// Tutorials defines the metadata for each tutorial
	Panes = map[string]Tutorial{
		"logon": {"Logon", "", logonScreen, true},
    "messages": {"Messages", "", messageScreen, true},
		
		"apptabs": {"AppTabs",
			"A container to help divide up an application into functional areas.",
			makeAppTabsTab,
			true,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"logon", "messages"},

	}
)
