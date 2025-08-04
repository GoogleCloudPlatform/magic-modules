package apigee

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/googleapi"
)

func RetryWithContext(ctx context.Context, opt transport_tpg.RetryOptions) error {
	doneCh := make(chan error, 1)

	go func() {
		doneCh <- transport_tpg.Retry(opt)
	}()

	select {
	case err := <-doneCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func sendRequestRawBodyFramework(ctx context.Context, opts SendRequestRawBodyOptions) (map[string]interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	tflog.Trace(ctx, "Executing raw body request", map[string]interface{}{
		"url":    opts.RawURL,
		"method": opts.Method,
	})

	reqHeaders := make(http.Header)
	reqHeaders.Set("User-Agent", opts.UserAgent)
	reqHeaders.Set("Content-Type", opts.ContentType)

	if opts.Config.UserProjectOverride && opts.Project != "" {
		reqHeaders.Set("X-Goog-User-Project", opts.Project)
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 2 * time.Minute
	}

	var httpResp *http.Response

	err := RetryWithContext(ctx, transport_tpg.RetryOptions{
		RetryFunc: func() error {
			client := opts.Config.Client

			if bodySeeker, ok := opts.Body.(io.ReadSeeker); ok {
				bodySeeker.Seek(0, io.SeekStart)
			}

			httpReq, err := http.NewRequestWithContext(ctx, opts.Method, opts.RawURL, opts.Body)
			if err != nil {
				return fmt.Errorf("error creating HTTP request: %w", err)
			}
			httpReq.Header = reqHeaders

			httpResp, err = client.Do(httpReq)
			if err != nil {
				return fmt.Errorf("error sending request: %w", err)
			}

			if err := googleapi.CheckResponse(httpResp); err != nil {
				return err
			}

			return nil
		},
		Timeout: timeout,
	})

	if err != nil {
		diags.AddError("API Request Failed", fmt.Sprintf("request to %s failed with retries: %s", opts.RawURL, err.Error()))
		return nil, diags
	}

	if httpResp == nil {
		diags.AddError("API Response Error", "Request was successful, but the HTTP response was nil.")
		return nil, diags
	}
	defer googleapi.CloseBody(httpResp)

	if httpResp.StatusCode == http.StatusNoContent {
		return nil, diags
	}

	result := make(map[string]interface{})
	if err := json.NewDecoder(httpResp.Body).Decode(&result); err != nil {
		diags.AddError("API Response Decode Error", fmt.Sprintf("failed to decode JSON response body: %s", err.Error()))
		return nil, diags
	}

	tflog.Trace(ctx, "Raw body request successful")
	return result, diags
}

type SendRequestRawBodyOptions struct {
	Config      *transport_tpg.Config
	Method      string
	Project     string
	RawURL      string
	UserAgent   string
	Body        io.Reader
	ContentType string
	Timeout     time.Duration
}
