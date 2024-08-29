package files

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type fileInfo struct {
	name       string
	path       string
	size       int64
	modifiedAt time.Time
}

var (
	srcTree  = make(map[string]fileInfo)
	destTree = make(map[string]fileInfo)
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

func ScanDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			fmt.Println(path + "/")
			err = ScanDir(path)
			if err != nil {
				return err
			}
		} else {
			fmt.Println(path)
		}
	}

	return nil
}

func SyncDirs(src, dest string) error {
	return nil
}
