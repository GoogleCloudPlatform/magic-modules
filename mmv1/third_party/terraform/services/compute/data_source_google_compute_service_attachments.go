package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/registry"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeServiceAttachments() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleComputeServiceAttachmentsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: tpgresource.DatasourceSchemaFromResourceSchema(ResourceComputeServiceAttachment().Schema),
				},
			},
		},
	}
}

func datasourceGoogleComputeServiceAttachmentsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("projects/%s/regions/%s/serviceAttachments", project, region)
	d.SetId(id)

	call := NewClient(config, userAgent).ServiceAttachments.List(project, region)
	if filter, ok := d.GetOk("filter"); ok {
		call = call.Filter(filter.(string))
	}

	serviceAttachments := make([]map[string]interface{}, 0)

	for {
		list, err := call.Do()
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Service Attachments Not Found in project %s, region %s", project, region))
		}

		for _, sa := range list.Items {
			consumerAcceptLists := make([]map[string]interface{}, 0, len(sa.ConsumerAcceptLists))
			for _, cal := range sa.ConsumerAcceptLists {
				consumerAcceptLists = append(consumerAcceptLists, map[string]interface{}{
					"project_id_or_num": cal.ProjectIdOrNum,
					"connection_limit":  cal.ConnectionLimit,
					"network_url":       cal.NetworkUrl,
				})
			}

			connectedEndpoints := make([]map[string]interface{}, 0, len(sa.ConnectedEndpoints))
			for _, ce := range sa.ConnectedEndpoints {
				connectedEndpoints = append(connectedEndpoints, map[string]interface{}{
					"endpoint":          ce.Endpoint,
					"status":            ce.Status,
					"psc_connection_id": ce.PscConnectionId,
					"consumer_network":  ce.ConsumerNetwork,
				})
			}

			var pscId []map[string]interface{}
			if sa.PscServiceAttachmentId != nil {
				pscId = []map[string]interface{}{
					{
						"high": sa.PscServiceAttachmentId.High,
						"low":  sa.PscServiceAttachmentId.Low,
					},
				}
			}

			mapped := map[string]interface{}{
				"name":                        sa.Name,
				"description":                 sa.Description,
				"self_link":                   sa.SelfLink,
				"target_service":              sa.TargetService,
				"connection_preference":       sa.ConnectionPreference,
				"nat_subnets":                 sa.NatSubnets,
				"enable_proxy_protocol":       sa.EnableProxyProtocol,
				"domain_names":                sa.DomainNames,
				"fingerprint":                 sa.Fingerprint,
				"region":                      sa.Region,
				"reconcile_connections":       sa.ReconcileConnections,
				"propagated_connection_limit": sa.PropagatedConnectionLimit,
				"consumer_reject_lists":       sa.ConsumerRejectLists,
				"consumer_accept_lists":       consumerAcceptLists,
				"connected_endpoints":         connectedEndpoints,
				"psc_service_attachment_id":   pscId,
			}
			serviceAttachments = append(serviceAttachments, mapped)
		}

		if list.NextPageToken == "" {
			break
		}
		call = call.PageToken(list.NextPageToken)
	}

	if err := d.Set("service_attachments", serviceAttachments); err != nil {
		return fmt.Errorf("error setting service_attachments: %s", err)
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("error setting project: %s", err)
	}

	if err := d.Set("region", region); err != nil {
		return fmt.Errorf("error setting region: %s", err)
	}

	return nil
}

func init() {
	registry.Schema{
		Name:        "google_compute_service_attachments",
		ProductName: "compute",
		Type:        registry.SchemaTypeDataSource,
		Schema:      DataSourceGoogleComputeServiceAttachments(),
	}.Register()
}
