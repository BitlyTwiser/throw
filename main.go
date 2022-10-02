package main

import (
	"fmt"
	"log"
	"time"
  "github.com/BitlyTwiser/throw/src/toolbar"
  "github.com/BitlyTwiser/throw/src/pufs_client"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
  a := app.New() 
  w := a.NewWindow("Filesystem")
  w.SetMaster()
  w.Resize(fyne.NewSize(500, 500))  
  
  // Unsubscribes client when application closes 
  defer pufs_client.UnsubscribeOnClose()
  data := []string{"One", "Two", "Three"}

  toolbar := widget.NewToolbar(
    widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { toolbar.UploadFile(w) } ),
        widget.NewToolbarSeparator(),
        widget.NewToolbarAction(theme.SettingsIcon(), func() { toolbar.Settings() }),
        widget.NewToolbarSpacer(),
        widget.NewToolbarAction(theme.HelpIcon(), func() { toolbar.HelpWindow() }),
        )
  
  fileList := widget.NewList(
      func() int { return len(data) },
      func() fyne.CanvasObject { return container.NewGridWithColumns(3, container.NewPadded(widget.NewLabel("")), container.NewPadded(widget.NewButtonWithIcon("Download", theme.DownloadIcon(), nil)), container.NewPadded(widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), nil))) },
      func(i widget.ListItemID, o fyne.CanvasObject) {
        o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(data[i])
        o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
          fmt.Printf("Hello from download. Data: %v", data[i])
        }
        o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
          fmt.Printf("Hello from delete. Data: %v", data[i])
        }
      },
    )

  go func() {
    time.Sleep(time.Second*5)
    log.Println("Adding element")
    data = append(data, "Nine")
    fileList.Refresh()
  }()

  //c := container.NewMax(fileList)
  content := container.NewBorder(toolbar, nil, nil, nil, fileList)
  
  w.SetContent(content)
  //w.SetContent(container.New(layout.NewAdaptiveGridLayout(2), c, uploadLayout))

  //w.Resize(fyne.NewSize(500, 500))

  w.ShowAndRun()
}

