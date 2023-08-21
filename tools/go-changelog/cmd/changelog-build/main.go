// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/go-changelog"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var lastRelease, thisRelease, repoDir, entriesDir, noteTmpl, changelogTmpl string
	flag.StringVar(&lastRelease, "last-release", "", "a git ref to the last commit in the previous release")
	flag.StringVar(&thisRelease, "this-release", "", "a git ref to the last commit to include in this release")
	flag.StringVar(&repoDir, "git-dir", pwd, "the directory of the git repo being released")
	flag.StringVar(&entriesDir, "entries-dir", "", "the directory within the repo containing changelog entry files")
	flag.StringVar(&noteTmpl, "note-template", "", "the path of the file holding the template to use for each item in the changelog")
	flag.StringVar(&changelogTmpl, "changelog-template", "", "the path of the file holding the template to use for the entire changelog")
	flag.Parse()

	if lastRelease == "" {
		fmt.Fprintln(os.Stderr, "Must specify last commit in the previous release.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	if thisRelease == "" {
		fmt.Fprintln(os.Stderr, "Must specify last commit in the release.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	if repoDir == "" {
		fmt.Fprintln(os.Stderr, "Must specify directory of the git repository being released.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	if entriesDir == "" {
		fmt.Fprintln(os.Stderr, "Must specify directory of the changelog entries within the repository being released.")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	if noteTmpl == "" {
		fmt.Fprintln(os.Stderr, "Must specify path to the file holding the template to use for each item in the changelog")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	if changelogTmpl == "" {
		fmt.Fprintln(os.Stderr, "Must specify path to the file holding the template to use for the entire changelog")
		fmt.Fprintln(os.Stderr, "")
		flag.Usage()
		os.Exit(1)
	}

	tmpl := template.New(filepath.Base(changelogTmpl)).Funcs(template.FuncMap{
		"sort": func(in []changelog.Note) []changelog.Note {
			sort.Slice(in, changelog.SortNotes(in))
			return in
		},
		"sortByDate": func(in []changelog.Note) []changelog.Note {
			sort.Slice(in, func(i, j int) bool {
				return in[i].Date.Before(in[j].Date)
			})
			return in
		},
		"combineTypes": func(in ...[]changelog.Note) []changelog.Note {
			count := 0
			for _, i := range in {
				count += len(i)
			}
			res := make([]changelog.Note, 0, count)
			for _, i := range in {
				res = append(res, i...)
			}
			return res
		},
		"stringHasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
	})
	tmpl, err = tmpl.ParseFiles(noteTmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %q as a Go template: %s\n", noteTmpl, err)
		os.Exit(1)
	}

	tmpl, err = tmpl.ParseFiles(changelogTmpl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %q as a Go template: %s\n", changelogTmpl, err)
		os.Exit(1)
	}

	entries, err := changelog.Diff(repoDir, lastRelease, thisRelease, entriesDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var notes []changelog.Note
	notesByType := map[string][]changelog.Note{}
	for i := 0; i < entries.Len(); i++ {
		entry := entries.Get(i)
		if strings.HasSuffix(entry.Issue, ".txt") {
			entry.Issue = strings.TrimSuffix(entry.Issue, ".txt")
		}
		notes = append(notes, changelog.NotesFromEntry(*entry)...)
	}
	for _, note := range notes {
		notesByType[note.Type] = append(notesByType[note.Type], note)
	}
	for _, n := range notesByType {
		sort.Slice(n, changelog.SortNotes(n))
	}
	sort.Slice(notes, changelog.SortNotes(notes))
	type renderData struct {
		Notes       []changelog.Note
		NotesByType map[string][]changelog.Note
	}
	err = tmpl.Execute(os.Stdout, renderData{
		Notes:       notes,
		NotesByType: notesByType,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing templates: %s\n", err)
		os.Exit(1)
	}
}
