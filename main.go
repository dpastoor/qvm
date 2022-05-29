package main

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	"github.com/dpastoor/qvm/cmd"
)

// https://goreleaser.com/cookbooks/using-main.version
var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

func main() {
	// normalize the config home on osx to linux and get rid of the path spacing
	if runtime.GOOS == "darwin" {
		hd, _ := os.UserHomeDir()
		os.Setenv("XDG_CONFIG_HOME", filepath.Join(hd, ".config"))
		os.Setenv("XDG_DATA_HOME", filepath.Join(hd, ".config"))
		xdg.Reload()
	}
	cmd.Execute(version, os.Args[1:])
}
