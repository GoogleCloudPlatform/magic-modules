package datalineage

import (
	"github.com/OpenLineage/openlineage/client/go/pkg/facets"
	"github.com/OpenLineage/openlineage/client/go/pkg/openlineage"
)

func buildRunEvent(in map[string]interface{}) *openlineage.RunEvent {
	runEvent := openlineage.NewNamespacedRunEvent(openlineage.EventTypeComplete, openlineage.NewRunID(), in["name"].(string), in["namespace"].(string), "_PRODUCER_")
	buildJobType(runEvent)
	buildOwners(in, runEvent)
	buildInputs(in, runEvent)
	buildOutputs(in, runEvent)
	return runEvent
}

func buildOutputs(in map[string]interface{}, runEvent *openlineage.RunEvent) {
	if v, ok := in["output"]; ok {
		outputs := make([]openlineage.OutputElement, 0)
		v := v.([]interface{})
		for _, item := range v {
			m := item.(map[string]interface{})
			element := openlineage.OutputElement{
				Name:      m["name"].(string),
				Namespace: m["namespace"].(string),
			}
			facets := getCommonDatasetFacets(m)
			element.WithFacets(facets...)
			if e, ok := m["columnLineage"].([]interface{}); ok {
				element.WithFacets(buildColumnLineage(e))
			}

			outputs = append(outputs, element)
		}
		runEvent.WithOutputs(outputs...)
	}
}

func buildInputs(in map[string]interface{}, runEvent *openlineage.RunEvent) {
	if v, ok := in["input"]; ok {
		v := v.([]interface{})
		inputs := make([]openlineage.InputElement, 0)
		for _, item := range v {
			m := item.(map[string]interface{})
			element := openlineage.InputElement{
				Name:      m["name"].(string),
				Namespace: m["namespace"].(string),
			}
			facets := getCommonDatasetFacets(m)
			element.WithFacets(facets...)
			inputs = append(inputs, element)
		}
		runEvent.WithInputs(inputs...)
	}
}

func buildOwners(in map[string]interface{}, runEvent *openlineage.RunEvent) {
	if v, ok := in["owner"]; ok {
		v := v.([]interface{})
		owners := make([]facets.OwnershipJobFacetOwner, 0, len(v))
		for _, item := range v {
			m := item.(map[string]interface{})

			t := m["type"].(string)
			owners = append(owners, facets.OwnershipJobFacetOwner{
				Name: m["name"].(string),
				Type: &t,
			})
		}
		runEvent.WithJobFacets(facets.NewOwnershipJobFacet("_PRODUCER_").WithOwners(owners))
	}
}

func buildJobType(runEvent *openlineage.RunEvent) {
	runEvent.WithJobFacets(
		facets.NewJobTypeJobFacet("_PRODUCER_", "TERRAFORM", "BYOL"),
	)
}

func buildColumnLineage(v []interface{}) *facets.ColumnLineageDatasetFacet {
	if len(v) == 0 {
		return nil
	}
	cll := v[0].(map[string]interface{})
	fields := make(map[string]facets.ColumnLineageDatasetFacetFieldsValue)
	for _, item := range cll["fields"].([]interface{}) {
		m := item.(map[string]interface{})
		name := m["name"].(string)
		i := m["input"].([]interface{})
		fields[name] = facets.ColumnLineageDatasetFacetFieldsValue{
			InputFields: buildCllInputs(i),
		}
	}

	di := cll["dataset_input"].([]interface{})
	facet := facets.NewColumnLineageDatasetFacet("_PRODUCER_", fields).WithDataset(buildCllInputs(di))
	return facet
}

func buildCllInputs(i []interface{}) []facets.InputField {
	in := make([]facets.InputField, 0)

	for _, it := range i {
		m2 := it.(map[string]interface{})
		iname := m2["name"].(string)
		inamespace := m2["namespace"].(string)
		field := m2["field"].(string)
		transformations := make([]facets.InputFieldTransformation, 0)
		for _, transformation := range m2["transformation"].([]interface{}) {
			tr := transformation.(map[string]interface{})
			subtype := tr["subtype"].(string)
			transformations = append(transformations, facets.InputFieldTransformation{
				Type:    tr["type"].(string),
				Subtype: &subtype,
			})
		}
		in = append(in, facets.InputField{
			Name:            iname,
			Namespace:       inamespace,
			Field:           field,
			Transformations: transformations,
		})

	}
	return in
}

func getCommonDatasetFacets(m map[string]interface{}) []facets.DatasetFacet {
	facets := make([]facets.DatasetFacet, 0)
	if v, ok := m["catalog"]; ok {
		facets = append(facets, buildCatalog(v.(map[string]interface{})))
	}

	if v, ok := m["symlink"]; ok {
		facets = append(facets, buildSymlinks(v.([]interface{})))
	}
	return facets
}

func buildSymlinks(v []interface{}) *facets.SymlinksDatasetFacet {
	if len(v) == 0 {
		return nil
	}

	symlinks := make([]facets.SymlinksDatasetFacetIdentifier, 0, len(v))
	for _, item := range v {
		m := item.(map[string]interface{})
		symlinks = append(symlinks, facets.SymlinksDatasetFacetIdentifier{
			Name:      m["name"].(string),
			Namespace: m["namespace"].(string),
			Type:      m["type"].(string),
		})
	}
	return facets.NewSymlinksDatasetFacet("_PRODUCER_").WithIdentifiers(symlinks)
}

func buildCatalog(v map[string]interface{}) *facets.CatalogDatasetFacet {
	return facets.NewCatalogDatasetFacet("_PRODUCER_", v["framework"].(string), v["name"].(string), v["type"].(string))
}

func flattenKnowledgeCatalog(
	process, run string) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"process": process,
			"run":     run,
		},
	}
}
