package cfg

// MorpheConfig is the main configuration for PostgreSQL compilation
type MorpheConfig struct {
	MorpheModelsConfig
	MorpheEnumsConfig
	MorpheStructuresConfig
	MorpheEntitiesConfig
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

	structuresErr := config.MorpheStructuresConfig.Validate()
	if structuresErr != nil {
		return structuresErr
	}

	entitiesErr := config.MorpheEntitiesConfig.Validate()
	if entitiesErr != nil {
		return entitiesErr
	}

	return nil
}

// DefaultMorpheConfig returns a default configuration
func DefaultMorpheConfig() MorpheConfig {
	return MorpheConfig{
		MorpheModelsConfig: MorpheModelsConfig{
			Schema: DefaultSchema,
		},
		MorpheEnumsConfig: MorpheEnumsConfig{
			Schema: DefaultSchema,
		},
		MorpheStructuresConfig: MorpheStructuresConfig{
			Schema:            DefaultSchema,
			EnablePersistence: false,
		},
		MorpheEntitiesConfig: MorpheEntitiesConfig{
			Schema:         DefaultSchema,
			ViewNameSuffix: "_entities",
		},
	}
}
