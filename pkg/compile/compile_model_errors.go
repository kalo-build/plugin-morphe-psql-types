package compile

import "errors"

var ErrNoModelTables = errors.New("no model tables provided")
var ErrNoModelTable = errors.New("no model table provided")
var ErrNoStructureTable = errors.New("no structure table provided")
var ErrNoStructureWriter = errors.New("structure writer must be provided when structure persistence is enabled")
