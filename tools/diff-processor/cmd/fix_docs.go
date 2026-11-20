package cmd

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/docfixer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cobra"

	newProvider "google/provider/new/google/provider"
)

const fixDocsDesc = `Add nested block anchors and links to handwritten resource docs so the missing doc detector can read them.

DOCS_DIR is a directory of resource docs, such as
mmv1/third_party/terraform/website/docs/r. Every resource with a doc there is
fixed, or only the given RESOURCEs. Prints what it changed and what it couldn't
resolve; pass --write to update the files.`

type fixDocsOptions struct {
	rootOptions *rootOptions
	resources   func() map[string]*schema.Resource
	stdout      io.Writer
	write       bool
}

func newFixDocsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &fixDocsOptions{
		rootOptions: rootOptions,
		resources:   newProvider.ResourceMap,
		stdout:      os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "fix-docs DOCS_DIR [RESOURCE...]",
		Short: "Add nested block anchors and links to handwritten resource docs",
		Long:  fixDocsDesc,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args[0], args[1:])
		},
	}
	cmd.Flags().BoolVar(&o.write, "write", false, "write the fixed docs back")
	return cmd
}

func (o *fixDocsOptions) run(dir string, names []string) error {
	resources := o.resources()
	all := len(names) == 0
	if all {
		names = slices.Sorted(maps.Keys(resources))
	}
	changed := 0
	for _, name := range names {
		r, ok := resources[name]
		if !ok {
			return fmt.Errorf("unknown resource %q", name)
		}
		path, err := resourceDoc(dir, name)
		if errors.Is(err, fs.ErrNotExist) && all {
			continue // generated or undocumented resource
		}
		if err != nil {
			return err
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		res := docfixer.Fix(string(src), r)
		base := filepath.Base(path)
		for _, c := range res.Changes {
			fmt.Fprintf(o.stdout, "%s: changed: %s\n", base, c)
		}
		for _, n := range res.Notes {
			fmt.Fprintf(o.stdout, "%s: note: %s\n", base, n)
		}
		if res.Text == string(src) {
			continue
		}
		changed++
		if o.write {
			if err := os.WriteFile(path, []byte(res.Text), 0644); err != nil {
				return err
			}
		}
	}
	verb := "would change"
	if o.write {
		verb = "changed"
	}
	fmt.Fprintf(o.stdout, "%s %d docs\n", verb, changed)
	return nil
}

// resourceDoc returns the path of the doc for resource name in dir.
func resourceDoc(dir, name string) (string, error) {
	for _, base := range []string{strings.TrimPrefix(name, "google_"), name} {
		path := filepath.Join(dir, base+".html.markdown")
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no doc for %s in %s: %w", name, dir, fs.ErrNotExist)
}
