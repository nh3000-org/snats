package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsScreen(_ fyne.Window) fyne.CanvasObject {
	MyJson("LOAD")
	pllabel := widget.NewLabel("Password Length")
	pl := widget.NewRadioGroup([]string{"6", "8", "12"}, func(string) {})
	pl.Horizontal = true
	pl.SetSelected(PasswordMinimumSize)

	mcletterlabel := widget.NewLabel("Password Must Contain Letter")
	mcletter := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	mcletter.Horizontal = true
	mcletter.SetSelected(PasswordMustContainLetter)

	mcnumberlabel := widget.NewLabel("Password Must Contain Number")
	mcnumber := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	mcnumber.Horizontal = true
	mcnumber.SetSelected(PasswordMustContainNumber)

	mcspeciallabel := widget.NewLabel("Password Must Contain Special")
	mcspecial := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	mcspecial.Horizontal = true
	mcspecial.SetSelected(PasswordMustContainSpecial)

	ssbutton := widget.NewButton("Save Settings", func() {
		PasswordMustContainNumber = mcnumber.Selected
		PasswordMinimumSize = pl.Selected
		PasswordMustContainLetter = mcletter.Selected
		PasswordMustContainSpecial = mcspecial.Selected
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
		ssbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""),
	))

}
