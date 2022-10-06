package settings

import "log"

type Settings struct {
	Host         string
	Port         string
	DownloadPath string
	Encrypted    bool
	Password     string
}

func SaveSettings(settings Settings) bool {
	log.Println("Saving settings")
	log.Println(settings)

	return true
}

func LoadSettings() {
}
