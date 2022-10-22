package pufs_client

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Create fyne table, insert file data within.
func FileMetadata(fileData FileData) {
	w := fyne.CurrentApp().NewWindow("File Metadata")
	w.Resize(fyne.NewSize(400, 400))

	fileSize := strconv.Itoa(int(fileData.FileSize))

	data := [][]string{[]string{"FileName", "FileSize", "UpdatedAt", "IpfsHash"}, []string{
		fileData.FileName, fileSize, fileData.UploadedAt.String(), fileData.IpfsHash,
	}}

	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("File Metadata")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(data[i.Row][i.Col])
		})

	w.SetContent(table)
	w.Show()
}
