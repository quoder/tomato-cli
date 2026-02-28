package storage

import (
	"os"
	"path/filepath"
)

func AppDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "tomato-cli"), nil
}
