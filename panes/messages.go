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
	//"encoding/json"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/goccy/go-json"
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
			return container.NewHBox(widget.NewIcon(theme.CheckButtonCheckedIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var prefix = ""
			if !NodeIsValid(NatsMessages[id].MSnodeuuid) {
				prefix = "-"
			}

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(prefix + NatsMessages[id].MSalias + " - " + NatsMessages[id].MSmessage)

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
	if PasswordValid == true {
		//nc, err := nats.Connect(Server, nats.RootCAs(DataStore("ca-root.pem").Path()), nats.ClientCert(DataStore("client-cert.pem").Path(), DataStore("client-key.pem").Path()))
		nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
		if err != nil {
			log.Println("natsconnect", err, " pv ", PasswordValid)
		}

		js, _ := nc.JetStream()

		smbutton := widget.NewButton("Send Message", func() {

			var formatMessage = FormatMessage(mymessage.Text)

			js.Publish(strings.ToLower(Queue)+"."+NodeUUID, []byte(formatMessage))

		})

		topbox := container.NewBorder(

			widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			smbutton,
			nil,
			nil,
			mymessage,
		)
		// recieve messages
		recbutton := widget.NewButton("Recieve Messages", func() {
			//nc, err := nats.Connect(Server, nats.RootCAs(DataStore("ca-root.pem").Path()), nats.ClientCert(DataStore("client-cert.pem").Path(), DataStore("client-key.pem").Path()))
			nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
			if err != nil {
				log.Println("Recieve Messages ", err, " pv ", PasswordValid)
			}

			js, _ := nc.JetStream()
			js.AddStream(&nats.StreamConfig{
				Name:     Queue,
				Subjects: []string{strings.ToLower(Queue) + NodeUUID},
			})

			sub, errsub := js.PullSubscribe(strings.ToLower(Queue)+".*", "", nats.BindStream(Queue))
			if errsub != nil {
				log.Println("pullsubscribe sub ", errsub)
			}
			msgs, errfetch := sub.Fetch(100)
			if errfetch != nil {
				log.Println("pullsubscribe fetch ", errfetch)
			}
			fmt.Printf("got %d messages\n", len(msgs))
			if len(msgs) > 0 {

				for i := 0; i < len(msgs); i++ {
					msgs[i].Nak()

					HandleMessage(msgs[i])
					//fmt.Printf("fetch message %d  ", msgs[i].Data)
				}

			}
			List.Refresh()
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
	return container.NewBorder(

		widget.NewLabelWithStyle("Logon First", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,

		//widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
		// natsmessages is message q
		//		container.NewVScroll(
		container.NewHSplit(List, container.NewCenter(hbox)),
	)
}
func HandleMessage(m *nats.Msg) {
	ms := MessageStore{}
	var unq = true // unique message id
	ejson, _ := Decrypt(string(m.Data), MySecret)
	err := json.Unmarshal([]byte(ejson), &ms)
	if err != nil {
		log.Println("HandleMessage Unmarshall: ", err)
	}

	for x := 0; x < len(NatsMessages); x++ {
		log.Println("HandleMessage store ", ms.MSiduuid, " - ", ms.MSmessage, " messages ", NatsMessages[x].MSiduuid, " - ", NatsMessages[x].MSmessage)
		if ms.MSiduuid == NatsMessages[x].MSiduuid {
			unq = false
		}
	}
	if unq == true {
		NatsMessages = append(NatsMessages, ms)
	}

}
func NodeIsValid(node string) bool {
	if node == NodeUUID {
		return true
	}

	return false

}
