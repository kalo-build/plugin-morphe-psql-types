# Polymorphism Implementation Plan

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
- **Database integrity**: All polymorphic tables must maintain referential integrity via constraints
- **Planning**: After task completion, reevaluate upcoming tasks and revisit adding new / updating existing

## 1. Codebase ⇄ Spec Gap Analysis

### Current State
- **Supported Relations**: `ForOne`, `ForMany`, `HasOne`, `HasMany` (non-polymorphic)
- **Junction Tables**: Implemented for `ForMany` relationships
- **Column Generation**: Foreign key columns for `ForOne` relationships  
- **Test Coverage**: Comprehensive tests for existing relationship types

### Missing Components for Polymorphism

#### Core Logic Modules
- `pkg/compile/compile_models.go`: No polymorphic relationship handling
- **External Dependency**: `morphe-go/pkg/yamlops/`: Missing polymorphic relation detection functions  
- `pkg/psqldef/`: No polymorphic column types (type + id fields)

#### Database Schema Components
- **ForOnePoly**: Missing polymorphic type + id column generation
- **ForManyPoly**: Missing polymorphic junction table creation
- **HasOnePoly/HasManyPoly**: Missing inverse relationship handling
- **Polymorphic constraints**: No validation for `for` property in YAML

#### Test Infrastructure
- No test models with polymorphic relationships
- No expected SQL output for polymorphic tables
- Missing validation tests for polymorphic configurations

#### Documentation & Examples
- No polymorphic model examples in `testdata/`
- Missing polymorphic schema generation in `ground-truth/`

### Impacted Files
```
# External Dependencies (morphe-go repo) - BATCHED CHANGES
morphe-go/pkg/yamlops/relation.go      - Polymorphic relation detection functions  
morphe-go/pkg/yamlops/relation_test.go - Tests for polymorphic detection
morphe-go/pkg/yaml/model_relation.go   - Add 'for' and 'through' fields
morphe-go/pkg/yaml/entity_relation.go  - Add 'for' and 'through' fields (if needed)
morphe-go/pkg/yaml/entity.go           - Add polymorphic types to validation
morphe-go/pkg/yaml/model.go            - Add polymorphic validation logic

# This Repository (plugin-morphe-psql-types)  
pkg/compile/compile_models.go          - Core polymorphic compilation logic
pkg/compile/compile_models_test.go     - Test coverage for new features
pkg/compile/naming.go                  - Polymorphic column/table naming
pkg/psqldef/                          - Polymorphic column definitions
testdata/registry/                    - Polymorphic model examples  
testdata/ground-truth/                - Expected polymorphic SQL output
```

## 2. Implementation Plan

**Statuses:**

- Ready: Ready to implement (in theory), practical prerequisites not checked
- In progress: Implementing
- Review: Ready for review by human
- Done: Successfully implemented & reviewed by human

### Task 1: Complete morphe-go Polymorphism Support  
**Status**: Done
**Scope**: BATCHED cross-repo changes - all morphe-go modifications for polymorphism
**Cross-Repo Dependency**: Single comprehensive morphe-go update + one dependency sync
**Sub-tasks**:
- ✅ **1.1** Polymorphic relation detection functions (yamlops) 
- ✅ **1.2** ModelRelation structure extensions (`for`, `through` fields)
- ✅ **1.3** EntityRelation validation updates (add polymorphic types)  
- ✅ **1.4** Model/Entity validation for polymorphic properties
- ✅ **1.5** Deep clone method updates
- ✅ **1.6** Comprehensive morphe-go testing

### Task 2: ForOnePoly Column Generation  
**Status**: Done
**Scope**: Generate polymorphic type + id columns for ForOnePoly relationships
**Acceptance Criteria**:
- Creates `{relation_name}_type` TEXT column 
- Creates `{relation_name}_id` TEXT column
- No foreign key constraints (polymorphic nature)
- Proper column naming conventions

### Task 3: ForManyPoly Junction Tables
**Status**: Done
**Scope**: Generate polymorphic junction tables for ForManyPoly relationships  
**Acceptance Criteria**:
- ✅ Junction table with polymorphic type + id columns
- ✅ Proper unique constraints on (source_id, target_type, target_id)
- ✅ No foreign key constraints on polymorphic columns

### Task 4: HasOnePoly/HasManyPoly Inverse Relations
**Status**: Done
**Scope**: Handle inverse polymorphic relationships using `through` property
**Acceptance Criteria**:
- ✅ Validates `through` property points to valid polymorphic relation
- ✅ No additional column generation (relies on forward relation)
- ✅ Proper relationship validation

