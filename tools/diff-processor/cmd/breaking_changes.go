package cmd
import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/GoogleCloudPlatform/magic-modules/.ci/diff-processor/diff"
)
const breakingChangesDesc = `Check for breaking changes between the new / old Terraform provider versions.`
type breakingChangesOptions struct {
	rootOptions           *rootOptions
}
func newBreakingChangesCmd(rootOptions *rootOptions) *cobra.Command {
	o := &breakingChangesOptions{
		rootOptions:           rootOptions,
	}
	cmd := &cobra.Command{
		Use:   "breaking-changes",
		Short: breakingChangesDesc,
		Long:  breakingChangesDesc,
		RunE: func(c *cobra.Command, args []string) error {
			return o.run()
		},
	}
	return cmd
}
func (o *breakingChangesOptions) run() error {
	breakages := diff.Compare()
	sort.Strings(breakages)
	for _, breakage := range breakages {
		fmt.Println(breakage)
	}
	return nil
}