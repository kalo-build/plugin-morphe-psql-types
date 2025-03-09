package cfg

import "fmt"

// MorpheModelsConfig holds configuration specific to PostgreSQL model tables
type MorpheModelsConfig struct {
	// Schema to use for model tables
	Schema string

	// Whether to use BIGSERIAL instead of SERIAL for auto-increment fields
	UseBigSerial bool
}

// Validate checks if the models configuration is valid
func (config MorpheModelsConfig) Validate() error {
	if config.Schema == "" {
		return fmt.Errorf("models %w", ErrNoModelSchema)
	}

	return nil
}
