package agentregistry

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

type AgentRegistryOperationWaiter struct {
	Config    *transport_tpg.Config
	UserAgent string
	Project   string
	tpgresource.CommonOperationWaiter
}

func (w *AgentRegistryOperationWaiter) QueryOp() (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Cannot query operation, it's unset or nil.")
	}
	url := fmt.Sprintf("%s%s", transport_tpg.BaseUrl(Product, w.Config), w.OpName())

	return transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    w.Config,
		Method:    "GET",
		Project:   w.Project,
		RawURL:    url,
		UserAgent: w.UserAgent,
	})
}

func createAgentRegistryWaiter(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string) (*AgentRegistryOperationWaiter, error) {
	w := &AgentRegistryOperationWaiter{
		Config:    config,
		UserAgent: userAgent,
		Project:   project,
	}
	if err := w.SetOp(op); err != nil {
		return nil, err
	}
	return w, nil
}

func AgentRegistryOperationWaitTime(config *transport_tpg.Config, op map[string]interface{}, project, activity, userAgent string, timeout time.Duration) error {
	if val, ok := op["name"]; !ok || val == "" {
		// This was a synchronous call - there is no operation to wait for.
		return nil
	}
	w, err := createAgentRegistryWaiter(config, op, project, activity, userAgent)
	if err != nil {
		return err
	}
	return tpgresource.OperationWait(w, activity, timeout, config.PollInterval)
}
