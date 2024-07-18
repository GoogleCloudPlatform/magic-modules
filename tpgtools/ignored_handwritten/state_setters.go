package google

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setStateForCollapsedObject sets nested fields as top-level state values for an object that should be collapsed.
// It takes in the flattened value of the object (i.e. converted from *obj.FooBar to []interface{}{map[string]interface{}}).
// Example Usage:
// err := setStateForFieldsInFlattenedObject(d, flattenMyTopLevelNestedObject(obj.MyTopLevelNestedObject))
func setStateForCollapsedObject(d *schema.ResourceData, v interface{}) error {
	if v == nil {
		return nil
	}

	ls, ok := v.([]interface{})
	if !ok {
		return fmt.Errorf("expected nested object value to be flattened to []interface{}")
	}
	if len(ls) == 0 {
		return nil
	}

	nestedObj := ls[0].(map[string]interface{})
	for k, kv := range nestedObj {
		if err := d.Set(k, kv); err != nil {
			return fmt.Errorf("error setting %s in state: %s", k, err)
		}
	}
	return nil
}

func setStateForTime(d *schema.ResourceData, v time.Time, name string) error {
	if !v.IsZero() {
		return d.Set(name, fmt.Sprintf(v.Format(time.RFC3339)))
	} else {
		return d.Set(name, nil)
	}
}

func generateIfNotSet(d *schema.ResourceData, field, prefix string) (string, error) {
	if _, ok := d.GetOkExists(field); !ok {
		if prefix == "" {
			prefix = "tf-generated-"
		}
		v := id.PrefixedUniqueId(prefix)
		if len(v) > 30 {
			v = v[:30]
		}

		if err := d.Set(field, v); err != nil {
			return "", err
		}
	}
	return d.Get(field).(string), nil
}
