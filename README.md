# plugin-morphe-psql-types

Generates PostgreSQL DDL schema from Morphe schema definitions (`KA:MO1:YAML1`). Produces `CREATE TABLE`, `CREATE VIEW`, indexes, foreign keys, and enum lookup tables with seed data.

## What it generates

| Morphe artifact | SQL output                                                                    |
|-----------------|-------------------------------------------------------------------------------|
| **Model**       | `CREATE TABLE` with columns, foreign keys, indexes, unique constraints        |
| **Enum**        | Lookup table with `INSERT` seed data                                          |
| **Structure**   | Standard `morphe_structures` table with JSONB storage (optional)              |
| **Entity**      | `CREATE OR REPLACE VIEW` with `SELECT` / `LEFT JOIN`                          |
| **Relationships** | Foreign key columns, indexes, junction tables for many-to-many              |

### Example output

**Model table** (`people.sql`):

```sql
CREATE TABLE IF NOT EXISTS public.people (
    first_name TEXT NOT NULL,
    id SERIAL PRIMARY KEY,
    last_name TEXT NOT NULL,
    nationality_id INTEGER NOT NULL,
    company_id INTEGER NOT NULL,
    CONSTRAINT fk_people_nationality_id FOREIGN KEY (nationality_id)
        REFERENCES public.nationalities (id) ON DELETE CASCADE,
    CONSTRAINT fk_people_company_id FOREIGN KEY (company_id)
        REFERENCES public.companies (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_people_nationality_id ON public.people (nationality_id);
CREATE INDEX IF NOT EXISTS idx_people_company_id ON public.people (company_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_people_first_name_last_name
    ON public.people (first_name, last_name);
```

**Enum lookup table** (`nationalities.sql`):

```sql
CREATE TABLE IF NOT EXISTS public.nationalities (
    id SERIAL PRIMARY KEY,
    key TEXT NOT NULL,
    value TEXT NOT NULL,
    value_type TEXT NOT NULL,
    UNIQUE (key)
);

INSERT INTO public.nationalities (key, value, value_type) VALUES ('DE', 'German', 'String');
INSERT INTO public.nationalities (key, value, value_type) VALUES ('FR', 'French', 'String');
INSERT INTO public.nationalities (key, value, value_type) VALUES ('US', 'American', 'String');
```

**Entity view** (`person_entities.sql`):

```sql
CREATE OR REPLACE VIEW public.person_entities AS
SELECT
    contact_infos.email,
    people.id,
    people.last_name,
    people.nationality
FROM public.people
LEFT JOIN public.contact_infos
    ON people.id = contact_infos.id;
```

### Relationship handling

| Relationship type | SQL output                                                |
|-------------------|-----------------------------------------------------------|
| `ForOne`          | Foreign key column + index                                |
| `ForMany`         | Junction table with composite unique constraint           |
| `HasOne`          | Foreign key column + index                                |
| `HasMany`         | Junction table with composite unique constraint           |
| Polymorphic       | `_type TEXT` + `_id` columns, composite unique constraint |

### Type mappings

| Morphe type     | PostgreSQL type | BigSerial variant |
|-----------------|-----------------|-------------------|
| `UUID`          | `UUID`          | `UUID`            |
| `AutoIncrement` | `SERIAL`        | `BIGSERIAL`       |
| `String`        | `TEXT`          | `TEXT`            |
| `Integer`       | `INTEGER`       | `INTEGER`         |
| `Float`         | `DOUBLE PRECISION` | `DOUBLE PRECISION` |
| `Boolean`       | `BOOLEAN`       | `BOOLEAN`         |
| `Time`          | `TIMESTAMPTZ`   | `TIMESTAMPTZ`     |
| `Date`          | `DATE`          | `DATE`            |
| `Protected`     | `TEXT`          | `TEXT`            |
| `Sealed`        | `TEXT`          | `TEXT`            |

## Input / output

| Direction | Format         | Store suggestion | Description                          |
|-----------|----------------|------------------|--------------------------------------|
| Input     | `KA:MO1:YAML1` | `KA_MO_YAML`   | Morphe registry (models, enums, structures, entities) |
| Output    | `KA:MO1:PSQL1` | `KA_MO_PSQL`   | PostgreSQL DDL `.sql` files          |

Output is organized into subdirectories: `enums/`, `models/`, `structures/`, `entities/`.
Table names are snake_case and pluralized.

## Configuration

| Key                  | Type    | Default    | Description                                               |
|----------------------|---------|------------|-----------------------------------------------------------|
| `orderedMigrations`  | boolean | `true`     | Prefix output files with numeric order (e.g., `001_`)     |
| `structures.Schema`  | string  | `"public"` | PostgreSQL schema name                                    |
| `structures.UseBigSerial` | boolean | `false` | Use `BIGSERIAL` instead of `SERIAL` for auto-increment    |
| `structures.EnablePersistence` | boolean | `true` | Generate the `morphe_structures` table                   |

The `Schema` and `UseBigSerial` options also apply to models, enums, and entities.

## Pipeline context

This plugin generates the **base schema** DDL from the current Morphe definitions.
For incremental schema **migrations** between versions, use
[`plugin-morphediff-psql`](../plugin-morphediff-psql) instead.

```yaml
stores:
  KA_MO_YAML:
    format: "KA:MO1:YAML1"
    type: "localFileSystem"
    options:
      path: "./morphe"

  KA_MO_PSQL:
    format: "KA:MO1:PSQL1"
    type: "localFileSystem"
    options:
      path: "./schema"

plugins:
  "@kalo-build/plugin-morphe-psql-types":
    version: "v1.0.0"
    inputs:
      morphe:
        format: "KA:MO1:YAML1"
        store: "KA_MO_YAML"
    output:
      format: "KA:MO1:PSQL1"
      store: "KA_MO_PSQL"
    config:
      orderedMigrations: true
```

## Project structure

```
plugin-morphe-psql-types/
├── cmd/plugin/             # WASM entry point
├── pkg/
│   ├── compile/            # Compilation pipeline
│   │   ├── compile.go      # MorpheToPSQL entry point
│   │   ├── compile_models.go
│   │   ├── compile_enums.go
│   │   ├── compile_structures.go
│   │   ├── compile_entities.go
│   │   ├── naming.go       # Snake_case, pluralization, identifier abbreviation
│   │   ├── dependency_sort.go  # Table creation order
│   │   ├── cfg/            # Configuration structs
│   │   ├── hook/           # Extensibility hooks
│   │   └── write/          # Table and view writers
│   ├── psqldef/            # PostgreSQL definition types (Table, View, Index, ForeignKey, etc.)
│   └── typemap/            # Morphe → PostgreSQL type mappings
├── testdata/
│   ├── registry/           # Sample Morphe registry input
│   └── ground-truth/       # Expected SQL output for integration tests
├── dist/                   # WASM output
└── plugin.yaml             # Kalo plugin manifest
```

## Building

```bash
# Native binary
go build ./cmd/plugin

# WASM (for Kalo CLI)
GOOS=wasip1 GOARCH=wasm go build -o dist/plugin.wasm cmd/plugin/main.go
```

## Testing

```bash
go test ./...
```
