package notifications

import (
	"fyne.io/fyne/v2"
)

func SendSuccessNotification(message string) {
	displayNotification(fyne.NewNotification("Success", message))
}

func SendErrorNotification(message string) {
	displayNotification(fyne.NewNotification("Error", message))
}

func displayNotification(n *fyne.Notification) {
	fyne.CurrentApp().SendNotification(n)
}
