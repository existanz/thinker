package main

import (
	"testing"
)

func TestValidateOptions(t *testing.T) {
	source := "src"
	dest := "dest"
	err := validateOptions(source, dest)
	if err != nil {
		t.Fatal(err)
	}

	source = ""
	err = validateOptions(source, dest)
	if err == nil {
		t.Fatal("expected error")
	}

	source = "dest"
	err = validateOptions(source, dest)
	if err == nil {
		t.Fatal("expected error")
	}
}
