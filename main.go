package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	// "context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/BitlyTwiser/throw/src/pufs_client"
	"github.com/BitlyTwiser/throw/src/toolbar"

	pufs_pb "github.com/BitlyTwiser/pufs-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initializeUI(w fyne.Window, client pufs_pb.IpfsFileSystemClient) {
	// Load files here.
	// Spawn alternate thread to handle the file watching.
	data := []string{"One", "Two", "Three"}

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() { toolbar.UploadFile(w) }),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() { toolbar.Settings() }),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() { toolbar.HelpWindow() }),
	)

	fileList := widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject {
			return container.NewGridWithColumns(3, container.NewPadded(widget.NewLabel("")), container.NewPadded(widget.NewButtonWithIcon("Download", theme.DownloadIcon(), nil)), container.NewPadded(widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), nil)))
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("Name: %v", data[i]))
			o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				fmt.Printf("Hello from download. Data: %v", data[i])
			}
			o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*widget.Button).OnTapped = func() {
				fmt.Printf("Hello from delete. Data: %v", data[i])
			}
		},
	)

	go func() {
		time.Sleep(time.Second * 5)
		log.Println("Adding element")
		data = append(data, "Nine")
		fileList.Refresh()
	}()

	//c := container.NewMax(fileList)
	content := container.NewBorder(toolbar, nil, nil, nil, fileList)

	w.SetContent(content)
}

var id int64

func main() {
	a := app.New()
	w := a.NewWindow("Filesystem")
	w.SetMaster()
	w.Resize(fyne.NewSize(500, 500))
  
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
  
  client := pufs_client.IpfsClient{Id: id, Client: c}
  // Remove  client after connection ends
	defer client.UnsubscribeClient()

  go client.SubscribeFileStream()
  
  // Initialize the UI elements.
	initializeUI(w, c)

	w.ShowAndRun()
}
