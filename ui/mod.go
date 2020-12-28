package ui

import (
	"archive/zip"
	"bufio"
	"fmt"
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/filipwtf/filips-installer/config"
	"io/ioutil"
	"os"
	"strings"
)

// Mod component
type Mod struct {
	Name       string `yaml:"name"`
	SimpleName string `yaml:"simplename"`
	Path       string `yaml:"path"`
	Version    string `yaml:"version"`
	Latest     bool   `yaml:"latest"`
}

// MCModInfo forge mod info, unreliable due to mod authors not using this
type MCModInfo struct {
	modid        string
	name         string
	description  string
	version      string
	mcversion    string
	url          string
	updateUrl    string
	authorList   []string
	credits      string
	logoFile     string
	screenshots  []string
	dependencies []string
}

func (mod *Mod) DeleteMod() {
	err := os.Remove(mod.Path)
	if err != nil {
		Log(fmt.Sprintf("Failed to delete file %s", mod.Name))
	}
	Log(fmt.Sprintf("Succesffuly deleted %s", mod.Name))
}

func (mod *Mod) modWidget(th *material.Theme) layout.Widget {
	return func(gtx ctx) dim {
		return layout.Flex{}.Layout(gtx,
			layout.Flexed(0.2, func(gtx ctx) dim {
				return material.Body1(th, mod.SimpleName).Layout(gtx)
			}),
			layout.Flexed(0.4, func(gtx ctx) dim {
				return material.Body1(th, mod.Version).Layout(gtx)
			}),
			// TODO Delete Button
			// TODO Integrity Check
			// TODO Update check? Not very feasible
		)
	}
}

func GetAllMods(config config.Config) []Mod {
	if config.IsDirSet() == false {
		Log("Not a valid path")
		return []Mod{}
	}
	dir := config.MCPath

	files := getJarFiles(dir)

	var mods []Mod

	for _, file := range files {
		path := fmt.Sprintf("%s/%s", dir, file.Name())
		version := getModMeta(path)
		m := Mod{
			Name:       file.Name(),
			SimpleName: strings.Replace(file.Name(), ".jar", "", 1),
			Path:       path,
			Version:    version,
		}
		Log(fmt.Sprintf("Found mod: %s, version: %s", file.Name(), version))
		mods = append(mods, m)
	}
	return mods
}

func getModMeta(path string) string {
	def := "Unable to determine version"
	dir, err := zip.OpenReader(path)
	if err != nil || dir == nil {
		Log("Unable to open jar")
		return def
	}

	defer dir.Close()

	if strings.Contains(path, "OptiFine") {
		s := strings.Split(path, "OptiFine_")[1]
		v := strings.Split(s, ".jar")[0]
		return v
	}

	for _, file := range dir.File {
		if strings.Contains(file.Name, ".class") {
			jItem, err := file.Open()
			if err != nil {
				Log("Failed to open class")
			}
			str := bufio.NewScanner(jItem)
			var s strings.Builder
			for str.Scan() {
				readLine := str.Text()
				if strings.Contains(readLine, "Lnet/minecraftforge/fml/common/Mod;") && strings.Contains(readLine, "version") {
					for _, h := range str.Bytes() {
						if h >= 0x20 && h <= 0xf1 {
							s.WriteString(string(h))
						} else {
							s.WriteString(".")
						}
					}
					final := s.String()
					vVar := strings.Split(final, "version...")[1]
					vVer := strings.Split(vVar, "..")[0]
					return vVer
				}
			}
		}
	}
	return def
}

func getJarFiles(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		Log(err.Error())
	}
	return files
}
