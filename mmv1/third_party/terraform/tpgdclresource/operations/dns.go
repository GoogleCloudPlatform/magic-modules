package operations

import (
	"bytes"
	"context"
	"fmt"
	"time"

	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// DNSOperation can be parsed from the returned API operation and waited on.
// This is used for Changes only.
// Project and ManagedZone must be set ahead of time.
type DNSOperation struct {
	Status      string `json:"status"`
	ID          string `json:"id"`
	Project     string
	ManagedZone string
	// other irrelevant fields omitted

	config *dcl.Config
}

// Wait waits for an DNSOperation to complete by fetching the operation until it completes.
func (op *DNSOperation) Wait(ctx context.Context, c *dcl.Config, project, managedZone string) error {
	c.Logger.Infof("Waiting on operation: %v", op)
	op.config = c
	op.ManagedZone = managedZone
	op.Project = project

	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *DNSOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := fmt.Sprintf("https://dns.googleapis.com/dns/v1/projects/%s/managedZones/%s/changes/%s", op.Project, op.ManagedZone, op.ID)
	resp, err := dcl.SendRequest(ctx, op.config, "GET", u, &bytes.Buffer{}, nil)
	if err != nil {
		if dcl.IsRetryableRequestError(op.config, err, false, time.Now()) {
			return nil, dcl.OperationNotDone{}
		}
		return nil, err
	}
	if err := dcl.ParseResponse(resp.Response, op); err != nil {
		return nil, err
	}
	if op.Status != "done" {
		return nil, dcl.OperationNotDone{}
	}
	return resp, nil
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *DNSOperation) FirstResponse() (map[string]any, bool) {
	return make(map[string]any), false
}
