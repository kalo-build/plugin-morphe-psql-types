# Relationship Aliasing Implementation Plan

## 0. Rules & Principles

- **Consistency**: Follow existing codebase patterns (naming conventions, error handling, test structure)
- **Minimal scope**: Each task focuses on one specific feature implementation
- **Minimal redundancy**: New code attempts to leverage existing code when there is domain- / usage-level overlap to mitigate code fragmentation and manual update friction
- **Cohesion <> Coupling**: Apply minimal coupling and maximum cohesion software principles
- **Idiomatic Go**: Clean, readable code following Go best practices
- **Clean Code**: Clean code practices like "no deep nesting", early return, SRP, ...
- **TDD approach**: Red-Green-Refactor cycle with 1-2 refactor loops per task, and tests for realistic edge cases (including error cases)
- **Testing style**: Use testify assertions; one test per case for simple functions (boolean/regex), arrange-act-assert blocks for complex tests; minimize comments
- **Human review**: Mandatory pause after each green run for testing and feedback
- **Backward compatibility**: All changes must not break existing functionality
- **Database integrity**: All aliased relationships must maintain referential integrity via proper constraint naming
- **Planning**: After task completion, reevaluate upcoming tasks and revisit adding new / updating existing

## 1. Codebase ⇄ Spec Gap Analysis

### Current State
- **Supported Relations**: `ForOne`, `ForMany`, `HasOne`, `HasMany`, `ForOnePoly`, `ForManyPoly`, `HasOnePoly`, `HasManyPoly`
- **Relationship Names**: Direct mapping between relationship name and target model name
- **SQL Generation**: Foreign keys, junction tables, constraints use relationship names directly
- **Entity References**: Field paths reference relationships by their direct names

### Missing Components for Aliasing

#### Core Logic Modules
- `morphe-go/pkg/yaml/model_relation.go`: ✅ **DONE** - Added `aliased` field in ModelRelation structure
- `morphe-go/pkg/yaml/entity_relation.go`: ✅ **DONE** - Added `aliased` field in EntityRelation structure
- `morphe-go/pkg/yamlops/relation.go`: ✅ **DONE** - Added aliasing detection functions (IsRelationAliased, GetRelationTargetName)
- `morphe-go/pkg/yaml/normalize.go`: ✅ **DONE** - Added whitespace normalization for aliased fields
- `pkg/compile/compile_models.go`: ⚠️ **PENDING** - Multiple functions need updates:
  - `getColumnsForModelRelations` (lines 207-279): Uses relatedModelName directly
  - `getForeignKeysForModelRelations` (lines 280-333): Uses relatedModelName directly  
  - `getTablesForModelManyRelations` (lines 390-528): Uses relatedModelName for junction tables
- `pkg/compile/compile_entities.go`: ⚠️ **PENDING** - Entity traversal issues:
  - `traverseRelationshipChain` (line 198): Uses relationName as modelName
  - `setupJoinsForRegularRelationships` (lines 302-314): Uses relationship name for joins
  - `addJoinClause` (lines 317-381): Needs correct target model resolution

#### Database Schema Components
- **ForOne Aliasing**: Missing aliased foreign key column generation (`{alias_name}_id`)
- **ForMany Aliasing**: Missing aliased junction table creation (`{source}_{alias_name}`)
- **HasOne/HasMany Aliasing**: No validation for aliased inverse relationships
- **Polymorphic Aliasing**: Missing aliased polymorphic column generation (`{alias_name}_type`, `{alias_name}_id`)
- **Polymorphic Inverse Aliasing**: Missing validation for `HasOnePoly`/`HasManyPoly` with `aliased` + `through` pattern
- **Constraint Naming**: No aliasing support in foreign key and unique constraint names

#### Validation Components
- **Aliased Target Validation**: ✅ **DONE** - Basic validation that `aliased` property references valid models/entities
- **Polymorphic Inverse Validation**: ✅ **DONE** - Complex validation for `HasOnePoly`/`HasManyPoly` with `through` + `aliased` pattern
- **Whitespace Normalization**: ✅ **DONE** - Aliased fields are normalized to trim whitespace
- **Duplicate Alias Detection**: ⚠️ **PENDING** - No prevention of multiple relationships with same alias to same target
- **Circular Alias Detection**: ⚠️ **PENDING** - No prevention of circular aliasing references
- **Polymorphic Alias Validation**: ✅ **DONE** - Validates polymorphic relationships can use aliasing with `for` property

#### Test Infrastructure
- ✅ **DONE** - Comprehensive test coverage for basic aliasing (31 passing tests)
- **PENDING** - Test models with polymorphic inverse aliased relationships
- **PENDING** - Expected SQL output for aliased tables and constraints
- Missing validation tests for complex polymorphic aliasing configurations
- No entity field path tests with aliased relationships

