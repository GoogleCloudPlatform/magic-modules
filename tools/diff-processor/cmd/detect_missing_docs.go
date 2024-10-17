package cmd

import (
	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"
	"slices"
	"sort"

	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const detectMissingDocDesc = `Compute list of fields missing documents`

type MissingDocsInfo struct {
	Name     string
	FilePath string
	Fields   []detector.MissingDocField
}

type detectMissingDocsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	stdout            io.Writer
}

func newDetectMissingDocsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingDocsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.ResourceMap(), newProvider.ResourceMap())
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
	detectedResources, err := detector.DetectMissingDocs(schemaDiff, args[0])
	if err != nil {
		return err
	}
	resources := maps.Keys(detectedResources)
	slices.Sort(resources)
	info := []MissingDocsInfo{}
	for _, r := range resources {
		details := detectedResources[r]
		sort.Slice(details.Fields, func(i, j int) bool {
			return details.Fields[i].Field < details.Fields[j].Field
		})
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
