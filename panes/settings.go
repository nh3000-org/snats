/*
 *	PROGRAM		: settings.go
 *	DESCRIPTION		:
 *
 *		This program handles option definitions.
 *
 *	PARAMETERS		:
  *
 *	RETURNS			:
 *		Canvas
*/

package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
 *	FUNCTION		: settingsScreen
 *	DESCRIPTION		:
 *		Interface for changing settings
 *
 *	PARAMETERS		:
 *		        	:
 *
 *	RETURNS			:
 *
 */
func settingsScreen(_ fyne.Window) fyne.CanvasObject {

	pllabel := widget.NewLabel("Password Length")
	pl := widget.NewRadioGroup([]string{"6", "8", "12"}, func(string) {})
	pl.Horizontal = true
	pl.SetSelected("6")

	mcletterlabel := widget.NewLabel("Password Must Contain Letter")
	mcletter := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	mcletter.Horizontal = true
	mcletter.SetSelected("False")

	mcnumberlabel := widget.NewLabel("Password Must Contain Number")
	mcnumber := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	mcnumber.Horizontal = true
	mcnumber.SetSelected("False")

	mcspeciallabel := widget.NewLabel("Password Must Contain Special")
	mcspecial := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	mcspecial.Horizontal = true
	mcspecial.SetSelected("False")

	usetlslabel := widget.NewLabel("Use TLS Authorization")
	usetls := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	usetls.Horizontal = true
	usetls.SetSelected("False")

	usejslabel := widget.NewLabel("Use Jetstream")
	usejs := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	usejs.Horizontal = true
	usejs.SetSelected("True")

	UseJetstream = editEntry("TRUEFALSE", usejs.Selected)
	UseTLS = editEntry("TRUEFALSE", usetls.Selected)
	PasswordMustContainNumber = editEntry("TRUEFALSE", mcnumber.Selected)
	PasswordMinimumSize = pl.Selected
	PasswordMustContainLetter = editEntry("TRUEFALSE", mcletter.Selected)
	PasswordMustContainSpecial = editEntry("TRUEFALSE", mcspecial.Selected)

	ssbutton := widget.NewButton("Save Settings", func() {

		if PasswordValid {

			MyJson("SAVE")

		}
	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		pllabel,
		pl,
		mcletterlabel,
		mcletter,
		mcnumberlabel,
		mcnumber,
		mcspeciallabel,
		mcspecial,
		usetlslabel,
		usetls,
		usejslabel,
		usejs,

		ssbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
