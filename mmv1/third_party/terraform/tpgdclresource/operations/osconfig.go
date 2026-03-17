package operations

import (
	"bytes"
	"context"

	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// OSPolicyAssignmentDeleteOperation can be parsed from the returned API operation and waited on.
type OSPolicyAssignmentDeleteOperation struct {
	Name string `json:"name"`

	config *dcl.Config
}

// Wait waits for an OSPolicyAssignmentDeleteOperation to complete by waiting until the operation returns a 404.
func (op *OSPolicyAssignmentDeleteOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	c.Logger.Infof("Waiting on: %q", op.Name)
	op.config = c

	return dcl.Do(ctx, op.operate, c.RetryProvider)
}

func (op *OSPolicyAssignmentDeleteOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := dcl.URL(op.Name, "https://osconfig.googleapis.com/v1alpha", op.config.BasePath, nil)
	resp, err := dcl.SendRequest(ctx, op.config, "GET", u, &bytes.Buffer{}, nil)
	if dcl.IsNotFound(err) {
		return nil, nil
	}
	return resp, dcl.OperationNotDone{}
}
