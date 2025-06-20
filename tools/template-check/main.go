package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/template-check/gotemplate"
)

func isValidTemplate(filename string) (bool, error) {
	results, err := gotemplate.CheckVersionGuardsForFile(filename)
	if err != nil {
		return false, err
	}

	if len(results) > 0 {
		fmt.Fprintf(os.Stderr, "error: invalid version checks found in %s:\n", filename)
		for _, result := range results {
			fmt.Fprintf(os.Stderr, "  %s\n", result)
		}
		return false, nil
	}

	return true, nil
}

func checkTemplate(filename string) bool {
	valid, err := isValidTemplate(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}
	return valid
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "template-check - check that a template file is valid\n    template-check [file]\n")
	}

	flag.Parse()

	// Handle file as a positional argument
	if flag.Arg(0) != "" {
		if !checkTemplate(flag.Arg(0)) {
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Handle files coming from a linux pipe
	fileInfo, _ := os.Stdin.Stat()
	if fileInfo.Mode()&os.ModeCharDevice == 0 {
		exitStatus := 0
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for scanner.Scan() {
			if !checkTemplate(scanner.Text()) {
				exitStatus = 1
			}
		}

		os.Exit(exitStatus)
	}
}
