package settings

import (
	"encoding/json"
	"log"
	"os"
)

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

	return true
}

// Load settings from file.
func LoadSettings() *Settings {
	log.Println("Loading settings")
	s := &Settings{}
	var settingsFilePath string = "../../settings.json"

	if _, err := os.Stat(settingsFilePath); os.IsNotExist(err) {
		file, err := os.OpenFile(settingsFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0600)

		if err != nil {
			log.Println("Error loading settings..")
		}

		defer file.Close()

		j, err := json.MarshalIndent(Settings{}, "", "")

		if err != nil {
			log.Printf("error storing settings file content. Error: %v", err)
		}

		n, err := file.Write(j)

		if n == 0 || err != nil {
			log.Println("Zero bytes written to settings file or there was an error")
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
