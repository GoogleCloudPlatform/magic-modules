var bigqueryAccessRoleToPrimitiveMap =  map[string]string {
    "roles/bigquery.dataOwner": "OWNER",
    "roles/bigquery.dataEditor": "WRITER",
    "roles/bigquery.dataViewer": "READER",
}

func resourceBigQueryDatasetAccessRoleDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    if primitiveRole, ok := bigqueryAccessRoleToPrimitiveMap[new]; ok {
        return primitiveRole == old
    }
    return false
}
