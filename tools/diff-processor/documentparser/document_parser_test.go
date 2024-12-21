package documentparser

import (
	"os"
	"path/filepath"
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
	wantArguments := []string{
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
		"params",
		// "params.resource_manager_tags", // params text does not include a nested tag
		"zone",
	}
	wantAttributes := []string{
		"id",
		"network_interface.access_config.nat_ip",
	}
	gotArguments := parser.Arguments()
	gotAttributes := parser.Attributes()
	if diff := cmp.Diff(wantArguments, gotArguments); diff != "" {
		t.Errorf("Parse returned diff in arguments(-want, +got): %s", diff)
	}
	if diff := cmp.Diff(wantAttributes, gotAttributes); diff != "" {
		t.Errorf("Parse returned diff in attributes(-want, +got): %s", diff)
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
	traverse(&paths, "", func(*node) bool { return true }, root)

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

func resourceToDocFile(resource string, repoPath string) string {
	fileBaseName := strings.TrimPrefix(resource, "google_") + ".html.markdown"
	return filepath.Join(repoPath, "website", "docs", "r", fileBaseName)
}
