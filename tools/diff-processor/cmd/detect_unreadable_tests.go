package cmd

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/test-reader/reader"
	"github.com/spf13/cobra"
)

const detectUnreadableTestsDesc = "Report acceptance tests in the given services that the test reader can't read and that aren't known to be unreadable"

type detectUnreadableTestsOptions struct {
	rootOptions      *rootOptions
	readServiceTests func(serviceDir string) ([]string, map[string]error, error)
	knownTests       []string
	stdout           io.Writer
}

func newDetectUnreadableTestsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectUnreadableTestsOptions{
		rootOptions:      rootOptions,
		readServiceTests: readServiceTests,
		knownTests:       detector.KnownUnreadableTests,
		stdout:           os.Stdout,
	}
	return &cobra.Command{
		Use:   "detect-unreadable-tests SERVICES_DIR [SERVICE...]",
		Short: detectUnreadableTestsDesc,
		Long:  detectUnreadableTestsDesc,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
}

func (o *detectUnreadableTestsOptions) run(args []string) error {
	servicesDir, services := args[0], args[1:]
	var allTestNames []string
	allErrs := make(map[string]error)
	for _, service := range services {
		testNames, errs, err := o.readServiceTests(filepath.Join(servicesDir, service))
		if err != nil {
			return fmt.Errorf("error reading tests for service %s: %w", service, err)
		}
		allTestNames = append(allTestNames, testNames...)
		for name, err := range errs {
			allErrs[name] = err
		}
	}
	report := detector.DetectUnreadableTests(allTestNames, allErrs, o.knownTests)
	if err := json.NewEncoder(o.stdout).Encode(report); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}

// Read the acceptance tests in a service directory, returning the names of the
// acceptance test functions and the errors from reading them. A service
// directory that doesn't exist, such as one removed by the change under review,
// has no tests.
//
// The names come from the function declarations rather than the tests the
// reader returns, because the reader names the tests in a serial test map
// after their subtest functions.
func readServiceTests(serviceDir string) ([]string, map[string]error, error) {
	files, err := os.ReadDir(serviceDir)
	if os.IsNotExist(err) {
		return nil, nil, nil
	} else if err != nil {
		return nil, nil, err
	}
	var testFileNames []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), "_test.go") {
			testFileNames = append(testFileNames, filepath.Join(serviceDir, file.Name()))
		}
	}
	var testNames []string
	fset := token.NewFileSet()
	for _, fileName := range testFileNames {
		f, err := parser.ParseFile(fset, fileName, nil, 0)
		if err != nil {
			// The reader reports the same error, keyed by file name.
			continue
		}
		for _, decl := range f.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && strings.HasPrefix(funcDecl.Name.Name, "TestAcc") {
				testNames = append(testNames, funcDecl.Name.Name)
			}
		}
	}
	_, errs := reader.ReadTestFiles(testFileNames)
	return testNames, errs, nil
}
