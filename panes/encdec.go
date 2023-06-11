package panes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func encdecScreen(win fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	password := widget.NewEntry()
	password.SetPlaceHolder("Enter Password For Encryption")

	myinputtext := widget.NewMultiLineEntry()
	myinputtext.SetPlaceHolder("Enter Value For Enc/Dec")
	myinputtext.SetMinRowsVisible(6)

	myinputtext.SetText(win.Clipboard().Content())
	myoutputtext := widget.NewMultiLineEntry()
	myoutputtext.SetPlaceHolder("Output Shows Up Here")
	myoutputtext.SetMinRowsVisible(6)
	var iserrors = false
	errors.SetText("...")
	encbutton := widget.NewButton("Encrypt Message", func() {
		iserrors = editEntry("STRING", password.Text)
		if iserrors == true {
			errors.SetText("Error Invalid Password")
			iserrors = true
		}
		if len(password.Text) != 24 {
			iserrors = true
			errors.SetText("Error Password Length is " + strconv.Itoa(len(password.Text)) + " shlould be length of 24")
		}
		iserrors = editEntry("STRING", myinputtext.Text)
		if iserrors == true {
			errors.SetText("Error Input Text Misssing")
			iserrors = true
		}
		if iserrors == false {
			t, err := Encrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText("Error Input Text " + err.Error())
			} else {
				myoutputtext.SetText(string(t))
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}
	})

	decbutton := widget.NewButton("Decrypt Message", func() {
		iserrors = editEntry("STRING", password.Text)

		iserrors = editEntry("STRING", myinputtext.Text)
		if iserrors == true {
			errors.SetText("Error Input Text Misssing")
		}
		if len(password.Text) != 24 {
			iserrors = true
			errors.SetText("Error Password Length is " + strconv.Itoa(len(password.Text)) + " shlould be length of 24")
		}
		if iserrors == false {
			t, err := Decrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText("Error Input Text " + err.Error())
			} else {
				myoutputtext.SetText(t)
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}

	})

	if iserrors == true {
		//encbutton.Disable()
		//decbutton.Disable()
	}
	keybox := container.NewBorder(
		password,
		nil,
		nil,
		nil,
		nil,
	)
	inputbox := container.NewBorder(
		widget.NewLabelWithStyle("Input", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myinputtext,
		nil,
		nil,
		nil,
	)
	outputbox := container.NewBorder(
		widget.NewLabelWithStyle("Output", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myoutputtext,
		nil,
		nil,
		nil,
	)
	buttonbox := container.NewBorder(
		nil,
		nil,
		nil,
		encbutton,
		decbutton,
	)
	c0box := container.NewBorder(
		keybox,
		nil,
		nil,
		nil,
		nil,
	)
	c1box := container.NewBorder(
		inputbox,
		buttonbox,
		nil,
		nil,
		nil,
	)
	c2box := container.NewBorder(
		c0box,
		c1box,
		nil,package panes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func encdecScreen(win fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	password := widget.NewEntry()
	password.SetPlaceHolder("Enter Password For Encryption")

	myinputtext := widget.NewMultiLineEntry()
	myinputtext.SetPlaceHolder("Enter Value For Enc/Dec")
	myinputtext.SetMinRowsVisible(6)

	myinputtext.SetText(win.Clipboard().Content())
	myoutputtext := widget.NewMultiLineEntry()
	myoutputtext.SetPlaceHolder("Output Shows Up Here")
	myoutputtext.SetMinRowsVisible(6)
	var iserrors = false
	errors.SetText("...")
	encbutton := widget.NewButton("Encrypt Message", func() {
		iserrors = editEntry("STRING", password.Text)
		if iserrors == true {
			errors.SetText("Error Invalid Password")
			iserrors = true
		}
		if len(password.Text) != 24 {
			iserrors = true
			errors.SetText("Error Password Length is " + strconv.Itoa(len(password.Text)) + " shlould be length of 24")
		}
		iserrors = editEntry("STRING", myinputtext.Text)
		if iserrors == true {
			errors.SetText("Error Input Text Misssing")
			iserrors = true
		}
		if iserrors == false {
			t, err := Encrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText("Error Input Text " + err.Error())
			} else {
				myoutputtext.SetText(string(t))
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}
	})

	decbutton := widget.NewButton("Decrypt Message", func() {
		iserrors = editEntry("STRING", password.Text)

		iserrors = editEntry("STRING", myinputtext.Text)
		if iserrors == true {
			errors.SetText("Error Input Text Misssing")
		}
		if len(password.Text) != 24 {
			iserrors = true
			errors.SetText("Error Password Length is " + strconv.Itoa(len(password.Text)) + " shlould be length of 24")
		}
		if iserrors == false {
			t, err := Decrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText("Error Input Text " + err.Error())
			} else {
				myoutputtext.SetText(t)
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}

	})

	if iserrors == true {
		//encbutton.Disable()
		//decbutton.Disable()
	}
	keybox := container.NewBorder(
		password,
		nil,
		nil,
		nil,
		nil,
	)
	inputbox := container.NewBorder(
		widget.NewLabelWithStyle("Input", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myinputtext,
		nil,
		nil,
		nil,
	)
	outputbox := container.NewBorder(
		widget.NewLabelWithStyle("Output", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myoutputtext,
		nil,
		nil,
		nil,
	)
	buttonbox := container.NewBorder(
		nil,
		nil,
		nil,
		encbutton,
		decbutton,
	)
	c0box := container.NewBorder(
		keybox,
		nil,
		nil,
		nil,
		nil,
	)
	c1box := container.NewBorder(
		inputbox,
		buttonbox,
		nil,
		nil,
		nil,
	)
	c2box := container.NewBorder(
		c0box,
		c1box,
		nil,
		nil,
		nil,
	)
	c3box := container.NewBorder(
		c2box,
		outputbox,
		nil,
		nil,
		nil,
	)

	return container.NewBorder(
		c3box,
		errors,
		nil,
		nil,
		nil,
	)

}

		nil,
		nil,
	)
	c3box := container.NewBorder(
		c2box,
		outputbox,
		nil,
		nil,
		nil,
	)

	return container.NewBorder(
		c3box,
		errors,
		nil,
		nil,
		nil,
	)

}
