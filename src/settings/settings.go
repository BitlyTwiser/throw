package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"fyne.io/fyne/v2"
)

const settingsFilePath string = "../../settings.json"

type Settings struct {
	Host         string
	Port         string
	DownloadPath string
	Encrypted    bool
	Password     string
}

func (s Settings) CurrentSettings() Settings {
	return s
}

func (s *Settings) SaveSettings(settings Settings) bool {
	// See if this works
	s = &settings

	file, err := os.OpenFile(settingsFilePath, os.O_TRUNC|os.O_RDWR, 0600)

	if err != nil {
		notification := fyne.NewNotification("Error", fmt.Sprintf("error saving settings. Error: %v", err))
		fyne.CurrentApp().SendNotification(notification)
	}

	j, err := json.MarshalIndent(&settings, "", "")

	if err != nil {
		notification := fyne.NewNotification("Error", fmt.Sprintf("error saving settings. Error: %v", err))
		fyne.CurrentApp().SendNotification(notification)
	}

	_, err = file.Write(j)

	if err != nil {
		notification := fyne.NewNotification("Error", fmt.Sprintf("error writing settings to file. Error: %v", err))
		fyne.CurrentApp().SendNotification(notification)
	}

	return true
}

// Load settings from file.
func LoadSettings() *Settings {
	log.Println("Loading settings")
	s := &Settings{}

	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		file, err := os.OpenFile(settingsFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

		if err != nil {
			notification := fyne.NewNotification("Error", fmt.Sprintf("error loading settings. Error: %v", err))
			fyne.CurrentApp().SendNotification(notification)
		}

		defer file.Close()

		j, err := json.MarshalIndent(Settings{}, "", "")

		if err != nil {
			notification := fyne.NewNotification("Error", fmt.Sprintf("error loading settings. Error: %v", err))
			fyne.CurrentApp().SendNotification(notification)
		}

		n, err := file.Write(j)

		if n == 0 || err != nil {
			notification := fyne.NewNotification("Error", "Zero bytes written or error")
			fyne.CurrentApp().SendNotification(notification)
		}

		// At the end of it all, we return empty settings after making the file
		return &Settings{}
	}

	fileData, err := os.ReadFile(settingsFilePath)

	if err != nil {
		log.Println("Death reading file data")
	}

	err = json.Unmarshal(fileData, &s)

	if err != nil {
		log.Printf("Error unmarshalling data. Error: %v", err)
	}

	return s
}
