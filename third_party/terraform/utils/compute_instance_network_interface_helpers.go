package google

import (
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	computeBeta "google.golang.org/api/compute/v0.beta"
	"google.golang.org/api/compute/v1"
)

type networkInterfaceHelper struct {
	WaitForOperation       func(interface{}, string) error
	GetSubnetworks         func(string, string, string) *compute.SubnetworksGetCall
	DeleteAccessConfig     func(string, string, string, string, string) *compute.InstancesDeleteAccessConfigCall
	AddAccessConfig        func(string, string, string, string, *computeBeta.AccessConfig) *computeBeta.InstancesAddAccessConfigCall
	UpdateNetworkInterface func(string, string, string, string, *computeBeta.NetworkInterface) *computeBeta.InstancesUpdateNetworkInterfaceCall
	HasChange              func(string) bool
	instanceName           string
	prefix                 string
	project                string
	zone                   string
}

type networkInterfaceHelperInterface interface {
	InferNetworkFromSubnetwork(d *schema.ResourceData, config *Config, instNetworkInterface, networkInterface *computeBeta.NetworkInterface) (string, error)
	DeleteAccessConfigs(instNetworkInterface *computeBeta.NetworkInterface) error
	CreateAccessConfigs(instNetworkInterface, networkInterface *computeBeta.NetworkInterface) error
	DeleteAliasIPRanges(instNetworkInterface *computeBeta.NetworkInterface) error
	CreateAliasIPRanges(instNetworkInterface, networkInterface *computeBeta.NetworkInterface) error
	CreateUpdateWhileStoppedCall(instNetworkInterface, networkInterface *computeBeta.NetworkInterface, index int) (func(inst *computeBeta.Instance) error, error)
}

func (niH *networkInterfaceHelper) InferNetworkFromSubnetwork(d *schema.ResourceData, config *Config, instNetworkInterface, networkInterface *computeBeta.NetworkInterface) (string, error) {
	subnetwork := networkInterface.Subnetwork
	subnetProjectField := niH.prefix + ".subnetwork_project"
	sf, err := ParseSubnetworkFieldValueWithProjectField(subnetwork, subnetProjectField, d, config)
	if err != nil {
		return "", fmt.Errorf("Cannot determine self_link for subnetwork %q: %s", subnetwork, err)
	}
	resp, err := niH.GetSubnetworks(sf.Project, sf.Region, sf.Name).Do()
	if err != nil {
		return "", errwrap.Wrapf("Error getting subnetwork value: {{err}}", err)
	}
	nf, err := ParseNetworkFieldValue(resp.Network, d, config)
	if err != nil {
		return "", fmt.Errorf("Cannot determine self_link for network %q: %s", resp.Network, err)
	}
	return nf.RelativeLink(), nil
}

func (niH *networkInterfaceHelper) DeleteAccessConfigs(instNetworkInterface *computeBeta.NetworkInterface) error {
	// Delete any accessConfig that currently exists in instNetworkInterface
	for _, ac := range instNetworkInterface.AccessConfigs {
		op, err := niH.DeleteAccessConfig(niH.project, niH.zone, niH.instanceName, ac.Name, instNetworkInterface.Name).Do()
		if err != nil {
			return fmt.Errorf("Error deleting old access_config: %s", err)
		}
		opErr := niH.WaitForOperation(op, "old access_config to delete")
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func (niH *networkInterfaceHelper) CreateAccessConfigs(instNetworkInterface, networkInterface *computeBeta.NetworkInterface) error {
	// Create new ones
	for _, ac := range networkInterface.AccessConfigs {
		op, err := niH.AddAccessConfig(niH.project, niH.zone, niH.instanceName, instNetworkInterface.Name, ac).Do()
		if err != nil {
			return fmt.Errorf("Error adding new access_config: %s", err)
		}
		opErr := niH.WaitForOperation(op, "new access_config to add")
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func (niH *networkInterfaceHelper) DeleteAliasIPRanges(instNetworkInterface *computeBeta.NetworkInterface) error {
	if len(instNetworkInterface.AliasIpRanges) > 0 {
		ni := &computeBeta.NetworkInterface{
			Fingerprint:     instNetworkInterface.Fingerprint,
			ForceSendFields: []string{"AliasIpRanges"},
		}
		op, err := niH.UpdateNetworkInterface(niH.project, niH.zone, niH.instanceName, instNetworkInterface.Name, ni).Do()
		if err != nil {
			return errwrap.Wrapf("Error removing alias_ip_range: {{err}}", err)
		}
		opErr := niH.WaitForOperation(op, "updating alias ip ranges")
		if opErr != nil {
			return opErr
		}
	}
	return nil
}

func (niH *networkInterfaceHelper) CreateAliasIPRanges(instNetworkInterface, networkInterface *computeBeta.NetworkInterface) error {
	// Lets be explicit about what we are changing in the patch call
	networkInterfacePatchObj := &computeBeta.NetworkInterface{}
	networkInterfacePatchObj.AliasIpRanges = networkInterface.AliasIpRanges
	networkInterfacePatchObj.Fingerprint = instNetworkInterface.Fingerprint
	op, err := niH.UpdateNetworkInterface(niH.project, niH.zone, niH.instanceName, instNetworkInterface.Name, networkInterfacePatchObj).Do()
	if err != nil {
		return errwrap.Wrapf("Error updating network interface: {{err}}", err)
	}
	opErr := niH.WaitForOperation(op, "network interface to update")
	if opErr != nil {
		return opErr
	}
	return nil
}

func (niH *networkInterfaceHelper) CreateUpdateWhileStoppedCall(instNetworkInterface, networkInterface *computeBeta.NetworkInterface, index int) (func(inst *computeBeta.Instance) error, error) {
	// Lets be explicit about what we are changing in the patch call
	networkInterfacePatchObj := &computeBeta.NetworkInterface{
		Network:       networkInterface.Network,
		Subnetwork:    networkInterface.Subnetwork,
		AliasIpRanges: networkInterface.AliasIpRanges,
	}

	// network_ip can be inferred if not declared. Let's only patch if it's being changed by user
	// otherwise this could fail if the network ip is not compatible with the new Subnetwork/Network.
	if niH.HasChange(niH.prefix + ".network_ip") {
		networkInterfacePatchObj.NetworkIP = networkInterface.NetworkIP
	}

	// Access config can run into some issues since we can't derive users original intent due to
	// terraform limitation. Lets only update it if we need to.
	if niH.HasChange(niH.prefix + ".access_config") {
		err := niH.DeleteAccessConfigs(instNetworkInterface)
		if err != nil {
			return nil, err
		}
	}

	// Access configs' ip changes when the instance stops invalidating our fingerprint
	// expect caller to re-validate instance before calling patch
	updateCall := func(instance *computeBeta.Instance) error {
		networkInterfacePatchObj.Fingerprint = instance.NetworkInterfaces[index].Fingerprint
		op, err := niH.UpdateNetworkInterface(niH.project, niH.zone, niH.instanceName, instNetworkInterface.Name, networkInterfacePatchObj).Do()
		if err != nil {
			return errwrap.Wrapf("Error updating network interface: {{err}}", err)
		}
		opErr := niH.WaitForOperation(op, "network interface to update")
		if opErr != nil {
			return opErr
		}
		if niH.HasChange(niH.prefix + ".access_config") {
			err := niH.CreateAccessConfigs(instance.NetworkInterfaces[index], networkInterface)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return updateCall, nil
}
