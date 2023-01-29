package panes

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var EncMessage MessageStore   // message store
const QueueCheckInterval = 30 // check interval in seconds


func messagesScreen(_ fyne.Window) fyne.CanvasObject {

	log.Println("messagesScreen")
	SaveCarootToFS()
	mymessage := widget.NewEntry()
	mymessage.SetPlaceHolder("Enter Message For Encryption")

	// try the password
	smbutton := widget.NewButton("Send Message", func() {

		EncMessage = FormatMessage(mymessage.Text)
		//AddMessage()
		log.Println("messagesScreen publish" + mymessage.Text)
		NATSPublish(EncMessage)

	})

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(icon, label)
	List := widget.NewList(
		func() int {
			log.Println("list size" + strconv.Itoa(len(NatsMessages)))
			return len(NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(NatsMessages[id].MSalias)

			//item.(*fyne.Container).Objects[1].(*widget.Label).SetText(NatsMessages[id])

		},
	)
	List.OnSelected = func(id widget.ListItemID) {
		var mytext = NatsMessages[id].MSmessage + "\n" + NatsMessages[id].MShostname + "\n" + NatsMessages[id].MSipadrs
		label.SetText(mytext)
		icon.SetResource(theme.DocumentIcon())
	}
	List.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	//list.Select(125)

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()

	vertbox := container.NewVBox(

		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		mymessage,

		smbutton,
	)
	// save the server
	sebutton := widget.NewButton("Security Erase", func() {
		NATSErase()
		NATSConnect()
	})
	if !LoggedOn {
		mymessage.Disable()
		smbutton.Disable()
		sebutton.Disable()
	}
	return container.NewBorder(
		//return container.NewCenter(container.NewVBox(
		//return container.NewCenter(container.NewGridWithRows(
		//widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		//mymessage,

		//smbutton,
		vertbox,
		sebutton,
		nil,
		nil,

		//widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
		// natsmessages is message q
		//		container.NewVScroll(
		container.NewHSplit(List, container.NewCenter(hbox)),
	)
	//)
}
