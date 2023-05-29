
package main

import (
	//"log"

	"fyne.io/fyne/v2"
	//"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	//"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	//"panes/data"
	"github.com/nh3000-org/snats/panes"
)

const preferenceCurrentApplication = "currentApplication"

var TopWindow fyne.Window

var ErrorMessage = "..."

//var MyApp fyne.App

func main() {
	//MyApp := app.NewWithID("org.nh3000.snats")
	panes.MyJson("LOAD")
	panes.MyAppDup = panes.GetMyApp()
	panes.MyAppDup.SetIcon(theme.FyneLogo())
	makeTray(panes.MyAppDup)
	logLifecycle(panes.MyAppDup)
	w := panes.MyAppDup.NewWindow("Secure NATS BETA")
	TopWindow = w

	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("SNATS")

	intro := widget.NewLabel("Secure Communications using NATS\nVisit nats.io for additional info.")
	intro.Wrapping = fyne.TextWrapWord
	setTutorial := func(t panes.MyPane) {
		if fyne.CurrentDevice().IsMobile() {
			child := panes.MyAppDup.NewWindow(t.Title)
			TopWindow = child
			child.SetContent(t.View(TopWindow))
			child.Show()
			child.SetOnClosed(func() {
				TopWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(setTutorial, false))
	} else {
		split := container.NewHSplit(makeNav(setTutorial, true), tutorial)
		split.Offset = 0.2
		w.SetContent(split)
	}

	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {

	a.Lifecycle().SetOnStopped(func() {

		//log.Println("DeleteCarootFS Deleting")

	})

}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Secure", func() {})
		menu := fyne.NewMenu("Encryption", h)
		h.Action = func() {
			//log.Println("System tray menu tapped")
			h.Label = "Secure"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}


func unsupportedApplication(t panes.MyPane) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}


func makeNav(setTutorial func(panes panes.MyPane), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return panes.MyPanesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := panes.MyPanesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := panes.MyPanes[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedApplication(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := panes.MyPanes[uid]; ok {
				if unsupportedApplication(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentApplication, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentApplication, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}
func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
