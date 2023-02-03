/*
 *	PROGRAM		: settings.go
 *	DESCRIPTION		:
 *
 *		This program handles loging on
 *
 *	PARAMETERS		:
  *
 *	RETURNS			:
 *		Canvas
*/

package panes

import (
	"log"

	"os"
	"os/exec"

	"fyne.io/fyne/v2"

	"golang.org/x/crypto/bcrypt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
 *	FUNCTION		: logonScren
 *	DESCRIPTION		:
 *		This function returns a logonwindow
 *
 *	PARAMETERS		:
 *
 *
 *	RETURNS			:
 *
 */
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

		if !iserrors {
			PasswordValid = true
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

	SSbutton := widget.NewButton("Logon", func() {
		if PasswordValid {
			var iserrors bool
			iserrors = false
			if !iserrors == false {
				iserrors = editEntry("URL", server.Text)
			}
			if !iserrors == false {
				iserrors = editEntry("CERTIFICATE", caroot.Text)
			}
			if !iserrors && PasswordValid {

				uuid, err := exec.Command("uuidgen").Output()
				if err != nil {
					log.Fatal(err)
				}
				NodeUUID = string(uuid)

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

				// dont disable and alLoggedOnlow for multiple queue entries
				//queue.Disable()
				//queuepassword.Disable()

				MyJson("SAVE")

				go NATSConnect()

				LoggedOn = true
			}
		}
	})
	// security erase
	SEbutton := widget.NewButton("Security Erase", func() {
		if PasswordValid {
			NATSErase()
			os.Exit(1)
		}

	})

	// check for logon
	if !PasswordValid {
		password.Enable()
		server.Disable()
		caroot.Disable()
		alias.Disable()
		queue.Disable()
		queuepassword.Disable()

	}
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		alias,
		server,
		caroot,
		queue,
		queuepassword,
		SSbutton,
		SEbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", parseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
