package psqldef

import (
	"github.com/barkimedes/go-deepcopy"
	"github.com/kalo-build/clone"
)

// DeepClonePSQLTypeMap attempts to deep clone a map of PSQLTypes.
//
// (Potentially unsafe, see docs for `DeepClonePSQLType`)
func DeepClonePSQLTypeMap[TType PSQLType](original map[string]TType) map[string]TType {
	if original == nil {
		return nil
	}
	newMap := make(map[string]TType, len(original))
	for key, ttype := range original {
		newMap[key] = DeepClonePSQLType(ttype)
	}
	return newMap
}

// DeepClonePSQLTypeSlice attempts to deep clone a slice of PSQLTypes.
//
// (Potentially unsafe, see docs for `DeepClonePSQLType`)
func DeepClonePSQLTypeSlice[TType PSQLType](original []TType) []TType {
	if original == nil {
		return nil
	}
	newSlice := make([]TType, len(original))
	for idx, ttype := range original {
		newSlice[idx] = DeepClonePSQLType(ttype)
	}
	return newSlice
}

// DeepClonePSQLType attempts to deep clone a PSQLType.
//
// If the passed PSQLType type implements `clone.DeepCloneable[TType]`, the type method's clone itself is used. This is the preferred method for
// all deep clones.
//
// However, if this fails we attempt to make a deepcopy (excluding functions, channels, and unsafe pointers)
// and then cast the result to the target type.
//
// If all else fails, we do not deep clone and instead return the input, potentially leading to side-effects.
func DeepClonePSQLType[TType PSQLType](original TType) TType {
	var originalAny any = original
	deepCloneable, isCloneable := originalAny.(clone.DeepCloneable[TType])
	if isCloneable {
		return deepCloneable.DeepClone()
	}
	deepClone, deepCloneErr := deepcopy.Anything(original)
	if deepCloneErr != nil {
		return original
	}
	typedClone, isTType := deepClone.(TType)
	if !isTType {
		return original
	}
	return typedClone
}
