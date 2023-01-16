package panes

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var EncMessage string         // message strung
const QueueCheckInterval = 30 // check interval in secinds

func messagesScreen(_ fyne.Window) fyne.CanvasObject {

	log.Println("messagesScreen")

	mymessage := widget.NewEntry()
	mymessage.SetPlaceHolder("Enter Message For Encryption")

	// try the password
	smbutton := widget.NewButton("Send Message", func() {

		EncMessage += FormatMessage(mymessage.Text)
		// pub the message to queue

	})

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(icon, label)
	list := widget.NewList(
		func() int {
			return len(NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			if id == 5 || id == 6 {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(NatsMessages[id] + "\ntaller")
			} else {
				item.(*fyne.Container).Objects[1].(*widget.Label).SetText(NatsMessages[id])
			}
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(NatsMessages[id])
		icon.SetResource(theme.DocumentIcon())
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)
	//list.SetItemHeight(5, 50)
	//list.SetItemHeight(6, 50)
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		mymessage,

		smbutton,

		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
		// natsmessages is message q
		//		container.NewVScroll(
		container.NewHSplit(list, container.NewCenter(hbox)),
	))
	//)
}
