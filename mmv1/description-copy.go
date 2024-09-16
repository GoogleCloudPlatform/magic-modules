package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

func CopyAllDescriptions(tempMode bool) {
	identifiers := []string{
		"description:",
		"note:",
		"set_hash_func:",
		"warning:",
		"required_properties:",
		"optional_properties:",
		"attributes:",
	}

	for i, id := range identifiers {
		CopyText(id, len(identifiers)-1 == i, tempMode)
	}
}

// Used to copy/paste text from Ruby -> Go YAML files
func CopyText(identifier string, last, tempMode bool) {
	var allProductFiles []string = make([]string, 0)
	glob := "products/**/go_product.yaml"
	if tempMode {
		glob = "products/**/*.temp"
	}
	files, err := filepath.Glob(glob)
	if err != nil {
		return
	}
	for _, filePath := range files {
		dir := filepath.Dir(filePath)
		productPath := fmt.Sprintf("products/%s", filepath.Base(dir))
		if !slices.Contains(allProductFiles, productPath) {
			allProductFiles = append(allProductFiles, productPath)
		}
	}

	for _, productPath := range allProductFiles {
		if strings.Contains(productPath, "healthcare") || strings.Contains(productPath, "memorystore") {
			continue
		}
		// Gather go and ruby file pairs
		yamlMap := make(map[string][]string)
		yamlPaths, err := filepath.Glob(fmt.Sprintf("%s/*", productPath))
		if err != nil {
			log.Fatalf("Cannot get yaml files: %v", err)
		}
		for _, yamlPath := range yamlPaths {
			if strings.HasSuffix(yamlPath, "_new") {
				continue
			}

			if tempMode {
				cutName, found := strings.CutSuffix(yamlPath, ".temp")
				if !found {
					continue
				}

				baseName := filepath.Base(yamlPath)
				yamlMap[baseName] = make([]string, 2)
				yamlMap[baseName][1] = yamlPath
				yamlMap[baseName][0] = cutName
				continue
			}

			fileName := filepath.Base(yamlPath)
			baseName, found := strings.CutPrefix(fileName, "go_")
			if yamlMap[baseName] == nil {
				yamlMap[baseName] = make([]string, 2)
			}
			if found {
				yamlMap[baseName][1] = yamlPath
			} else {
				yamlMap[baseName][0] = yamlPath
			}
		}

		for _, files := range yamlMap {
			rubyPath := files[0]
			goPath := files[1]
			var text []string
			currText := ""
			recording := false

			if strings.Contains(rubyPath, "product.yaml") {
				// log.Printf("skipping %s", rubyPath)
				continue
			}

			// Ready Ruby yaml
			file, _ := os.Open(rubyPath)
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, identifier) && !strings.HasPrefix(strings.TrimSpace(line), "#") {
					currText = strings.SplitAfter(line, identifier)[1]
					recording = true
				} else if recording {
					if terminateText(line) {
						text = append(text, currText)
						currText = ""
						recording = false
					} else {
						currText = fmt.Sprintf("%s\n%s", currText, line)
					}
				}
			}
			if recording {
				text = append(text, currText)
			}

			// Read Go yaml while writing to a temp file
			index := 0
			firstLine := true
			newFilePath := fmt.Sprintf("%s_new", goPath)
			fo, _ := os.Create(newFilePath)
			w := bufio.NewWriter(fo)
			file, _ = os.Open(goPath)
			defer file.Close()
			scanner = bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if firstLine {
					if strings.Contains(line, "NOT CONVERTED - RUN YAML MODE") {
						firstLine = false
						if !last {
							w.WriteString(fmt.Sprintf("NOT CONVERTED - RUN YAML MODE\n"))
						}
						continue
					} else {
						break
					}
				}
				if strings.Contains(line, identifier) {
					if index >= len(text) {
						log.Printf("did not replace %s correctly! Is the file named correctly?", goPath)
						w.Flush()
						break
					}
					line = fmt.Sprintf("%s%s", line, text[index])
					index += 1
				}
				w.WriteString(fmt.Sprintf("%s\n", line))
			}

			if !firstLine {
				if index != len(text) {
					log.Printf("potential issue with %s, only completed %d index out of %d replacements", goPath, index, len(text))
				}
				if err = w.Flush(); err != nil {
					panic(err)
				}

				// Overwrite original file with temp
				os.Rename(newFilePath, goPath)
			} else {
				os.Remove(newFilePath)
			}
		}

	}

}

// quick and dirty logic to determine if a description/note is terminated
func terminateText(line string) bool {
	terminalStrings := []string{
		"!ruby/",
	}

	for _, t := range terminalStrings {
		if strings.Contains(line, t) {
			return true
		}
	}

	if regexp.MustCompile(`^\s*https:[\s$]*`).MatchString(line) {
		return false
	}

	// Whole line comments
	if regexp.MustCompile(`^\s*#.*?`).MatchString(line) {
		return true
	}

	return regexp.MustCompile(`^\s*[a-z_]+:[\s$]*`).MatchString(line)
}
