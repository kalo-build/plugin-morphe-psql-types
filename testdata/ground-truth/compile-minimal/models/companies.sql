-- Table definition for companies

CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.companies (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	tax_id TEXT NOT NULL
);

-- Indices
CREATE UNIQUE INDEX IF NOT EXISTS idx_companies_name ON public.companies ("name");

