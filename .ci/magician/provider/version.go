package provider

type Version int

const (
	None Version = iota
	GA
	Beta
<<<<<<< HEAD
)

const NumVersions = 2
=======
	Alpha
)

const NumVersions = 3
>>>>>>> 2fdda66097e2c96688e59f7c58c1f717c7785856

func (v Version) String() string {
	switch v {
	case GA:
		return "ga"
	case Beta:
		return "beta"
<<<<<<< HEAD
=======
	case Alpha:
		return "alpha"
>>>>>>> 2fdda66097e2c96688e59f7c58c1f717c7785856
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
<<<<<<< HEAD
=======
	case Alpha:
		return "terraform-next"
>>>>>>> 2fdda66097e2c96688e59f7c58c1f717c7785856
	}
	return "unknown"
}
