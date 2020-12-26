package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/filipwtf/filips-installer/config"
	"github.com/filipwtf/filips-installer/ui"
	"gopkg.in/yaml.v2"
)

var std string

func main() {
	configFile := config.GetConfig(os.O_CREATE | os.O_RDONLY)
	var config config.Config
	if err := yaml.NewDecoder(configFile).Decode(&config); err != nil {
		log.Println(err)
		config.Version = "0.0.0"
		config.MCPath = "Enter mc path"
		config.ShowLogs = true
	}
	defer configFile.Close()

	ui := ui.NewUI(config)

	go func() {
		title := fmt.Sprintf("Filip's Mod Manager")
		window := app.NewWindow(
			app.Title(title),
			app.MinSize(unit.Dp(1260), unit.Dp(640)),
			app.MaxSize(unit.Dp(640), unit.Dp(640)),
		)
		if err := ui.Run(window); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}
