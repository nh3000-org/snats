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
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var prefix = ""
			if !NodeIsValid(NatsMessages[id].MSnodeuuid) {
				prefix = "-"
			}

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(prefix + NatsMessages[id].MSalias + " - " + NatsMessages[id].MShostname)

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
		js.AddStream(&nats.StreamConfig{
			Name:     Queue,
			Subjects: []string{Queue},
		})
		// Create a Consumer
		ac, se7 := js.AddConsumer(Queue, &nats.ConsumerConfig{
			Durable:       Queue,
			AckPolicy:     nats.AckExplicitPolicy,
			DeliverPolicy: nats.DeliverAllPolicy,
			//		ReplayPolicy: nats.ReplayInstantPolicy,
		})
		if se7 != nil {
			log.Println("natsconnect se7 ", se7, ac)
		}
		log.Println("natsconnect se8 ")
		sub, errsub := js.PullSubscribe(Queue, Queue, nats.BindStream(Queue))
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
				fmt.Printf("fetch message %d  ", msgs[i].Data)
			}

		}
		//sub1, errsub := js.Subscribe(Queue, func (msg nats.Msg) {

		//		})
		//		nc.Close()
		//ephemeralName := <-js.ConsumerNames(Queue)
		//fmt.Printf("ephemeral name is %q\n", ephemeralName)

		smbutton := widget.NewButton("Send Message", func() {

			var formatMessage = FormatMessage(mymessage.Text)

			js.Publish(Queue, []byte(formatMessage))

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
				Subjects: []string{Queue},
			})
			// Create a Consumer
			ac, se7 := js.AddConsumer(Queue, &nats.ConsumerConfig{
				Durable:       Queue,
				AckPolicy:     nats.AckExplicitPolicy,
				DeliverPolicy: nats.DeliverAllPolicy,
				//		ReplayPolicy: nats.ReplayInstantPolicy,
			})
			if se7 != nil {
				log.Println("Recieve Messages AddConsumer ", se7, ac)
			}

			_, errsub := js.QueueSubscribe(Queue, Alias, func(msg *nats.Msg) {
				log.Println("Recieve Messages Receive ", msg.Data)
				HandleMessage(msg)
				msg.Ack()

			})
			if errsub != nil {
				log.Println("Recieve Messages Subscribe ", errsub)
			}

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
	mm := MessageStore{}
	err := json.Unmarshal([]byte(m.Data), &mm)
	if err != nil {
		log.Println("HandleMessage Unmarshall: ", err)
	}
	NatsMessages = append(NatsMessages, mm)
}
func NodeIsValid(node string) bool {
	if node == NodeUUID {
		return true
	}

	return false

}
