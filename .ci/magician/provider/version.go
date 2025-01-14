package provider

type Version int

const (
	None Version = iota
	GA
	Beta
	Private
)

const NumVersions = 3

func (v Version) String() string {
	switch v {
	case GA:
		return "ga"
	case Beta:
		return "beta"
	case Private:
		return "private"
	}
	return "unknown"
}

func (v Version) ProviderName() string {
	switch v {
	case GA:
		return "google"
	case Beta:
		return "google-beta"
	case Private:
		return "google-private"
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
	case Private:
		return "terraform-next"
	}
	return "unknown"
}
