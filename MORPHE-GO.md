# Morphe-Go Entity Aliasing Validation Implementation Guide

## Overview

This document provides implementation instructions for adding entity aliasing validation support to morphe-go. The plugin-morphe-psql-types has already implemented aliasing support for SQL generation, but morphe-go's entity validation needs to be updated to understand aliased relationships in entity field paths.

## Problem Statement

### Current Behavior
When an entity field references a path through an aliased relationship (e.g., `Person.WorkContact.Email`), the entity validation fails with an error like:
```
morphe entity PersonProfile field workEmail references unknown model: WorkContact in path Person.WorkContact.Email
```

This happens because the validation tries to find a model named "WorkContact" but doesn't know it's an alias for "Contact".

### Root Cause
The entity validation in morphe-go occurs before the compilation phase, so it doesn't have access to the alias resolution logic that's already implemented in plugin-morphe-psql-types.

## Required Implementation

### 1. Core Changes to Entity Validation

#### 1.1 Update Field Path Validation Logic

The entity field path validation needs to resolve aliases when traversing relationships. Here's the algorithm:

```go
// Pseudo-code for the updated validation logic
func validateEntityFieldPath(path string, models map[string]Model) error {
    segments := strings.Split(path, ".")
    if len(segments) < 2 {
        return fmt.Errorf("invalid field path: %s", path)
    }
    
    // Start with the root model
    currentModelName := segments[0]
    currentModel, exists := models[currentModelName]
    if !exists {
        return fmt.Errorf("model %s not found in field path %s", currentModelName, path)
    }
    
    // Traverse through relationships
    for i := 1; i < len(segments)-1; i++ {
        relationName := segments[i]
        relation, exists := currentModel.Related[relationName]
        if !exists {
            return fmt.Errorf("relationship %s not found in model %s", 
                relationName, currentModel.Name)
        }
        
        // CRITICAL: Resolve the aliased target using yamlops
        targetModelName := yamlops.GetRelationTargetName(relationName, relation.Aliased)
        
        // Load the target model
        nextModel, exists := models[targetModelName]
        if !exists {
            return fmt.Errorf("aliased target model %s not found for relationship %s in model %s", 
                targetModelName, relationName, currentModel.Name)
        }
        
        currentModel = nextModel
        currentModelName = targetModelName
    }
    
    // Validate the final field exists
    fieldName := segments[len(segments)-1]
    if _, exists := currentModel.Fields[fieldName]; !exists {
        return fmt.Errorf("field %s not found in model %s (reached via path %s)", 
            fieldName, currentModelName, path)
    }
    
    return nil
}
```

#### 1.2 Integration Point

This logic should be integrated into the `Entity.Validate()` method, specifically in the section that validates entity field types that contain dots (field paths).

### 2. Special Cases to Handle

#### 2.1 Polymorphic Relationships

Polymorphic relationships with aliases need special handling:

```go
// For polymorphic relationships, check if the alias contains a dot
// Example: "Comment.Commentable" where Commentable is the inverse
if yamlops.IsRelationPoly(relation.Type) && strings.Contains(relation.Aliased, ".") {
    // Extract just the model name part
    parts := strings.Split(relation.Aliased, ".")
    targetModelName = parts[0]
    // Note: The inverse relationship validation should happen separately
}
```

#### 2.2 Direct Relationship References

The current validation correctly prevents direct references to relationships in entity fields. This should continue to work:

```go
// This should still be invalid:
fields:
  myRelation:
    type: Model.RelationshipName  # Invalid - can't reference relationship directly
```

### 3. Error Message Guidelines

Error messages should be clear and helpful:

```go
// Good error messages that distinguish between different failure modes:

// When relationship doesn't exist
"relationship WorkContact not found in model Person"

// When aliased target doesn't exist
"aliased target model Contact not found for relationship WorkContact in model Person"

// When field doesn't exist in the resolved model
"field Email not found in model Contact (reached via path Person.WorkContact.Email)"

// For polymorphic relationships
"cannot traverse through polymorphic relationship Commentable in path Post.Commentable.Text"
```

### 4. Implementation Checklist

- [ ] Locate the `Entity.Validate()` method in morphe-go
- [ ] Find the section that validates entity field types containing dots
- [ ] Import or ensure access to `yamlops.GetRelationTargetName()`
- [ ] Implement the alias resolution logic in the field path traversal
- [ ] Update error messages to be more informative about aliasing
- [ ] Add unit tests for aliased field paths (see Testing section)
- [ ] Ensure backward compatibility - non-aliased paths must continue to work

### 5. Testing

#### 5.1 Basic Aliasing Test Case

```yaml
# Models
---
name: Person
fields:
  ID:
    type: AutoIncrement
related:
  WorkContact:
    type: ForOne
    aliased: Contact
  PersonalContact:
    type: ForOne  
    aliased: Contact
---
name: Contact
fields:
  ID:
    type: AutoIncrement
  Email:
    type: String
  Phone:
    type: String

# Entity - should validate successfully
---
name: PersonProfile
fields:
  id:
    type: Person.ID
  workEmail:
    type: Person.WorkContact.Email      # Should resolve to Contact.Email
  personalPhone:
    type: Person.PersonalContact.Phone  # Should resolve to Contact.Phone
```

#### 5.2 Error Cases to Test

```yaml
# Test 1: Aliased target doesn't exist
related:
  BadAlias:
    type: ForOne
    aliased: NonExistentModel  # Should fail validation

# Test 2: Field doesn't exist in aliased target
fields:
  badField:
    type: Person.WorkContact.NonExistentField  # Should fail

# Test 3: Attempting to traverse through polymorphic (should fail appropriately)
fields:
  invalid:
    type: Post.Comments.Author  # Where Comments is polymorphic
```

#### 5.3 Complex Traversal Test

```yaml
# Test multi-hop traversal through aliases
name: Company
related:
  CEO:
    type: ForOne
    aliased: Person

# Should be able to traverse: Company.CEO.WorkContact.Email
```

### 6. Dependencies

- The implementation requires access to `yamlops.GetRelationTargetName()` function
- No changes needed to the YAML structure or schema
- Must maintain backward compatibility with existing entity definitions

### 7. Expected Outcome

After implementing these changes:

1. Entities can reference fields through aliased relationships
2. Validation errors clearly indicate whether issues are with relationships, aliases, or fields
3. All existing entity validations continue to work
4. The validation phase and compilation phase use the same alias resolution logic

### 8. Additional Notes

- The plugin-morphe-psql-types already has working aliasing implementation for the compilation phase
- The key is to use the same `yamlops.GetRelationTargetName()` function for consistency
- Entity validation happens in `entity.Validate()` before compilation, so it needs its own alias resolution
- The validation should only resolve aliases for traversal - it should not change how the compiled SQL is generated (that's already handled correctly)

## Questions or Issues?

If you encounter any issues or need clarification:

1. Check the existing aliasing implementation in plugin-morphe-psql-types (pkg/compile/compile_models.go and compile_entities.go)
2. Look for uses of `yamlops.GetRelationTargetName()` for examples
3. The test files in plugin-morphe-psql-types show expected behavior for aliased relationships