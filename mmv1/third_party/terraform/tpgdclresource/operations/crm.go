package operations

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// CRMOperation can be parsed from the returned API operation and waited on.
// This is the typical GCP operation.
type CRMOperation struct {
	Name     string             `json:"name"`
	Error    *CRMOperationError `json:"error"`
	Done     bool               `json:"done"`
	Response map[string]any     `json:"response"`
	Metadata map[string]any     `json:"metadata"`
	// other irrelevant fields omitted

	config   *dcl.Config
	basePath string
	verb     string
	version  string

	response map[string]any
}

// CRMOperationError is the GCP operation's Error body.
type CRMOperationError struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Errors  []*CRMOperationErrorError `json:"errors"`
}

// String formats the CRMOperationError as an error string.
func (e *CRMOperationError) String() string {
	if e == nil {
		return "nil"
	}
	var b strings.Builder
	for _, err := range e.Errors {
		fmt.Fprintf(&b, "error code %q, message: %s\n", err.Code, err.Message)
	}
	if e.Code != 0 || e.Message != "" {
		fmt.Fprintf(&b, "error code %d, message: %s\n", e.Code, e.Message)
	}

	return b.String()
}

// CRMOperationErrorError is a singular error in a GCP operation.
type CRMOperationErrorError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Wait waits for an CRMOperation to complete by fetching the operation until it completes.
func (op *CRMOperation) Wait(ctx context.Context, c *dcl.Config, basePath, verb string) error {
	c.Logger.Infof("Waiting on operation: %v", op)
	op.config = c
	op.basePath = basePath
	op.verb = verb

	if len(op.Response) > 0 {
		op.response = op.Response
	}

	// base CRM resources use the v1 endpoint
	op.version = "v1"

	// Tags resources require the v3 endpoint, and DCL merges the two into one Operation handler. Identify
	// the operation kind by the "type" returned.
	if t, ok := op.Metadata["@type"].(string); ok && strings.HasPrefix(t, "type.googleapis.com/google.cloud.resourcemanager.v3") {
		op.version = "v3"
	}

	if op.Done {
		c.Logger.Infof("Completed operation: %v", op)
		return nil
	}

	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *CRMOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := dcl.URL(op.version+"/"+op.Name, op.basePath, op.config.BasePath, nil)
	resp, err := dcl.SendRequest(ctx, op.config, op.verb, u, &bytes.Buffer{}, nil)
	if err != nil {
		if dcl.IsRetryableRequestError(op.config, err, false, time.Now()) {
			return nil, dcl.OperationNotDone{}
		}
		return nil, err
	}

	if err := dcl.ParseResponse(resp.Response, op); err != nil {
		return nil, err
	}

	if !op.Done {
		return nil, dcl.OperationNotDone{}
	}

	if op.Error != nil {
		return nil, fmt.Errorf("operation received error: %+v", op.Error)
	}

	if len(op.response) == 0 && len(op.Response) > 0 {
		op.response = op.Response
	}

	return resp, nil
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *CRMOperation) FirstResponse() (map[string]any, bool) {
	return op.response, len(op.response) > 0
}
