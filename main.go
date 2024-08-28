package main

import (
	"flag"
	"fmt"
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
	fmt.Println(opts)

}
