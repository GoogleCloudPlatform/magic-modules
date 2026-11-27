package compute

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// DataSourceGoogleComputeInterconnectAttachments returns the schema of the google_compute_interconnect_attachments data source.
func DataSourceGoogleComputeInterconnectAttachments() *schema.Resource {
	attachmentSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeInterconnectAttachment().Schema)
	// These are arguments of the resource that the API does not return, or that only describe its own configuration.
	tpgresource.DeleteFieldsFromSchema(attachmentSchema, "project", "deletion_policy", "params", "candidate_subnets", "subnet_length", "terraform_labels", "effective_labels")

	return &schema.Resource{
		Read: dataSourceGoogleComputeInterconnectAttachmentsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the project whose interconnect attachments are listed. If it is not provided, the provider project is used.`,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The region whose interconnect attachments are listed, as a name or a self link. If it is not provided, the attachments of every region are listed, and the provider region is not used.`,
			},
			"interconnect_attachments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The interconnect attachments (VLAN attachments) of the project, sorted by region.`,
				Elem:        &schema.Resource{Schema: attachmentSchema},
			},
		},
	}
}

// dataSourceGoogleComputeInterconnectAttachmentsRead lists the attachments of one region, or of every region when none is given, following pages.
func dataSourceGoogleComputeInterconnectAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	billingProject := project
	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	// Only the region argument selects a region: the provider region is not used, so leaving it out lists every region.
	id := fmt.Sprintf("projects/%s/aggregated/interconnectAttachments", project)
	if region, ok := d.GetOk("region"); ok {
		id = fmt.Sprintf("projects/%s/regions/%s/interconnectAttachments", project, tpgresource.GetResourceNameFromSelfLink(region.(string)))
	}

	params := map[string]string{}
	attachments := make([]map[string]interface{}, 0)
	for {
		url, err := transport_tpg.AddQueryParams(transport_tpg.BaseUrl(Product, config)+id, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error retrieving list of interconnect attachments: %w", err)
		}

		items, err := computeInterconnectAttachmentsPageItems(res["items"])
		if err != nil {
			return err
		}
		for _, item := range items {
			attachments = append(attachments, flattenComputeInterconnectAttachmentsItem(item, d, config))
		}

		pageToken, _ := res["nextPageToken"].(string)
		if pageToken == "" {
			break
		}
		params["pageToken"] = pageToken
	}

	// The list of every region comes as an object keyed by region, which has no order, so the result is sorted by region.
	sort.SliceStable(attachments, func(i, j int) bool {
		a, _ := attachments[i]["region"].(string)
		b, _ := attachments[j]["region"].(string)
		return a < b
	})

	if err := d.Set("interconnect_attachments", attachments); err != nil {
		return fmt.Errorf("Error setting interconnect_attachments: %w", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %w", err)
	}

	d.SetId(id)

	return nil
}

// computeInterconnectAttachmentsPageItems returns the attachments of one page: a list for one region, or an object keyed by region for every region.
func computeInterconnectAttachmentsPageItems(items interface{}) ([]map[string]interface{}, error) {
	var raw []interface{}
	switch v := items.(type) {
	case nil:
	case []interface{}:
		raw = v
	case map[string]interface{}:
		for scope, entry := range v {
			scoped, ok := entry.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Error reading list of interconnect attachments: expected an object for %s, got %T", scope, entry)
			}
			// A region with no attachments carries a warning and no list.
			switch list := scoped["interconnectAttachments"].(type) {
			case nil:
			case []interface{}:
				raw = append(raw, list...)
			default:
				return nil, fmt.Errorf("Error reading list of interconnect attachments: expected a list for %s, got %T", scope, list)
			}
		}
	default:
		return nil, fmt.Errorf("Error reading list of interconnect attachments: expected a list or an object, got %T", v)
	}

	attachments := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		attachment, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Error reading list of interconnect attachments: expected an object, got %T", item)
		}
		attachments = append(attachments, attachment)
	}
	return attachments, nil
}

