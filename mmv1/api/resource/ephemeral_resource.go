package resource

type EphemeralResource struct {
	// boolean to determine whether the datasource file should be generated
	Generate bool `yaml:"generate"`
	// boolean to determine whether tests should be generated for a ephemeral resource
	ExcludeTest bool `yaml:"exclude_test"`
}