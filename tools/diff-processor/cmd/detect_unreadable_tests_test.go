package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/google/go-cmp/cmp"
)

const readableServiceTestFile = `package readable_test

func TestAccReadable(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{Config: testAccReadableConfig()},
		},
	})
}

func testAccReadableConfig() string {
	return ` + "`" + `
resource "google_readable" "default" {
  name = "readable"
}
` + "`" + `
}
`

// The reader names the tests in a serial test map after their subtest
// functions, not after TestAccSerial.
const serialServiceTestFile = `package serial_test

func TestAccSerial(t *testing.T) {
	testCases := map[string]func(t *testing.T){
		"basic": testAccSerial_basic,
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc(t)
		})
	}
}

func testAccSerial_basic(t *testing.T) {
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{Config: testAccSerialConfig()},
		},
	})
}

func testAccSerialConfig() string {
	return ` + "`" + `
resource "google_serial" "default" {
  name = "serial"
}
` + "`" + `
}
`

const unreadableServiceTestFile = `package unreadable_test

func TestAccKnownUnreadable(t *testing.T) {
	config := testAccUnreadableConfig()
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{Config: config},
		},
	})
}

func TestAccNewUnreadable(t *testing.T) {
	config := testAccUnreadableConfig()
	acctest.VcrTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			{Config: config},
		},
	})
}

func testAccUnreadableConfig() string {
	return ` + "`" + `
resource "google_unreadable" "default" {
  name = "unreadable"
}
` + "`" + `
}
`

func writeServiceTestFile(t *testing.T, servicesDir, service, contents string) {
	t.Helper()
	dir := filepath.Join(servicesDir, service)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "resource_"+service+"_test.go"), []byte(contents), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectUnreadableTests(t *testing.T) {
	servicesDir := t.TempDir()
	writeServiceTestFile(t, servicesDir, "readable", readableServiceTestFile)
	writeServiceTestFile(t, servicesDir, "unreadable", unreadableServiceTestFile)
	writeServiceTestFile(t, servicesDir, "serial", serialServiceTestFile)
	// A service the change didn't touch: its unreadable tests aren't reported.
	writeServiceTestFile(t, servicesDir, "untouched", strings.ReplaceAll(unreadableServiceTestFile, "TestAccNewUnreadable", "TestAccUntouchedUnreadable"))

	cases := map[string]struct {
		services []string
		want     detector.UnreadableTestsReport
	}{
		"no services": {
			want: detector.UnreadableTestsReport{
				NewUnreadableTests:   map[string]string{},
				StaleUnreadableTests: []string{},
			},
		},
		"only the given services are read": {
			// "removed" has no directory, as for a service the change deletes.
			services: []string{"readable", "unreadable", "serial", "removed"},
			want: detector.UnreadableTestsReport{
				NewUnreadableTests: map[string]string{
					"TestAccNewUnreadable": "errors reading test func TestAccNewUnreadable: [errors reading test steps: [failed to find variable declaration config]]",
				},
				StaleUnreadableTests: []string{"TestAccReadable", "TestAccSerial"},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			o := detectUnreadableTestsOptions{
				readServiceTests: readServiceTests,
				knownTests:       []string{"TestAccKnownUnreadable", "TestAccReadable", "TestAccSerial"},
				stdout:           &buf,
			}
			if err := o.run(append([]string{servicesDir}, tc.services...)); err != nil {
				t.Fatalf("run() returned error: %v", err)
			}
			var got detector.UnreadableTestsReport
			if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
				t.Fatalf("failed to unmarshal output %q: %v", buf.String(), err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("output diff (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDetectUnreadableTestsServiceReadError(t *testing.T) {
	servicesDir := t.TempDir()
	// A file where a service directory is expected can't be listed.
	if err := os.WriteFile(filepath.Join(servicesDir, "notadir"), nil, 0644); err != nil {
		t.Fatal(err)
	}
	o := detectUnreadableTestsOptions{
		readServiceTests: readServiceTests,
		stdout:           &bytes.Buffer{},
	}
	err := o.run([]string{servicesDir, "notadir"})
	if err == nil || !strings.Contains(err.Error(), "error reading tests for service notadir") {
		t.Errorf("run() error = %v, want an error reading tests for service notadir", err)
	}
}
