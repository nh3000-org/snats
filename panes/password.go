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
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/bcrypt"
)

/*
 *	FUNCTION		: passwordScren
 *	DESCRIPTION		:
 *		This function returns a settings window
 *
 *	PARAMETERS		:
 *
 *
 *	RETURNS			:
 *
 */
func passwordScreen(_ fyne.Window) fyne.CanvasObject {

	_, configfileerr := os.Stat("config.json")
	if configfileerr != nil {

		MyJson("CREATE")
	}
	MyJson("LOAD")
	password := widget.NewPasswordEntry()

	password.SetPlaceHolder("Enter Password For Encryption")
	password.SetText(Password)

	passwordc1 := widget.NewPasswordEntry()
	passwordc1.SetPlaceHolder("Enter Change Password")
	passwordc1.SetText(Password)
	passwordc1.Disable()
	passwordc2 := widget.NewPasswordEntry()
	passwordc2.SetPlaceHolder("Enter Change Password Again")
	passwordc1.SetText(Password)
	passwordc2.Disable()

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
			MyJson("LOAD")
			log.Println("save password config", Password)
			log.Println("save password gui", password.Text)

			password.Disable()
			passwordc1.Enable()
			passwordc2.Enable()

		}

	})

	cpbutton := widget.NewButton("Change Password", func() {
		var iserrors bool
		iserrors = false

		if editEntry("STRING", passwordc1.Text) != true {
			iserrors = true
		}
		if editEntry("STRING", passwordc2.Text) != true {
			iserrors = true
		}
		if passwordc1.Text != passwordc2.Text != true {
			iserrors = true
		}
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
			Passwordhash = string(pwh)
			if err != nil {
				log.Fatal(err)
			}
			_, confighasherr := os.Stat("config.hash")
			if confighasherr == nil {

				MyHash("CREATE", Passwordhash)
				MyJson("SAVE")
			}
		}
		Password = passwordc1.Text
		MyHash("LOAD", "NONE")
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {
			// TODO: Properly handle error
			iserrors = true
		}
		if !iserrors {
			MyJson("SAVE")
			log.Println("save password config", Password)
			log.Println("save password gui", password.Text)

		}

	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Enter Password To Reset", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		tpbutton,
		passwordc1,
		passwordc2,
		cpbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", parseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
