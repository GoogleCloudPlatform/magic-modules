package vcr

import (
	"fmt"
	"magician/provider"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/rand"
)

func TestFSNotify(t *testing.T) {
	folderPath, err := os.MkdirTemp(t.TempDir(), "testdir")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	defer w.Close()
	if err := w.Add(folderPath); err != nil {
		t.Fatalf("failed to add folder: %v", err)
	}

	var mu sync.Mutex
	visited := make(map[string]int)
	tester := &Tester{
		watcher: w,
		uploadFunc: func(head string, version provider.Version, filename string) error {
			mu.Lock()
			defer mu.Unlock()
			visited[filename] = visited[filename] + 1
			// simulate gsutil takes time to upload
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			return nil
		},
	}

	go tester.asyncUploadCassettes(provider.Beta, "branch", w)

	var wgWrite sync.WaitGroup
	for i := 0; i < 5; i++ {
		wgWrite.Add(1)
		go func(i int) {
			defer wgWrite.Done()
			f, err := os.OpenFile(filepath.Join(folderPath, fmt.Sprintf("%d.log", i+1)), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				fmt.Println(err)
			}
			defer f.Close()
			for j := 0; j < 1; j++ {
				if _, err := f.WriteString("abcdefg\n"); err != nil {
					fmt.Println("error writing file:", err)
				}
			}
		}(i)
	}
	wgWrite.Wait()
	// wait a bit so all events are received
	synced := false
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		mu.Lock()
		if got, want := len(visited), 5; want == got {
			synced = true
			break
		}
		mu.Unlock()
	}
	if !synced {
		mu.Lock()
		t.Errorf("want callback called for 5 files, got = %d", len(visited))
		mu.Unlock()
	}
}

func TestCollectResults(t *testing.T) {
	for _, test := range []struct {
		name     string
		output   string
		expected Result
	}{
		{
			name: "no compound tests",
			output: `--- FAIL: TestAccServiceOneResourceOne (100.00s)
--- PASS: TestAccServiceOneResourceTwo (100.00s)
--- PASS: TestAccServiceTwoResourceOne (100.00s)
--- PASS: TestAccServiceTwoResourceTwo (100.00s)
`,
			expected: Result{
				PassedTests: []string{"TestAccServiceOneResourceTwo", "TestAccServiceTwoResourceOne", "TestAccServiceTwoResourceTwo"},
				FailedTests: []string{"TestAccServiceOneResourceOne"},
			},
		},
		{
			name: "compound tests",
			output: `--- FAIL: TestAccServiceOneResourceOne (100.00s)
--- FAIL: TestAccServiceOneResourceTwo (100.00s)
    --- PASS: TestAccServiceOneResourceTwo/test_one (100.00s)
    --- FAIL: TestAccServiceOneResourceTwo/test_two (100.00s)
--- PASS: TestAccServiceTwoResourceOne (100.00s)
    --- PASS: TestAccServiceTwoResourceOne/test_one (100.00s)
    --- PASS: TestAccServiceTwoResourceOne/test_two (100.00s)
--- PASS: TestAccServiceTwoResourceTwo (100.00s)
`,
			expected: Result{
				PassedTests: []string{
					"TestAccServiceTwoResourceOne",
					"TestAccServiceTwoResourceTwo",
				},
				FailedTests: []string{"TestAccServiceOneResourceOne", "TestAccServiceOneResourceTwo"},
				PassedSubtests: []string{
					"TestAccServiceOneResourceTwo__test_one",
					"TestAccServiceTwoResourceOne__test_one",
					"TestAccServiceTwoResourceOne__test_two",
				},
				FailedSubtests: []string{"TestAccServiceOneResourceTwo__test_two"},
			},
		},
	} {
		if diff := cmp.Diff(test.expected, collectResult(test.output)); diff != "" {
			t.Errorf("collectResult(%q) got unexpected diff (-want +got):\n%s", test.output, diff)
		}
	}

}
