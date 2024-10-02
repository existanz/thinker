package files

import (
	"os"
	"testing"
	"time"
)

func TestSyncDirs(t *testing.T) {
	src := t.TempDir()
	dest := t.TempDir()

	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = copyFile(src, dest, fileInfo{
		name:       "test.txt",
		path:       "test.txt",
		size:       10,
		hash:       []byte("test"),
		modifiedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}

	err = SyncDirs(src, dest)
	if err != nil {
		t.Fatal(err)
	}
}
