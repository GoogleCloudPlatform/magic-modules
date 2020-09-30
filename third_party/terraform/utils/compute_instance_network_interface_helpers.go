package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

type networkInterfaceHelper struct {
	InferNetworkFromSubnetwork   func() error
	RefreshInstance              func() error
	DeleteAccessConfigs          func() error
	CreateAccessConfigs          func() error
	DeleteAliasIPRanges          func() error
	CreateAliasIPRanges          func() error
	CreateUpdateWhileStoppedCall func() (func(inst *computeBeta.Instance) error, error)
}

func networkInterfaceHelperFactory(d *schema.ResourceData, config *Config, instance *computeBeta.Instance, networkInterfaces []*computeBeta.NetworkInterface, index int, project, zone, userAgent string) (networkInterfaceHelper, error) {
	prefix := fmt.Sprintf("network_interface.%d", index)
	networkInterface := networkInterfaces[index]
	instNetworkInterface := instance.NetworkInterfaces[index]
	networkName := d.Get(prefix + ".name").(string)

	if networkName != instNetworkInterface.Name {
		return networkInterfaceHelper{}, fmt.Errorf("Instance networkInterface had unexpected name: %s", instNetworkInterface.Name)
	}

	inferNetworkFromSubnetwork := func() error {
		subnetwork := networkInterface.Subnetwork
		subnetProjectField := prefix + ".subnetwork_project"
		sf, err := ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
		if err != nil {
			return fmt.Errorf("Cannot determine self_link for subnetwork %q: %s", subnetwork, err)
		}
		resp, err := config.NewComputeClient(userAgent).Subnetworks.Get(sf.Project, sf.Region, sf.Name).Do()
		if err != nil {
			return errwrap.Wrapf("Error getting subnetwork value: {{err}}", err)
		}
		nf, err := ParseNetworkFieldValue(resp.Network, d, config)
		if err != nil {
			return fmt.Errorf("Cannot determine self_link for network %q: %s", resp.Network, err)
		}
		networkInterface.Network = nf.RelativeLink()
		return nil
	}

	refreshInstance := func() error {
		// re-read fingerprint
		inst, err := config.clientComputeBeta.Instances.Get(project, zone, instance.Name).Do()
		if err != nil {
			return nil
		}

		instance = inst
		instNetworkInterface = instance.NetworkInterfaces[index]
		return nil
	}

	deleteAccessConfigs := func() error {
		// Delete any accessConfig that currently exists in instNetworkInterface
		for _, ac := range instNetworkInterface.AccessConfigs {
			op, err := config.NewComputeClient(userAgent).Instances.DeleteAccessConfig(
				project, zone, instance.Name, ac.Name, networkName).Do()
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

	createAccessConfigs := func() error {
		// Create new ones
		for _, ac := range networkInterface.AccessConfigs {
			op, err := config.clientComputeBeta.Instances.AddAccessConfig(
				project, zone, instance.Name, networkName, ac).Do()
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

	deleteAliasIPRanges := func() error {
		if len(instNetworkInterface.AliasIpRanges) > 0 {
			ni := &computeBeta.NetworkInterface{
				Fingerprint:     instNetworkInterface.Fingerprint,
				ForceSendFields: []string{"AliasIpRanges"},
			}
			op, err := config.clientComputeBeta.Instances.UpdateNetworkInterface(project, zone, instance.Name, networkName, ni).Do()
			if err != nil {
				return errwrap.Wrapf("Error removing alias_ip_range: {{err}}", err)
			}
			opErr := computeOperationWaitTime(config, op, project, "updating alias ip ranges", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
		}
		return nil
	}

	createAliasIPRanges := func() error {
		// Lets be explicit about what we are changing in the patch call
		networkInterfacePatchObj := &computeBeta.NetworkInterface{}
		networkInterfacePatchObj.AliasIpRanges = networkInterface.AliasIpRanges
		networkInterfacePatchObj.Fingerprint = instNetworkInterface.Fingerprint
		op, err := config.clientComputeBeta.Instances.UpdateNetworkInterface(project, zone, instance.Name, networkName, networkInterfacePatchObj).Do()
		if err != nil {
			return errwrap.Wrapf("Error updating network interface: {{err}}", err)
		}
		opErr := computeOperationWaitTime(config, op, project, "network interface to update", userAgent, d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
		return nil
	}

	createUpdateWhileStoppedCall := func() (func(inst *computeBeta.Instance) error, error) {
		// Lets be explicit about what we are changing in the patch call
		networkInterfacePatchObj := &computeBeta.NetworkInterface{
			Network:       networkInterface.Network,
			Subnetwork:    networkInterface.Subnetwork,
			AliasIpRanges: networkInterface.AliasIpRanges,
		}

		// network_ip can be inferred if not declared. Let's only patch if it's being changed by user
		// otherwise this could fail if the network ip is not compatible with the new Subnetwork/Network.
		if d.HasChange(prefix + ".network_ip") {
			networkInterfacePatchObj.NetworkIP = networkInterface.NetworkIP
		}

		// Access config can run into some issues since we can't derive users original intent due to
		// terraform limitation. Lets only update it if we need to.
		if d.HasChange(prefix + ".access_config") {
			err := deleteAccessConfigs()
			if err != nil {
				return nil, err
			}
		}

		// Access configs' ip changes when the instance stops invalidating our fingerprint
		// expect caller to re-validate instance before calling patch
		updateCall := func(instance *computeBeta.Instance) error {
			networkInterfacePatchObj.Fingerprint = instance.NetworkInterfaces[index].Fingerprint
			op, err := config.clientComputeBeta.Instances.UpdateNetworkInterface(project, zone, instance.Name, networkName, networkInterfacePatchObj).Do()
			if err != nil {
				return errwrap.Wrapf("Error updating network interface: {{err}}", err)
			}
			opErr := computeOperationWaitTime(config, op, project, "network interface to update", userAgent, d.Timeout(schema.TimeoutUpdate))
			if opErr != nil {
				return opErr
			}
			if d.HasChange(prefix + ".access_config") {
				err := createAccessConfigs()
				if err != nil {
					return err
				}
			}
			return nil
		}
		return updateCall, nil
	}

	return networkInterfaceHelper{
		InferNetworkFromSubnetwork:   inferNetworkFromSubnetwork,
		RefreshInstance:              refreshInstance,
		DeleteAccessConfigs:          deleteAccessConfigs,
		CreateAccessConfigs:          createAccessConfigs,
		DeleteAliasIPRanges:          deleteAliasIPRanges,
		CreateAliasIPRanges:          createAliasIPRanges,
		CreateUpdateWhileStoppedCall: createUpdateWhileStoppedCall,
	}, nil
}
