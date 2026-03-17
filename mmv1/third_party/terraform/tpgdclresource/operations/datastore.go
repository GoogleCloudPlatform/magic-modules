package operations

import (
	"bytes"
	"context"
	"fmt"
	"time"

	dcl "github.com/hashicorp/terraform-provider-google/google/tpgdclresource"
)

// DatastoreOperation can be parsed from the returned API operation and waited on.
type DatastoreOperation struct {
	Name     string                      `json:"name"`
	Done     bool                        `json:"done"`
	Metadata *DatastoreOperationMetadata `json:"metadata"`
	Error    *DatastoreOperationError    `json:"error"`
	config   *dcl.Config
}

// DatastoreOperationMetadata is an error in a datastore operation.
type DatastoreOperationMetadata struct {
	IndexID string `json:"indexId"`
}

// DatastoreOperationError is an error in a datastore operation.
type DatastoreOperationError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// Wait waits for an DatastoreOperation to complete by fetching the operation until it completes.
func (op *DatastoreOperation) Wait(ctx context.Context, c *dcl.Config, _, _ string) error {
	c.Logger.Infof("Waiting on operation: %v", op)
	op.config = c
	err := dcl.Do(ctx, op.operate, c.RetryProvider)
	c.Logger.Infof("Completed operation: %v", op)
	return err
}

func (op *DatastoreOperation) operate(ctx context.Context) (*dcl.RetryDetails, error) {
	u := dcl.URL(op.Name, "https://datastore.googleapis.com/v1/", op.config.BasePath, nil)
	resp, err := dcl.SendRequest(ctx, op.config, "GET", u, &bytes.Buffer{}, nil)
	if err != nil {
		if dcl.IsRetryableRequestError(op.config, err, true, time.Now()) {
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
	return resp, nil
}

// FirstResponse returns the first response that this operation receives with the resource.
// This response may contain special information.
func (op *DatastoreOperation) FirstResponse() (map[string]any, bool) {
	return map[string]any{
		"indexId": op.Metadata.IndexID,
	}, false
}
