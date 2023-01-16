/*
* Modify cipherkey for your installation
 */

package panes

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsScreen(_ fyne.Window) fyne.CanvasObject {

	_, configfileerr := os.Stat("config.json")
	if configfileerr != nil {

		MyJson("CREATE")
	}
	MyJson("LOAD")

	server := widget.NewEntry()
	server.SetPlaceHolder("URL: nats://xxxxxx:4332")
	server.Disable()
	caroot := widget.NewMultiLineEntry()
	caroot.SetPlaceHolder("CAROOT For nats://xxxxxx:4332")
	caroot.Disable()
	queue := widget.NewEntry()
	queue.SetPlaceHolder("Message Queue for Pub/Sub")
	queue.Disable()
	queuepassword := widget.NewEntry()
	queuepassword.SetPlaceHolder("Message Queue Password")
	queuepassword.Disable()
	MyJson("LOAD")

	server.SetText(Server)
	caroot.SetText(Caroot)
	queue.SetText(Queue)
	queuepassword.SetText(Queuepassword)

	server.Enable()
	caroot.Enable()
	queue.Enable()
	queuepassword.Enable()

	ssbutton := widget.NewButton("Connect To Server", func() {
		var iserrors bool
		iserrors = false
		if !iserrors == false {
			iserrors = editEntry("URL", server.Text)
		}
		if !iserrors == false {
			iserrors = editEntry("CERTIFICATE", caroot.Text)
		}
		if !iserrors {
			Server = server.Text
			Caroot = caroot.Text
			Queue = queue.Text
			Queuepassword = queuepassword.Text

			server.Disable()
			caroot.Disable()
			// dont disable and allow for multiple queue entries
			//queue.Disable()
			//queuepassword.Disable()

			MyJson("SAVE")

			go NATSConnect()
		}
	})
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		server,
		caroot,
		queue,
		queuepassword,
		ssbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.og", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
