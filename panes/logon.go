/*
 *	PROGRAM		: logon.go
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

	"github.com/google/uuid"

	"fyne.io/fyne/v2"

	"golang.org/x/crypto/bcrypt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
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
func logonScreen(MyWin fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")
	configbool, configboolerr := storage.Exists(DataStore("config.json"))
	if configboolerr != nil {
		log.Println(configboolerr)
	}

	if configbool == false {

		if MyJson("CREATE") {
			log.Println("logon.go Error Creating Password Hash ")
			errors.SetText("Error Creating Password Hash ")
		}

	}

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

	// try the password
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

			if MyHash("CREATE", Passwordhash) {
				log.Println("logon.go Error Creating Password Hash")
				errors.SetText("Error Creating Password Hash")
			}
		}

		if MyHash("LOAD", "NONE") {
			log.Println("logon.go Error Loading Password Hash")
			errors.SetText("Error Loading Password Hash")
		}
		// Comparing the password with the hash
		errpw := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password))
		// TODO: Properly handle error
		if errpw != nil {
			iserrors = true
			log.Println("logon.go Error Invalid Password")
			errors.SetText("Error Invalid Password")
		}

		if !iserrors {
			errors.SetText("...")
			PasswordValid = true
			if MyJson("LOAD") {
				iserrors = true
				log.Println("logon.go Error Cannot Load JSON")
				errors.SetText("Error Cannot Load JSON")
			}
			alias.SetText(Alias)
			server.SetText(Server)
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

		//var iserrors bool = false

		var iserrors = editEntry("URL", server.Text)
		if iserrors == true {
			log.Println("logon.go Error URL Incorrect Format ")

			errors.SetText("Error Invalid Password")

		}

		if !iserrors && PasswordValid {

			NodeUUID = uuid.New().String()

			Alias = alias.Text
			Server = server.Text
			//Caroot = caroot.Text
			//Clientcert = clientcert.Text
			//Clientkey = clientkey.Text
			Queue = queue.Text
			Queuepassword = queuepassword.Text
			password.Disable()
			server.Disable()
			//caroot.Disable()
			//clientcert.Disable()
			//clientkey.Disable()
			alias.Disable()
			queue.Disable()
			queuepassword.Disable()

			// dont disable and alLoggedOnlow for multiple queue entries
			//queue.Disable()
			//queuepassword.Disable()

			if MyJson("SAVE") {

				log.Println("logon.go Error Cannot Save JSON")
				errors.SetText("Error Cannot Save JSON")
			}

			LoggedOn = true/*
2
 *      PROGRAM         : settings.go
3
 *      DESCRIPTION             :
4
 *
5
 *              This program handles loging on
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
​
13
package panes
14
​
15
import (
16
        "log"
17
        "os/exec"
18
​
19
        "os"
20
        //"os/exec"
21
​
22
        "fyne.io/fyne/v2"
23
​
24
        "golang.org/x/crypto/bcrypt"
25
​
26
        "fyne.io/fyne/v2/container"
27
        "fyne.io/fyne/v2/storage"
28
        "fyne.io/fyne/v2/widget"
29
)
30
​
31
/*
32
 *      FUNCTION                : logonScren
33
 *      DESCRIPTION             :
34
 *              This function returns a logonwindow
35
 *
36
 *      PARAMETERS              :
37
 *
38
 *
39
 *      RETURNS                 :
40
 *
41
 */
42
func logonScreen(_ fyne.Window) fyne.CanvasObject {
43
​
44
        configbool, _ := storage.Exists(DataStore("config.json"))
45
        if configbool == false {
46
​
47
                MyJson("CREATE")
48
​
49
        }
			//log.Println("loggedOn ", LoggedOn)
			errors.SetText("...")
		}

	})
	// security erase
	SEbutton := widget.NewButton("Security Erase", func() {
		if PasswordValid {
			NATSErase()

		}

	})

	// check for logon
	if !PasswordValid {
		password.Enable()
		server.Disable()
		//caroot.Disable()
		//clientcert.Disable()
		//clientkey.Disable()
		alias.Disable()
		queue.Disable()
		queuepassword.Disable()

	}
	//return container.NewCenter(container.NewVBox(
/*
2
 *      PROGRAM         : settings.go
3
 *      DESCRIPTION             :
4
 *
5
 *              This program handles loging on
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
​
13
package panes
14
​
15
import (
16
        "log"
17
        "os/exec"
18
​
19
        "os"
20
        //"os/exec"
21
​
22
        "fyne.io/fyne/v2"
23
​
24
        "golang.org/x/crypto/bcrypt"
25
​
26
        "fyne.io/fyne/v2/container"
27
        "fyne.io/fyne/v2/storage"
28
        "fyne.io/fyne/v2/widget"
29
)
30
​
31
/*
32
 *      FUNCTION                : logonScren
33
 *      DESCRIPTION             :
34
 *              This function returns a logonwindow
35
 *
36
 *      PARAMETERS              :
37
 *
38
 *
39
 *      RETURNS                 :
40
 *
41
 */
42
func logonScreen(_ fyne.Window) fyne.CanvasObject {
43
​
44
        configbool, _ := storage.Exists(DataStore("config.json"))
45
        if configbool == false {
46
​
47
                MyJson("CREATE")
48
​
49
        }
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		alias,
		server,
		//caroot,
		//clientcert,
		//clientkey,
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
