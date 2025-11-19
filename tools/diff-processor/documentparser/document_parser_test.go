package documentparser

import (
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	b, err := os.ReadFile("../testdata/resource.html.markdown")
	if err != nil {
		t.Fatal(err)
	}
	parser := NewParser()
	if err := parser.Parse(b); err != nil {
		t.Fatal(err)
	}
	want := []string{
		// The below are from arguments section.
		"boot_disk",
		"boot_disk.auto_delete",
		"boot_disk.device_name",
		"boot_disk.disk_encryption_key_raw",
		"boot_disk.initialize_params",
		"boot_disk.initialize_params.enable_confidential_compute",
		"boot_disk.initialize_params.image",
		"boot_disk.initialize_params.labels",
		"boot_disk.initialize_params.provisioned_iops",
		"boot_disk.initialize_params.provisioned_throughput",
		"boot_disk.initialize_params.resource_manager_tags",
		"boot_disk.initialize_params.size",
		"boot_disk.initialize_params.storage_pool",
		"boot_disk.initialize_params.type",
		"boot_disk.kms_key_self_link",
		"boot_disk.mode",
		"boot_disk.source",
		"name",
		"network_interface",
		"network_interface.access_config",
		"network_interface.access_config.nat_ip",
		"network_interface.access_config.network_tier",
		"network_interface.access_config.public_ptr_domain_name",
		"network_interface.alias_ip_range",
		"network_interface.alias_ip_range.ip_cidr_range",
		"network_interface.alias_ip_range.subnetwork_range_name",
		"network_interface.ipv6_access_config",
		"network_interface.ipv6_access_config.external_ipv6",
		"network_interface.ipv6_access_config.external_ipv6_prefix_length",
		"network_interface.ipv6_access_config.name",
		"network_interface.ipv6_access_config.network_tier",
		"network_interface.ipv6_access_config.public_ptr_domain_name",
		"network_interface.network",
		"network_interface.network_attachment",
		"network_interface.network_ip",
		"network_interface.nic_type",
		"network_interface.queue_count",
		"network_interface.security_policy",
		"network_interface.stack_type",
		"network_interface.subnetwork",
		"network_interface.subnetwork_project",
		"params",
		// "params.resource_manager_tags", // params text does not include a nested tag
		"zone",
		"labels",
		"description",
		"traffic_port_selector",
		"traffic_port_selector.ports",
		"project",
		// The below are from attributes section.
		"id",
		"network_interface.access_config.nat_ip",
		"workload_identity_config",
		"errors",
		"workload_identity_config.identity_provider",
		"workload_identity_config.issuer_uri",
		"workload_identity_config.workload_pool",
		"errors.message",
		// The below are from the ephemeral attributes section.
		"shared_secret_wo",
		"sensitive_params",
		"sensitive_params.secret_access_key_wo",
	}
	got := parser.FlattenFields()
	// gotAttributes := parser.Attributes()
	for _, arr := range [][]string{got, want} {
		sort.Strings(arr)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Parse returned diff in arguments(-want, +got): %s", diff)
	}
}

func TestTraverse(t *testing.T) {
	n1 := &node{name: "n1"}
	n2 := &node{name: "n2"}
	n3 := &node{name: "n3"}
	n4 := &node{name: "n4"}
	root := &node{
		children: []*node{n1, n2, n3},
	}
	n1.children = []*node{n4}
	n2.children = []*node{n4}

	var paths []string
	traverse(&paths, "", root)

	wantPaths := []string{
		"n1",
		"n1.n4",
		"n2",
		"n2.n4",
		"n3",
	}
	if diff := cmp.Diff(wantPaths, paths); diff != "" {
		t.Errorf("traverse returned diff(-want, +got): %s", diff)
	}
}

func TestSplitWithRegexp(t *testing.T) {
	paragraph := []string{
		"Lorem ipsum",
		"*   `name` - (Required) Resource name.",
		"",
		"* `os_policies` - (Required) List of OS policies to be applied to the VMs. Structure is [documented below](#nested_os_policies).	",
		"-   `some_field` - (Required) Lorem ipsum.	",
	}

	got := splitWithRegexp(strings.Join(paragraph, "\n"), fieldNameRegex)
	want := []string{
		"Lorem ipsum\n",
		"*   `name` - (Required) Resource name.\n\n",
		"* `os_policies` - (Required) List of OS policies to be applied to the VMs. Structure is [documented below](#nested_os_policies).	\n",
		"-   `some_field` - (Required) Lorem ipsum.	",
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("splitWithRegexp returned diff(-want, +got): %s", diff)
	}
}
