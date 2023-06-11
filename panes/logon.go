package panes

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func logonScreen(MyWin fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

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

	queue := widget.NewEntry()
	queue.SetPlaceHolder("Message Queue for Pub/Sub")
	queue.Disable()

	queuepassword := widget.NewEntry()
	queuepassword.SetPlaceHolder("Message Queue Password")
	queuepassword.Disable()

	tpbutton := widget.NewButton("Try Password", func() {
		errors.SetText("...")
		var iserrors = false
		Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		Passwordhash = string(pwh)
		if err != nil {
			iserrors = true
			log.Println("logon.go Error Creating Password Hash ")
			errors.SetText(err.Error())
		}
		confighasherr, _ := storage.Exists(DataStore("config.hash"))
		//log.Println("hash logon ", confighasherr, " at ", DataStore("config.hash"))
		if confighasherr == false {
			if MyHash("CREATE") {
				log.Println("logon.go Error Creating Password Hash")
				errors.SetText("Error Creating Password Hash")
			}
		}

		if MyHash("LOAD") {
			log.Println("logon.go Error Loading Password Hash")
			errors.SetText("Error Loading Password Hash")
		}
		// Comparing the password with the hash
		errpw := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password))

		if errpw != nil {
			iserrors = true
			log.Println("logon.go Error Invalid Password")
			errors.SetText("Error Invalid Password")
		}

		if !iserrors {
			errors.SetText("...")
			PasswordValid = true
			MyJson("LOAD")
			alias.SetText(Alias)
			server.SetText(Server)
			//log.Println("logon.go Server " + Server)
			queue.SetText(Queue)
			queuepassword.SetText(Queuepassword)
			password.Disable()
			server.Enable()
			queue.Enable()
			alias.Enable()
			queuepassword.Enable()
		}
	})

	SSbutton := widget.NewButton("Logon", func() {

		var iserrors = editEntry("URL", server.Text)
		if iserrors == true {
			log.Println("logon.go Error URL Incorrect Format ")
			errors.SetText("Error Invalid Password")
		}
		iserrors = editEntry("STRING", queuepassword.Text)
		if iserrors == true {
			errors.SetText("Error Invalid Queue Password")
			iserrors = true
		}
		if len(queuepassword.Text) != 24 {
			iserrors = true
			errors.SetText("Error Queue Password Length is " + strconv.Itoa(len(queuepassword.Text)) + " shlould be length of 24")
		}

		if !iserrors && PasswordValid {
			NodeUUID = uuid.New().String()
			Alias = alias.Text
			Server = server.Text
			Queue = queue.Text
			Queuepassword = queuepassword.Text
			password.Disable()
			server.Disable()

			alias.Disable()
			queue.Disable()
			queuepassword.Disable()

			MyJson("SAVE")

			LoggedOn = true

			errors.SetText("...")
		}

	})
	// security erase
	SEbutton := widget.NewButton("Security Erase", func() {
		if PasswordValid {
			NATSErase()
		}

	})

	if !PasswordValid {
		password.Enable()
		server.Disable()
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
		queue,
		queuepassword,
		SSbutton,
		SEbutton,
		errors,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", parseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}
