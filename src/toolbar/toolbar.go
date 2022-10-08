package toolbar

import (
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BitlyTwiser/throw/src/notifications"
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
			notifications.SendErrorNotification(fmt.Sprintf("Error uploading File. Error: %v", err))

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
	tg := widget.NewTextGrid()

	tg.Resize(fyne.NewSize(400, 400))

	tg.SetText("Throw is a vitualized file system utilizing the distributed power of IPFS to host your files.")

	// Neet to size this and add text lol
	helpWindow.SetContent(tg)

	helpWindow.Show()
}

// Set the values from the settings on load
func Settings(s *settings.Settings) {
	var downloadPath string
	settingsWindow := fyne.CurrentApp().NewWindow("Settings")

	settingsWindow.Resize(fyne.NewSize(400, 400))

	host := widget.NewEntry()
	if s.Host != "" {
		host.SetText(s.Host)
	}
	host.SetPlaceHolder("Enter Host Address...")

	port := widget.NewEntry()
	if s.Port != "" {
		port.SetText(s.Port)
	}
	port.SetPlaceHolder("Enter Host Port...")

	password := widget.NewPasswordEntry()
	if s.Password != "" {
		password.SetText(s.Password)
	}
	password.SetPlaceHolder("password...")

	downloadFolder := dialog.NewFolderOpen(func(f fyne.ListableURI, _ error) {
		if f == nil {
			return
		}

		downloadPath = f.Path()
	}, settingsWindow)

	button := widget.NewButtonWithIcon("Download Path", theme.FolderIcon(), nil)
	button.OnTapped = func() { downloadFolder.Show() }

	if s.Encrypted {
		password.Enable()
	} else {
		password.Disable()
	}
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
			newSettings := settings.Settings{
				Host:         host.Text,
				Port:         port.Text,
				Encrypted:    checkBox.Checked,
				Password:     password.Text,
				DownloadPath: downloadPath,
			}

			saved := s.SaveSettings(newSettings)

			if saved {
				notifications.SendSuccessNotification("Settings saved")
			} else {
				notifications.SendErrorNotification("Error saving settings")
			}
		},
		OnCancel: func() {
			settingsWindow.Close()
		},
	}

	tg := widget.NewTextGrid()
	tg.Resize(fyne.NewSize(100, 200))
	tg.SetText("Warning: Changes will apply after the client has been restarted")
	tg.SetStyleRange(0, 0, 0, len(tg.Text()), &widget.CustomTextGridStyle{color.White, color.RGBA{255, 0, 0, 0}})

	dp := widget.NewTextGrid()
	dp.SetText(fmt.Sprintf("%v", s.DownloadPath))

	// Append form elements
	form.Append("Host Address", host)
	form.Append("Host Port", port)
	form.Append("Encrypt Files", checkBox)
	form.Append("Encryption Password", password)
	form.Append("File Download Path", button)
	if s.DownloadPath != "" {
		form.Append("Curent Download Path", dp)
	}

	c := container.NewGridWithRows(2, form, tg)
	settingsWindow.SetContent(c)

	settingsWindow.Show()
}
