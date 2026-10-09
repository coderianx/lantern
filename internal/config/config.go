package config

import (
	"os"
	"path/filepath"
)

var lanternPath string

func Init() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	lanternPath = filepath.Join(home, ".lantern")
	return os.MkdirAll(lanternPath, 0755)
}

func GetLanternPath() string {
	return lanternPath
}
