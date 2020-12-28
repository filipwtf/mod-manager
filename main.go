package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"github.com/filipwtf/filips-installer/config"
	"github.com/filipwtf/filips-installer/ui"
	"gopkg.in/yaml.v2"
)

func main() {
	configFile := config.GetConfig(os.O_CREATE | os.O_RDONLY)
	var uiCfg *config.Config
	if err := yaml.NewDecoder(configFile).Decode(&uiCfg); err != nil {
		log.Println(err)
		uiCfg.Version = "0.0.0"
		uiCfg.MCPath = "Mods Directory"
		uiCfg.ShowLogs = true
		uiCfg.Installer = false
	}
	defer configFile.Close()

	modUI := ui.NewUI(uiCfg)

	go func() {
		title := fmt.Sprintf("Filip's Mod Manager")
		window := app.NewWindow(
			app.Title(title),
			app.MinSize(config.Width, config.Height),
			app.MaxSize(config.Width, config.Height),
		)
		if err := modUI.Run(window); err != nil {
			log.Println(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	app.Main()
}
