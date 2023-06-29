package main

import (
	"flag"
	"fmt"
	"sort"
)

func main() {
	flag.Parse()
	breakages := compare()
	sort.Strings(breakages)
	for _, breakage := range breakages {
		fmt.Println(breakage)
	}
}
