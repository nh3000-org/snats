package panes

import (
	"log"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/goccy/go-json"
	"github.com/nats-io/nats.go"
)

var EncMessage MessageStore   // message store
const QueueCheckInterval = 30 // check interval in seconds

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
				}

			}
			List.Refresh()
		})

		if !LoggedOn {
			mymessage.Disable()
			smbutton.Disable()
			recbutton.Disable()
			ErrorMessage = "Please Logon First"
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
			container.NewHSplit(List, container.NewCenter(hbox)),
		)

	}
	return container.NewBorder(

		widget.NewLabelWithStyle("Logon First", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
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

func NodeMap(node string) bool {
	_, x := MyMap[node]
	return x
}
