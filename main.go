package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
  a := app.New() 
  w := a.NewWindow("Filesystem")
  w.SetMaster()
  w.Resize(fyne.NewSize(500, 500))  

  data := []string{"One", "Two", "Three"}

  uploadFile := func() {
      dialog.NewFileOpen(func(f fyne.URIReadCloser, _ error) { 
        if f == nil {
          log.Println("No file selected")
          return
        }

        data, err := io.ReadAll(f)

        if err != nil {
          log.Printf("Error in reading file data. Error: %v", err)

          return
        }

        result := fyne.NewStaticResource("File", data)
        
        entry := widget.NewMultiLineEntry()
        entry.SetText(string(result.StaticContent))

        w := fyne.CurrentApp().NewWindow(string(result.StaticName))
        w.SetContent(container.NewScroll(entry))

        w.Resize(fyne.NewSize(400,400))
        w.Show()

      }, w).Show()}
  
  toolbar := widget.NewToolbar(
    widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { uploadFile() } ),
        widget.NewToolbarSeparator(),
        )
  
  fileList := widget.NewList(
      func() int { return len(data) },
      func() fyne.CanvasObject { return container.NewPadded(container.NewPadded(widget.NewLabel("")), container.NewPadded(widget.NewButtonWithIcon("Download", theme.DownloadIcon(), nil)), container.NewPadded(widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), nil))) },
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

