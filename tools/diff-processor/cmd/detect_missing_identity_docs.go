package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/detector"
	"github.com/GoogleCloudPlatform/magic-modules/tools/diff-processor/diff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/spf13/cobra"
	"golang.org/x/exp/maps"
)

const detectMissingIdentityDocsDesc = `Compute list of resources with undocumented ResourceIdentity fields`

type detectMissingIdentityDocsOptions struct {
	rootOptions       *rootOptions
	computeSchemaDiff func() diff.SchemaDiff
	resourceMap       func() map[string]*schema.Resource
	stdout            io.Writer
}

func newDetectMissingIdentityDocsCmd(rootOptions *rootOptions) *cobra.Command {
	o := &detectMissingIdentityDocsOptions{
		rootOptions: rootOptions,
		computeSchemaDiff: func() diff.SchemaDiff {
			return schemaDiff
		},
		resourceMap: func() map[string]*schema.Resource {
			return newResourceMap
		},
		stdout: os.Stdout,
	}
	return &cobra.Command{
		Use:   "detect-missing-identity-docs REPO_PATH",
		Short: detectMissingIdentityDocsDesc,
		Long:  detectMissingIdentityDocsDesc,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return o.run(args)
		},
	}
}

func (o *detectMissingIdentityDocsOptions) run(args []string) error {
	sd := o.computeSchemaDiff()
	detected, err := detector.DetectMissingIdentityDocs(sd, o.resourceMap(), args[0])
	if err != nil {
		return fmt.Errorf("error detecting missing identity docs: %w", err)
	}

	itemNames := maps.Keys(detected)
	slices.Sort(itemNames)
	result := make([]detector.MissingDocDetails, 0, len(itemNames))
	for _, name := range itemNames {
		result = append(result, detected[name])
	}

	if err := json.NewEncoder(o.stdout).Encode(result); err != nil {
		return fmt.Errorf("error encoding json: %w", err)
	}
	return nil
}
