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
			return container.NewGridWithColumns(3, container.NewPadded(widget.NewLabel("")), container.NewPadded(widget.NewButtonWithIcon("Delete", theme.DownloadIcon(), nil)), container.NewPadded(widget.NewButtonWithIcon("Download", theme.DeleteIcon(), nil)))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("Name: %v", client.Files[i]))
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
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
			o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				fileName := client.Files[i]
				var err error

				if client.ChunkFile(fileName) {
					err = client.DownloadCappedFile(fileName, client.Settings.DownloadPath)
				} else {
					err = client.DownloadFile(fileName, client.Settings.DownloadPath)
				}

				if err != nil {
					notifications.SendErrorNotification(fmt.Sprintf("Error downloading file: %v", fileName))
				} else {
					notifications.SendSuccessNotification(fmt.Sprintf("File %v downloaded", fileName))
				}
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
				for i, v := range client.Files {
					if v == file {
						client.Files = append(client.Files[:i], client.Files[i+1:]...)
					}
				}
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
	w := a.NewWindow("Filesystem")
	w.SetMaster()
	w.Resize(fyne.NewSize(500, 500))

	rand.Seed(time.Now().UTC().UnixNano())
	//Note: Look to add validation server side that ID is unique.
	id = int64(rand.Intn(100))

	// Must load values for address and server port from storage.
	// The settings page will store these values.
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", "127.0.0.1", 9000), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error connection to server: %v", err)
	} else {
		log.Println("Connected to gRPC server")
	}

	defer conn.Close()

	c := pufs_pb.NewIpfsFileSystemClient(conn)

	client := pufs_client.IpfsClient{
		Id:          id,
		Client:      c,
		Files:       []string{},
		FileUpload:  make(chan string, 1),
		DeletedFile: make(chan string, 1),
		FileDeleted: make(chan bool, 1),
		Settings:    settings.LoadSettings(),
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
