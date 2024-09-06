package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"thinker/internal/files"
)

func main() {
	var (
		source string
		dest   string
	)
	flag.StringVar(&source, "src", "", "Source folder")
	flag.StringVar(&dest, "dest", "", "Destination folder")
	flag.Parse()
	err := validateOptions(source, dest)
	if err != nil {
		log.Fatal(err)
	}
	absSource, err := files.GetAbsPath(source)
	if err != nil {
		log.Fatal(err)
	}
	err = files.CheckPath(absSource)
	if err != nil {
		log.Fatal(err)
	}
	absDest, err := files.GetAbsPath(dest)
	if err != nil {
		log.Fatal(err)
	}
	err = files.CheckPath(absDest)
	if err != nil {
		log.Fatal(err)
	}

	err = files.SyncDirs(absSource, absDest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Sinchronised")

}

func validateOptions(source, dest string) error {
	if source == "" || dest == "" {
		return errors.New("empty path")
	}
	if source == dest {
		return errors.New("source and destination are the same")
	}
	return nil
}
