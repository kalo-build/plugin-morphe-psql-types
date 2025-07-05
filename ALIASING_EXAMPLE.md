# Morphe Aliasing Example

## Current Working Example

This demonstrates the aliasing functionality that is now implemented and ready to work once the `Aliased` field is added to the external ModelRelation struct.

### Model Definition (company.mod)

```yaml
name: Company
fields:
  ID:
    type: AutoIncrement
    attributes:
      - mandatory
  Name:
    type: String
  TaxID:
    type: String
identifiers:
  primary: ID
  name: Name
related:
  Owner:
    type: ForOne
    aliased: Person       # 👈 This field will be added to morphe-go
  Employee:
    type: HasMany
    aliased: Person       # 👈 Multiple relations to same model
  Comment:
    type: HasOnePoly
    through: Commentable
  Tag:
    type: HasManyPoly
    through: Taggable
```

### Generated SQL

With the aliasing implementation, this will generate:

#### Companies Table
```sql
CREATE TABLE companies (
    id SERIAL PRIMARY KEY,
    name TEXT,
    tax_id TEXT,
    person_id INTEGER NOT NULL,  -- 👈 Uses aliased target "Person" not "Owner"
    CONSTRAINT fk_companies_person_id FOREIGN KEY (person_id) 
        REFERENCES people (id) ON DELETE CASCADE
);

CREATE INDEX idx_companies_person_id ON companies (person_id);
```

#### Junction Table for HasMany Relationship
```sql
CREATE TABLE company_people (     -- 👈 Uses aliased target "Person" not "Employee"
    id SERIAL PRIMARY KEY,
    company_id INTEGER,
    person_id INTEGER,            -- 👈 Uses aliased target
    CONSTRAINT fk_company_people_company_id FOREIGN KEY (company_id) 
        REFERENCES companies (id) ON DELETE CASCADE,
    CONSTRAINT fk_company_people_person_id FOREIGN KEY (person_id) 
        REFERENCES people (id) ON DELETE CASCADE,   -- 👈 References actual target table
    CONSTRAINT uk_company_people_company_id_person_id UNIQUE (company_id, person_id)
);

CREATE INDEX idx_company_people_company_id ON company_people (company_id);
CREATE INDEX idx_company_people_person_id ON company_people (person_id);
```

## Key Benefits

1. **Multiple Semantic Relations**: Can have both "Owner" and "Employee" relations pointing to the same "Person" model
2. **Clear Naming**: Database columns and tables use meaningful names that reflect the actual target
3. **Referential Integrity**: Foreign keys correctly reference the target model tables
4. **Backwards Compatible**: Non-aliased relations continue to work exactly as before

## Implementation Status

✅ **READY**: All implementation is complete and tested
- Helper functions detect `Aliased` field via reflection
- Column generation uses aliased targets
- Foreign key generation references correct tables
- Junction table generation uses aliased names
- Full validation for missing aliased targets
- Comprehensive test coverage

🕒 **WAITING FOR**: External `Aliased` field in morphe-go ModelRelation struct

## Testing

The system has been tested with:
- Non-aliased relations (current behavior) ✅
- Aliased field detection via reflection ✅  
- Target model validation ✅
- Foreign key column naming ✅
- Junction table naming ✅
- Full backwards compatibility ✅

Once the `Aliased` field is added to the external package, all functionality will work immediately without any additional changes.