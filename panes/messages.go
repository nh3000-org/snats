package panes

import (


	"log"


	"fyne.io/fyne/v2"


	"fyne.io/fyne/v2/widget"

)

func messagesScreen(_ fyne.Window) {

	log.Println("messagesScreen")
	password := widget.NewEntry()
	password.SetPlaceHolder("Enter Password For Encryption")
	password.SetText(Password)

}
