/*
 *	PROGRAM		: password.go
 *	DESCRIPTION		:
 *		This program handles password entry, json decription, 
 *      password changing and placing the CAROOT on the file system
 *      from the config.json file.
 *
 *      The caroot is deleted from the file system when program closes.
 *
 *	PARAMETERS		:
  *
 *	RETURNS			:
 *		Fyne window
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

func passwordScreen(_ fyne.Window) fyne.CanvasObject {

	_, configfileerr := os.Stat("config.json")
	if configfileerr != nil {

		myjson("CREATE")
	}
	myjson("LOAD")
	password := widget.NewPasswordEntry()

	password.SetPlaceHolder("Enter Password For Encryption")
	password.SetText(Password)

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
			myjson("LOAD")
			log.Println("save password config", Password)
			log.Println("save password gui", password.Text)

			password.Disable()

		}

	})
	// save the server

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Enter Password", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		tpbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.og", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
