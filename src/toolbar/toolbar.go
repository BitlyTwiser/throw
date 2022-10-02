package toolbar

import (
  "log"
  "io"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func UploadFile(window fyne.Window) {
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

      }, window).Show()
}

func HelpWindow() {
  log.Println("Hello from settings")
}

func Settings() {
  log.Println("Hello from settings")
}
