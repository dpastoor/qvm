package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
	log "github.com/sirupsen/logrus"
)

func GetConfigPath() (string, error) {
	return xdg.ConfigFile("qvm/config.yml")
}

func GetRootConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "qvm")
}
func GetRootDataPath() string {
	return filepath.Join(xdg.DataHome, "qvm")
}

// GetPathToActiveBinDir returns the path to the active version's bin directory
// if that directory does not exist, it will be created with perms 700
func GetPathToActiveBinDir() string {
	path := filepath.Join(xdg.ConfigHome, "qvm", "bin")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0700)
		if err != nil {
			fmt.Println("unable to create active bin dir, there could be issues when attempting to access!")
		}
	}
	return path
}

func GetPathToVersionsDir() string {
	return filepath.Join(xdg.DataHome, "qvm", "versions")
}

// GetInstalledVersions returns a map of installed versions where they key is the version
// and the value is the path to the quarto executable
func GetInstalledVersions() (map[string]string, error) {
	iv := make(map[string]string)
	entries, err := os.ReadDir(GetPathToVersionsDir())
	if err != nil {
		return iv, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			quartoExe := "quarto"
			if runtime.GOOS == "windows" {
				quartoExe = "quarto.exe"
			}
			quartoPath := filepath.Join(GetPathToVersionsDir(), entry.Name(), "bin", quartoExe)
			if _, err := os.Stat(quartoPath); err == nil {
				iv[entry.Name()] = quartoPath
			} else {
				log.Warn("could not find expected quarto executable for version: ", entry.Name())
			}
		}
	}
	return iv, nil
}

func Read() (Config, string, error) {
	// for now don't check on error until consider what type of
	// logging - this should be completely optional anyway
	globalConfigPath, err := xdg.SearchConfigFile("qvm/config.yml")
	if err != nil {
		// this message is and print out all the paths that were searched
		// for example, by mangling the name to `onfig.yml` to see what the
		// error was, got the following result:
		//
		// could not locate `config.yml` in any of the following paths: /Users/devinp/Library/Application Support/qvm, /Users/devinp/Library/Preferences/qvm, /Library/Application Support/qvm, /Library/Preferences/qvm
		// exit status 1
		return Config{}, "", err
	}
	cfg, err := read(globalConfigPath)
	return cfg, globalConfigPath, err
}
