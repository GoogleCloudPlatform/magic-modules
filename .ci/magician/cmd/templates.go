package cmd

import (
	"fmt"
)

func color(color, text string) string {
	if color == "" || text == "" {
		return text
	}
	var emoji string
	switch color {
	case "red":
		emoji = "🔴"
	case "yellow":
		emoji = "🟡"
	case "green":
		emoji = "🟢"
	default:
		emoji = ""
	}
	return fmt.Sprintf("%s%s%s", emoji, text, emoji)
}