### Impacted Files
```
# External Dependencies (morphe-go repo) - ✅ COMPLETED
morphe-go/pkg/yaml/model_relation.go     - ✅ Added 'aliased' field to ModelRelation
morphe-go/pkg/yaml/entity_relation.go    - ✅ Added 'aliased' field to EntityRelation  
morphe-go/pkg/yamlops/relation.go        - ✅ Aliasing detection functions
morphe-go/pkg/yamlops/relation_test.go   - ✅ Tests for aliasing detection
morphe-go/pkg/yaml/model.go              - ✅ Basic aliasing validation logic
morphe-go/pkg/yaml/entity.go             - ✅ Basic aliased relationship validation

# This Repository (plugin-morphe-psql-types)
pkg/compile/compile_models.go            - Core aliasing compilation logic
pkg/compile/compile_models_test.go       - Test coverage for aliased relationships
pkg/compile/compile_entities.go          - Aliased field path resolution
pkg/compile/compile_entities_test.go     - Entity aliasing test coverage
pkg/compile/naming.go                    - Aliased table/column/constraint naming
testdata/registry/                       - Aliased relationship examples
testdata/ground-truth/                   - Expected aliased SQL output
```

## 2. Implementation Plan

**Statuses:**

- Ready: Ready to implement (in theory), practical prerequisites not checked
- In progress: Implementing
- Review: Ready for review by human
- Done: Successfully implemented & reviewed by human

### Task 1: Complete morphe-go Aliasing Support
**Status**: ✅ **Done**
**Scope**: BATCHED cross-repo changes - all morphe-go modifications for aliasing
**Cross-Repo Dependency**: Single comprehensive morphe-go update + one dependency sync
**Sub-tasks**:
- **1.1** ✅ Add `aliased` field to ModelRelation and EntityRelation structures
- **1.2** ✅ Implement aliasing detection functions in yamlops (`IsRelationAliased`, `GetRelationTargetName`)
- **1.3** ✅ Add aliasing validation to Model and Entity validation logic
- **1.4** ✅ Update deep clone methods for new aliased field
- **1.5** ✅ Comprehensive morphe-go aliasing testing (31 passing tests)

**Acceptance Criteria**: ✅ All completed
- ✅ `aliased` field properly added to relation structures with YAML marshaling
- ✅ Aliasing detection functions work correctly for all relationship types
- ✅ Validation ensures aliased targets exist in registry
- ✅ All existing tests continue to pass + new aliasing tests added

### Task 2: ForOne/HasOne Aliasing Support
**Status**: Ready
**Scope**: Generate aliased foreign key columns and constraints for direct relationships
**Implementation Details**:
- **Files to modify**: 
  - `compile_models.go`: `getColumnsForModelRelations` (line 238), `getForeignKeysForModelRelations` (line 292)
  - Pattern: `targetModelName := yamlops.GetRelationTargetName(relatedModelName, modelRelation.Aliased)`
**Acceptance Criteria**:
- Uses relationship name for column naming (backward compatibility): `{relation_name}_id`
- Foreign key constraints reference correct aliased target table
- HasOne inverse relationships validate through aliased ForOne relationships
- Proper constraint naming maintains relationship name (not aliased target)

### Task 3: ForMany/HasMany Aliasing Support  
**Status**: Ready
**Scope**: Generate aliased junction tables and constraints for many-to-many relationships
**Implementation Details**:
- **Files to modify**:
  - `compile_models.go`: `getTablesForModelManyRelations` (lines 408-528)
  - Pattern: Use GetRelationTargetName for model lookups, maintain relationship name for table/column naming
**Acceptance Criteria**:
- Junction tables use relationship name: `{source_table}_{relation_name}` (backward compatibility)
- Junction table columns use relationship name: `{relation_name}_id`
- Foreign key constraints reference correct aliased target tables
- HasMany inverse relationships validate through aliased ForMany relationships

### Task 4: Polymorphic Aliasing Support
**Status**: Ready
**Scope**: Support aliasing for polymorphic relationships (ForOnePoly, ForManyPoly, Has*Poly)
**Acceptance Criteria**:
- ForOnePoly generates `{alias_name}_type` and `{alias_name}_id` columns
- ForManyPoly junction tables use aliased naming (`{source_table}_{alias_name}`)
- HasOnePoly/HasManyPoly validate through aliased polymorphic relationships
- **Advanced**: Polymorphic inverse aliasing (`HasOnePoly` with `through` + `aliased`)
- Polymorphic constraints and indexes use aliased names

### Task 5: Entity-Level Aliasing Support
**Status**: Ready
**Scope**: Support referencing aliased model relationships in entities + entity-level aliasing
**Implementation Details**:
- **Files to modify**:
  - `compile_entities.go`: `traverseRelationshipChain` (line 198), `setupJoinsForRegularRelationships`, `addJoinClause`
  - Pattern: Resolve aliased targets during traversal while maintaining relationship names for SQL
