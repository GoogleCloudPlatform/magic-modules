// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package changelog

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Entry struct {
	Issue string
	Body  string
	Date  time.Time
	Hash  string
}

// EntryList provides thread-safe operations on a list of Entry values
type EntryList struct {
	mu sync.RWMutex
	es []*Entry
}

type EntryErrorCode string

const (
	EntryErrorNotFound                            EntryErrorCode = "NOT_FOUND"
	EntryErrorUnknownTypes                        EntryErrorCode = "UNKNOWN_TYPES"
	EntryErrorInvalidNewReourceOrDatasourceFormat EntryErrorCode = "INVALID_NEW_RESOURCE_OR_DATASOURCE_FORMAT"
	EntryErrorMultipleLines                       EntryErrorCode = "MULTIPLE_LINES"
	EntryErrorInvalidEnhancementOrBugFixFormat    EntryErrorCode = "INVALID_ENHANCEMENT_OR_BUGFIX_FORMAT"
)

type EntryValidationError struct {
	message string
	Code    EntryErrorCode
	Details map[string]interface{}
}

func (e *EntryValidationError) Error() string {
	return e.message
}

// Validates that an Entry body contains properly formatted changelog notes
func (e *Entry) Validate() []*EntryValidationError {
	notes := NotesFromEntry(*e)

	var errors []*EntryValidationError

	if len(notes) < 1 {
		errors = append(errors, &EntryValidationError{
			message: fmt.Sprintf("no changelog entry found in: %s", string(e.Body)),
			Code:    EntryErrorNotFound,
		})
		return errors
	}

	for _, note := range notes {
		err := note.Validate()
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

// NewEntryList returns an EntryList with capacity c
func NewEntryList(c int) *EntryList {
	return &EntryList{
		es: make([]*Entry, 0, c),
	}
}

// Append appends entries to the EntryList
func (el *EntryList) Append(entries ...*Entry) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.es = append(el.es, entries...)
}

// Get returns the Entry at index i
func (el *EntryList) Get(i int) *Entry {
	el.mu.RLock()
	defer el.mu.RUnlock()
	if i >= len(el.es) || i < 0 {
		return nil
	}
	return el.es[i]
}

// Set sets the Entry at index i. The list will be resized if i is larger than
// the current list capacity.
func (el *EntryList) Set(i int, e *Entry) {
	if i < 0 {
		panic("invalid slice index")
	}
	el.mu.Lock()
	defer el.mu.Unlock()

	if i > (cap(el.es) - 1) {
		// resize the slice
		newEntries := make([]*Entry, i)
		copy(newEntries, el.es)
		el.es = newEntries
	}
	el.es[i] = e
}

// Len returns the number of items in the EntryList
func (el *EntryList) Len() int {
	el.mu.RLock()
	defer el.mu.RUnlock()
	return len(el.es)
}

// SortByIssue does an in-place sort of the entries by their issue number.
func (el *EntryList) SortByIssue() {
	el.mu.Lock()
	defer el.mu.Unlock()
	sort.Slice(el.es, func(i, j int) bool {
		return el.es[i].Issue < el.es[j].Issue
	})
}

type changelog struct {
	content []byte
	hash    string
	date    time.Time
}

// Diff returns the slice of Entry values that represent the difference of
// entries in the dir directory within repo from ref1 revision to ref2 revision.
// ref1 and ref2 should be valid git refs as strings and dir should be a valid
// directory path in the repository.
//
// The function calculates the diff by first checking out ref2 and collecting
// the set of all entries in dir. It then checks out ref1 and subtracts the
// entries found in dir. The resulting set of entries is then filtered to
// exclude any entries that came before the commit date of ref1.
//
// Along the way, if any git or filesystem interactions fail, an error is returned.
func Diff(repo, ref1, ref2, dir string) (*EntryList, error) {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL: repo,
	})
	if err != nil {
		return nil, err
	}
	rev2, err := r.ResolveRevision(plumbing.Revision(ref2))
	if err != nil {
		return nil, fmt.Errorf("could not resolve revision %s: %w", ref2, err)
	}
	var rev1 *plumbing.Hash
	if ref1 != "-" {
		rev1, err = r.ResolveRevision(plumbing.Revision(ref1))
		if err != nil {
			return nil, fmt.Errorf("could not resolve revision %s: %w", ref1, err)
		}
	}
	wt, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	if err := wt.Checkout(&git.CheckoutOptions{
		Hash:  *rev2,
		Force: true,
	}); err != nil {
		return nil, fmt.Errorf("could not checkout repository at %s: %w", ref2, err)
	}
	entriesAfterFI, err := wt.Filesystem.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("could not read repository directory %s: %w", dir, err)
	}
	// a set of all entries at rev2 (this release); the set of entries at ref1
	// will then be subtracted from it to arrive at a set of 'candidate' entries.
	entryCandidates := make(map[string]bool, len(entriesAfterFI))
	for _, i := range entriesAfterFI {
		entryCandidates[i.Name()] = true
	}
	if rev1 != nil {
		err = wt.Checkout(&git.CheckoutOptions{
			Hash:  *rev1,
			Force: true,
		})
		if err != nil {
			return nil, err
		}
		entriesBeforeFI, err := wt.Filesystem.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("could not read repository directory %s: %w", dir, err)
		}
		for _, i := range entriesBeforeFI {
			delete(entryCandidates, i.Name())
		}
		// checkout rev2 so that we can read files later
		if err := wt.Checkout(&git.CheckoutOptions{
			Hash:  *rev2,
			Force: true,
		}); err != nil {
			return nil, fmt.Errorf("could not checkout repository at %s: %w", ref2, err)
		}
	}

	entries := NewEntryList(len(entryCandidates))
	errg := new(errgroup.Group)
	for name := range entryCandidates {
		name := name // https://golang.org/doc/faq#closures_and_goroutines
		errg.Go(func() error {
			fp := filepath.Join(dir, name)
			f, err := wt.Filesystem.Open(fp)
			if err != nil {
				return fmt.Errorf("error opening file at %s: %w", name, err)
			}
			contents, err := ioutil.ReadAll(f)
			f.Close()
			if err != nil {
				return fmt.Errorf("error reading file at %s: %w", name, err)
			}
			log, err := r.Log(&git.LogOptions{FileName: &fp})
			if err != nil {
				return fmt.Errorf("error fetching git log for %s: %w", name, err)
			}
			lastChange, err := log.Next()
			if err != nil {
				return fmt.Errorf("error fetching next git log: %w", err)
			}
			entries.Append(&Entry{
				Issue: name,
				Body:  string(contents),
				Date:  lastChange.Author.When,
				Hash:  lastChange.Hash.String(),
			})
			return nil
		})
	}
	if err := errg.Wait(); err != nil {
		return nil, err
	}
	entries.SortByIssue()
	return entries, nil
}
