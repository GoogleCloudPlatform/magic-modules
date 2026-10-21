package publicca

import (
	"encoding/base64"
	"fmt"
)

// Public CA returns an outer base64url encoding around the ACME MAC key.
func flattenEphemeralPublicCAMacKey(v interface{}) (interface{}, error) {
	if v == nil {
		return "", nil
	}
	encoded, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("expected a string, got %T", v)
	}
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	return string(decoded), nil
}

func flattenEphemeralPublicCAB64URLMacKey(v interface{}) (interface{}, error) {
	decoded, err := flattenEphemeralPublicCAMacKey(v)
	if err != nil {
		return nil, err
	}
	return base64.RawURLEncoding.EncodeToString([]byte(decoded.(string))), nil
}