**Acceptance Criteria**:
- Entity field paths correctly traverse through aliased model relationships
- Entity relationships can themselves be aliased (`aliased` field in EntityRelation)
- Entity view joins use correct aliased target tables
- Proper validation for entity-level aliased references

### Task 6: Aliasing Validation & Error Handling
**Status**: Ready
**Scope**: Comprehensive validation for aliased relationship configurations
**Acceptance Criteria**:
- Validates `aliased` property contains valid model/entity names
- **Advanced**: Validates polymorphic inverse aliasing pattern (`HasOnePoly` + `through` + `aliased`)
- Prevents duplicate aliases pointing to same target within same model/entity
- Clear error messages for invalid aliased references
- Prevents circular aliasing references

### Task 7: Integration Testing & Documentation
**Status**: Ready
**Scope**: End-to-end testing with comprehensive aliasing examples
**Acceptance Criteria**:
- Complete aliased relationship examples in testdata (all relationship types)
- **Advanced**: Polymorphic inverse aliasing examples and validation
- Generated SQL matches expected aliased schema (tables, columns, constraints)
- All relationship + aliasing combinations tested at model and entity levels
- Integration test covers full aliasing pipeline

## 3. Technical Possibilities & Trade-offs

### Alias Naming Strategy
**Options**:
1. Direct alias substitution: `{alias_name}_id` (Chosen)
2. Prefixed aliasing: `alias_{alias_name}_id`
3. Suffixed aliasing: `{alias_name}_alias_id`

**Trade-offs**:
- **Chosen**: Clean, intuitive naming that matches relationship alias directly
- **Prefixed**: Clear alias indication but verbose
- **Suffixed**: Distinguishable but less intuitive

### Junction Table Naming
**Options**:
1. `{source_table}_{alias_name}` (Chosen)
2. `{source_table}_{alias_name}_{target_table}`
3. `{source_table}_to_{alias_name}`

**Trade-offs**:
- **Chosen**: Consistent with current ForMany pattern, concise
- **With target**: More explicit but potentially very long names
- **With preposition**: Clear but breaks existing conventions

### Aliased Target Resolution
**Options**:
1. Runtime target resolution from `aliased` field (Chosen)
2. Preprocessing to expand aliases into direct references
3. Dual storage of both alias and target names

**Trade-offs**:
- **Chosen**: Flexible, maintains alias semantics throughout pipeline
- **Preprocessing**: Simpler processing but loses alias context
- **Dual storage**: Redundant but faster lookups

### Validation Strategy
**Options**:
1. Registry-based validation during model loading (Chosen)
2. Lazy validation during SQL generation
3. Two-phase validation (structure + references)

**Trade-offs**:
- **Chosen**: Early error detection, consistent with existing patterns
- **Lazy**: Simpler but errors discovered late in process
- **Two-phase**: Comprehensive but complex implementation

### Entity Field Path Resolution
**Options**:
1. Alias-aware path traversal during entity compilation (Chosen)
2. Alias expansion before field path processing
3. Virtual field mapping for aliased relationships

**Trade-offs**:
- **Chosen**: Maintains full alias semantics in entity definitions
- **Expansion**: Simpler but loses alias context in entities
- **Virtual mapping**: Flexible but adds complexity layer

## 4. Active Working Log

| Task | Status | Commit/PR | Notes |
|------|--------|-----------|-------|
| **Task 1** Complete morphe-go Aliasing | ✅ **Done** | morphe-go v0.0.0-20250824082856 | **COMPLETED**: All tests passing, includes polymorphic inverse validation |
| **Task 2** ForOne/HasOne Aliasing | Ready | - | **READY**: Specific code locations identified |
| **Task 3** ForMany/HasMany Aliasing | Ready | - | **READY**: Junction table handling identified |
| **Task 4** Polymorphic Aliasing | Ready | - | **READY**: Includes polymorphic inverse pattern |
| **Task 5** Entity-Level Aliasing | Ready | - | **READY**: Join traversal updates identified |
| **Task 6** Aliasing Validation | Partial | - | **PARTIAL**: Polymorphic inverse ✅, duplicates/circular pending |
| **Task 7** Integration Testing | Ready | - | **READY**: Requires implementation first |

## 5. Status Updates

### 2025-01-02 - Initial Analysis Completed
- Analyzed existing relationship processing architecture
- Identified gap between current implementation and aliasing requirements
- Created comprehensive implementation plan with 7 focused tasks
- Established cross-repo dependency strategy with batched morphe-go changes

