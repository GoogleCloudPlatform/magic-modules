package common

// ConverterMap is a configuration object to map assets to converters.
type ConverterMap struct {
	// AssetTypeToConverterName maps asset type to converter name.
	AssetTypeToConverterName map[string]string

	// ConverterNameToConverter maps converter name to converter instance.
	ConverterNameToConverter map[string]Converter
}
