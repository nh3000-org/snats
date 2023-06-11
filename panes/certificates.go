package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func certificatesScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")
	if PasswordValid == true {
		calabel := widget.NewLabel("CAROOT Certificate")
		ca := widget.NewMultiLineEntry()
		ca.SetText(Caroot)

		cclabel := widget.NewLabel("CLIENT Certificate")
		cc := widget.NewMultiLineEntry()
		cc.SetText(Clientcert)

		cklabel := widget.NewLabel("CLIENT Key")
		ck := widget.NewMultiLineEntry()
		ck.SetText(Clientkey)

		ssbutton := widget.NewButton("Save Settings", func() {
			errors.SetText("...")
			if PasswordValid {
				var iserrors = editEntry("CERTIFICATE", ca.Text)
				if iserrors {
					errors.SetText("Error CAROOT is invalid ")
				}
				iserrors = editEntry("CERTIFICATE", cc.Text)
				if iserrors {
					errors.SetText("Error CLIENT CERTIFICATE is invalid ")
				}
				iserrors = editEntry("KEY", ck.Text)
				if iserrors {
					errors.SetText("Error CLIENT KEY is invalid ")
				}
				if !iserrors {
					MyJson("SAVE")
				}
			}
		})

		return container.NewCenter(container.NewVBox(
			widget.NewLabelWithStyle("New Horizons 3000 Secure Communications ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			calabel,
			ca,
			cclabel,
			cc,
			cklabel,
			cklabel,
			ck,

			ssbutton,
			errors,
			container.NewHBox(
				widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			),
			widget.NewLabel(""),
		))
	}
	errors.SetText("Logon First")
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle("New Horizons 3000 Secure Communications ", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		errors,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""),
	))

}