### 2025-01-02 - Task 1 Complete ✅
- **COMPLETED**: All morphe-go aliasing support with comprehensive testing
- Added `aliased` field to ModelRelation and EntityRelation structures
- Implemented aliasing detection functions in yamlops (IsRelationAliased, GetRelationTargetName)
- Added comprehensive validation including polymorphic inverse validation
- Added whitespace normalization for all string fields including aliased
- **DISCOVERED**: Advanced polymorphic inverse aliasing pattern for type generation

### 2025-01-15 - Integration Audit Complete
- **AUDITED**: Complete scan of plugin-morphe-psql-types for integration points
- **IDENTIFIED**: 6 key functions requiring updates across 2 main files
- **PATTERN**: Consistent fix using `yamlops.GetRelationTargetName(relationName, relation.Aliased)`
- **STRATEGY**: Use aliased target for model/entity lookups, maintain relationship name for SQL objects
- **COMPATIBILITY**: All changes maintain backward compatibility by preserving column/table naming

---

## Relationship Aliasing Specification Reference

### Supported Aliased Relation Types

* **ForOne Aliased**: Direct relationship with alias - current model belongs to one instance of target model via alias
  * Implementation: Aliased foreign key column on model table
  * Example: Person ForOne WorkContact (aliased: ContactInfo), Person ForOne HomeContact (aliased: ContactInfo)

* **ForMany Aliased**: Many-to-many relationship with alias via junction table
  * Implementation: Aliased junction table with proper naming
  * Example: Person ForMany WorkProjects (aliased: Project), Person ForMany HobbyProjects (aliased: Project)

* **HasOne Aliased**: Inverse of ForOne with alias - current model can have one model referring to it via alias
  * Implementation: Validates through aliased ForOne relationship
  * Example: ContactInfo HasOne WorkPerson (through: WorkContact), ContactInfo HasOne HomePerson (through: HomeContact)

* **HasMany Aliased**: Inverse of ForMany with alias - current model can have multiple models referring to it via alias
  * Implementation: Validates through aliased ForMany relationship
  * Example: Project HasMany WorkPeople (through: WorkProjects), Project HasMany HobbyPeople (through: HobbyProjects)

* **Polymorphic Aliased**: All polymorphic types support aliasing with `for` and `through` properties
  * Implementation: Aliased polymorphic columns and junction tables
  * Example: Comment ForOnePoly WorkCommentable (aliased: Commentable, for: [WorkDocument, WorkTask])

* **🆕 Polymorphic Inverse Aliased**: Advanced pattern for semantic field naming in generated types
  * **Pattern**: `HasOnePoly`/`HasManyPoly` with `through` + `aliased` properties
  * **Purpose**: Create semantic field names for inverse polymorphic relationships
  * **Example**: 
    ```yaml
    # Comment model defines polymorphic interface
    Comment:
      related:
        Commentable:
          type: ForOnePoly
          for: [Post, Article, Task]
    
    # Post model creates semantic field name
    Post:
      related:
        Note:  # ← Field name in generated types (semantic alias)
          type: HasOnePoly
          through: Commentable  # ← References polymorphic relationship name
          aliased: Comment      # ← Actual model type
    ```
  * **Generated Types**:
    ```go
    type Post struct {
        Note *Comment  // Field: "Note", Type: Comment
    }
    
    type Task struct {
        StatusUpdate *Comment  // Field: "StatusUpdate", Type: Comment
    }
    ```
  * **Validation Requirements**:
    - `aliased` model must exist in registry
    - `aliased` model must have polymorphic relationship named `through`
    - Current model must be in `for` list of that polymorphic relationship
  * **SQL Impact**: Minimal - `HasOnePoly` doesn't create database columns, only inverse semantics

### Configuration Format

```yaml
# Model-level aliasing
related:
  WorkContact:
    type: ForOne
    aliased: ContactInfo
  HomeContact:
    type: ForOne
    aliased: ContactInfo
  
  WorkProject:
    type: ForMany
    aliased: Project
  HobbyProject:
    type: ForMany
    aliased: Project

  # 🆕 Advanced: Polymorphic inverse aliasing
  Note:
    type: HasOnePoly
    through: Commentable
    aliased: Comment
  StatusUpdate:
    type: HasOnePoly
    through: Commentable
    aliased: Comment

# Entity-level aliasing (referencing aliased model relationships)
related:
  WorkContact:
    type: ForOne
    # Inherits aliasing from model relationship
    
  # Entity-level aliasing itself
  PrimaryContact:
    type: ForOne
    aliased: ContactInfo
```

### Database Implementation

- **Aliased Columns**: Foreign key columns use alias names (`work_contact_id`, `home_contact_id`)
- **Aliased Tables**: Junction tables use alias names (`people_work_projects`, `people_hobby_projects`)  
- **Aliased Constraints**: Foreign keys and unique constraints incorporate alias names
- **Target Resolution**: Runtime resolution of aliased target models for validation and SQL generation
- **Polymorphic Inverse**: No additional SQL schema - purely for type generation semantics 