package toolbar

import (
	"io"
	"log"

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
  helpWindow := fyne.CurrentApp().NewWindow("Help Window")
  helpWindow.Resize(fyne.NewSize(400, 400))
  
  button := widget.NewButton("Close", func() { helpWindow.Close() })
  
  // Neet to size this and add text lol
  helpWindow.SetContent(button)

  helpWindow.Show()
}

func Settings() {
  settingsWindow := fyne.CurrentApp().NewWindow("Settings")

  settingsWindow.Resize(fyne.NewSize(400, 400))

  form := &widget.Form{
    Items: []*widget.FormItem{},
    OnSubmit: func() {
      log.Println("Form submitted")
    },
    OnCancel: func() {
      log.Println("Close Form")
      settingsWindow.Close()
    },
  }
  
  host := widget.NewEntry()
  host.SetPlaceHolder("Enter Host Address...")
  port := widget.NewEntry()
  port.SetPlaceHolder("Enter Host Port...")


  form.Append("Host Address", host)
  form.Append("Host Port", port)

  settingsWindow.SetContent(form)

  settingsWindow.Show()
}
