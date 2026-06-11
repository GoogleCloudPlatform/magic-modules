package main

import (
	"flag"
)

func main() {
	// Workaround for "ERROR: logging before flag.Parse"
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	_ = fs.Parse([]string{})
	flag.CommandLine = fs

	Execute()
}
