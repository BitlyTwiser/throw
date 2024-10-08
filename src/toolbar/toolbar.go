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
	tg.SetText(`
	Throw is a vitualized file system utilizing the distributed power of IPFS to host your files.
	--------------------------------------------------------------------------------------------------------------------
	File Upload: 
		One can upload files using the Pencil Icon from the main page that is initially loaded upon start of the application.
		All files adde to the application, will be displayed in real time, when upload/delete actions commence.
	--------------------------------------------------------------------------------------------------------------------
	Settings:
		Using the Gear icon from within the toolbar, the user can set adjust server settings, download path, and if the data is to be encrypted in transit.
	--------------------------------------------------------------------------------------------------------------------
	Encryption:
		Data is encrypted using the tinycrypt golang library. https://github.com/BitlyTwiser/tinycrypt	
		Tinycrypt uses AES enryption algorithms to protect your data while the files are stored on IPFS.
	--------------------------------------------------------------------------------------------------------------------
	IPFS:
		IPFS is an amazing project known as the InterPlanetary File system. I do reccomend reading the docs to know more about this project.
		https://docs.ipfs.tech/concepts/what-is-ipfs/
		For the purposes of this project however, the core element is that your files are distributed to an IPFS node (Pufs gRPC server hosts/runs the node)
		Those files are then chunked and distributed to numerous peers, no data is stored directly on a disk.
	--------------------------------------------------------------------------------------------------------------------
	Server:
		gRPC is the cornerstone of the realtime transmission elements of Throw. Using the power of the Pufs project, the files are distrubuted to all clients in near real time.
	`)

	// Neet to size this and add text lol
	helpWindow.SetContent(tg)

	helpWindow.Show()
}

// Set the values from the settings on load
func Settings(s *settings.Settings) {
	var downloadPath string
	settingsWindow := fyne.CurrentApp().NewWindow("Settings")

	settingsWindow.Resize(fyne.NewSize(400, 400))

	if s.DownloadPath != "" {
		downloadPath = s.DownloadPath
	}

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

	checkBox.SetChecked(s.Encrypted)

	selectedFolder := widget.NewEntry()
	selectedFolder.SetText(downloadPath)

	if downloadPath != "" {
		selectedFolder.Enable()
	} else {
		selectedFolder.Disabled()
	}

	downloadFolder := dialog.NewFolderOpen(func(f fyne.ListableURI, _ error) {
		if f == nil {
			return
		}

		downloadPath = f.Path()
		selectedFolder.SetText(downloadPath)
		selectedFolder.Enable()
	}, settingsWindow)

	downloadFolderButton := widget.NewButtonWithIcon("Download Path", theme.FolderIcon(), nil)
	downloadFolderButton.OnTapped = func() { downloadFolder.Show() }

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
	tg.SetText("Warning: Any Host/IP changes will apply after the client has been restarted")
	tg.SetStyleRange(0, 0, 0, len(tg.Text()), &widget.CustomTextGridStyle{color.White, color.RGBA{255, 0, 0, 0}})

	// Append form elements
	form.Append("Host Address", host)
	form.Append("Host Port", port)
	form.Append("Encrypt Files", checkBox)
	form.Append("Encryption Password", password)
	form.Append("File Download Path", downloadFolderButton)
	if downloadPath != "" {
		form.Append("Curent Download Path", selectedFolder)
	}

	c := container.NewGridWithRows(2, form, tg)
	settingsWindow.SetContent(c)

	settingsWindow.Show()
}
