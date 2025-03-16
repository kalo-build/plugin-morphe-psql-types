package cfg

import "errors"

var ErrNoSchema = errors.New("schema cannot be empty")
var ErrNoModelSchema = errors.New("model schema cannot be empty")
var ErrNoEnumSchema = errors.New("enum schema cannot be empty")
