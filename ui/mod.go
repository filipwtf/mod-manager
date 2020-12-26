package ui

// Mod component
type Mod struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	latest  bool   `yaml:"latest"`
	Path    string `yaml:"path"`
}

// TODO Handle deletion of mods etc
