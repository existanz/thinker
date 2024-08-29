package files

import (
	"errors"
	"os"
	"path/filepath"
)

func CheckPath(path string) error {
	dir, err := os.Stat(path)
	if err != nil {
		return err
	}
	if dir != nil && !dir.IsDir() {
		return errors.New("not a directory")
	}
	return nil
}

func GetAbsPath(path string) (string, error) {
	if path == "" {
		return "", errors.New("empty path")
	}
	return filepath.Abs(path)
}
