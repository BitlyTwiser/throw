package toolbar

import (
	"log"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BitlyTwiser/throw/src/pufs_client"
	"github.com/BitlyTwiser/throw/src/settings"
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
	var downloadPath string
	settingsWindow := fyne.CurrentApp().NewWindow("Settings")

	settingsWindow.Resize(fyne.NewSize(400, 400))

	host := widget.NewEntry()
	host.SetPlaceHolder("Enter Host Address...")

	port := widget.NewEntry()
	port.SetPlaceHolder("Enter Host Port...")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("password...")

	downloadFolder := dialog.NewFolderOpen(func(f fyne.ListableURI, _ error) {
		if f == nil {
			return
		}

		downloadPath = f.Path()

		log.Println(f.Path())
	}, settingsWindow)

	button := widget.NewButtonWithIcon("Download Path", theme.FolderIcon(), nil)
	button.OnTapped = func() { downloadFolder.Show() }

	checkBox := widget.NewCheck("", func(changed bool) {
		if changed {
			password.Enable()
		} else {
			password.Disable()
		}
	})

	form := &widget.Form{
		Items: []*widget.FormItem{},
		OnSubmit: func() {
			s := settings.Settings{
				Host:         host.Text,
				Port:         port.Text,
				Encrypted:    checkBox.Checked,
				Password:     password.Text,
				DownloadPath: downloadPath,
			}

			settings.SaveSettings(s)
		},
		OnCancel: func() {
			settingsWindow.Close()
		},
	}

	// Append form elements
	form.Append("Host Address", host)
	form.Append("Host Port", port)
	form.Append("Encrypt Files", checkBox)
	form.Append("Encryption Password", password)
	form.Append("File Download Path", button)

	settingsWindow.SetContent(form)

	settingsWindow.Show()
}
