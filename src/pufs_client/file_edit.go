package pufs_client

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func FileEditor(data []byte) *fyne.Container {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentSaveIcon(), func() {}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.CancelIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {}),
	)

	fileEditor := widget.NewMultiLineEntry()
	fileEditor.SetText(string(data))

	return container.NewBorder(toolbar, nil, nil, nil, fileEditor)
}
