package config

import (
	"fmt"
	"log"
	"os"
)

// Config stores application version
type Config struct {
	Version  string `yaml:"version"`
	ShowLogs bool   `yaml:"mods"`
	MCPath   string `yaml:"mcpath"`
}

// GetConfig returns the config file
func GetConfig(flags int) *os.File {
	userConfig, _ := os.UserConfigDir()
	f, err := os.OpenFile(fmt.Sprintf("%s/installer-config.yaml", userConfig), flags, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	return f
}