// flattenComputeInterconnectAttachmentsItem maps one attachment from the API response with the flatteners of google_compute_interconnect_attachment.
func flattenComputeInterconnectAttachmentsItem(attachment map[string]interface{}, d *schema.ResourceData, config *transport_tpg.Config) map[string]interface{} {
	region, _ := attachment["region"].(string)
	l2Forwarding := flattenComputeInterconnectAttachmentL2Forwarding(attachment["l2Forwarding"], d, config)
	sortComputeInterconnectAttachmentsApplianceMappings(l2Forwarding)

	return map[string]interface{}{
		"name":        flattenComputeInterconnectAttachmentName(attachment["name"], d, config),
		"description": flattenComputeInterconnectAttachmentDescription(attachment["description"], d, config),
		"self_link":   flattenComputeInterconnectAttachmentsLink(attachment["selfLink"]),
		// The API returns the region as a link, and its name is what a list over every region is grouped by.
		"region":                                 tpgresource.GetResourceNameFromSelfLink(region),
		"router":                                 flattenComputeInterconnectAttachmentRouter(attachment["router"], d, config),
		"interconnect":                           flattenComputeInterconnectAttachmentsLink(attachment["interconnect"]),
		"attachment_group":                       flattenComputeInterconnectAttachmentsLink(attachment["attachmentGroup"]),
		"type":                                   flattenComputeInterconnectAttachmentType(attachment["type"], d, config),
		"state":                                  flattenComputeInterconnectAttachmentState(attachment["state"], d, config),
		"admin_enabled":                          flattenComputeInterconnectAttachmentAdminEnabled(attachment["adminEnabled"], d, config),
		"bandwidth":                              flattenComputeInterconnectAttachmentBandwidth(attachment["bandwidth"], d, config),
		"edge_availability_domain":               flattenComputeInterconnectAttachmentEdgeAvailabilityDomain(attachment["edgeAvailabilityDomain"], d, config),
		"mtu":                                    flattenComputeInterconnectAttachmentMtu(attachment["mtu"], d, config),
		"vlan_tag8021q":                          flattenComputeInterconnectAttachmentVlanTag8021q(attachment["vlanTag8021q"], d, config),
		"stack_type":                             flattenComputeInterconnectAttachmentStackType(attachment["stackType"], d, config),
		"encryption":                             flattenComputeInterconnectAttachmentEncryption(attachment["encryption"], d, config),
		"ipsec_internal_addresses":               flattenComputeInterconnectAttachmentIpsecInternalAddresses(attachment["ipsecInternalAddresses"], d, config),
		"cloud_router_ip_address":                flattenComputeInterconnectAttachmentCloudRouterIpAddress(attachment["cloudRouterIpAddress"], d, config),
		"customer_router_ip_address":             flattenComputeInterconnectAttachmentCustomerRouterIpAddress(attachment["customerRouterIpAddress"], d, config),
		"cloud_router_ipv6_address":              flattenComputeInterconnectAttachmentCloudRouterIpv6Address(attachment["cloudRouterIpv6Address"], d, config),
		"customer_router_ipv6_address":           flattenComputeInterconnectAttachmentCustomerRouterIpv6Address(attachment["customerRouterIpv6Address"], d, config),
		"candidate_cloud_router_ip_address":      flattenComputeInterconnectAttachmentCandidateCloudRouterIpAddress(attachment["candidateCloudRouterIpAddress"], d, config),
		"candidate_customer_router_ip_address":   flattenComputeInterconnectAttachmentCandidateCustomerRouterIpAddress(attachment["candidateCustomerRouterIpAddress"], d, config),
		"candidate_cloud_router_ipv6_address":    flattenComputeInterconnectAttachmentCandidateCloudRouterIpv6Address(attachment["candidateCloudRouterIpv6Address"], d, config),
		"candidate_customer_router_ipv6_address": flattenComputeInterconnectAttachmentCandidateCustomerRouterIpv6Address(attachment["candidateCustomerRouterIpv6Address"], d, config),
		"pairing_key":                            flattenComputeInterconnectAttachmentPairingKey(attachment["pairingKey"], d, config),
		"partner_asn":                            flattenComputeInterconnectAttachmentPartnerAsn(attachment["partnerAsn"], d, config),
		"private_interconnect_info":              flattenComputeInterconnectAttachmentPrivateInterconnectInfo(attachment["privateInterconnectInfo"], d, config),
		"google_reference_id":                    flattenComputeInterconnectAttachmentGoogleReferenceId(attachment["googleReferenceId"], d, config),
		"creation_timestamp":                     flattenComputeInterconnectAttachmentCreationTimestamp(attachment["creationTimestamp"], d, config),
		// Every label is returned: the labels flattener of the resource keeps only the labels found in its configuration.
		"labels":            attachment["labels"],
		"label_fingerprint": flattenComputeInterconnectAttachmentLabelFingerprint(attachment["labelFingerprint"], d, config),
		"l2_forwarding":     l2Forwarding,
	}
}

// flattenComputeInterconnectAttachmentsLink returns a link to a Compute resource in its v1 form, whichever API version sent it.
func flattenComputeInterconnectAttachmentsLink(v interface{}) interface{} {
	if s, ok := v.(string); ok {
		return tpgresource.ConvertSelfLinkToV1(s)
	}
	return v
}

// sortComputeInterconnectAttachmentsApplianceMappings orders the appliance mappings by VLAN id: the API returns them as an object, which has no order.
func sortComputeInterconnectAttachmentsApplianceMappings(l2Forwarding interface{}) {
	blocks, _ := l2Forwarding.([]interface{})
	for _, block := range blocks {
		fields, _ := block.(map[string]interface{})
		mappings, _ := fields["appliance_mappings"].([]interface{})
		sort.Slice(mappings, func(i, j int) bool {
			a, _ := mappings[i].(map[string]interface{})["vlan_id"].(string)
			b, _ := mappings[j].(map[string]interface{})["vlan_id"].(string)
			aTag, aErr := strconv.Atoi(a)
			bTag, bErr := strconv.Atoi(b)
			if aErr == nil && bErr == nil {
				return aTag < bTag
			}
			return a < b
		})
	}
}

func init() {
	registry.Schema{
		Name:        "google_compute_interconnect_attachments",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeInterconnectAttachments(),
	}.Register()
}