### Task 5: Polymorphic Validation & Error Handling
**Status**: Done
**Scope**: Comprehensive validation for polymorphic relationship configurations
**Acceptance Criteria**:
- ✅ Validates `for` property contains valid model names
- ✅ Validates `through` property for Has* relationships  
- ✅ Clear error messages for misconfigurations
- ✅ Prevents circular polymorphic references

### Task 6: Integration Testing & Documentation  
**Status**: Done
**Scope**: End-to-end testing with comprehensive polymorphic examples
**Acceptance Criteria**:
- ✅ Complete polymorphic model examples in testdata (ForOnePoly + ForManyPoly)
- ✅ Generated SQL matches expected polymorphic schema
- ✅ All relationship combinations tested at model level

### Task 7: Entity-Level Polymorphic View Generation
**Status**: Done  
**Scope**: Generate appropriate entity views for polymorphic relationships
**Acceptance Criteria**:
- ✅ ForOnePoly entity views include raw polymorphic columns (type + id)
- ✅ ForManyPoly entity views are simple (no junction table materialization)
- ✅ HasOnePoly/HasManyPoly inverse relationships not materialized in entity views
- ✅ Entity view generation maintains existing patterns

### Task 8: Entity-Level Integration Testing
**Status**: Done
**Scope**: Comprehensive testing of polymorphic entity view generation
**Acceptance Criteria**:
- ✅ Entity views generated correctly for all polymorphic relationship types
- ✅ Integration tests cover entity-level polymorphic scenarios  
- ✅ Generated entity SQL matches expected polymorphic view schema
- ✅ No regressions in existing entity view functionality

## 3. Technical Possibilities & Trade-offs

### Column Naming Strategy
**Options**:
1. `{relation_name}_type` + `{relation_name}_id` (Chosen)
2. `polymorphic_type` + `polymorphic_id` (Generic)

**Trade-offs**:
- **Chosen**: Clear relation association, supports multiple polymorphic relations per model
- **Alternative**: Simpler but limits to one polymorphic relation per model

### ID Storage Type
**Options**:
1. TEXT columns for polymorphic IDs (Chosen - per spec)
2. Separate typed columns per target model
3. JSONB column with structured data

**Trade-offs**:
- **Chosen**: Flexible, spec-compliant, no schema coupling
- **Alternatives**: Type safety vs. flexibility trade-off

### Junction Table Strategy  
**Options**:
1. Dedicated polymorphic junction tables (Chosen)
2. Single polymorphic relations table
3. Hybrid approach with table-per-relation-type

**Trade-offs**:
- **Chosen**: Clear separation, follows existing ForMany pattern
- **Single table**: Centralized but harder to query efficiently
- **Hybrid**: Complex implementation, unclear benefits

### Migration Strategy
**Options**:
1. Additive-only changes (Chosen)
2. Schema versioning with migrations
3. Breaking changes with major version bump

**Trade-offs**:
- **Chosen**: Backward compatible, safe for existing systems
- **Versioning**: More complex but enables schema evolution
- **Breaking**: Simpler implementation but adoption barrier

## 4. Active Working Log

| Task | Status | Commit/PR | Notes |
|------|--------|-----------|-------|
| **Task 1** Complete morphe-go Polymorphism | Done | morphe-go batched | **SYNCED**: All morphe-go changes complete and synced |
| **Task 2** ForOnePoly Column Generation | Done | ✅ | **COMPLETE**: ForOnePoly columns generated correctly |
| **Task 3** ForManyPoly Junction Tables | Done | ✅ | **COMPLETE**: Polymorphic junction tables with proper constraints |
| **Task 4** HasOnePoly/HasManyPoly Relations | Done | ✅ | **COMPLETE**: Inverse polymorphic relationships with through validation |
| **Task 5** Polymorphic Validation | Done | ✅ | **COMPLETE**: Comprehensive validation with clear error messages |
| **Task 6** Integration Testing | Done | ✅ | **MODEL-LEVEL COMPLETE**: ForOnePoly + ForManyPoly integration tests passing |
| **Task 7** Entity-Level Polymorphic Views | Done | ✅ | **COMPLETE**: Entity views correctly handle polymorphic relationships |
| **Task 8** Entity-Level Integration Testing | Done | ✅ | **COMPLETE**: Comprehensive entity-level polymorphic testing |

## 5. Status Updates

### 2025-06-23 - Initial Analysis Completed
- Analyzed existing codebase architecture
- Identified gap between current implementation and polymorphism spec
- Created comprehensive implementation plan with 6 focused tasks
- Established TDD workflow with human review checkpoints

