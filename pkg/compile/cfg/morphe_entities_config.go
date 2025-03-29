package cfg

import "fmt"

// MorpheEntitiesConfig defines configuration options for compiling Morphe entities to PostgreSQL views
type MorpheEntitiesConfig struct {
	// Schema is the PostgreSQL schema name to use for generated views
	Schema string

	// ViewNameSuffix is appended to view names (default: "_entities")
	ViewNameSuffix string
}

// Validate validates the MorpheEntitiesConfig
func (c MorpheEntitiesConfig) Validate() error {
	if c.Schema == "" {
		return fmt.Errorf("schema is required")
	}
	return nil
}
