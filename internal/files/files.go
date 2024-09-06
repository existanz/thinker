package files

import (
	"errors"
	"fmt"
	"io"
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

type dirTree map[string]fileInfo

func CheckPath(path string) error {
	dir, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("check path: %w", err)
	}
	if dir != nil && !dir.IsDir() {
		return errors.New("not a directory")
	}
	return nil
}

func GetAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func ScanDir(base, dir string, dt *dirTree) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}
	relPath, err := filepath.Rel(base, dir)
	if err != nil {
		return fmt.Errorf("rel path: %w", err)
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			err = ScanDir(base, path, dt)
			if err != nil {
				return err
			}
		} else {
			info, err := file.Info()
			if err != nil {
				return fmt.Errorf("file info: %w", err)
			}
			(*dt)[path] = fileInfo{
				name:       file.Name(),
				path:       relPath,
				size:       info.Size(),
				modifiedAt: info.ModTime(),
			}
		}
	}

	return nil
}

func SyncDirs(source, dest string) error {
	absSource, err := GetAbsPath(source)
	if err != nil {
		return fmt.Errorf("get abs path: %w", err)
	}
	err = CheckPath(absSource)
	if err != nil {
		return err
	}
	absDest, err := GetAbsPath(dest)
	if err != nil {
		return fmt.Errorf("get abs path: %w", err)
	}
	err = CheckPath(absDest)
	if err != nil {
		return err
	}
	srcTree := dirTree{}
	destTree := dirTree{}

	err = ScanDir(absSource, absSource, &srcTree)
	if err != nil {
		return err
	}
	err = ScanDir(absDest, absDest, &destTree)
	if err != nil {
		return err
	}
	for path, info := range srcTree {
		destInfo, ok := destTree[path]
		if !ok {
			err = CopyFile(absSource, absDest, info)
			if err != nil {
				return err
			}
		} else {
			if info.modifiedAt.After(destInfo.modifiedAt) {
				err = CopyFile(path, absDest, info)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func CopyFile(src, dest string, info fileInfo) error {
	if info.path != "." {
		src = filepath.Join(src, info.path)
		dest = filepath.Join(dest, info.path)
		fmt.Println(filepath.Dir(dest), dest)
		err := os.MkdirAll(dest, os.ModePerm)
		if err != nil {
			return fmt.Errorf("make dir: %w", err)
		}
	}
	fmt.Println(info.path, info.name)
	srcFilePath := filepath.Join(src, info.name)
	destFilePath := filepath.Join(dest, info.name)
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	return nil
}
