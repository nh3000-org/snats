package panes

import (
	"log"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var Encmessage string // message strung
func messagesScreen(_ fyne.Window) fyne.CanvasObject {

	log.Println("messagesScreen")
	mymessage := widget.NewEntry()
	mymessage.SetPlaceHolder("Enter Message For Encryption")

	// try the password
	smbutton := widget.NewButton("Send Message", func() {

		Encmessage = mymessage.Text

	})
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		mymessage,

		smbutton,

		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}
