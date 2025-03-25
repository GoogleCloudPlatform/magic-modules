package cmd

import (
	"slices"

	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"

	newProvider "google/provider/new/google/provider"
	oldProvider "google/provider/old/google/provider"
)

const detectMissingDocDesc = `Compute list of fields missing documents`

type MissingDocsInfo struct {
	Name     string
	FilePath string
	Fields   []string
}

type detectMissingDocsOptions struct {
	rootOptions                 *rootOptions
	computeSchemaDiff           func() diff.SchemaDiff
	computeDatasourceSchemaDiff func() diff.SchemaDiff
	stdout                      io.Writer
}

type MissingDocsSummary struct {
	Resource   []detector.MissingDocDetails
	DataSource []detector.MissingDocDetails
}

func newDetectMissingDocsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingDocsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		computeDatasourceSchemaDiff: func() diff.SchemaDiff {
			return diff.ComputeSchemaDiff(oldProvider.DatasourceMap(), newProvider.DatasourceMap())
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

	datasourceSchemaDiff := o.computeDatasourceSchemaDiff()
	detectedDataSources, err := detector.DetectMissingDocsForDatasource(datasourceSchemaDiff, args[0])
	if err != nil {
		return err
	}

	sum := MissingDocsSummary{
		Resource:   sortMissingDocDetails(detectedResources),
		DataSource: sortMissingDocDetails(detectedDataSources),
	}

	if err := json.NewEncoder(o.stdout).Encode(sum); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}

	return nil
}

func sortMissingDocDetails(m map[string]detector.MissingDocDetails) []detector.MissingDocDetails {
	itemNames := maps.Keys(m)
	slices.Sort(itemNames)
	arr := []detector.MissingDocDetails{}
	for _, itemName := range itemNames {
		details := m[itemName]
		arr = append(arr, detector.MissingDocDetails{
			Name:     itemName,
			FilePath: details.FilePath,
			Fields:   details.Fields,
		})
	}
	return arr
}
