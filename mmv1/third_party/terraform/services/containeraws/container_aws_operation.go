package containeraws

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type ContainerAwsOperationWaiter struct {
	Config    *transport_tpg.Config
	UserAgent string
	Project   string
	tpgresource.CommonOperationWaiter
}

func (w *ContainerAwsOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}

	region := tpgresource.GetRegionFromRegionalSelfLink(w.CommonOperationWaiter.Op.Name)

	// Returns the proper get.
	url := fmt.Sprintf("%s%s", w.Config.ContainerAwsBasePath, w.Op.Name)
	if strings.Contains(w.Config.ContainerAwsBasePath, "https://{{location}}") {
		url = fmt.Sprintf("https://%s-gkemulticloud.googleapis.com/v1/%s", region, w.CommonOperationWaiter.Op.Name)
	}

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    url,
		UserAgent: w.UserAgent,
	})
}

func createContainerAwsWaiter(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string) (*ContainerAwsOperationWaiter, error) {
	w := &ContainerAwsOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Project:   project,
	}
	if err := w.CommonOperationWaiter.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

func ContainerAwsOperationWaitTimeWithResponse(config *transport_tpg.Config, op map[string]interface{}, response *map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	w, err := createContainerAwsWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	if err := tpgresource.OperationWait(w, activity, timeout, config.PollInterval); err != nil {
		return err
	}
	return json.Unmarshal([]byte(w.CommonOperationWaiter.Op.Response), response)
}

func ContainerAwsOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createContainerAwsWaiter(config, op, project, activity, userAgent)
	if err != nil {
		// If w is nil, the op was synchronous.
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
