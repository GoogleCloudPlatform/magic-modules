package vcr

import (
	"fmt"
	"io/fs"
	"regexp"
	"strings"
)

func (vt *Tester) filterTraceFromLogFiles(logPath string) error {
	return vt.rnr.Walk(logPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		// Only process .log files
		if !strings.HasSuffix(info.Name(), ".log") {
			return nil
		}

		// Read the log file
		content, err := vt.rnr.ReadFile(path)
		if err != nil {
			return nil // Skip files we can't read
		}

		// Filter TRACE lines
		filteredContent := filterTraceLines(content)

		// Write back the filtered content
		if err := vt.rnr.WriteFile(path, filteredContent); err != nil {
			fmt.Printf("Warning: could not filter log file %s: %v\n", path, err)
		}

		return nil
	})
}

func filterTraceLines(output string) string {
	lines := strings.Split(output, "\n")
	var filtered []string
	inTraceBlock := false
	tracePattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z \[TRACE\]`)
	timestampPattern := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z`)

	for _, line := range lines {
		if timestampPattern.MatchString(line) {
			inTraceBlock = tracePattern.MatchString(line)
			if !inTraceBlock {
				filtered = append(filtered, line)
			}
		} else {
			if !inTraceBlock {
				filtered = append(filtered, line)
			}
		}
	}

	return strings.Join(filtered, "\n")
}
