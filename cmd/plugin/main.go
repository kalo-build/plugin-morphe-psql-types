package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kalo-build/plugin-morphe-psql-types/pkg/compile"
)

type CompileConfig struct {
	InputPath  string         `json:"inputPath"`
	OutputPath string         `json:"outputPath"`
	Config     map[string]any `json:"config,omitempty"`
	Verbose    bool           `json:"verbose,omitempty"`
}

const (
	ErrMissingConfig      = 3
	ErrInvalidConfig      = 4
	ErrInputPathRequired  = 12
	ErrOutputPathRequired = 13
	ErrCompileFailed      = 1
)

// logInfo prints info messages only when verbose mode is enabled
func logInfo(verbose bool, format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: plugin-morphe-psql-types <config>")
		fmt.Fprintln(os.Stderr, "  config: JSON string with inputPath, outputPath, and optional config parameters")
		os.Exit(ErrMissingConfig)
	}

	rawConfig := os.Args[1]
	var compileConfig CompileConfig
	if err := json.Unmarshal([]byte(rawConfig), &compileConfig); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing config JSON:", err)
		fmt.Fprintln(os.Stderr, "Expected format: {\"inputPath\":\"...\",\"outputPath\":\"...\",\"config\":{...},\"verbose\":false}")
		os.Exit(ErrInvalidConfig)
	}

	if compileConfig.InputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Input path is required")
		os.Exit(ErrInputPathRequired)
	}

	if compileConfig.OutputPath == "" {
		fmt.Fprintln(os.Stderr, "Error: Output path is required")
		os.Exit(ErrOutputPathRequired)
	}

	inputAbs, err := filepath.Abs(compileConfig.InputPath)
	if err == nil {
		compileConfig.InputPath = inputAbs
	}

	outputAbs, err := filepath.Abs(compileConfig.OutputPath)
	if err == nil {
		compileConfig.OutputPath = outputAbs
	}

	logInfo(compileConfig.Verbose, "Processing Morphe registry from: '%s'", compileConfig.InputPath)
	logInfo(compileConfig.Verbose, "Output PostgreSQL types to: '%s'", compileConfig.OutputPath)

	morpheConfig := compile.DefaultMorpheCompileConfig(
		compileConfig.InputPath,
		compileConfig.OutputPath,
	)

	// Check for ordered migrations config option
	if orderedMigrations, ok := compileConfig.Config["orderedMigrations"].(bool); ok && orderedMigrations {
		morpheConfig.EnableOrderedMigrations = true
		logInfo(compileConfig.Verbose, "Ordered migrations enabled - files will have numeric prefixes")
	}

	logInfo(compileConfig.Verbose, "Starting compilation process...")
	compileErr := compile.MorpheToPSQL(morpheConfig)
	if compileErr != nil {
		fmt.Fprintln(os.Stderr, "Compilation failed:", compileErr)
		os.Exit(ErrCompileFailed)
	}

	logInfo(compileConfig.Verbose, "Compilation completed successfully")
	os.Exit(0)
}
