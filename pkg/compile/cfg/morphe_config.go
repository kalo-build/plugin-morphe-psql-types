package cfg

// MorpheConfig is the main configuration for PostgreSQL compilation
type MorpheConfig struct {
	MorpheModelsConfig
	MorpheEnumsConfig
}

// Default schema
const (
	DefaultSchema = "public"
)

// Validate checks if the configuration is valid
func (config MorpheConfig) Validate() error {
	// Validate each component config
	modelsErr := config.MorpheModelsConfig.Validate()
	if modelsErr != nil {
		return modelsErr
	}

	enumsErr := config.MorpheEnumsConfig.Validate()
	if enumsErr != nil {
		return enumsErr
	}

	return nil
}

// DefaultMorpheConfig returns a default configuration
func DefaultMorpheConfig() MorpheConfig {
	return MorpheConfig{
		MorpheModelsConfig: MorpheModelsConfig{
			Schema: DefaultSchema,
		},
	}
}
