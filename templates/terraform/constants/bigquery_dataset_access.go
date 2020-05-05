var bigqueryAccessRoleToPrimitiveMap =  map[string]string {
    "roles/bigQuery.dataOwner": "OWNER",
    "roles/bigQuery.dataEditor": "EDITOR",
    "roles/bigQuery.dataViewer": "VIEWER",
}

func resourceBigQueryDatasetAccessRoleDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    if primitiveRole, ok := bigqueryAccessRoleToPrimitiveMap[new]; ok {
        return primitiveRole == old
    }
    return false
}
