package google

import (
	"os"
)

// MultiEnvDefaultFunc is a helper function that returns the value of the first
// environment variable in the given list that returns a non-empty value. If
// none of the environment variables return a value, the default value is
// returned.
func MultiEnvDefault(ks []string, dv interface{}) interface{} {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return dv
}
