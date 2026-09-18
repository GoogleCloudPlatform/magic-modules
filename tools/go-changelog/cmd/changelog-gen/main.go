// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	changelog "github.com/hashicorp/go-changelog"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage: changelog-gen <repo> <ref1> <ref2>\n")
	}

	repo := os.Args[1]
	ref1 := os.Args[2]
	ref2 := os.Args[3]

	entries, err := changelog.Diff(repo, ref1, ref2, ".changelog")
	if err != nil {
		log.Fatalf("Error getting changelog diff: %s", err)
	}

	var unknown, notes, breaking, deprecations, features, improvements, bugs []string

	for i := 0; i < entries.Len(); i++ {
		entry := entries.Get(i)
		notes_from_entry := changelog.NotesFromEntry(*entry)
		for _, note := range notes_from_entry {
			line := fmt.Sprintf("* %s", note.Body)
			switch note.Type {
			case "note":
				notes = append(notes, line)
			case "breaking-change":
				breaking = append(breaking, line)
			case "deprecation":
				deprecations = append(deprecations, line)
			case "new-resource", "new-datasource", "new-list-resource", "feature":
				features = append(features, line)
			case "improvement", "enhancement":
				improvements = append(improvements, line)
			case "bug":
				bugs = append(bugs, line)
			case "none":
				// skip
			default:
				unknown = append(unknown, line)
			}
		}
	}

	var sb strings.Builder

	if len(unknown) > 0 {
		sb.WriteString("UNKNOWN CHANGELOG TYPE:\n")
		for _, l := range unknown {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}
	if len(notes) > 0 {
		sb.WriteString("NOTES:\n")
		for _, l := range notes {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}
	if len(deprecations) > 0 {
		sb.WriteString("DEPRECATIONS:\n")
		for _, l := range deprecations {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}
	if len(breaking) > 0 {
		sb.WriteString("BREAKING CHANGES:\n")
		for _, l := range breaking {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}
	if len(features) > 0 {
		sb.WriteString("FEATURES:\n")
		for _, l := range features {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}
	if len(improvements) > 0 {
		sb.WriteString("IMPROVEMENTS:\n")
		for _, l := range improvements {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}
	if len(bugs) > 0 {
		sb.WriteString("BUG FIXES:\n")
		for _, l := range bugs {
			sb.WriteString(l + "\n")
		}
		sb.WriteString("\n")
	}

	fmt.Print(sb.String())
}
