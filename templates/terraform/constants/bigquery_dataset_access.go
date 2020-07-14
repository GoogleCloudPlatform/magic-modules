var bigqueryAccessRoleToPrimitiveMap =  map[string]string {
    "roles/bigquery.dataOwner": "OWNER",
    "roles/bigquery.dataEditor": "WRITER",
    "roles/bigquery.dataViewer": "READER",
}

var bigqueryAccessIamMemberToTypeMap = map[string]string{
	"serviceAccount": "user_by_email",
	"user":           "user_by_email",
	"group":          "group_by_email",
	"domain":         "domain",
	"specialGroup":   "special_group",
}

func resourceBigQueryDatasetAccessRoleDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    if primitiveRole, ok := bigqueryAccessRoleToPrimitiveMap[new]; ok {
        return primitiveRole == old
    }
    return false
}

// iam_member can be passed into the request, but the response will have the value in one of
// UserByEmail, GroupByEmail, Domain, or SpecialGroup fields. The key is determined by the
// prefix of the iam_member value, and the value follows the : of the prefix.
// Instead of dealing with the issues in the response, we'll do the translation before we
// request.
func customDiffBigQueryDatasetAccess(d *schema.ResourceDiff, meta interface{}) error {
	if !d.NewValueKnown("iam_member") {
		return nil
	}

	_, configValue := d.GetChange("iam_member")

	parts := strings.Split(configValue.(string), ":")
	if len(parts) != 2 {
		return nil
	}

	var key string
	if k, ok := bigqueryAccessIamMemberToTypeMap[parts[0]]; !ok {
		return nil
	} else {
		key = k

		if err := d.Clear("iam_member"); err != nil {
			return err
		}
	}
	value := parts[1]

	return d.SetNew(key, value)
}