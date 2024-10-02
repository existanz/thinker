package hash

import (
	"encoding/hex"
	"testing"
)

func TestMD5(t *testing.T) {
	data := []byte("test")
	hash := MD5(data)
	want := "098f6bcd4621d373cade4e832627b4f6"
	wantLen := 16

	if len(hash) != wantLen {
		t.Errorf("wrong hash length expected %d, got %d", wantLen, len(hash))
	}

	if string(hex.EncodeToString(hash)) != want {
		t.Errorf("wrong hash expected %s, got %s", want, hex.EncodeToString(hash))
	}
}
