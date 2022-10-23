package pufs_client

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Create fyne table, insert file data within.
func FileMetadata(fileData FileData) {
	w := fyne.CurrentApp().NewWindow("File Metadata")
	w.Resize(fyne.NewSize(400, 200))

	fileSize := strconv.Itoa(int(fileData.FileSize))

	fileNameLabel := widget.NewLabel(fmt.Sprintf("File Name: %v\n--------------", fileData.FileName))
	fileNameLabel.Wrapping = 1

	fileSizeLabel := widget.NewLabel(fmt.Sprintf("File Size: %vMB\n------------------------", fileSize))
	fileSizeLabel.Wrapping = 1

	fileUploadedAt := widget.NewLabel(fmt.Sprintf("Uploaded at: %v\n------------", fileData.UploadedAt))
	fileUploadedAt.Wrapping = 1

	ipfsHash := widget.NewLabel(fmt.Sprintf("Ipfs Hash: %v\n--------------------", fileData.IpfsHash))
	ipfsHash.Wrapping = 1

	table := container.NewGridWithRows(
		2,
		fileNameLabel,
		fileSizeLabel,
		fileUploadedAt,
		ipfsHash,
	)

	w.SetContent(table)
	w.Show()
}
