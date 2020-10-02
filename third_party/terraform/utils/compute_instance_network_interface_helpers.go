package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
)

type networkInterfaceHelper struct {
	d                    *schema.ResourceData
	config               *Config
	instNetworkInterface *computeBeta.NetworkInterface
	networkInterface     *computeBeta.NetworkInterface
	index                int
	instanceName         string
	prefix               string
	project              string
	zone                 string
	userAgent            string
}

func createNetworkInterfaceHelper(d *schema.ResourceData, config *Config, instance *computeBeta.Instance, networkInterfaces []*computeBeta.NetworkInterface, index int, project, zone, userAgent string) (networkInterfaceHelper, error) {
	prefix := fmt.Sprintf("network_interface.%d", index)
	networkInterface := networkInterfaces[index]
	instNetworkInterface := instance.NetworkInterfaces[index]
	networkName := d.Get(prefix + ".name").(string)

	if networkName != instNetworkInterface.Name {
		return networkInterfaceHelper{}, fmt.Errorf("Instance networkInterface had unexpected name: %s", instNetworkInterface.Name)
	}
	return networkInterfaceHelper{
		d:                    d,
		config:               config,
		instNetworkInterface: instNetworkInterface,
		networkInterface:     networkInterface,
		index:                index,
		instanceName:         instance.Name,
		prefix:               prefix,
		project:              project,
		zone:                 zone,
		userAgent:            userAgent,
	}, nil
}

func (niH *networkInterfaceHelper) InferNetworkFromSubnetwork() error {
	subnetwork := niH.networkInterface.Subnetwork
	subnetProjectField := niH.prefix + ".subnetwork_project"
	sf, err := ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, niH.d, niH.config)
	if err != nil {
		return fmt.Errorf("Cannot determine self_link for subnetwork %q: %s", subnetwork, err)
	}
	resp, err := niH.config.NewComputeClient(niH.userAgent).Subnetworks.Get(sf.Project, sf.Region, sf.Name).Do()
	if err != nil {
		return errwrap.Wrapf("Error getting subnetwork value: {{err}}", err)
	}
	nf, err := ParseNetworkFieldValue(resp.Network, niH.d, niH.config)
	if err != nil {
		return fmt.Errorf("Cannot determine self_link for network %q: %s", resp.Network, err)
	}
	niH.networkInterface.Network = nf.RelativeLink()
	return nil
}

func (niH *networkInterfaceHelper) RefreshInstance() error {
	// re-read fingerprint
	inst, err := niH.config.NewComputeBetaClient(niH.userAgent).Instances.Get(niH.project, niH.zone, niH.networkInterface.Name).Do()
	if err != nil {
		return nil
	}

	instance := inst
	niH.instNetworkInterface = instance.NetworkInterfaces[niH.index]
	return nil
}

