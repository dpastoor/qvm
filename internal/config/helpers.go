package config

import (
	"errors"
	"os"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return errors.Is(err, os.ErrNotExist)
}
