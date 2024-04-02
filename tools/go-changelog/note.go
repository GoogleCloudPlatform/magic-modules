// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package changelog

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Note struct {
	Type  string
	Body  string
	Issue string
	Hash  string
	Date  time.Time
}

var TypeValues = []string{
	"enhancement",
	"bug",
	"note",
	"none",
	"new-resource",
	"new-datasource",
	"deprecation",
	"breaking-change",
}

var textInBodyREs = []*regexp.Regexp{
	regexp.MustCompile("(?ms)^```release-note\r?\n(?P<note>.+?)\r?\n```"),
	regexp.MustCompile("(?ms)^```releasenote\r?\n(?P<note>.+?)\r?\n```"),
	regexp.MustCompile("(?ms)^```release-note:(?P<type>[^\r\n]*)\r?\n?(?P<note>.*?)\r?\n?```"),
	regexp.MustCompile("(?ms)^```releasenote:(?P<type>[^\r\n]*)\r?\n?(?P<note>.*?)\r?\n?```"),
}

var enhancementOrBugFixRegexp = regexp.MustCompile(`^[a-z0-9]+: .+$`)
var newResourceOrDatasourceRegexp = regexp.MustCompile("`google_[a-z0-9_]+`")
var newlineRegexp = regexp.MustCompile(`\n`)

func NotesFromEntry(entry Entry) []Note {
	var res []Note
	for _, re := range textInBodyREs {
		matches := re.FindAllStringSubmatch(entry.Body, -1)
		if len(matches) == 0 {
			continue
		}

		for _, match := range matches {
			note := ""
			typ := ""
			for i, name := range re.SubexpNames() {
				switch name {
				case "note":
					note = match[i]
				case "type":
					typ = match[i]
				}
				if note != "" && typ != "" {
					break
				}
			}

			typ = strings.TrimSpace(typ)

			if note == "" && typ == "" {
				continue
			}

			res = append(res, Note{
				Type:  typ,
				Body:  note,
				Issue: entry.Issue,
				Hash:  entry.Hash,
				Date:  entry.Date,
			})
		}
	}
	sort.Slice(res, SortNotes(res))
	return res
}

// Validates if a changelog note is properly formatted
func (n *Note) Validate() *EntryValidationError {
	typ := n.Type
	content := n.Body

	if !TypeValid(typ) {
		return &EntryValidationError{
			message: fmt.Sprintf("unknown changelog types %v: please use only the configured changelog entry types: %v", typ, content),
			Code:    EntryErrorUnknownTypes,
			Details: map[string]interface{}{
				"type": typ,
				"note": content,
			},
		}
	}

	if newlineRegexp.MatchString(content) {
		return &EntryValidationError{
			message: fmt.Sprintf("multiple lines are found in changelog entry %v: Please only have one CONTENT line per release note block. Use multiple blocks if there are multiple related changes in a single PR.", content),
			Code:    EntryErrorMultipleLines,
			Details: map[string]interface{}{
				"type": typ,
				"note": content,
			},
		}
	}

	if typ == "new-resource" || typ == "new-datasource" {
		if !newResourceOrDatasourceRegexp.MatchString(content) {
			return &EntryValidationError{
				message: fmt.Sprintf("invalid resource/datasource format in changelog entry %v: Please follow format in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/#type-specific-guidelines-and-examples", content),
				Code:    EntryErrorInvalidNewReourceOrDatasourceFormat,
				Details: map[string]interface{}{
					"type": typ,
					"note": content,
				},
			}
		}
	}

	if typ == "enhancement" || typ == "bug" {
		if !enhancementOrBugFixRegexp.MatchString(content) {
			return &EntryValidationError{
				message: fmt.Sprintf("invalid enhancement/bug fix format in changelog entry %v: Please follow format in https://googlecloudplatform.github.io/magic-modules/contribute/release-notes/#type-specific-guidelines-and-examples", content),
				Code:    EntryErrorInvalidEnhancementOrBugFixFormat,
				Details: map[string]interface{}{
					"type": typ,
					"note": content,
				},
			}
		}
	}
	return nil
}

func SortNotes(res []Note) func(i, j int) bool {
	return func(i, j int) bool {
		if res[i].Type < res[j].Type {
			return true
		} else if res[j].Type < res[i].Type {
			return false
		} else if res[i].Body < res[j].Body {
			return true
		} else if res[j].Body < res[i].Body {
			return false
		} else if res[i].Issue < res[j].Issue {
			return true
		} else if res[j].Issue < res[i].Issue {
			return false
		}
		return false
	}
}

func TypeValid(Type string) bool {
	for _, a := range TypeValues {
		if a == Type {
			return true
		}
	}
	return false
}
