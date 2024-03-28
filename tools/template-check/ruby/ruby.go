package ruby

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

// Note: this is allowlisted to combat other issues like `=` instead of `==`, but it is possible we
// need to add more options to this list in the future, like `private`. The main thing we want to
// prevent is targeting `beta` in version guards, because it mishandles either `ga` or `private`.
var allowedGuards = []string{
	"<% unless version == 'ga' -%>",
	"<% if version == 'ga' -%>",
	"<% unless version == \"ga\" -%>",
	"<% if version == \"ga\" -%>",
}

// Note: this does not account for _every_ possible use of a version guard (for example, those
// starting with `version.nil?`), because the logic would start to get more complicated. Instead,
// the goal is to capture (and validate) all "standard" version guards that would be added for new
// resources/fields.
func isVersionGuard(line string) bool {
	re := regexp.MustCompile("<% [a-z]+ version ")
	return re.MatchString(line)
}

func isValidVersionGuard(line string) bool {
	for _, g := range allowedGuards {
		if strings.Contains(line, g) {
			return true
		}
	}
	return false
}

func CheckVersionGuards(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	lineNum := 1
	var invalidGuards []string
	for scanner.Scan() {
		if isVersionGuard(scanner.Text()) && !isValidVersionGuard(scanner.Text()) {
			invalidGuards = append(invalidGuards, fmt.Sprintf("%d: %s", lineNum, scanner.Text()))
		}
		lineNum++
	}
	return invalidGuards
}

func CheckVersionGuardsForFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return CheckVersionGuards(file), nil
}
