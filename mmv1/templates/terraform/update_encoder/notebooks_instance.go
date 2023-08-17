// Update requests use "items" as the api name instead of "metadata"
// https://cloud.google.com/vertex-ai/docs/workbench/reference/rest/v1/projects.locations.instances/updateMetadataItems
if metadata, ok := obj["metadata"]; ok {
	obj["items"] = metadata
	delete(obj, "metadata")
}
return obj, nil
