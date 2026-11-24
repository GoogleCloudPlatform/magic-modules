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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"testing/fstest"
	"text/template"
)

func noFuncs() template.FuncMap { return template.FuncMap{} }

func execute(t *testing.T, tmpl *template.Template, name string, data any) string {
	t.Helper()
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		t.Fatalf("ExecuteTemplate(%s): %v", name, err)
	}
	return buf.String()
}

func TestParseTemplatesCached_MemoizesPerKey(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.tmpl"), []byte("a={{.}}"), 0644); err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(dir) // comparable, so memoization applies

	first, err := ParseTemplatesCached(fsys, "k1", "a.tmpl", noFuncs, "a.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	second, err := ParseTemplatesCached(fsys, "k1", "a.tmpl", noFuncs, "a.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Errorf("expected the same *Template for repeated key, got distinct instances")
	}
	if got := execute(t, second, "a.tmpl", 1); got != "a=1" {
		t.Errorf("unexpected output %q", got)
	}

	// A different key must not reuse the entry, even for the same files.
	calls := 0
	countingFuncs := func() template.FuncMap { calls++; return template.FuncMap{} }
	if _, err := ParseTemplatesCached(fsys, "k2", "a.tmpl", countingFuncs, "a.tmpl"); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseTemplatesCached(fsys, "k2", "a.tmpl", countingFuncs, "a.tmpl"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Errorf("funcs should be invoked once per key on a cache miss, got %d calls", calls)
	}
}

func TestParseTemplatesCached_DistinctFSDoNotCollide(t *testing.T) {
	dir1, dir2 := t.TempDir(), t.TempDir()
	for dir, body := range map[string]string{dir1: "one", dir2: "two"} {
		if err := os.WriteFile(filepath.Join(dir, "a.tmpl"), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}

	t1, err := ParseTemplatesCached(os.DirFS(dir1), "same-key", "a.tmpl", noFuncs, "a.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	t2, err := ParseTemplatesCached(os.DirFS(dir2), "same-key", "a.tmpl", noFuncs, "a.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if got := execute(t, t1, "a.tmpl", nil); got != "one" {
		t.Errorf("fs1: got %q, want %q", got, "one")
	}
	if got := execute(t, t2, "a.tmpl", nil); got != "two" {
		t.Errorf("fs2: got %q, want %q", got, "two")
	}
}

func TestParseTemplatesCached_UncomparableFSBypassesCache(t *testing.T) {
	// fstest.MapFS is a map and cannot be used as a map key; the cache must
	// degrade to parsing every time rather than panicking.
	fsys := fstest.MapFS{"a.tmpl": &fstest.MapFile{Data: []byte("x")}}

	first, err := ParseTemplatesCached(fsys, "k", "a.tmpl", noFuncs, "a.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	second, err := ParseTemplatesCached(fsys, "k", "a.tmpl", noFuncs, "a.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if first == second {
		t.Errorf("expected uncomparable FS to bypass memoization")
	}
}

func TestParseTemplatesCached_ReturnsParseErrors(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "bad.tmpl"), []byte("{{"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := ParseTemplatesCached(os.DirFS(dir), "bad", "bad.tmpl", noFuncs, "bad.tmpl"); err == nil {
		t.Errorf("expected a parse error")
	}
}

func TestParseTemplatesCached_ConcurrentCallersShareOneParse(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.tmpl"), []byte("{{.}}"), 0644); err != nil {
		t.Fatal(err)
	}
	fsys := os.DirFS(dir)

	var mu sync.Mutex
	calls := 0
	funcs := func() template.FuncMap {
		mu.Lock()
		calls++
		mu.Unlock()
		return template.FuncMap{}
	}

	var wg sync.WaitGroup
	results := make([]*template.Template, 32)
	for i := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			tmpl, err := ParseTemplatesCached(fsys, "concurrent", "a.tmpl", funcs, "a.tmpl")
			if err != nil {
				t.Error(err)
				return
			}
			results[i] = tmpl
			// Parsed templates may be executed concurrently.
			if got, want := execute(t, tmpl, "a.tmpl", i), fmt.Sprint(i); got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		}(i)
	}
	wg.Wait()

	if calls != 1 {
		t.Errorf("expected exactly one parse for concurrent callers, got %d", calls)
	}
	for i := 1; i < len(results); i++ {
		if results[i] != results[0] {
			t.Fatalf("result %d is a different template instance", i)
		}
	}
}
