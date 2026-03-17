package operations

import (
	"context"
	"fmt"
	"strings"

	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// MonitoringOperation can be parsed from the returned API operation and waited on.
type MonitoringOperation struct {
	Name string `json:"name"`
}

// Wait waits for an MonitoringOperation to complete by fetching the operation until it completes.
func (op *MonitoringOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	if op.Name != "" {
		// Names come in the form "accessPolicies/{{name}}"
		parts := strings.Split(op.Name, "/")
		op.Name = parts[len(parts)-1]
	}
	return nil
}

// FetchName will fetch the operation and return the name of the resource created.
// Monitoring creates resources with machine generated names.
// It must be called after the resource has been created.
func (op *MonitoringOperation) FetchName() (*string, error) {
	if op.Name == "" {
		return nil, fmt.Errorf("this operation (%s) has no name and probably hasn't been run before", op.Name)
	}
	return &op.Name, nil
}
