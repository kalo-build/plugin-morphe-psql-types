-- Table definition for morphe_structures

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.morphe_structures (
	id SERIAL PRIMARY KEY,
	"type" TEXT NOT NULL,
	"data" JSONB NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indices
CREATE INDEX IF NOT EXISTS idx_morphe_structures_type ON public.morphe_structures ("type");
CREATE INDEX IF NOT EXISTS idx_morphe_structures_data ON public.morphe_structures USING GIN ("data");

