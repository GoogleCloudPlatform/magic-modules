package resource

// Sweeper provides configuration for the test sweeper
type Sweeper struct {
	// The field checked by sweeper to determine
	// eligibility for deletion for generated resources
	IdentifierField string   `yaml:"identifier_field"`
	Regions                  []string          `yaml:"regions,omitempty"`
	Prefixes                 []string          `yaml:"prefixes,omitempty"`
	URLSubstitutions         []URLSubstitution `yaml:"url_substitutions,omitempty"`
}

// URLSubstitution represents a region-zone pair for URL substitution
type URLSubstitution struct {
	Region string `yaml:"region,omitempty"`
	Zone   string `yaml:"zone,omitempty"`
>>>>>>> Stashed changes
}
