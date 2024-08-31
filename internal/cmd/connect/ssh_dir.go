package connect

import (
	"os"
	"path/filepath"
)

// GetUserSSHDir returns the path to the user's SSH directory(e.g. ~/.ssh)
func GetUserSSHDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".ssh"), nil
}
