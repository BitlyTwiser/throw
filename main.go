package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BitlyTwiser/throw/src/notifications"
	"github.com/BitlyTwiser/throw/src/pufs_client"
	"github.com/BitlyTwiser/throw/src/settings"
	"github.com/BitlyTwiser/throw/src/toolbar"

	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initializeUI(w fyne.Window, client pufs_client.IpfsClient) {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { toolbar.UploadFile(w, client) }),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() { toolbar.Settings(client.Settings) }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() { toolbar.HelpWindow() }),
	)

	fileList := widget.NewList(
		func() int { return len(client.Files) },
		func() fyne.CanvasObject {
			deleteButton := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			deleteButton.Resize(fyne.NewSize(5, 5))

			downloadButton := widget.NewButtonWithIcon("", theme.DownloadIcon(), nil)
			downloadButton.Resize(fyne.NewSize(5, 5))

			editButton := widget.NewButtonWithIcon("", theme.ContentAddIcon(), nil)
			editButton.Resize(fyne.NewSize(5, 5))

			fileNameLabel := widget.NewLabel("")
			fileNameLabel.Wrapping = 1

			return container.NewGridWithColumns(
				4,
				container.NewPadded(fileNameLabel),
				container.NewPadded(editButton),
				container.NewPadded(downloadButton),
				container.NewPadded(deleteButton),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(client.Files[i])
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				fileName := client.Files[i]
				w := fyne.CurrentApp().NewWindow(fmt.Sprintf("Edit %v", fileName))
				w.Resize(fyne.NewSize(300, 400))

				err := client.Download(fileName)

				if err != nil {
					notifications.SendErrorNotification("Error opening file for editing.")

					return
				}

				data, err := client.DownloadedFileContent(fileName)

				if err != nil {
					notifications.SendErrorNotification("Error loading file data for editing..")

					return
				}

				// Open File editor
				w.SetContent(pufs_client.FileEditor(*data, client, fileName, w))

				w.Show()
			}
			o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				client.Download(client.Files[i])
			}
			o.(*fyne.Container).Objects[3].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				var message *fyne.Notification
				fileName := client.Files[i]
				err := client.DeleteFile(fileName)

				if err != nil {
					message = fyne.NewNotification("Error", fmt.Sprintf("Error deleting file: %v", fileName))
				} else {
					message = fyne.NewNotification("Success", fmt.Sprintf("File %v delete", fileName))
				}

				fyne.CurrentApp().SendNotification(message)
			}
		},
	)

	go func() {
		for {
			select {
			case file := <-client.FileUpload:
				client.Files = append(client.Files, file)
				log.Println("Refreshing.")
				fileList.Refresh()
			case file := <-client.DeletedFile:
				var refresh []string
				for _, v := range client.Files {
					if v != file {
						refresh = append(refresh, v)
					}
				}
				client.Files = refresh
				log.Println("Refreshing after deleting file.")
				fileList.Refresh()
			}
		}
	}()

	content := container.NewBorder(toolbar, nil, nil, nil, fileList)

	w.SetContent(content)
}

var id int64

func main() {
	a := app.New()
	w := a.NewWindow("Throw")
	w.SetMaster()
	w.Resize(fyne.NewSize(750, 500))

	rand.Seed(time.Now().UTC().UnixNano())
	//Note: Look to add validation server side that ID is unique.
	id = int64(rand.Intn(100))

	s := settings.LoadSettings()

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", s.Host, s.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error connection to server: %v", err)
	} else {
		log.Println("Connected to gRPC server")
	}

	defer conn.Close()

	c := pufs_pb.NewIpfsFileSystemClient(conn)

	client := pufs_client.IpfsClient{
		Id:                id,
		Client:            c,
		Files:             []string{},
		FileUpload:        make(chan string, 1),
		DeletedFile:       make(chan string, 1),
		FileDeleted:       make(chan bool, 2),
		FileUploadedInApp: make(chan bool, 2),
		Settings:          s,
		InvalidFileTypes:  []string{"ELF", "EXE"},
	}
	// Remove  client after connection ends
	defer client.UnsubscribeClient()

	// Load existing files from server on application start
	client.LoadFiles()

	// Initialize the UI elements.
	initializeUI(w, client)

	go client.SubscribeFileStream()

	w.ShowAndRun()
}
