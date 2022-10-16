package pufs_client

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BitlyTwiser/throw/src/notifications"
)

func FileEditor(data []byte, client IpfsClient, fileName string, w fyne.Window) *fyne.Container {

	fileEditor := widget.NewMultiLineEntry()
	fileEditor.Wrapping = 1

	fileEditor.SetText(string(data))

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {
			err := saveFile(
				fileEditor.Text,
				client.Settings.DownloadPath,
				fileName)

			if err != nil {
				notifications.SendErrorNotification(fmt.Sprintf("File data failed to save. Error: %v", err.Error()))
			} else {
				notifications.SendSuccessNotification("File data saved")
			}

			client.UploadFile(fmt.Sprintf("%v/%v", client.Settings.DownloadPath, fileName), fileName)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.CancelIcon(), func() {
			w.Close()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {}),
	)

	return container.NewBorder(toolbar, nil, nil, nil, fileEditor)
}

func saveFile(data, path, filename string) error {
	file, err := os.OpenFile(fmt.Sprintf("%v/%v", path, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if err != nil {
		return err
	}

	defer file.Close()

	// Write file from 0th byte, replacing file with newfound content.
	_, err = file.WriteAt([]byte(data), 0)

	if err != nil {
		return err
	}

	return nil
}
