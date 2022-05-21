package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/dpastoor/qvm/internal/config"
	log "github.com/sirupsen/logrus"

	"github.com/adrg/xdg"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}

func getDefaultConfigPath() (string, error) {
	return xdg.ConfigFile("qvm/config.yml")
}

func readConfig() (config.Config, string, error) {
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
		return config.Config{}, "", err
	}
	cfg, err := config.Read(globalConfigPath)
	return cfg, globalConfigPath, err
}

func setLogLevel(level string) {
	// We want the log to be reset whenever it is initialized.
	logLevel := strings.ToLower(level)

	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.Fatalf("Invalid log level: %s", logLevel)
	}
}
