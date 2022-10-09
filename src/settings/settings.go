package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/BitlyTwiser/throw/src/notifications"
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
	s.SaveSettingsMemory(&settings)

	settings.Password = Base64EncodeString([]byte(settings.Password))

	file, err := os.OpenFile(settingsFilePath, os.O_TRUNC|os.O_RDWR, 0600)

	if err != nil {
		notifications.SendErrorNotification(fmt.Sprintf("error saving settings. Error: %v", err))
	}

	j, err := json.MarshalIndent(&settings, "", "")

	if err != nil {
		notifications.SendErrorNotification(fmt.Sprintf("error saving settings. Error: %v", err))
	}

	_, err = file.Write(j)

	if err != nil {
		notifications.SendErrorNotification(fmt.Sprintf("error writing settings to file. Error: %v", err))
	}

	return true
}

func (s *Settings) SaveSettingsMemory(settings *Settings) {
	*s = *settings
}

// Load settings from file.
func LoadSettings() *Settings {
	log.Println("Loading settings")
	s := &Settings{}

	// If the settings file does now exist, write the generic struct outline to the file.
	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		file, err := os.OpenFile(settingsFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

		if err != nil {
			notifications.SendErrorNotification(fmt.Sprintf("error loading settings. Error: %v", err))
		}

		defer file.Close()

		j, err := json.MarshalIndent(Settings{}, "", "")

		if err != nil {
			notifications.SendErrorNotification(fmt.Sprintf("error loading settings. Error: %v", err))
		}

		n, err := file.Write(j)

		if n == 0 || err != nil {
			notifications.SendErrorNotification("Zero bytes written or error")
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

	// Decode password here.
	if s.Password != "" {
		pass, err := DecodeString(s.Password)

		if err == nil {
			s.Password = pass
		}
	}

	return s
}