func (niH *networkInterfaceHelper) DeleteAccessConfigs() error {
	// Delete any accessConfig that currently exists in instNetworkInterface
	for _, ac := range niH.instNetworkInterface.AccessConfigs {
		op, err := niH.config.NewComputeClient(niH.userAgent).Instances.DeleteAccessConfig(
			niH.project, niH.zone, niH.instanceName, ac.Name, niH.networkInterface.Name).Do()
		if err != nil {
			return fmt.Errorf("Error deleting old access_config: %s", err)
		}
		opErr := computeOperationWaitTime(niH.config, op, niH.project, "old access_config to delete", niH.userAgent, niH.d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func (niH *networkInterfaceHelper) CreateAccessConfigs() error {
	// Create new ones
	for _, ac := range niH.networkInterface.AccessConfigs {
		op, err := niH.config.NewComputeBetaClient(niH.userAgent).Instances.AddAccessConfig(
			niH.project, niH.zone, niH.instanceName, niH.networkInterface.Name, ac).Do()
		if err != nil {
			return fmt.Errorf("Error adding new access_config: %s", err)
		}
		opErr := computeOperationWaitTime(niH.config, op, niH.project, "new access_config to add", niH.userAgent, niH.d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func (niH *networkInterfaceHelper) DeleteAliasIPRanges() error {
	if len(niH.instNetworkInterface.AliasIpRanges) > 0 {
		ni := &computeBeta.NetworkInterface{
			Fingerprint:     niH.instNetworkInterface.Fingerprint,
			ForceSendFields: []string{"AliasIpRanges"},
		}
		op, err := niH.config.NewComputeBetaClient(niH.userAgent).Instances.UpdateNetworkInterface(niH.project, niH.zone, niH.instanceName, niH.networkInterface.Name, ni).Do()
		if err != nil {
			return errwrap.Wrapf("Error removing alias_ip_range: {{err}}", err)
		}
		opErr := computeOperationWaitTime(niH.config, op, niH.project, "updating alias ip ranges", niH.userAgent, niH.d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func (niH *networkInterfaceHelper) CreateAliasIPRanges() error {
	// Lets be explicit about what we are changing in the patch call
	networkInterfacePatchObj := &computeBeta.NetworkInterface{}
	networkInterfacePatchObj.AliasIpRanges = niH.networkInterface.AliasIpRanges
	networkInterfacePatchObj.Fingerprint = niH.instNetworkInterface.Fingerprint
	op, err := niH.config.NewComputeBetaClient(niH.userAgent).Instances.UpdateNetworkInterface(niH.project, niH.zone, niH.instanceName, niH.networkInterface.Name, networkInterfacePatchObj).Do()
	if err != nil {
		return errwrap.Wrapf("Error updating network interface: {{err}}", err)
	}
	opErr := computeOperationWaitTime(niH.config, op, niH.project, "network interface to update", niH.userAgent, niH.d.Timeout(schema.TimeoutUpdate))
	if opErr != nil {
		return opErr
	}
	return nil
}

func (niH *networkInterfaceHelper) CreateUpdateWhileStoppedCall() (func(inst *computeBeta.Instance) error, error) {
	// Lets be explicit about what we are changing in the patch call
	networkInterfacePatchObj := &computeBeta.NetworkInterface{
		Network:       niH.networkInterface.Network,
		Subnetwork:    niH.networkInterface.Subnetwork,
		AliasIpRanges: niH.networkInterface.AliasIpRanges,
	}

	// network_ip can be inferred if not declared. Let's only patch if it's being changed by user
	// otherwise this could fail if the network ip is not compatible with the new Subnetwork/Network.
	if niH.d.HasChange(niH.prefix + ".network_ip") {
		networkInterfacePatchObj.NetworkIP = niH.networkInterface.NetworkIP
	}

	// Access config can run into some issues since we can't derive users original intent due to
	// terraform limitation. Lets only update it if we need to.
	if niH.d.HasChange(niH.prefix + ".access_config") {
		err := niH.DeleteAccessConfigs()
		if err != nil {
			return nil, err
		}
	}

	// Access configs' ip changes when the instance stops invalidating our fingerprint
	// expect caller to re-validate instance before calling patch
	updateCall := func(instance *computeBeta.Instance) error {
		networkInterfacePatchObj.Fingerprint = instance.NetworkInterfaces[niH.index].Fingerprint
		op, err := niH.config.NewComputeBetaClient(niH.userAgent).Instances.UpdateNetworkInterface(niH.project, niH.zone, niH.instanceName, niH.networkInterface.Name, networkInterfacePatchObj).Do()
		if err != nil {
			return errwrap.Wrapf("Error updating network interface: {{err}}", err)
		}
		opErr := computeOperationWaitTime(niH.config, op, niH.project, "network interface to update", niH.userAgent, niH.d.Timeout(schema.TimeoutUpdate))
		if opErr != nil {
			return opErr
		}
		if niH.d.HasChange(niH.prefix + ".access_config") {
			err := niH.CreateAccessConfigs()
			if err != nil {
				return err
			}
		}
		return nil
	}
	return updateCall, nil
}
