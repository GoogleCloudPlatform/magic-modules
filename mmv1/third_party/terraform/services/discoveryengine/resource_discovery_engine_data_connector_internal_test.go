package discoveryengine

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceDiscoveryEngineDataConnectorDecoder(t *testing.T) {
	resourceSchema := ResourceDiscoveryEngineDataConnector().Schema

	t.Run("outlook root placeholder filtered when not in config", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"entities": []interface{}{
				map[string]interface{}{"entity_name": "mail"},
				map[string]interface{}{"entity_name": "calendar"},
			},
		})

		apiRes := map[string]interface{}{
			"dataSource": "outlook",
			"entities": []interface{}{
				map[string]interface{}{
					"entityName": "mail",
					"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_mail",
				},
				map[string]interface{}{
					"entityName": "calendar",
					"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_calendar",
				},
				map[string]interface{}{
					"entityName": "root",
					"params":     map[string]interface{}{},
				},
			},
		}

		got, err := resourceDiscoveryEngineDataConnectorDecoder(d, nil, apiRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []interface{}{
			map[string]interface{}{
				"entityName": "mail",
				"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_mail",
			},
			map[string]interface{}{
				"entityName": "calendar",
				"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_calendar",
			},
		}
		if !reflect.DeepEqual(got["entities"], expected) {
			t.Errorf("expected %#v, got %#v", expected, got["entities"])
		}
	})

	t.Run("sharepoint site placeholder filtered when not in config", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"entities": []interface{}{
				map[string]interface{}{"entity_name": "file"},
				map[string]interface{}{"entity_name": "page"},
			},
		})

		apiRes := map[string]interface{}{
			"dataSource": "sharepoint",
			"entities": []interface{}{
				map[string]interface{}{
					"entityName": "file",
					"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_file",
				},
				map[string]interface{}{
					"entityName": "page",
					"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_page",
				},
				map[string]interface{}{
					"entityName": "site",
					"params":     map[string]interface{}{},
				},
			},
		}

		got, err := resourceDiscoveryEngineDataConnectorDecoder(d, nil, apiRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []interface{}{
			map[string]interface{}{
				"entityName": "file",
				"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_file",
			},
			map[string]interface{}{
				"entityName": "page",
				"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_page",
			},
		}
		if !reflect.DeepEqual(got["entities"], expected) {
			t.Errorf("expected %#v, got %#v", expected, got["entities"])
		}
	})

	t.Run("placeholder preserved when explicitly declared in config", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
			"entities": []interface{}{
				map[string]interface{}{"entity_name": "mail"},
				map[string]interface{}{"entity_name": "root"},
			},
		})

		entities := []interface{}{
			map[string]interface{}{
				"entityName": "mail",
				"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_mail",
			},
			map[string]interface{}{
				"entityName": "root",
				"params":     map[string]interface{}{},
			},
		}
		apiRes := map[string]interface{}{
			"dataSource": "outlook",
			"entities":   entities,
		}

		got, err := resourceDiscoveryEngineDataConnectorDecoder(d, nil, apiRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got["entities"], entities) {
			t.Errorf("expected %#v, got %#v", entities, got["entities"])
		}
	})

	t.Run("import filters empty placeholder but preserves placeholder with parent params", func(t *testing.T) {
		dEmpty := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{})

		outlookRes := map[string]interface{}{
			"dataSource": "outlook",
			"entities": []interface{}{
				map[string]interface{}{
					"entityName": "mail",
					"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_mail",
				},
				map[string]interface{}{
					"entityName": "root",
					"params":     map[string]interface{}{},
				},
			},
		}
		gotOutlook, err := resourceDiscoveryEngineDataConnectorDecoder(dEmpty, nil, outlookRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotOutlook["entities"].([]interface{})) != 1 {
			t.Errorf("expected 1 entity on outlook import, got %#v", gotOutlook["entities"])
		}

		servicenowRes := map[string]interface{}{
			"dataSource": "servicenow",
			"entities": []interface{}{
				map[string]interface{}{
					"entityName": "catalog",
					"dataStore":  "projects/p/locations/global/collections/default_collection/dataStores/ds_catalog",
				},
				map[string]interface{}{
					"entityName": "knowledge_base",
					"params": map[string]interface{}{
						"inclusion_filters": map[string]interface{}{
							"knowledgeBaseSysId": []interface{}{"123"},
						},
					},
				},
			},
		}
		gotServicenow, err := resourceDiscoveryEngineDataConnectorDecoder(dEmpty, nil, servicenowRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(gotServicenow["entities"].([]interface{})) != 2 {
			t.Errorf("expected 2 entities on servicenow import with parent params, got %#v", gotServicenow["entities"])
		}
	})

	t.Run("gcs_asp entities without dataStore are preserved", func(t *testing.T) {
		entities := []interface{}{
			map[string]interface{}{
				"entityName": "gcs_store",
			},
		}
		apiRes := map[string]interface{}{
			"dataSource": "gcs_asp",
			"entities":   entities,
		}

		got, err := resourceDiscoveryEngineDataConnectorDecoder(nil, nil, apiRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got["entities"], entities) {
			t.Errorf("expected %#v, got %#v", entities, got["entities"])
		}
	})
}
