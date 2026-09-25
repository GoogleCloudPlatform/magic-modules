package tpgiamresource

import (
	"fmt"
	"time"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"google.golang.org/api/cloudresourcemanager/v1"
)

const (
	batchKeyTmplModifyIamPolicy = "%s modifyIamPolicy"
)

// A single batched change: how to modify the policy, and an optional check that the write is allowed.
type iamPolicyChange struct {
	modify     iamPolicyModifyFunc
	allowWrite iamPolicyWriteAllowedFunc
}

func BatchRequestModifyIamPolicy(updater ResourceIamUpdater, modify iamPolicyModifyFunc, allowWrite iamPolicyWriteAllowedFunc, config *transport_tpg.Config, reqDesc string) error {
	batchKey := fmt.Sprintf(batchKeyTmplModifyIamPolicy, updater.GetMutexKey())

	request := &transport_tpg.BatchRequest{
		ResourceName: updater.GetResourceId(),
		Body:         []iamPolicyChange{{modify: modify, allowWrite: allowWrite}},
		CombineF:     combineBatchIamPolicyModifiers,
		SendF:        sendBatchModifyIamPolicy(updater),
		DebugId:      reqDesc,
	}

	_, err := config.RequestBatcherIam.SendRequestWithTimeout(batchKey, request, time.Minute*30)
	return err
}

func combineBatchIamPolicyModifiers(currV interface{}, toAddV interface{}) (interface{}, error) {
	currChanges, ok := currV.([]iamPolicyChange)
	if !ok {
		return nil, fmt.Errorf("provider error in batch combiner: expected data to be type []iamPolicyChange, got %v with type %T", currV, currV)
	}

	newChanges, ok := toAddV.([]iamPolicyChange)
	if !ok {
		return nil, fmt.Errorf("provider error in batch combiner: expected data to be type []iamPolicyChange, got %v with type %T", currV, currV)
	}

	return append(currChanges, newChanges...), nil
}

func sendBatchModifyIamPolicy(updater ResourceIamUpdater) transport_tpg.BatcherSendFunc {
	return func(resourceName string, body interface{}) (interface{}, error) {
		changes, ok := body.([]iamPolicyChange)
		if !ok {
			return nil, fmt.Errorf("provider error: expected data to be type []iamPolicyChange, got %v with type %T", body, body)
		}
		modify := func(policy *cloudresourcemanager.Policy) error {
			for _, c := range changes {
				if err := c.modify(policy); err != nil {
					return err
				}
			}
			return nil
		}

		var allowWrite iamPolicyWriteAllowedFunc
		for _, c := range changes {
			if c.allowWrite != nil {
				allowWrite = allowBatchedWrite(changes)
				break
			}
		}
		return nil, iamPolicyReadModifyWrite(updater, modify, allowWrite)
	}
}

// Each change's allowWrite must only see its own effect, not the rest of the batch's, so replay the
// batch on a copy of before and check each change against the policy just before it was applied.
func allowBatchedWrite(changes []iamPolicyChange) iamPolicyWriteAllowedFunc {
	return func(before, _ *cloudresourcemanager.Policy) error {
		current, err := copyIamPolicy(before)
		if err != nil {
			return err
		}
		for _, c := range changes {
			if c.allowWrite == nil {
				if err := c.modify(current); err != nil {
					return err
				}
				continue
			}
			prev, err := copyIamPolicy(current)
			if err != nil {
				return err
			}
			if err := c.modify(current); err != nil {
				return err
			}
			if err := c.allowWrite(prev, current); err != nil {
				return err
			}
		}
		return nil
	}
}
