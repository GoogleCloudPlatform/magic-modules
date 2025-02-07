package cmd

import (
	"slices"
	"sort"

	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const detectMissingDocDesc = `Compute list of fields missing documents`

type MissingDocsInfo struct {
	Name     string
	FilePath string
	Fields   []string
}

type detectMissingDocsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	newResourceSchema map[string]*schema.Resource
	stdout            io.Writer
}

func newDetectMissingDocsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingDocsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "detect-missing-docs",
		Short: detectMissingDocDesc,
		Long:  detectMissingDocDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
	return cmd
}
func (o *detectMissingDocsOptions) run(args []string) error {
	schemaDiff := o.computeSchemaDiff()
	detectedResources, err := detector.DetectMissingDocs(schemaDiff, args[0], o.newResourceSchema)
	if err != nil {
		return err
	}
	resources := maps.Keys(detectedResources)
	slices.Sort(resources)
	info := []MissingDocsInfo{}
	for _, r := range resources {
		details := detectedResources[r]
		sort.Strings(details.Fields)
		info = append(info, MissingDocsInfo{
			Name:     r,
			FilePath: details.FilePath,
			Fields:   details.Fields,
		})
	}

	if err := json.NewEncoder(o.stdout).Encode(info); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}

	return nil
}