### 2025-06-23 - Cross-Repository Dependency Identified
- **Critical Discovery**: yamlops package is in separate `morphe-go` repository
- Updated implementation plan to reflect cross-repo dependencies
- Task 1 requires changes to `morphe-go/pkg/yamlops/relation.go` + manual dependency sync
- All polymorphic detection functions must be implemented in external dependency first

### 2025-06-23 - Task 1 Progress: Complete morphe-go Polymorphism Support
- ✅ **Sub-task 1.1**: Polymorphic relation detection functions (yamlops)
  - RED-GREEN-REFACTOR cycle complete with testify assertions
  - 5 compositional functions: `IsRelationPoly`, `IsRelationPolyFor`, `IsRelationPolyHas`, `IsRelationPolyOne`, `IsRelationPolyMany`
  - All tests pass: 9/9 test functions in yamlops package
- ✅ **Sub-task 1.2**: ModelRelation structure extensions (`for`, `through` fields)
- ✅ **Sub-task 1.3**: EntityRelation validation updates (add polymorphic types)
- ✅ **Sub-task 1.4**: Model/Entity polymorphic validation logic
- ✅ **Sub-task 1.5**: Deep clone method updates
- ✅ **Sub-task 1.6**: Comprehensive morphe-go testing
- ✅ **BATCHED APPROACH**: All morphe-go changes complete, ready for single dependency sync

### 2025-06-23 - Task 1 COMPLETE: All morphe-go Polymorphism Support
- ✅ **All Sub-tasks Complete**: Successfully implemented comprehensive polymorphic support in morphe-go
- 🔧 **Data Structures**: Extended ModelRelation and EntityRelation with `for` and `through` fields
- 🛡️ **Validation**: Added polymorphic type validation with proper error handling for missing properties
- 🧪 **Testing**: All existing tests pass + new polymorphic tests added and passing
- 📦 **Batched Changes**: Single comprehensive update prevents multiple dependency syncs
- ✅ **SYNCED**: Complete morphe-go implementation synced to plugin-morphe-psql-types

### 2025-06-23 - Task 2 COMPLETE: ForOnePoly Column Generation
- ✅ **TDD RED-GREEN-REFACTOR**: Followed proper TDD cycle with failing test → implementation → refactor
- 🏗️ **Column Generation**: Successfully generates `{relation_name}_type` and `{relation_name}_id` TEXT columns
- 🚫 **No Foreign Keys**: Correctly skips foreign key constraint generation for polymorphic relationships
- 🐍 **Naming Conventions**: Proper snake_case column naming from camelCase relation names
- 🧪 **Test Coverage**: Two comprehensive tests covering basic and edge cases (long relation names)
- 🔧 **Clean Implementation**: Added polymorphic logic without breaking existing functionality

### 2025-06-23 - Task 3 COMPLETE: ForManyPoly Junction Tables
- ✅ **TDD RED-GREEN-REFACTOR**: Followed proper TDD cycle with failing test → implementation → refactor
- 🔍 **Root Cause Analysis**: Identified that ForManyPoly was incorrectly matching regular ForMany pattern
- 🛠️ **Fix Applied**: Added `!yamlops.IsRelationPoly()` check to exclude polymorphic relations from regular ForMany processing
- 🏗️ **Junction Tables**: Successfully generates polymorphic junction tables with type + id columns
- 🔐 **Proper Constraints**: Unique constraints on (source_id, target_type, target_id) with no foreign key constraints on polymorphic columns
- 🧪 **Test Coverage**: Comprehensive test covering junction table structure, columns, constraints, and foreign keys
- 🔧 **Clean Implementation**: Added polymorphic junction table logic without breaking existing functionality

### 2025-06-23 - Task 4 COMPLETE: HasOnePoly/HasManyPoly Inverse Relations
- ✅ **TDD RED-GREEN-REFACTOR**: Followed proper TDD cycle with failing validation tests → implementation → passing tests
- 🔍 **No Additional Columns**: HasOnePoly/HasManyPoly relationships correctly don't generate database columns (rely on forward relations)
- 🛡️ **Through Validation**: Added cross-model validation that ensures `through` property references valid polymorphic relations
- 🏗️ **Registry-Based Lookup**: Implemented proper validation that searches across all models in registry for through relationships
- 🐛 **Initial Bug Fix**: Fixed validation logic that was incorrectly looking for through relations only in same model
- 🧪 **Test Coverage**: Added comprehensive tests for both valid scenarios and error cases (invalid through properties)
- 🔧 **Clean Implementation**: Added validation without breaking existing functionality, all existing tests continue to pass

