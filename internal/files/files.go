package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"thinker/internal/hash"
	"time"

	"golang.org/x/sync/errgroup"
)

type fileInfo struct {
	name       string
	path       string
	size       int64
	hash       []byte
	modifiedAt time.Time
}

type dirTree struct {
	basePath string
	tree     map[string]fileInfo
}

func checkPath(path string) error {
	dir, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("check path: %w", err)
	}

	if dir != nil && !dir.IsDir() {
		return errors.New("not a directory")
	}

	return nil
}

func getAbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func scanDir(dir string, dt *dirTree) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir: %w", err)
	}

	relPath, err := filepath.Rel(dt.basePath, dir)
	if err != nil {
		return fmt.Errorf("rel path: %w", err)
	}

	for _, file := range files {
		path := filepath.Join(dir, file.Name())

		if file.IsDir() {
			err = scanDir(path, dt)

			if err != nil {
				return err
			}
		} else {
			info, err := file.Info()
			if err != nil {
				return fmt.Errorf("file info: %w", err)
			}

			hash, err := getFileHash(path)
			if err != nil {
				return fmt.Errorf("get file hash: %w", err)
			}

			relFullPath := filepath.Join(relPath, file.Name())

			dt.tree[relFullPath] = fileInfo{
				name:       file.Name(),
				path:       relPath,
				size:       info.Size(),
				hash:       hash,
				modifiedAt: info.ModTime(),
			}
		}
	}

	return nil
}

func SyncDirs(src, dest string) error {
	srcDT, destDT, err := getDirTrees(src, dest)
	if err != nil {
		return err
	}

	err = syncDirTrees(srcDT, destDT)
	if err != nil {
		return err
	}

	err = syncDirTrees(destDT, srcDT)
	if err != nil {
		return err
	}

	return nil
}

func getDirTrees(source, dest string) (*dirTree, *dirTree, error) {
	absSource, err := getAbsPath(source)
	if err != nil {
		return nil, nil, fmt.Errorf("get abs path: %w", err)
	}

	err = checkPath(absSource)
	if err != nil {
		return nil, nil, err
	}

	absDest, err := getAbsPath(dest)
	if err != nil {
		return nil, nil, fmt.Errorf("get abs path: %w", err)
	}

	err = checkPath(absDest)
	if err != nil {
		return nil, nil, err
	}

	srcTree := dirTree{basePath: absSource, tree: make(map[string]fileInfo)}
	destTree := dirTree{basePath: absDest, tree: make(map[string]fileInfo)}

	var eg errgroup.Group

	eg.Go(func() error { return scanDir(absSource, &srcTree) })
	eg.Go(func() error { return scanDir(absDest, &destTree) })

	err = eg.Wait()

	if err != nil {
		return nil, nil, err
	}

	return &srcTree, &destTree, nil
}

func syncDirTrees(src, dest *dirTree) error {
	var eg errgroup.Group
	eg.SetLimit(runtime.NumCPU())
	for path, info := range src.tree {
		destInfo, ok := dest.tree[path]
		if !ok {
			eg.Go(func() error { return copyFile(src.basePath, dest.basePath, info) })
		} else if string(info.hash) != string(destInfo.hash) && info.modifiedAt.After(destInfo.modifiedAt) {
			eg.Go(func() error { return copyFile(src.basePath, dest.basePath, info) })
		}
	}

	err := eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dest string, info fileInfo) error {

	src = filepath.Join(src, info.path)
	dest = filepath.Join(dest, info.path)

	err := os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("make dir: %w", err)
	}

	srcFilePath := filepath.Join(src, info.name)
	destFilePath := filepath.Join(dest, info.name)

	fmt.Printf("Copied file: %s of size: %d\n", srcFilePath, info.size)

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return fmt.Errorf("open file for copy: %w", err)
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

func getFileHash(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	return hash.GetHash(file)
}
