package cmd

import (
	oldConfig "google/provider/new/google/transport"
	newConfig "google/provider/old/google/transport"

	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const repDiffDesc = `Return a list of products where the regionalized endpoint default has changed.`

type repDiffOptions struct {
	rootOptions    *rootOptions
	repDefaultDiff func() (map[string]bool, map[string]bool)
	stdout         io.Writer
}

func newRepDefaultChangeCmd(rootOptions *rootOptions) *cobra.Command {
	o := &repDiffOptions{
		rootOptions: rootOptions,
		repDefaultDiff: func() (map[string]bool, map[string]bool) {
			return oldConfig.DefaultRepStatus(), newConfig.DefaultRepStatus()
		},
		stdout: os.Stdout,
	}
	cmd := &cobra.Command{
		Use:   "rep-diff",
		Short: repDiffDesc,
		Long:  repDiffDesc,
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *repDiffOptions) run() error {
	old, new := o.repDefaultDiff()
	var results []string
	for key, enabled := range new {
		oldEnabled, ok := old[key]
		if ok && oldEnabled != enabled {
			results = append(results, key)
		}
	}

	if err := json.NewEncoder(o.stdout).Encode(results); err != nil {
		return fmt.Errorf("Error encoding json: %w", err)
	}

	return nil
}
