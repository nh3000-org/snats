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

			if MyHash("CREATE") {
				errors.SetText("Error Creating Password Hash")
			}
		}

		Password = password.Text
		if MyHash("LOAD") {
			errors.SetText("Error Reading Password Hash")
		}
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {

			iserrors = true
			errors.SetText(err.Error())
		}
		if !iserrors {
			MyJson("LOAD")
			errors.SetText("Password Accepted ")

			password.Disable()
			passwordc1.Enable()
			passwordc2.Enable()

		}
	})

	cpbutton := widget.NewButton("Change Password", func() {
		var iserrors = false

		if editEntry("STRING", passwordc1.Text) == true {
			iserrors = true
			errors.SetText("Error Pasword 1 Invalid")
		}

		log.Println("password.go iserrors b4 pass ", iserrors)
		if editEntry("PASSWORD", passwordc1.Text) {
			iserrors = true
			errors.SetText("Error Pasword Does Not Meet Requirements")
		}
		if passwordc1.Text != passwordc2.Text {
			iserrors = true
			errors.SetText("Error Pasword 1 Does Not Match Password 2")
		}
		log.Println("password.go iserrors", iserrors)
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(passwordc1.Text), bcrypt.DefaultCost)
			Passwordhash = string(pwh)
			log.Println("password hash", Passwordhash)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("change save ", Passwordhash, " ", DataStore("config.hash"))
			if MyHash("SAVE") {
				errors.SetText("Error Saving Password Hash")
				log.Println("password.go Error Saving Password Hash ")
			}

		}
		Password = passwordc1.Text
		if MyHash("LOAD") {
			errors.SetText("Error Reading Password Hash")
			log.Println("password.go Error Reading Password Hash ")
			iserrors = true
		}
		// Comparing the password with the hash
		err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(passwordc1.Text))
		if err != nil {
			log.Println("Error Invalid Password ", err)
			errors.SetText("password.go Error Invalid Password ")
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
