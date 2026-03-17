package operations

import (
	"bytes"
	"context"
	"time"

	glog "github.com/golang/glog"
	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// SQLOperation can be parsed from the returned API operation and waited on.
type SQLOperation struct {
	ID         string `json:"id"`
	SelfLink   string `json:"selfLink"`
	Status     string `json:"status"`
	TargetLink string `json:"targetLink"`
	// other irrelevant fields omitted

	config *dcl.Config
}

// Wait waits for an Operation to complete by fetching the operation until it completes.
func (op *SQLOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	glog.Infof("Waiting on operation: %v", op)
	op.config = c

	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *SQLOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	resp, err := dcl.SendRequest(ctx, op.config, "GET", op.SelfLink, &bytes.Buffer{}, nil)
	if err != nil {
		if dcl.IsRetryableRequestError(op.config, err, true, time.Now()) {
			return nil, dcl.OperationNotDone{}
		}
		return nil, err
	}
	if err := dcl.ParseResponse(resp.Response, op); err != nil {
		return nil, err
	}
	if op.Status != "DONE" {
		return nil, dcl.OperationNotDone{}
	}
	return resp, nil
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *SQLOperation) FirstResponse() (map[string]any, bool) {
	return make(map[string]any), false
}

// SQLCreateCertOperation is the operation used for creating SSL certs.
// They have a different format from other resources and other methods.
type SQLCreateCertOperation struct {
	Operation  SQLOperation `json:"operation"`
	ClientCert struct {
		CertInfo map[string]any `json:"certInfo"`
	} `json:"clientCert"`
	response map[string]any
}

// Wait waits for an SQLOperation to complete by fetching the operation until it completes.
func (op *SQLCreateCertOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	return op.Operation.Wait(ctx, c, "", "")
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *SQLCreateCertOperation) FirstResponse() (map[string]any, bool) {
	if len(op.ClientCert.CertInfo) > 0 {
		return op.ClientCert.CertInfo, true
	}
	return make(map[string]any), false
}
