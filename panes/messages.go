/*
 *	PROGRAM		: messages.go
 *	DESCRIPTION		:
 *
 *		This program handles setting and recieving messages
 *
 *	PARAMETERS		:
  *
 *	RETURNS			:
 *		Canvas
*/
package panes

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/nats-io/nats.go"

	//	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var EncMessage MessageStore   // message store
const QueueCheckInterval = 30 // check interval in seconds
/*
 *	FUNCTION		: messagesScren
 *	DESCRIPTION		:
 *		This function returns a message window
 *
 *	PARAMETERS		:
 *
 *
 *	RETURNS			:
 *
 */
func messagesScreen(_ fyne.Window) fyne.CanvasObject {

	//SaveCarootToFS()
	mymessage := widget.NewMultiLineEntry()
	mymessage.SetPlaceHolder("Enter Message For Encryption")
	mymessage.SetMinRowsVisible(5)

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	hbox := container.NewHBox(icon, label)
	List := widget.NewList(
		func() int {
			return len(NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var prefix = ""
			if !NodeIsValid(NatsMessages[id].MSnodeuuid) {
				prefix = "[Unknown]"
			}

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(prefix + NatsMessages[id].MSalias)

		},
	)
	List.OnSelected = func(id widget.ListItemID) {
		var mytext = NatsMessages[id].MSmessage + "\n" + NatsMessages[id].MShostname + "\n" + NatsMessages[id].MSipadrs + "\n" + NatsMessages[id].MSnodeuuid
		label.SetText(mytext)
		icon.SetResource(theme.DocumentIcon())
	}
	List.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()
	if UseJetstream == true {

		NC, err := nats.Connect(Server, nats.RootCAs(DataStore("ca-nats.pem").Path()))
		if err != nil {
			log.Println("natsconnect", err, " pv ", PasswordValid)
		}
		c, errenc := nats.NewEncodedConn(NC, nats.JSON_ENCODER)
		if errenc != nil {
			log.Println("natsconnect enc ", errenc)
		}
		_, errqs := c.QueueSubscribe(Queue, Queue, func(msg MessageStore) {

			HandleMessage(msg)

			List.Refresh()
		})
		if errqs != nil {
			log.Println("natsconnect", errqs)
		}
		if err == nil && errenc == nil {
			//wg.Add(10)

			//qcnt, qerr := q.Fetch(100)
			//log.Println("queue fetch count ", qcnt)

			if errqs != nil {
				//log.Println("queue fetch ", qerr)
			}

		}

	}
	// try the password
	smbutton := widget.NewButton("Send Message", func() {

		EncMessage = FormatMessage(mymessage.Text)
		//AddMessage()
		//log.Println("messagesScreen publish" + mymessage.Text)
		NATSPublish(EncMessage)

	})

	topbox := container.NewBorder(

		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		smbutton,
		nil,
		nil,
		mymessage,
	)
	// save the server
	recbutton := widget.NewButton("Recieve Messages", func() {

		//NATSErase()
		//go NATSConnect()

	})

	if !LoggedOn {
		mymessage.Disable()
		smbutton.Disable()
		recbutton.Disable()
		ErrorMessage = "Please Logon"
		//ErrorScreen(TopWindow)
	}
	return container.NewBorder(

		topbox,
		recbutton,
		nil,
		nil,

		//widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
		// natsmessages is message q
		//		container.NewVScroll(
		container.NewHSplit(List, container.NewCenter(hbox)),
	)
	//)

}
func HandleMessage(m MessageStore) {

	NatsMessages = append(NatsMessages, m)
}
func NodeIsValid(node string) bool {
	return false

}
