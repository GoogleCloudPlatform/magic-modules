package artifactregistry

import "net/url"

// artifactResourceRef builds the trailing "{name}{sep}{version}" path component
// of an Artifact Registry resource URL from caller-supplied identifiers (an
// image_name/package_name/artifact_id and its tag or digest).
//
// The name keeps the query-escaping the data sources already applied. The
// version/digest was previously interpolated raw, so a value carrying "/", "?"
// or "#" could add path segments, a query string, or a fragment to the request
// URL and reach a different resource than the one configured. Path-escaping it
// leaves a valid semver or "sha256:<hex>" digest unchanged while neutralizing
// those characters.
func artifactResourceRef(name, sep, version string) string {
	if version == "" {
		return url.QueryEscape(name)
	}
	return url.QueryEscape(name) + sep + url.PathEscape(version)
}
