package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func dataSourceMonitoringServiceMeshIstio() *schema.Resource {
	istioSchema := map[string]*schema.Schema{
		"mesh_uid": {
			Type:     schema.TypeString,
			Required: true,
			Description: `Identifier for the mesh in which this Istio service is defined.
Corresponds to the meshUid metric label in Istio metrics.`,
		},
		"service_name": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The name of the Istio service underlying this service.
Corresponds to the destination_service_name metric label in
Istio metrics.`,
		},
		"service_namespace": {
			Type:     schema.TypeString,
			Required: true,
			Description: `The namespace of the Istio service underlying this service.
Corresponds to the destination_service_namespace metric label in
Istio metrics.`,
		},
	}

	filter := strings.Join([]string{
		`mesh_istio.mesh_uid="{{}}"`,
		`mesh_istio.service_namespace="{{}}"`,
		`mesh_istio.service_name="{{}}"`,
	}, "&")

	return dataSourceMonitoringServiceType(istioSchema, filter, dataSourceMonitoringServiceMeshIstioRead)
}

func dataSourceMonitoringServiceMeshIstioRead(res map[string]interface{}, d *schema.ResourceData, meta interface{}) error {
	var istio map[string]interface{}
	if v, ok := res["mesh_istio"]; ok {
		istio = v.(map[string]interface{})
	}
	if len(istio) == 0 {
		return nil
	}

	if err := d.Set("mesh_uid", istio["mesh_uid"]); err != nil {
		return err
	}
	if err := d.Set("service_name", istio["service_name"]); err != nil {
		return err
	}
	if err := d.Set("service_namespace", istio["service_namespace"]); err != nil {
		return err
	}

	return nil
}
