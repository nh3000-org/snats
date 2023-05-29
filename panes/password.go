package panes

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/bcrypt"
)

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
	//passwordc1.SetText(Password)
	passwordc1.Disable()

	passwordc2 := widget.NewPasswordEntry()
	passwordc2.SetPlaceHolder("Enter Change Password Again")
	//passwordc1.SetText(Password)
	passwordc2.Disable()
	errors := widget.NewLabel("...")
	// try the password
	tpbutton := widget.NewButton("Try Password", func() {
		var iserrors = false
		Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		Passwordhash = string(pwh)
		if err != nil {
			iserrors = true
			log.Println("Error Creating Password Hash ", err)
			errors.SetText(err.Error())
		}
		_, confighasherr := os.Stat("config.hash")
		if confighasherr != nil {

			if MyHash("CREATE", Passwordhash) {
				errors.SetText("Error Creating Password Hash")
			}
		}

		Password = password.Text
		if MyHash("LOAD", "NONE") {
			errors.SetText("Error Reading Password Hash")
		}
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {
			// TODO: Properly handle error
			iserrors = true
			errors.SetText(err.Error())
		}
		if !iserrors {
			MyJson("LOAD")
			errors.SetText("Password Changed ")

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
			errors.SetText("Error Pasword 1 Invalid")
		}
		if editEntry("STRING", passwordc2.Text) != true {
			iserrors = true
			errors.SetText("Error Pasword 2 Invalid")
		}
		if passwordc1.Text != passwordc2.Text != true {
			iserrors = true
			errors.SetText("Error Pasword 1 Dows Not Password 2")
		}
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
			Passwordhash = string(pwh)
			if err != nil {
				log.Fatal(err)
			}
			_, confighasherr := os.Stat("config.hash")
			if confighasherr == nil {

				if MyHash("CREATE", Passwordhash) {
					errors.SetText("Error Creating Password Hash")
				}
				MyJson("SAVE")
			}
		}
		Password = passwordc1.Text
		if MyHash("LOAD", "NONE") {
			errors.SetText("Error Reading Password Hash")
			iserrors = true
		}
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {
			errors.SetText("Error Invalid Password ")
			iserrors = true
		}
		if !iserrors {
			MyJson("SAVE")

		}

	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Enter Password To Reset", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("For Local Security", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		tpbutton,
		passwordc1,
		passwordc2,
		cpbutton,
		errors,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", parseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
