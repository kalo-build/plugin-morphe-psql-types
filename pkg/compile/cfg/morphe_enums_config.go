package cfg

// MorpheEnumsConfig holds configuration specific to PostgreSQL enum tables
type MorpheEnumsConfig struct {
	// Schema to use for enum tables
	Schema string

	// Whether to use BIGSERIAL instead of SERIAL for auto-increment fields
	UseBigSerial bool
}

// Validate checks if the models configuration is valid
func (config MorpheEnumsConfig) Validate() error {
	if config.Schema == "" {
		return ErrNoEnumSchema
	}

	return nil
}
