package main

import (
	"errors"
	"flag"
	"fmt"
	"thinker/internal/files"
)

type Options struct {
	Source string
	Dest   string
}

func main() {

	var opts Options
	flag.StringVar(&opts.Source, "src", "", "Source folder")
	flag.StringVar(&opts.Dest, "dest", "", "Destination folder")
	flag.Parse()
	err := ValidateOptions(&opts)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(opts)

}

func ValidateOptions(opts *Options) error {
	absSrc, err := files.GetAbsPath(opts.Source)
	if err != nil {
		return err
	}
	absDest, err := files.GetAbsPath(opts.Dest)
	if err != nil {
		return err
	}
	if err := files.CheckPath(absSrc); err != nil {
		return err
	}
	if err := files.CheckPath(absDest); err != nil {
		return err
	}
	if absSrc == absDest {
		return errors.New("source and destination are the same")
	}
	opts.Source = absSrc
	opts.Dest = absDest
	return nil
}
