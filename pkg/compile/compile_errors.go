package compile

import (
	"errors"
	"fmt"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

var ErrNoRegistry = errors.New("registry not initialized")

func ErrUnsupportedMorpheFieldType[TType yaml.ModelFieldType | yaml.StructureFieldType](unsupportedType TType) error {
	return fmt.Errorf("unsupported morphe field type for go conversion: '%s'", unsupportedType)
}

func ErrMissingMorpheIdentifierField(modelName string, identifierName string, fieldName string) error {
	return fmt.Errorf("morphe model '%s' has no field '%s' referenced in identifiers ('%s')", modelName, identifierName, fieldName)
}
