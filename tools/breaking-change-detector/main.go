package main

import (
	"flag"
	"fmt"
	"sort"
)

var docMode = flag.Bool("docs", false, "legacy flag to not break existing CI can be removed after 7/10")

func main() {
	flag.Parse()
	if !*docMode {
		breakages := compare()
		sort.Strings(breakages)
		for _, breakage := range breakages {
			fmt.Println(breakage)
		}
	}
}
