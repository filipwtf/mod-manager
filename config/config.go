package config

import (
	"fmt"
	"gioui.org/unit"
	"log"
	"os"
)

// Config stores application version
type Config struct {
	Version   string `yaml:"version"`
	ShowLogs  bool   `yaml:"mods"`
	MCPath    string `yaml:"mcpath"`
	Installer bool   `yaml:"installer"`
}

var (
	Width  = unit.Dp(1260)
	Height = unit.Dp(640)
)

// GetConfig returns the config file
func GetConfig(flags int) *os.File {
	userConfig, _ := os.UserConfigDir()
	f, err := os.OpenFile(fmt.Sprintf("%s/installer-config.yaml", userConfig), flags, os.ModePerm)
	if err != nil {
		log.Println(err)
	}
	return f
}

// IsDirSet checks if the user has set their mc mods directory
func (c Config) IsDirSet() bool {
	if c.MCPath == "Mods Directory" || c.MCPath == "" {
		return false
	}
	return true
}
