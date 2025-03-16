package cfg

// MorpheStructuresConfig holds configuration specific to PostgreSQL structure tables
type MorpheStructuresConfig struct {
	// Schema to use for structure tables
	Schema string

	// Whether to use BIGSERIAL instead of SERIAL for auto-increment fields
	UseBigSerial bool

	// Whether to enable structure persistence
	EnablePersistence bool
}

// Validate checks if the structures configuration is valid
func (config MorpheStructuresConfig) Validate() error {
	if config.EnablePersistence && config.Schema == "" {
		return ErrNoStructureSchema
	}

	return nil
}
