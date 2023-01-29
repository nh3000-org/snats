/*
* Modify cipherkey for your installation
 */

package panes

import (
	"log"

	"os"

	"fyne.io/fyne/v2"

	"golang.org/x/crypto/bcrypt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func logonScreen(_ fyne.Window) fyne.CanvasObject {

	_, configfileerr := os.Stat("config.json")
	if configfileerr != nil {

		MyJson("CREATE")
	}
	MyJson("LOAD")

	password := widget.NewEntry()
	password.SetPlaceHolder("Enter Password For Encryption")
	password.SetText(Password)
	alias := widget.NewEntry()
	alias.SetPlaceHolder("Enter User Alias")
	alias.SetText(Alias)
	alias.Disable()
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
	// try the password
	tpbutton := widget.NewButton("Try Password", func() {
		var iserrors bool
		iserrors = false
		Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		Passwordhash = string(pwh)
		if err != nil {
			log.Fatal(err)
		}
		_, confighasherr := os.Stat("config.hash")
		if confighasherr != nil {

			MyHash("CREATE", Passwordhash)
		}

		Password = password.Text
		MyHash("LOAD", "NONE")
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {
			// TODO: Properly handle error
			iserrors = true
		}
		GoodPassword = true
		if !iserrors {
			MyJson("LOAD")
			alias.SetText(Alias)
			server.SetText(Server)
			caroot.SetText(Caroot)
			queue.SetText(Queue)
			queuepassword.SetText(Queuepassword)
			password.Disable()
			server.Enable()
			caroot.Enable()
			queue.Enable()
			alias.Enable()
			queuepassword.Enable()
		}

	})
	// save the server
	if !LoggedOn {
		password.Enable()
		server.Disable()
		caroot.Disable()
		alias.Disable()
		queue.Disable()
		queuepassword.Disable()

	}
	ssbutton := widget.NewButton("Logon", func() {
		var iserrors bool
		iserrors = false
		if !iserrors == false {
			iserrors = editEntry("URL", server.Text)
		}
		if !iserrors == false {
			iserrors = editEntry("CERTIFICATE", caroot.Text)
		}
		if !iserrors && GoodPassword {
			Alias = alias.Text
			Server = server.Text
			Caroot = caroot.Text
			Queue = queue.Text
			Queuepassword = queuepassword.Text
			password.Disable()
			server.Disable()
			caroot.Disable()
			alias.Disable()
			queue.Disable()
			queuepassword.Disable()
			// dont disable and allow for multiple queue entries
			//queue.Disable()
			//queuepassword.Disable()

			MyJson("SAVE")

			go NATSConnect()

			LoggedOn = true
		}
	})
	// security erase
	sebutton := widget.NewButton("Security Erase", func() {
		NATSErase()
		NATSConnect()
	})
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		alias,
		server,
		caroot,
		queue,
		queuepassword,
		ssbutton,
		sebutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.og", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
