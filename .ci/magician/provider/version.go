package provider

type Version int

const (
	None Version = iota
	GA
	Beta
)

const NumVersions = 2

func (v Version) String() string {
	switch v {
	case GA:
		return "ga"
	case Beta:
		return "beta"
	}
	return "unknown"
}

func (v Version) BucketPath() string {
	if v == GA {
		return ""
	}
	return v.String() + "/"
}

func (v Version) RepoName() string {
	switch v {
	case GA:
		return "terraform-provider-google"
	case Beta:
		return "terraform-provider-google-beta"
	}
	return "unknown"
}
