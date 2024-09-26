package cmd

import (
	"fmt"
)

func color(color, text string) string {
	if color == "" || text == "" {
		return text
	}
	return fmt.Sprintf("$\\textcolor{%s}{\\textsf{%s}}$", color, text)
}
