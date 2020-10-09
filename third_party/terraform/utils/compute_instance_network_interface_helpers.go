package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

func computeInstanceDeleteAccessConfigs(d *schema.ResourceData, config *Config, instNetworkInterface *computeBeta.NetworkInterface, project, zone, userAgent, instanceName string) error {
	// TODO: This code deletes then recreates accessConfigs.  This is bad because it may
	// leave the machine inaccessible from either ip if the creation part fails (network
	// timeout etc).  However right now there is a GCE limit of 1 accessConfig so it is
	// the only way to do it.  In future this should be revised to only change what is
	// necessary, and also add before removing.

	// Delete any accessConfig that currently exists in instNetworkInterface
	for _, ac := range instNetworkInterface.AccessConfigs {
		op, err := config.NewComputeClient(userAgent).Instances.DeleteAccessConfig(
			project, zone, instanceName, ac.Name, instNetworkInterface.Name).Do()
		if err != nil {
			return fmt.Errorf("Error deleting old access_config: %s", err)
		}
		opErr := computeOperationWaitTime(config, op, project, "old access_config to delete", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func computeInstanceAddAccessConfigs(d *schema.ResourceData, config *Config, instNetworkInterface *computeBeta.NetworkInterface, accessConfigs []*computeBeta.AccessConfig, project, zone, userAgent, instanceName string) error {
	// Create new ones
	for _, ac := range accessConfigs {
		op, err := config.NewComputeBetaClient(userAgent).Instances.AddAccessConfig(project, zone, instanceName, instNetworkInterface.Name, ac).Do()
		if err != nil {
			return fmt.Errorf("Error adding new access_config: %s", err)
		}
		opErr := computeOperationWaitTime(config, op, project, "new access_config to add", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func computeInstanceCreateUpdateWhileStoppedCall(d *schema.ResourceData, config *Config, networkInterfacePatchObj *computeBeta.NetworkInterface, accessConfigs []*computeBeta.AccessConfig, accessConfigsHaveChanged bool, index int, project, zone, userAgent, instanceName string) func(inst *computeBeta.Instance) error {

	// Access configs' ip changes when the instance stops invalidating our fingerprint
	// expect caller to re-validate instance before calling patch this is why we expect
	// instance to be passed in
	return func(instance *computeBeta.Instance) error {

		instNetworkInterface := instance.NetworkInterfaces[index]
		networkInterfacePatchObj.Fingerprint = instNetworkInterface.Fingerprint

		// Access config can run into some issues since we can't derive users original intent due to
		// terraform limitation. Lets only update it if we need to.
		if accessConfigsHaveChanged {
			err := computeInstanceDeleteAccessConfigs(d, config, instNetworkInterface, project, zone, userAgent, instanceName)
			if err != nil {
				return err
			}
		}

		op, err := config.NewComputeBetaClient(userAgent).Instances.UpdateNetworkInterface(project, zone, instanceName, instNetworkInterface.Name, networkInterfacePatchObj).Do()
		if err != nil {
			return errwrap.Wrapf("Error updating network interface: {{err}}", err)
		}
		opErr := computeOperationWaitTime(config, op, project, "network interface to update", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}

		if accessConfigsHaveChanged {
			err := computeInstanceAddAccessConfigs(d, config, instNetworkInterface, accessConfigs, project, zone, userAgent, instanceName)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
