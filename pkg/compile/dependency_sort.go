package compile

import (
	"fmt"
	"sort"

	"github.com/kalo-build/plugin-morphe-psql-types/pkg/psqldef"
)

// TableDependencyGraph represents the dependency relationships between tables
type TableDependencyGraph struct {
	// tableDeps maps table name -> list of tables it depends on (via FK)
	tableDeps map[string][]string
	// allTables is the set of all known table names
	allTables map[string]bool
}

// NewTableDependencyGraph creates a new dependency graph
func NewTableDependencyGraph() *TableDependencyGraph {
	return &TableDependencyGraph{
		tableDeps: make(map[string][]string),
		allTables: make(map[string]bool),
	}
}

// AddTable adds a table and its FK dependencies to the graph
func (g *TableDependencyGraph) AddTable(table *psqldef.Table) {
	tableName := table.Name
	g.allTables[tableName] = true

	deps := []string{}
	for _, fk := range table.ForeignKeys {
		// Only add dependency if it's not self-referential
		if fk.RefTableName != tableName {
			deps = append(deps, fk.RefTableName)
			g.allTables[fk.RefTableName] = true
		}
	}
	g.tableDeps[tableName] = deps
}

// TopologicalSort returns tables sorted so that dependencies come before dependents
// Returns an error if there's a circular dependency
func (g *TableDependencyGraph) TopologicalSort() ([]string, error) {
	// Kahn's algorithm for topological sorting
	// Calculate in-degree for each node (considering only tables we're writing)
	inDegree := make(map[string]int)
	for table := range g.tableDeps {
		if _, exists := inDegree[table]; !exists {
			inDegree[table] = 0
		}
		for _, dep := range g.tableDeps[table] {
			// Only count dependencies on tables we're actually writing
			if _, exists := g.tableDeps[dep]; exists {
				inDegree[table]++
			}
		}
	}

	// Find all nodes with in-degree 0
	queue := []string{}
	for table := range g.tableDeps {
		if inDegree[table] == 0 {
			queue = append(queue, table)
		}
	}
	// Sort the initial queue for deterministic output
	sort.Strings(queue)

	result := []string{}
	for len(queue) > 0 {
		// Take first element (already sorted)
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)

		// Find tables that depend on current and reduce their in-degree
		newQueue := []string{}
		for table, deps := range g.tableDeps {
			for _, dep := range deps {
				if dep == current {
					inDegree[table]--
					if inDegree[table] == 0 {
						newQueue = append(newQueue, table)
					}
				}
			}
		}
		// Sort new entries before adding to queue
		sort.Strings(newQueue)
		queue = append(queue, newQueue...)
	}

	// Check for cycles
	if len(result) != len(g.tableDeps) {
		return nil, fmt.Errorf("circular dependency detected in table definitions")
	}

	return result, nil
}

// SortTablesByDependency sorts tables so that dependencies come before dependents
func SortTablesByDependency(tables []*psqldef.Table) ([]*psqldef.Table, error) {
	graph := NewTableDependencyGraph()
	tableMap := make(map[string]*psqldef.Table)

	for _, table := range tables {
		graph.AddTable(table)
		tableMap[table.Name] = table
	}

	sortedNames, err := graph.TopologicalSort()
	if err != nil {
		return nil, err
	}

	sortedTables := make([]*psqldef.Table, 0, len(sortedNames))
	for _, name := range sortedNames {
		if table, exists := tableMap[name]; exists {
			sortedTables = append(sortedTables, table)
		}
	}

	return sortedTables, nil
}
