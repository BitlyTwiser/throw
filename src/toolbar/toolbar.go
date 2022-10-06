package toolbar

import (
	"log"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/BitlyTwiser/throw/src/pufs_client"
)

func UploadFile(window fyne.Window, client pufs_client.IpfsClient) {
	dialog.NewFileOpen(func(f fyne.URIReadCloser, _ error) {
		if f == nil {
			log.Println("No file selected")

			return
		}

		err := client.UploadFile(f.URI().Path(), f.URI().Name())

		if err != nil {
			errorNotification := fyne.NewNotification("Error", "Error uploading File")
			fyne.CurrentApp().SendNotification(errorNotification)

			return
		}

	}, window).Show()
}

func EditFileWindow(data []byte) {
	result := fyne.NewStaticResource("File", data)

	entry := widget.NewMultiLineEntry()
	entry.SetText(string(result.StaticContent))

	w := fyne.CurrentApp().NewWindow(string(result.StaticName))
	w.SetContent(container.NewScroll(entry))

	w.Resize(fyne.NewSize(400, 400))
	w.Show()
}

func HelpWindow() {
	helpWindow := fyne.CurrentApp().NewWindow("Help Window")
	helpWindow.Resize(fyne.NewSize(400, 400))

	button := widget.NewButton("Close", func() { helpWindow.Close() })

	// Neet to size this and add text lol
	helpWindow.SetContent(button)

	helpWindow.Show()
}

func Settings() {
	settingsWindow := fyne.CurrentApp().NewWindow("Settings")

	settingsWindow.Resize(fyne.NewSize(400, 400))

	form := &widget.Form{
		Items: []*widget.FormItem{},
		OnSubmit: func() {
			log.Println("Form submitted")
		},
		OnCancel: func() {
			log.Println("Close Form")
			settingsWindow.Close()
		},
	}

	host := widget.NewEntry()
	host.SetPlaceHolder("Enter Host Address...")
	port := widget.NewEntry()
	port.SetPlaceHolder("Enter Host Port...")

	form.Append("Host Address", host)
	form.Append("Host Port", port)

	settingsWindow.SetContent(form)

	settingsWindow.Show()
}
