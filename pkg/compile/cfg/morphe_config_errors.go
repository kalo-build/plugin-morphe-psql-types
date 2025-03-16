package cfg

import "errors"

var ErrNoSchema = errors.New("schema cannot be empty")
var ErrNoModelSchema = errors.New("model schema cannot be empty")
var ErrNoEnumSchema = errors.New("enum schema cannot be empty")
var ErrNoStructureSchema = errors.New("structure schema cannot be empty when persistence is enabled")