### 2025-06-23 - Task 5 COMPLETE: Polymorphic Validation & Error Handling
- ✅ **TDD RED-GREEN-REFACTOR**: Followed proper TDD cycle with comprehensive failing validation tests → implementation → all tests passing
- 🔍 **For Property Validation**: Validates that all models referenced in `for` property exist in the registry
- 🛡️ **Missing Property Detection**: Clear error messages for missing `for` properties in ForOnePoly/ForManyPoly relationships
- 🚫 **Empty Array Validation**: Detects and reports when `for` property is empty (must have at least one target model)
- 🔄 **Circular Reference Detection**: Implemented DFS-based cycle detection to prevent circular polymorphic references
- 🏗️ **Registry Integration**: All validation logic properly integrates with the registry for cross-model validation
- 🧪 **Comprehensive Test Coverage**: Added 7 new validation tests covering all edge cases and error scenarios
- 🐛 **Test Fix**: Fixed existing LongRelationName test that was missing required model in registry
- 🔧 **Enhanced Error Messages**: All validation errors provide clear, actionable error messages for developers
- 🎯 **Detailed Circular Reference Detection**: Shows full cycle path and actionable guidance for resolving circular polymorphic references

### 2025-06-23 - Task 6 Progress: Model-Level Integration Testing Complete
- ✅ **ForOnePoly Integration**: Added Comment model with polymorphic relationship to Person and Company
- ✅ **ForManyPoly Integration**: Added Tag model with polymorphic many-to-many relationship via junction table  
- ✅ **HasOnePoly/HasManyPoly Integration**: Added inverse polymorphic relationships to Person and Company models
- 🔧 **Junction Table Discovery**: Identified actual naming pattern is `tag_taggables.sql` (not `tags_taggable.sql`)
- 📊 **Comprehensive Testing**: Integration test now covers all polymorphic relationship types at model level
- 🏗️ **Database Schema**: Generated SQL correctly includes polymorphic columns, junction tables, and proper constraints
- ✅ **Integration Tests Passing**: All model-level polymorphic functionality verified through end-to-end testing

### 2025-06-23 - New Tasks Added: Entity-Level Polymorphism Support
- 📋 **Task 7 Added**: Entity-Level Polymorphic View Generation with clear acceptance criteria
- 📋 **Task 8 Added**: Entity-Level Integration Testing for comprehensive polymorphic entity coverage
- 🎯 **Entity Strategy Defined**: Simple approach including raw polymorphic columns in entity views
- 🔄 **Next Phase**: Ready to implement entity-level polymorphic support following established patterns

---

## Polymorphism Specification Reference

### Supported Polymorphic Relation Types

* **ForOnePoly**: Polymorphic "belongs to" relationship - current model belongs to one instance of multiple possible model types
  * Implementation: Single polymorphic reference columns on model table
  * Example: Comment ForOnePoly Commentable (Post, Article, Video)

* **ForManyPoly**: Polymorphic "many-to-many" relationship via junction table  
  * Implementation: Junction table with polymorphic columns
  * Example: Tag ForManyPoly Taggable (Post, Product, User)

* **HasOnePoly**: Inverse of ForOnePoly - current model can have one model of different types referring to it
  * Implementation: Configured via `through` property, no additional columns
  * Example: Post HasOnePoly Comment (through Commentable)

* **HasManyPoly**: Inverse of ForManyPoly - current model can have multiple models of different types referring to it
  * Implementation: Configured via `through` property, no additional columns  
  * Example: Post HasManyPoly Tag (through Taggable)

### Configuration Format

```yaml
# ForOnePoly - Comment can belong to Post, Article, or Video
related:
  Commentable:
    type: ForOnePoly
    for:
      - Post
      - Article  
      - Video

# HasOnePoly - Post can have Comments referring to it
related:
  Comment:
    type: HasOnePoly
    through: Commentable
```

### Database Implementation

- **Polymorphic IDs**: Stored as TEXT, round-tripped as strings
- **Type Field**: Indicates which model type the relationship targets
- **No Foreign Keys**: Polymorphic nature prevents traditional FK constraints
- **Unique Constraints**: Applied to prevent duplicate polymorphic relationships

### Polymorphic Aliasing Support

As of 2025-01-15, polymorphic relationships fully support aliasing through the general aliasing feature:
- **ForOnePoly/ForManyPoly**: Can use `aliased` property to reference different target models
- **HasOnePoly/HasManyPoly**: Support the advanced pattern with `through` + `aliased` for semantic field naming
- See `aliasing.md` for complete aliasing documentation and the polymorphic inverse aliasing pattern 