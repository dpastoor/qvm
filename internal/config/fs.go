package config

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

func GetConfigPath() (string, error) {
	return xdg.ConfigFile("qvm/config.yml")
}

func GetRootConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "qvm")
}

func GetPathToActiveBinDir() string {
	return filepath.Join(xdg.DataHome, "qvm", "active", "bin")
}

func GetPathToVersionsDir() string {
	return filepath.Join(xdg.ConfigHome, "qvm", "versions")
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
