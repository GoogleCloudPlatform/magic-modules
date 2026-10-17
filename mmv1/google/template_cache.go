// Copyright 2026 Google Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package google

import (
	"io/fs"
	"reflect"
	"sync"
	"text/template"
)

// templateCache memoizes parsed templates for a single fs.FS.
//
// Template files never change during a generation run, but the generator
// historically re-read and re-parsed the same template set for every output
// file (and, via trimTemplate/customTemplate, for every property). Parsing
// accounted for roughly a third of total CPU time. Parsed templates are safe
// to execute concurrently, so parsing each distinct template set once and
// sharing the result is both safe and a large win.
type templateCache struct {
	mu      sync.Mutex
	entries map[string]*templateCacheEntry
}

type templateCacheEntry struct {
	once sync.Once
	tmpl *template.Template
	err  error
}

// templateCaches maps an fs.FS to its templateCache. Keyed by the FS value
// itself so that distinct filesystems (e.g. in tests) never share entries.
var templateCaches sync.Map

// cacheFor returns the templateCache for fsys. If fsys is not a comparable
// type (for example fstest.MapFS) it cannot be used as a map key, so a fresh,
// unshared cache is returned and no memoization takes place.
func cacheFor(fsys fs.FS) *templateCache {
	if fsys == nil || !reflect.TypeOf(fsys).Comparable() {
		return &templateCache{entries: map[string]*templateCacheEntry{}}
	}
	if c, ok := templateCaches.Load(fsys); ok {
		return c.(*templateCache)
	}
	c, _ := templateCaches.LoadOrStore(fsys, &templateCache{entries: map[string]*templateCacheEntry{}})
	return c.(*templateCache)
}

// ParseTemplatesCached parses files from fsys into a template named name with
// the functions returned by funcs, memoizing the result per (fsys, key).
//
// key must uniquely identify the parsed result: callers should derive it from
// the list of files and from anything that influences funcs (such as a
// closure over the primary template path). funcs is only invoked on a cache
// miss. The returned template must not be mutated by callers; it may be
// executed concurrently.
func ParseTemplatesCached(fsys fs.FS, key, name string, funcs func() template.FuncMap, files ...string) (*template.Template, error) {
	c := cacheFor(fsys)

	c.mu.Lock()
	e, ok := c.entries[key]
	if !ok {
		e = &templateCacheEntry{}
		c.entries[key] = e
	}
	c.mu.Unlock()

	e.once.Do(func() {
		e.tmpl, e.err = template.New(name).Funcs(funcs()).ParseFS(fsys, files...)
	})
	return e.tmpl, e.err
}
