package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/docfixer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cobra"

	newProvider "google/provider/new/google/provider"
)

const fixDocsDesc = "Add nested block anchors and links to a handwritten resource doc so the missing doc detector can read it. Prints what it changed and what it couldn't resolve; pass --write to update the file."

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
		Use:   "fix-docs RESOURCE DOC_PATH",
		Short: fixDocsDesc,
		Long:  fixDocsDesc,
		Args:  cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
	cmd.Flags().BoolVar(&o.write, "write", false, "write the fixed doc back to DOC_PATH")
	return cmd
}

func (o *fixDocsOptions) run(args []string) error {
	name, path := args[0], args[1]
	r, ok := o.resources()[name]
	if !ok {
		return fmt.Errorf("unknown resource %q", name)
	}
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	res := docfixer.Fix(string(src), r)
	for _, c := range res.Changes {
		fmt.Fprintln(o.stdout, "changed:", c)
	}
	for _, n := range res.Notes {
		fmt.Fprintln(o.stdout, "note:", n)
	}
	if o.write && res.Text != string(src) {
		if err := os.WriteFile(path, []byte(res.Text), 0644); err != nil {
			return err
		}
	}
	return nil
}
