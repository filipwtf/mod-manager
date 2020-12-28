package ui

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/filipwtf/filips-installer/config"
	"gopkg.in/yaml.v2"
)

type (
	ctx = layout.Context
	dim = layout.Dimensions
)

var (
	std  string
	list = &layout.List{
		Axis: layout.Vertical,
	}
	editor = new(widget.Editor)
	prefix = fmt.Sprintf("[%s] ", time.Now().Format("2-1-2006 15:04"))
	start  time.Time
)

// UI holds all of the application state.
type UI struct {
	theme   *material.Theme
	version string
	update  bool
	main    mainWidget
	cfg     configWidget
	log     logWidget
	mods    []Mod
}

type mainWidget struct {
	installer bool
	editor    *widget.Editor
	list      *layout.List
}

type configWidget struct {
	showLog    bool
	mcPath     string
	editor     *widget.Editor
	list       *layout.List
	showLogBtn *widget.Clickable
	setDirBtn  *widget.Clickable
	configRef  config.Config
}

type logWidget struct {
	editor *widget.Editor
	list   *layout.List
}

// NewUI creates a new UI using the Go Fonts.
func NewUI(config config.Config) *UI {
	ui := &UI{
		theme:   material.NewTheme(gofont.Collection()),
		version: config.Version,
		main: mainWidget{
			installer: config.Installer,
			editor:    new(widget.Editor),
			list: &layout.List{
				Axis: layout.Vertical,
			},
		},
		cfg: configWidget{
			showLog: config.ShowLogs,
			mcPath:  config.MCPath,
			editor:  new(widget.Editor),
			list: &layout.List{
				Axis: layout.Horizontal,
			},
			showLogBtn: new(widget.Clickable),
			setDirBtn:  new(widget.Clickable),
			configRef:  config,
		},
		log: logWidget{
			editor: new(widget.Editor),
			list: &layout.List{
				Axis: layout.Vertical,
			},
		},
	}

	go func() {
		start = time.Now()
		ui.mods = GetAllMods(config)
		end := time.Since(start)
		Log(fmt.Sprintf("Loading %d mods took %fs", len(ui.mods), end.Seconds()))
	}()

	go func() {
		ui.checkForUpdate()
		if ui.update {
			downloadUpdate()
		}
	}()

	return ui
}

func (ui *UI) checkForUpdate() {
	Log("Checking for updates")
	latest := getLatestVersion()
	if ui.version != latest {
		Log(fmt.Sprintf("Update available %s -> %s", ui.version, latest))
		ui.update = true
	} else {
		Log(fmt.Sprintf("No updates %s", ui.version))
		ui.update = false
	}
}

func downloadUpdate() {
	// TODO : Update
	Log("Downloading update")

	Log("Update complete please restart")
}

// TODO Fetch from url
func getLatestVersion() string {
	return "1.0.0"
}

func (ui *UI) saveConfig() {
	configFile := config.GetConfig(os.O_WRONLY)
	defer configFile.Close()

	cfg := config.Config{
		Version:   ui.version,
		ShowLogs:  ui.cfg.showLog,
		MCPath:    ui.cfg.mcPath,
		Installer: false,
	}
	if err := yaml.NewEncoder(configFile).Encode(cfg); err != nil {
		log.Println(err)
	}
}

// Run handles window events and renders the application.
func (ui *UI) Run(w *app.Window) error {
	var ops op.Ops
	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			ui.drawLayout(gtx)
			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			ui.saveConfig()
			return e.Err
		}
	}
	return nil
}

func (ui *UI) drawLayout(gtx ctx) dim {
	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx ctx) dim {
			return ui.log.logLayout(ui.theme, gtx, ui.cfg.showLog)
		}),
		layout.Stacked(func(gtx ctx) dim {
			return ui.cfg.configLayout(ui.theme, gtx, ui.version)
		}),
		layout.Stacked(func(gtx ctx) dim {
			if ui.main.installer {
				return ui.main.installerLayout(ui.theme, gtx, ui.cfg.showLog)
			}
			return ui.main.managerLayout(ui.theme, gtx, ui.mods, ui.cfg.showLog)
		}),
	)
}

func (cfg *configWidget) configLayout(th *material.Theme, gtx ctx, version string) dim {
	height := getConfigHeight(gtx, cfg.showLog)
	return layout.Inset{Top: unit.Dp(height)}.Layout(gtx, func(gtx ctx) dim {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(0.2, func(gtx ctx) dim {
				return layout.Inset{Left: unit.Dp(2)}.Layout(gtx, func(gtx ctx) dim {
					for cfg.showLogBtn.Clicked() {
						cfg.showLog = !cfg.showLog
					}
					return material.Button(th, cfg.showLogBtn, "Show logs").Layout(gtx)
				})
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Flexed(0.6, func(gtx ctx) dim {
				e := material.Editor(th, cfg.editor, cfg.mcPath)
				border := widget.Border{Color: color.NRGBA{A: 0xff}, CornerRadius: unit.Dp(8), Width: unit.Px(1)}
				return border.Layout(gtx, func(gtx ctx) dim {
					gtx.Constraints.Min.X = 260
					return layout.UniformInset(unit.Dp(8)).Layout(gtx, e.Layout)
				})
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
			layout.Flexed(0.2, func(gtx ctx) dim {
				return layout.Inset{Left: unit.Dp(2)}.Layout(gtx, func(gtx ctx) dim {
					for cfg.setDirBtn.Clicked() {
						if cfg.mcPath != cfg.editor.Text() {
							cfg.mcPath = cfg.editor.Text()
							Log(fmt.Sprintf("Mods directory set to: %s", cfg.mcPath))
							GetAllMods(cfg.configRef)
						}
					}
					return material.Button(th, cfg.setDirBtn, "Set Dir").Layout(gtx)
				})
			}),
			layout.Rigid(layout.Spacer{Width: unit.Dp(50)}.Layout),
			layout.Rigid(func(gtx ctx) dim {
				return material.Body1(th, fmt.Sprintf("Version - %s", version)).Layout(gtx)
			}),
		)
	})
}

func (log *logWidget) logLayout(th *material.Theme, gtx ctx, showLog bool) dim {
	if showLog {
		editor.SetText(std)
		width := float32(gtx.Constraints.Max.X)
		logHeight := float32(float64(gtx.Constraints.Max.Y) * 0.25)
		height := float32(float64(gtx.Constraints.Max.Y) * 0.75)

		return list.Layout(gtx, 1, func(gtx ctx, i int) dim {
			return layout.Inset{
				Top: unit.Dp(height),
			}.Layout(gtx, func(gtx ctx) dim {
				gtx.Constraints.Max.Y = gtx.Px(unit.Dp(logHeight))
				gtx.Constraints.Min.X = gtx.Px(unit.Dp(width))
				return material.Editor(th, editor, "Log is empty").Layout(gtx)
			})
		})
	}
	return layoutWidget(0, 0)
}

func layoutWidget(width, height int) dim {
	return dim{
		Size: image.Point{
			X: width,
			Y: height,
		},
	}
}

func getConfigHeight(gtx ctx, show bool) float32 {
	if show {
		return float32(float64(gtx.Constraints.Max.Y)*0.75) - 40
	}
	return float32(float64(gtx.Constraints.Max.Y)) - 40
}

// Log a message to editor
func Log(msg string) {
	std += fmt.Sprintf("%s %s\n", prefix, msg)
}
