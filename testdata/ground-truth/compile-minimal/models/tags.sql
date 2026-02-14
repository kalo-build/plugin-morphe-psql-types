-- Table definition for tags

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.tags (
	color TEXT NOT NULL,
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

-- Indices
CREATE UNIQUE INDEX IF NOT EXISTS idx_tags_name ON public.tags ("name");

