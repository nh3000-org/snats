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
	//	"fmt"
	"log"
	"strings"
	"time"

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
 *	FUNCTION		: messagesScreen
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
	errors := widget.NewLabel("...")

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

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(NatsMessages[id].MSalias + " - " + NatsMessages[id].MSmessage)

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
		nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
		if err != nil {
			log.Println("Messages", err)
		}

		js, _ := nc.JetStream()

		smbutton := widget.NewButton("Send Message", func() {

			//var formatedMessage = FormatMessage("NEW", mymessage.Text,nil)

			js.Publish(strings.ToLower(Queue)+"."+NodeUUID, FormatMessage(mymessage.Text))

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
			NatsMessages = nil
			nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
			if err != nil {
				errors.SetText("Receive Messaged " + err.Error())
				log.Println("messages.go Recieve Messages ", err, " pv ", PasswordValid)
			}

			js, _ := nc.JetStream()
			js.AddStream(&nats.StreamConfig{
				Name: Queue + NodeUUID,

				Subjects: []string{strings.ToLower(Queue) + ".>"},
			})
			var duration time.Duration = 604800000000
			ac, err1 := js.AddConsumer(Queue, &nats.ConsumerConfig{
				Durable:           NodeUUID,
				AckPolicy:         nats.AckExplicitPolicy,
				InactiveThreshold: duration,
				DeliverPolicy:     nats.DeliverAllPolicy,
				ReplayPolicy:      nats.ReplayInstantPolicy,
			})
			if err1 != nil {
				log.Println("messages.go AddConsumer ", err1, " ", ac)
				errors.SetText("Add Consumer " + err1.Error())
			}
			sub, errsub := js.PullSubscribe("", "", nats.BindStream(Queue))
			if errsub != nil {
				log.Println("messages.go PullSubscribe Sub ", errsub)
				errors.SetText("PullSubscribe Sub " + errsub.Error())
			}
			msgs, errfetch := sub.Fetch(100)
			if errfetch != nil {
				errors.SetText("PullSubscribe Fetch " + errfetch.Error())
				log.Println("messages.go PullSubscribe Fetch ", errfetch)
			}
			log.Println("messages: ", len(msgs))
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
			ErrorMessage = "Please Logon First"
			//ErrorScreen(TopWindow)
		}
		bottombox := container.NewBorder(

			recbutton,
			errors,
			nil,
			nil,
			nil,
		)
		return container.NewBorder(

			topbox,
			bottombox,
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
	var inmap = true // unique message id
	ejson, _ := Decrypt(string(m.Data), MySecret)
	err := json.Unmarshal([]byte(ejson), &ms)
	if err != nil {
		log.Println("HandleMessage Unmarshall: ", err)
	}

	inmap = NodeMap("MI" + ms.MSiduuid)
	if inmap == false {
		NatsMessages = append(NatsMessages, ms)
	}

}

/*
 *	FUNCTION		: NodeMap
 *	DESCRIPTION		:
 *		This function returns true if present
 */*
2
 *      PROGRAM         : messages.go
3
 *      DESCRIPTION             :
4
 *
5
 *              This program handles setting and recieving messages
6
 *
7
 *      PARAMETERS              :
8
  *
9
 *      RETURNS                 :
10
 *              Canvas
11
*/
12
package panes
13
​
14
import (
15
        //"encoding/json"
16
        //      "fmt"
17
        "log"
18
        "strings"
19
​
20
        "fyne.io/fyne/v2"
21
        "fyne.io/fyne/v2/container"
22
        "github.com/goccy/go-json"
23
        "github.com/nats-io/nats.go"
24
​
25
        //      "fyne.io/fyne/v2/dialog"
26
        "fyne.io/fyne/v2/theme"
27
        "fyne.io/fyne/v2/widget"
28
)
29
​
30
var EncMessage MessageStore   // message store
31
const QueueCheckInterval = 30 // check interval in seconds
32
/*
33
 *      FUNCTION                : messagesScreen
34
 *      DESCRIPTION             :
35
 *              This function returns a message window
 *	PARAMETERS		: action + node  to lookup
 *                    MI + IDuuid for message id
 *                    AL + Alias for user id
 *
 *
 *	RETURNS			: Array of indexes
 *
 */
func NodeMap(node string) bool {

	_, x := MyMap[node]

	return x

}
