package google

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// jsonMap converts a given value to a map[string]interface{} that
// matches its JSON format.
func jsonMap(x interface{}) (map[string]interface{}, error) {
	jsn, err := json.Marshal(x)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling")
	}

	m := make(map[string]interface{})
	if err := json.Unmarshal(jsn, &m); err != nil {
		return nil, errors.Wrap(err, "unmarshalling")
	}

	return m, nil
}
